// sheduler - пакет, содержащий объекты
// для запуска задач по расписанию
package scheduler

import (
	"context"
	"pora/internal/infrastructure/logger"
	"time"
)

// Sheduler раз в определенный интервал вызывает
// переданную в конструкторе функцию
type Scheduler struct {
	interval time.Duration
	job      func(ctx context.Context) error
}

// NewScheduler создает и возвращает новый объект Scheduler
func NewScheduler(interval time.Duration,
	job func(ctx context.Context) error) *Scheduler {
	return &Scheduler{
		interval: interval,
		job:      job,
	}
}

// Run запускает бесконечный цикл выполнения задач по расписанию
func (s *Scheduler) Run(ctx context.Context) {

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	if err := s.job(ctx); err != nil {
		logger.Log.Warn("Ошибка при выполнении задачи: ", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := s.job(ctx); err != nil {
				logger.Log.Warn("Ошибка при выполнении задачи: ", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
