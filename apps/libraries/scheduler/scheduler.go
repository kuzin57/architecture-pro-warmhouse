package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type Scheduler struct {
	scheduler *gocron.Scheduler
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewScheduler(lc fx.Lifecycle) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Scheduler{scheduler: gocron.NewScheduler(time.UTC), ctx: ctx, cancel: cancel}

	s.scheduler.SetMaxConcurrentJobs(3, gocron.WaitMode)

	lc.Append(fx.Hook{
		OnStart: s.Run,
		OnStop:  s.Stop,
	})

	return s
}

func (s *Scheduler) Run(ctx context.Context) error {
	s.scheduler.StartAsync()

	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.cancel()
	s.scheduler.Stop()

	return nil
}

func (s *Scheduler) AddJob(ctx context.Context, schedule string, id uuid.UUID, job func(ctx context.Context) error) error {
	cronJob, err := s.scheduler.Cron(schedule).Do(func() error {
		err := job(s.ctx)
		if err != nil {
			log.Println("failed to run job", err)
		}

		return err
	})
	if err != nil {
		return fmt.Errorf("failed to add job: %w", err)
	}

	cronJob.Tag(id.String())

	return nil
}

func (s *Scheduler) RemoveJob(ctx context.Context, id uuid.UUID) error {
	s.scheduler.RemoveByTag(id.String())

	return nil
}
