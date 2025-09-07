package port

import (
	"context"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
	Name() string
}

type Scheduler interface {
	ScheduleJob(job Job, interval time.Duration)
	Start(ctx context.Context)
	Stop() error
}
