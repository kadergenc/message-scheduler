package scheduler

import (
	"context"
	"message-scheduler/internal/port"
	"message-scheduler/log"
	"sync"
	"time"
)

type ScheduledJob struct {
	job      port.Job
	interval time.Duration
	ticker   *time.Ticker
	stop     chan bool
}

type SimpleScheduler struct {
	jobs   []*ScheduledJob
	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewSimpleScheduler() *SimpleScheduler {
	return &SimpleScheduler{
		jobs: make([]*ScheduledJob, 0),
	}
}

func (s *SimpleScheduler) ScheduleJob(job port.Job, interval time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	scheduledJob := &ScheduledJob{
		job:      job,
		interval: interval,
		stop:     make(chan bool, 1),
	}

	s.jobs = append(s.jobs, scheduledJob)

	log.Logger.Info().
		Str("job_name", job.Name()).
		Dur("interval", interval).
		Msg("Job scheduled successfully")
}

func (s *SimpleScheduler) Start(ctx context.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.ctx, s.cancel = context.WithCancel(ctx)

	log.Logger.Info().Int("job_count", len(s.jobs)).Msg("Starting scheduler...")

	for _, scheduledJob := range s.jobs {
		s.wg.Add(1)
		go s.runJob(scheduledJob)
	}

	log.Logger.Info().Msg("Scheduler started successfully")
}

func (s *SimpleScheduler) runJob(scheduledJob *ScheduledJob) {
	defer s.wg.Done()

	scheduledJob.ticker = time.NewTicker(scheduledJob.interval)
	defer scheduledJob.ticker.Stop()

	s.executeJob(scheduledJob.job)

	for {
		select {
		case <-s.ctx.Done():
			log.Logger.Info().
				Str("job_name", scheduledJob.job.Name()).
				Msg("Stopping job due to context cancellation")
			return

		case <-scheduledJob.stop:
			log.Logger.Info().
				Str("job_name", scheduledJob.job.Name()).
				Msg("Stopping job due to stop signal")
			return

		case <-scheduledJob.ticker.C:
			s.executeJob(scheduledJob.job)
		}
	}
}

func (s *SimpleScheduler) executeJob(job port.Job) {
	start := time.Now()

	log.Logger.Info().
		Str("job_name", job.Name()).
		Msg("Executing scheduled job")

	jobCtx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	err := job.Execute(jobCtx)

	duration := time.Since(start)

	if err != nil {
		log.Logger.Error().
			Err(err).
			Str("job_name", job.Name()).
			Dur("duration", duration).
			Msg("Job execution failed")
	} else {
		log.Logger.Info().
			Str("job_name", job.Name()).
			Dur("duration", duration).
			Msg("Job executed successfully")
	}
}

func (s *SimpleScheduler) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.cancel != nil {
		log.Logger.Info().Msg("Stopping scheduler...")

		s.cancel()

		for _, scheduledJob := range s.jobs {
			select {
			case scheduledJob.stop <- true:
			default:

			}
		}

		done := make(chan struct{})
		go func() {
			s.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Logger.Info().Msg("All scheduled jobs stopped successfully")
		case <-time.After(10 * time.Second):
			log.Logger.Warn().Msg("Timeout waiting for scheduled jobs to stop")
		}
	}

	return nil
}
