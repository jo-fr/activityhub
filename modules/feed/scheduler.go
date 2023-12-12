package feed

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func Schedule(lc fx.Lifecycle, logger *log.Logger, h *Handler) error {
	s := gocron.NewScheduler(time.UTC)

	s.RegisterEventListeners(
		gocron.WhenJobReturnsError(func(jobName string, err error) {
			logger.
				WithField("jobName", jobName).
				Errorf("scheduler job failed. Err %s", err.Error())
		}),
	)

	_, err := s.Every(20).Second().Name("feed fetcher").Do(h.FetchFeed)
	if err != nil {
		return errors.Wrap(err, "failed to setup scheduler job")
	}

	registerHooks(lc, s, logger)

	return nil
}

// registerHooks for uber fx
func registerHooks(lc fx.Lifecycle, scheduler *gocron.Scheduler, logger *log.Logger) {
	lc.Append(
		fx.Hook{
			OnStop: func(context.Context) error {
				logger.Info("stopping scheduler jobs")

				scheduler.Stop()
				return nil
			},

			OnStart: func(context.Context) error {
				logger.Info("starting scheduler job")
				scheduler.StartAsync()
				return nil
			},
		})
}

func (h *Handler) FetchFeed() error {
	ctx := context.Background()

	sources, err := h.store.ListSourceFeeds(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get source feeds")
	}

	for _, source := range sources {
		if err := h.FetchSourceFeedUpdates(ctx, source); err != nil {
			return errors.Wrap(err, "failed to fetch source feed")
		}
	}
	return nil
}
