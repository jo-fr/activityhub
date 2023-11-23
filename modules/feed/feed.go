package feed

import (
	"context"
	"strings"

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
	parser *gofeed.Parser
	store  *store.Store
}

func NewHandler(store *store.Store, log *log.Logger) *Handler {
	h := &Handler{
		parser: gofeed.NewParser(),
		store:  store,
	}

	return h
}

func (h *Handler) AddNewSourceFeed(ctx context.Context, feedurl string) (model.SourceFeed, error) {

	sanatizedURL, err := httputil.SanitizeURL(feedurl)
	if err != nil {
		return model.SourceFeed{}, err
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
		return model.SourceFeed{}, err
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
		return model.SourceFeed{}, err
	}

	return sourceFeed, nil

}
