package slogpretty

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

const (
	development = "development"
	production  = "production"
)

// Асинхронный handler для slog
type AsyncHandler struct {
	handler slog.Handler
	ch      chan *slog.Record
	wg      *sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// Создание нового асинхронного хендлера
func NewAsyncHandler(handler slog.Handler, bufferSize int) *AsyncHandler {
	ctx, cancel := context.WithCancel(context.Background())

	ah := &AsyncHandler{
		handler: handler,
		ch:      make(chan *slog.Record, bufferSize),
		wg:      &sync.WaitGroup{},
		ctx:     ctx,
		cancel:  cancel,
	}

	ah.wg.Add(1)
	go ah.worker()

	return ah
}

// Горутинa, обрабатывающая записи логов
func (ah *AsyncHandler) worker() {
	defer ah.wg.Done()
	defer close(ah.ch)

	for {
		select {
		case record := <-ah.ch:
			if record != nil {
				ah.handler.Handle(ah.ctx, *record)
			}
		case <-ah.ctx.Done():
			// Дочищаем буфер перед завершением
			for len(ah.ch) > 0 {
				record := <-ah.ch
				if record != nil {
					ah.handler.Handle(context.Background(), *record)
				}
			}
			return
		}
	}
}

// Обработка входящей записи
func (ah *AsyncHandler) Handle(ctx context.Context, record slog.Record) error {
	// Копируем, чтобы избежать переиспользования slog.Record
	recordCopy := record.Clone()

	select {
	case ah.ch <- &recordCopy:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Если буфер заполнен — ждём до 100мс, иначе дропаем
		timeout := time.NewTimer(100 * time.Millisecond)
		defer timeout.Stop()

		select {
		case ah.ch <- &recordCopy:
			return nil
		case <-timeout.C:
			os.Stderr.WriteString("async logger buffer full, dropping message\n")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Закрытие асинхронного хендлера
func (ah *AsyncHandler) Close() {
	ah.cancel()
	ah.wg.Wait()
}

// Поддержка slog.Handler API
func (ah *AsyncHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ah.handler.Enabled(ctx, level)
}

func (ah *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AsyncHandler{
		handler: ah.handler.WithAttrs(attrs),
		ch:      ah.ch,
		wg:      ah.wg,
		ctx:     ah.ctx,
		cancel:  ah.cancel,
	}
}

func (ah *AsyncHandler) WithGroup(name string) slog.Handler {
	return &AsyncHandler{
		handler: ah.handler.WithGroup(name),
		ch:      ah.ch,
		wg:      ah.wg,
		ctx:     ah.ctx,
		cancel:  ah.cancel,
	}
}

func SetupAsyncLogger(env string) (*slog.Logger, func()) {
	baseHandler := newBaseHandler(env)

	asyncHandler := NewAsyncHandler(baseHandler, 1000)
	logger := slog.New(asyncHandler)

	return logger, func() {
		asyncHandler.Close()
	}
}

func SetupLogger(env string) *slog.Logger {
	return slog.New(newBaseHandler(env))
}

func newBaseHandler(env string) slog.Handler {
	switch env {
	case production:
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case development:
		return setupPrettySlog()
	default:
		return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
}

// Красивый форматер для dev-режима
func setupPrettySlog() slog.Handler {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	return opts.NewPrettyHandler(os.Stdout)
}
