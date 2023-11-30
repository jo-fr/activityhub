package feed

import (
	"context"
	"fmt"
	"strings"

	"github.com/jo-fr/activityhub/modules/activitypub"
	"github.com/jo-fr/activityhub/modules/feed/internal/store"
	"github.com/jo-fr/activityhub/modules/feed/model"
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
	fx.Invoke(NewHandler),
)

// define errors
var (
	ErrSourceFeedAlreadyExists = errutil.NewError(errutil.TypeAlreadyExists, "source feed already exists")
)

type Handler struct {
	parser      *gofeed.Parser
	store       *store.Store
	activitypub *activitypub.Handler
}

func NewHandler(store *store.Store, log *log.Logger, activitypub *activitypub.Handler) *Handler {
	h := &Handler{
		parser:      gofeed.NewParser(),
		store:       store,
		activitypub: activitypub,
	}

	if err := h.FetchSourceFeed(context.Background(), "fef9533f-751e-49f1-bcf1-6b166167c67b"); err != nil {
		log.Fatal(err)
	}

	if err := h.FetchSourceFeed(context.Background(), "964ff6c2-0ae7-4707-b23b-b6a60fe17aab"); err != nil {
		log.Fatal(err)
	}
	log.Info("fetched source feed")

	return h
}

func (h *Handler) AddNewSourceFeed(ctx context.Context, feedurl string) (model.SourceFeed, error) {

	sanatizedURL, err := httputil.SanitizeURL(feedurl)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to sanitize url")
	}

	_, err = h.store.GetSourceFeedWithFeedURL(ctx, sanatizedURL)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.SourceFeed{}, errors.Wrap(err, "failed to get source feed")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.SourceFeed{}, ErrSourceFeedAlreadyExists
	}

	feed, err := h.parser.ParseURLWithContext(sanatizedURL, ctx)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to parse feed")
	}

	name := feed.Title
	description := strings.ReplaceAll(feed.Description, "\n", " ")
	url := sanatizedURL

	sourceFeed := model.SourceFeed{
		Name:        name,
		Type:        model.SourceFeedTypeRSS,
		URL:         url,
		Description: util.TrimStringLength(description, 500),
	}

	sourceFeed, err = h.store.CreateSourceFeed(ctx, sourceFeed)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to create source feed")
	}

	return sourceFeed, nil

}

func (h *Handler) FetchSourceFeed(ctx context.Context, sourceFeedID string) error {

	sourceFeed, err := h.store.GetSourceFeedWithID(ctx, sourceFeedID)
	if err != nil {
		return errors.Wrap(err, "failed to get source feed")
	}

	feed, err := h.parser.ParseURLWithContext(sourceFeed.URL, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to parse feed")
	}

	items := feed.Items
	if len(items) < 1 {
		return errors.New("no items found in feed")
	}

	newestItem := items[0]

	// todo do something with the code
	_ = builtPost(newestItem.Title, newestItem.Description, newestItem.Link)

	return nil

}

func builtPost(title string, description string, link string) string {

	// sanatize
	title = "<strong>" + util.RemoveHTMLTags(title) + "</strong>"
	description = strings.ReplaceAll(description, "\n", " ")
	description = util.RemoveHTMLTags(description)
	link = fmt.Sprintf("<a href=\"%s\" target=\"_blank\" rel=\"nofollow noopener noreferrer\" translate=\"no\">%s...</a>", link, link[:27])

	content := title + "\n" + description
	content = util.TrimStringLength(content, 500-30)

	return fmt.Sprintf("<p>%s</br>%s</p>", content, link)

}
