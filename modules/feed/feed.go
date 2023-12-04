package feed

import (
	"context"
	"fmt"
	"regexp"
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
	fx.Provide(NewHandler),
	fx.Invoke(Schedule),
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

	return h
}

func (h *Handler) AddNewSourceFeed(ctx context.Context, feedurl string) (model.SourceFeed, error) {

	sanatizedFeedURL, err := httputil.SanitizeURL(feedurl)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to sanitize url")
	}

	_, err = h.store.GetSourceFeedWithFeedURL(ctx, sanatizedFeedURL)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.SourceFeed{}, errors.Wrap(err, "failed to get source feed")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.SourceFeed{}, ErrSourceFeedAlreadyExists
	}

	feed, err := h.parser.ParseURLWithContext(sanatizedFeedURL, ctx)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to parse feed")
	}

	name := feed.Title
	description := strings.ReplaceAll(feed.Description, "\n", " ")
	authorsSlice := util.Map(feed.Authors, func(item *gofeed.Person, index int) string {
		if item == nil {
			return ""
		}

		return item.Name
	})
	author := strings.Join(authorsSlice, ", ")

	accountUsername := UsernameFromSourceFeedTitle(name)
	account, err := h.activitypub.CreateAccount(ctx, accountUsername)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to create account")
	}

	sourceFeed := model.SourceFeed{
		Name:        name,
		Type:        model.SourceFeedTypeRSS,
		FeedURL:     sanatizedFeedURL,
		HostURL:     feed.Link,
		Author:      author,
		Description: util.TrimStringLength(description, 500),
		ImageURL:    util.FromPointer(feed.Image).URL,
		AccountID:   account.ID,
	}

	sourceFeed, err = h.store.CreateSourceFeed(ctx, sourceFeed)
	if err != nil {
		return model.SourceFeed{}, errors.Wrap(err, "failed to create source feed")
	}

	return sourceFeed, nil

}

func (h *Handler) FetchSourceFeedUpdates(ctx context.Context, sourceFeed model.SourceFeed) error {

	feed, err := h.parser.ParseURLWithContext(sourceFeed.FeedURL, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to parse feed")
	}

	items := feed.Items
	if len(items) < 1 {
		return errors.New("no items found in feed")
	}

	newestItem := items[0]

	latestStatus, err := h.store.GetLatestStatusFromSourceFeed(ctx, sourceFeed.AccountID)
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

	status, err = h.store.CreateStatus(ctx, status)
	if err != nil {
		return errors.Wrap(err, "failed to create status")
	}

	// put status into queue

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

func UsernameFromSourceFeedTitle(title string) string {
	title = RemovePunctuation(title)
	title = strings.ReplaceAll(title, " ", "_")
	return title
}

func RemovePunctuation(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(s, "")
}
