package cleanup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	cleanAPIKeysQuery = `DELETE FROM api_key WHERE expire_at<=$1`
	cleanPasteQuery   = `DELETE FROM paste_info WHERE expire_at<=$1`
)

type Worker struct {
	pool        *pgxpool.Pool
	log         *zap.Logger
	cleanUpTick time.Duration
	wg          *sync.WaitGroup
	quite       chan struct{}
}

func NewWorker(pool *pgxpool.Pool, log *zap.Logger, cleanUpTick time.Duration) *Worker {
	return &Worker{
		pool:        pool,
		log:         log,
		cleanUpTick: cleanUpTick,
		wg:          &sync.WaitGroup{},
		quite:       make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		ticker := time.NewTicker(w.cleanUpTick)
		defer ticker.Stop()

		defer w.wg.Done()

		for {
			select {
			case <-ticker.C:
				err := w.cleanUp(ctx)
				if err != nil {
					w.log.Error(fmt.Sprintf("DB clean up error - %v", err))
					continue
				}
			case <-w.quite:
				w.log.Debug("quite")
				return
			case <-ctx.Done():
				w.log.Error(fmt.Sprintf("exit by context - %v", ctx.Err()))
				return
			}
		}
	}()
}

func (w *Worker) cleanUp(ctx context.Context) error {
	now := time.Now()
	_, err := w.pool.Exec(ctx, cleanAPIKeysQuery, now)
	if err != nil {
		return err
	}

	_, err = w.pool.Exec(ctx, cleanPasteQuery, now)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop(ctx context.Context) {
	close(w.quite)

	done := make(chan struct{})

	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		w.log.Error(fmt.Sprintf("exit by context - %v", ctx.Err()))
		return
	}
}
