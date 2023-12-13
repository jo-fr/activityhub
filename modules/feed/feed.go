package feed

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/feed/internal/store"
	"github.com/jo-fr/activityhub/modules/feed/model"
	"github.com/jo-fr/activityhub/pkg/database"
	"github.com/jo-fr/activityhub/pkg/errutil"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/jo-fr/activityhub/pkg/util"
	"github.com/jo-fr/activityhub/pkg/util/httputil"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	store.Module,
	fx.Provide(NewHandler),
	fx.Invoke(ScheduleFeedFetcher),
)

// define errors
var (
	ErrSourceFeedAlreadyExists = errutil.NewError(errutil.TypeAlreadyExists, "source feed already exists")
)

type Handler struct {
	parser      *gofeed.Parser
	store       *database.Store[*store.FeedRepository]
	activitypub *activitypub.Handler
}

func NewHandler(s *database.Store[*store.FeedRepository], log *log.Logger, activitypub *activitypub.Handler) *Handler {
	h := &Handler{
		parser:      gofeed.NewParser(),
		store:       s,
		activitypub: activitypub,
	}

	return h
}

func (h *Handler) AddNewSourceFeed(ctx context.Context, feedurl string) (sourceFeed model.SourceFeed, err error) {

	err = h.store.Execute(ctx, func(e *store.FeedRepository) error {
		sanatizedFeedURL, err := httputil.SanitizeURL(feedurl)
		if err != nil {
			return errors.Wrap(err, "failed to sanitize url")
		}

		_, err = e.GetSourceFeedWithFeedURL(sanatizedFeedURL)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get source feed")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSourceFeedAlreadyExists
		}

		feed, err := h.parser.ParseURLWithContext(sanatizedFeedURL, ctx)
		if err != nil {
			return errors.Wrap(err, "failed to parse feed")
		}

		title := feed.Title
		description := strings.ReplaceAll(feed.Description, "\n", " ")
		authorsSlice := util.Map(feed.Authors, func(item *gofeed.Person, index int) string {
			if item == nil {
				return ""
			}

			return item.Name
		})
		author := strings.Join(authorsSlice, ", ")

		accountUsername := usernameFromSourceFeedTitle(title)
		name := fmt.Sprintf("%s ActivityHub Bot", title)
		summary := fmt.Sprintf("This is the ActivityHub Bot of %s. This is NOT an offical account and is not related with the owners of the posted content. Posting entries of RSS feed.", title)

		ctx = context.WithValue(ctx, "tx", e.GetTX())
		account, err := h.activitypub.CreateAccount(ctx, accountUsername, name, summary)
		if err != nil {
			return errors.Wrap(err, "failed to create account")
		}

		sourceFeed = model.SourceFeed{
			Name:        title,
			Type:        model.SourceFeedTypeRSS,
			FeedURL:     sanatizedFeedURL,
			HostURL:     feed.Link,
			Author:      author,
			Description: util.TrimStringLength(description, 500),
			ImageURL:    util.FromPointer(feed.Image).URL,
			AccountID:   account.ID,
		}

		sourceFeed, err = e.CreateSourceFeed(sourceFeed)
		if err != nil {
			return errors.Wrap(err, "failed to create source feed")
		}

		return nil
	})

	return sourceFeed, err

}

func (h *Handler) FetchSourceFeedUpdates(ctx context.Context, sourceFeed model.SourceFeed) error {
	return h.store.Execute(ctx, func(e *store.FeedRepository) error {
		feed, err := h.parser.ParseURLWithContext(sourceFeed.FeedURL, ctx)
		if err != nil {
			return errors.Wrap(err, "failed to parse feed")
		}

		items := feed.Items
		if len(items) < 1 {
			return errors.New("no items found in feed")
		}

		newestItem := items[0]

		latestStatus, err := e.GetLatestStatusFromSourceFeed(sourceFeed.AccountID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get latest status")
		}

		if !newestItem.PublishedParsed.After(latestStatus.CreatedAt) {
			return nil
		}

		post := builtPost(newestItem.Title, newestItem.Description, newestItem.Link)
		status := model.Status{
			Content:   post,
			AccountID: sourceFeed.AccountID,
		}

		_, err = e.CreateStatus(status)
		if err != nil {
			return errors.Wrap(err, "failed to create status")
		}

		// put status into queue

		return nil
	})

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

// usernameFromSourceFeedTitle creates a username by removing punctuation and replacing spaces with underscores
// e.g. "Hello, World!" -> "Hello_World"
func usernameFromSourceFeedTitle(title string) string {
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
