package feed

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/jo-fr/activityhub/backend/modules/activitypub"

	"github.com/jo-fr/activityhub/backend/modules/feed/internal/repository"
	"github.com/jo-fr/activityhub/backend/modules/feed/model"
	"github.com/jo-fr/activityhub/backend/pkg/errutil"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"github.com/jo-fr/activityhub/backend/pkg/pubsub"
	"github.com/jo-fr/activityhub/backend/pkg/store"
	"github.com/jo-fr/activityhub/backend/pkg/util"
	"github.com/jo-fr/activityhub/backend/pkg/util/httputil"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	repository.Module,
	fx.Provide(NewHandler),
	fx.Invoke(ScheduleFeedFetcher),
)

// define errors
var (
	ErrFeedAlreadyExists = errutil.NewError(errutil.TypeAlreadyExists, "feed already exists")
	ErrFeedNotFound      = errutil.NewError(errutil.TypeNotFound, "feed not found")
)

type Handler struct {
	parser      *gofeed.Parser
	store       *store.Store[repository.FeedRepository]
	activitypub *activitypub.Handler
	pubsub      *pubsub.Client
	log         *log.Logger
	scheduler   *gocron.Scheduler
}

func NewHandler(s *store.Store[repository.FeedRepository], log *log.Logger, activitypub *activitypub.Handler, pubsub *pubsub.Client) *Handler {
	h := &Handler{
		parser:      gofeed.NewParser(),
		store:       s,
		log:         log,
		activitypub: activitypub,
		pubsub:      pubsub,
	}

	return h
}

func (h *Handler) AddNewFeed(ctx context.Context, feedurl string) (feed model.Feed, err error) {
	err = h.store.Execute(ctx, func(e *repository.FeedRepository) error {
		sanatizedFeedURL, err := httputil.SanitizeURL(feedurl)
		if err != nil {
			return errors.Wrap(err, "failed to sanitize url")
		}

		_, err = e.GetFeedWithFeedURL(sanatizedFeedURL)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get feed")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFeedAlreadyExists
		}

		sourceFeed, err := h.parser.ParseURLWithContext(sanatizedFeedURL, ctx)
		if err != nil {
			return errors.Wrap(err, "failed to parse feed")
		}

		title := sourceFeed.Title
		description := strings.ReplaceAll(util.RemoveHTMLTags(sourceFeed.Description), "\n", " ")
		authorsSlice := util.Map(sourceFeed.Authors, func(item *gofeed.Person, index int) string {
			if item == nil {
				return ""
			}

			return item.Name
		})
		author := strings.Join(authorsSlice, ", ")

		accountUsername := usernameFromFeedTitle(title)
		name := fmt.Sprintf("%s ActivityHub Bot", title)
		summary := fmt.Sprintf("This is the ActivityHub Bot of %s. This is NOT an official account and is not related with the owners of the posted content. Posting entries of RSS feed.", title)

		ctxWithTx := e.GetCtxWithTx(ctx)
		account, err := h.activitypub.CreateAccount(ctxWithTx, accountUsername, name, summary)
		if err != nil {
			return errors.Wrap(err, "failed to create account")
		}

		feed = model.Feed{
			Name:        title,
			Type:        model.FeedTypeRSS,
			FeedURL:     sanatizedFeedURL,
			HostURL:     sourceFeed.Link,
			Author:      author,
			Description: util.TrimStringLength(description, 500),
			ImageURL:    util.FromPointer(sourceFeed.Image).URL,
			AccountID:   account.ID,
		}

		feed, err = e.CreateFeed(feed)
		if err != nil {
			return errors.Wrap(err, "failed to create feed")
		}

		// fetch feed again to get full model
		feed, err = e.GetFeedWithID(feed.ID)
		if err != nil {
			return errors.Wrap(err, "failed to create feed")
		}

		if err := scheduleNewJob(ctx, h.scheduler, h.log, feed.Name, h.FetchFeed(context.Background(), feed)); err != nil {
			return errors.Wrap(err, "failed to schedule new job")
		}

		return nil
	})

	return feed, err
}

func (h *Handler) FetchFeedUpdates(ctx context.Context, feed model.Feed) error {
	return h.store.Execute(ctx, func(e *repository.FeedRepository) error {
		sourceFeed, err := h.parser.ParseURLWithContext(feed.FeedURL, ctx)
		if err != nil {
			return errors.Wrap(err, "failed to parse feed")
		}

		items := sourceFeed.Items
		if len(items) < 1 {
			return errors.New("no items found in feed")
		}

		newestItem := items[0]

		latestStatus, err := e.GetLatestStatusFromFeed(feed.AccountID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get latest status")
		}

		if !util.FromPointer(newestItem.PublishedParsed).After(latestStatus.CreatedAt) {
			return nil
		}

		post := builtPost(newestItem.Title, newestItem.Description, newestItem.Link)
		status := model.Status{
			Content:   post,
			AccountID: feed.AccountID,
		}

		status, err = e.CreateStatus(status)
		if err != nil {
			return errors.Wrap(err, "failed to create status")
		}

		if err := h.pubsub.Publish(ctx, pubsub.TopicOutbox, status); err != nil {
			return errors.Wrap(err, "failed to publish status")
		}

		return nil
	})

}

func (h *Handler) ListFeeds(ctx context.Context, offset int, limit int) (totalCount int, feeds []model.Feed, err error) {
	var sources []model.Feed
	err = h.store.Execute(ctx, func(e *repository.FeedRepository) error {

		count, err := e.FeedCount()
		if err != nil {
			return errors.Wrap(err, "failed to count feeds")
		}
		totalCount = int(count)

		sources, err = e.ListFeeds(offset, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get feeds")
		}
		return nil
	})

	return totalCount, sources, err
}

func (h *Handler) GetFeed(ctx context.Context, id string) (feed model.Feed, err error) {
	err = h.store.Execute(ctx, func(e *repository.FeedRepository) error {
		feed, err = e.GetFeedWithID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrFeedNotFound
			}
			return errors.Wrap(err, "failed to get feed")
		}
		return nil
	})

	return feed, err
}

func (h *Handler) GetFeedWithUsername(ctx context.Context, username string) (feed model.Feed, err error) {
	err = h.store.Execute(ctx, func(e *repository.FeedRepository) error {

		ctx = e.GetCtxWithTx(ctx)
		account, err := h.activitypub.GetActor(ctx, username)
		if err != nil {
			if errors.Is(err, activitypub.ErrActorNotFound) {
				return err
			}
			return errors.Wrap(err, "failed to get actor")
		}

		feed, err = e.GetFeedWithAccountID(account.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrFeedNotFound
			}
			return errors.Wrap(err, "failed to get feed")
		}
		return nil
	})

	return feed, err
}

func (h *Handler) ListFeedStatus(ctx context.Context, id string, offset int, limit int) (totalCount int, statuses []model.Status, err error) {
	err = h.store.Execute(ctx, func(e *repository.FeedRepository) error {
		feed, err := e.GetFeedWithID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrFeedNotFound
			}
			return errors.Wrap(err, "failed to get feed")
		}

		count, err := e.StatusCount(feed.AccountID)
		if err != nil {
			return errors.Wrap(err, "failed to count statuses")
		}
		totalCount = int(count)

		statuses, err = e.ListStatusFromAccount(feed.AccountID, offset, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get statuses")
		}
		return nil
	})

	return totalCount, statuses, err
}

func builtPost(title string, description string, link string) string {

	// sanatize
	title = "<strong>" + util.RemoveHTMLTags(title) + "</strong><br/>"
	description = strings.ReplaceAll(description, "\n", " ")
	description = util.RemoveHTMLTags(description)

	// Mastodon only displays the first 27 characters of a link
	if len(link) > 27 {
		link = fmt.Sprintf("<a href=\"%s\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">%s...</a>", link, link[:27])
	} else {
		link = fmt.Sprintf("<a href=\"%s\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">%s</a>", link, link)
	}

	content := title + description
	content = util.TrimStringLength(content, 500-30)

	return fmt.Sprintf("<p>%s</br>%s</p>", content, link)

}

// usernameFromFeedTitle creates a username by removing punctuation and replacing spaces with underscores
// e.g. "Hello, World!" -> "Hello_World"
func usernameFromFeedTitle(title string) string {
	title = removePunctuation(title)
	title = strings.ReplaceAll(title, " ", "_")
	title = CamelToSnake(title)
	return fmt.Sprintf("%s_activityhub", title)
}

// removePunctuation removes all punctuation from a string
// e.g. "Hello, World!" -> "Hello World"
func removePunctuation(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(s, "")
}

// CamelToSnake converts a string from camel case to snake case and replaces spaces with underscores
// e.g. "HelloWorld" -> "hello_world"
func CamelToSnake(camelCase string) string {
	// Use regular expression to match uppercase letters and add underscore before them
	snakeCase := regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(camelCase, "${1}_${2}")

	snakeCase = strings.ReplaceAll(snakeCase, " ", "_")
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}
