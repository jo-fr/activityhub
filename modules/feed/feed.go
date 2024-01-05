package feed

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/jo-fr/activityhub/modules/activitypub"

	"github.com/jo-fr/activityhub/modules/feed/internal/repository"
	"github.com/jo-fr/activityhub/modules/feed/model"
	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/jo-fr/activityhub/pkg/pubsub"
	"github.com/jo-fr/activityhub/pkg/store"
	"github.com/jo-fr/activityhub/pkg/util"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
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
		summary := fmt.Sprintf("This is the ActivityHub Bot of %s. This is NOT an offical account and is not related with the owners of the posted content. Posting entries of RSS feed.", title)

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

		count, err := e.CountFeeds()
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
				return errutil.NewError(errutil.TypeNotFound, "feed not found")
			}

			return errors.Wrap(err, "failed to get feed")
		}
		return nil
	})

	return feed, err
}

func builtPost(title string, description string, link string) string {

	// sanatize
	title = "<strong>" + util.RemoveHTMLTags(title) + "</strong><br/>"
	description = strings.ReplaceAll(description, "\n", " ")
	description = util.RemoveHTMLTags(description)
	link = fmt.Sprintf("<a href=\"%s\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">%s...</a>", link, link[:27])

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
