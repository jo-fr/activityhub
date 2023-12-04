package feed

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jo-fr/activityhub/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func Schedule(lc fx.Lifecycle, logger *log.Logger) error {
	s := gocron.NewScheduler(time.UTC)

	job, err := s.Every(1).Second().Name("feed fetcher").Do(func() error {

		logger.Info(time.Now().Second())

		return errors.New("error")

	})
	if err != nil {
		return errors.Wrap(err, "failed to setup scheduler job")
	}
	job.RegisterEventListeners(
		gocron.WhenJobReturnsError(func(jobName string, err error) {
			logger.
				WithField("jobName", jobName).
				Errorf("scheduler job failed. Err %s", err.Error())
		}),
	)

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
