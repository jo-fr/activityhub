package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jo-fr/activityhub/backend/modules/feed/internal/repository"
	"github.com/jo-fr/activityhub/backend/modules/feed/model"
	"github.com/jo-fr/activityhub/backend/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func ScheduleFeedFetcher(lc fx.Lifecycle, logger *log.Logger, h *Handler) error {
	ctx := context.Background()
	return h.store.Execute(ctx, func(e *repository.FeedRepository) error {
		s := gocron.NewScheduler(time.UTC)
		s.RegisterEventListeners(
			gocron.WhenJobReturnsError(func(jobName string, err error) {
				logger.
					WithField("jobName", jobName).
					Errorf("scheduler job failed. Err %s", err.Error())
			}),
		)

		sources, err := e.ListFeeds(0, 100000)
		if err != nil {
			return errors.Wrap(err, "failed to get feeds")
		}

		for _, source := range sources {
			if err := scheduleNewJob(ctx, s, logger, source.Name, h.FetchFeed(context.Background(), source)); err != nil {
				return errors.Wrap(err, "failed to schedule new job")
			}
		}

		registerHooks(lc, s, logger)

		h.scheduler = s
		return nil
	})

}

func scheduleNewJob(ctx context.Context, scheduler *gocron.Scheduler, logger *log.Logger, name string, job func() error) error {
	jobName := getSchedulerJobName(name)

	_, err := scheduler.Every(20).Second().Name(jobName).Do(job)
	if err != nil {
		return errors.Wrapf(err, "failed to setup scheduler job. source name: %s", name)
	}

	logger.Infof("%s successfully scheduled", jobName)
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

func (h *Handler) FetchFeed(ctx context.Context, source model.Feed) func() error {
	return func() error {
		if err := h.FetchFeedUpdates(ctx, source); err != nil {
			return errors.Wrap(err, "failed to fetch feed")
		}
		return nil
	}
}

func getSchedulerJobName(name string) string {
	return fmt.Sprintf("%s_scheduler_job", CamelToSnake(name))
}
