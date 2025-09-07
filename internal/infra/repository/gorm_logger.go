package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	gl "gorm.io/gorm/logger"
)

type GormLogger struct {
	l zerolog.Logger
}

func (gl *GormLogger) LogMode(gl.LogLevel) gl.Interface { return gl }
func (gl *GormLogger) Info(ctx context.Context, msg string, values ...interface{}) {
	ev := gl.l.Info()
	for i, v := range values {
		ev = ev.Any(fmt.Sprint(i), v)
	}
	ev.Msg(msg)
}
func (gl *GormLogger) Warn(ctx context.Context, msg string, values ...interface{}) {
	ev := gl.l.Warn()
	for i, v := range values {
		ev = ev.Any(fmt.Sprint(i), v)
	}
	ev.Msg(msg)
}
func (gl *GormLogger) Error(ctx context.Context, msg string, values ...interface{}) {
	ev := gl.l.Error()
	for i, v := range values {
		ev = ev.Any(fmt.Sprint(i), v)
	}
	ev.Msg(msg)
}
func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()

	if err != nil {
		gl.l.Error().
			Err(err).
			Dur("elapsed", elapsed).
			Str("sql", sql).
			Int64("rows_affected", rowsAffected).
			Msg("GORM SQL execution failed")
	} else {
		gl.l.Info().
			Dur("elapsed", elapsed).
			Str("sql", sql).
			Int64("rows_affected", rowsAffected).
			Msg("GORM SQL executed successfully")
	}
}
