package main

import (
	"sync"
	"time"
)

type Buffer[T any] struct {
	Events    []*T
	Mutex     sync.Mutex
	LastFlush time.Time
}

type BufferConfig[T any] struct {
	BatchSize    int
	BatchTimeout time.Duration
	FlusherFunc  func(evt []*T) error
}

type BufferedFlush[T any] struct {
	Buffer Buffer[T]
	Config BufferConfig[T]
}

func (b *BufferedFlush[T]) AddToBuffer(e *T) {
	b.Buffer.Mutex.Lock()
	defer b.Buffer.Mutex.Unlock()

	b.Buffer.Events = append(b.Buffer.Events, e)
	if len(b.Buffer.Events) >= b.Config.BatchSize || time.Since(b.Buffer.LastFlush) >= b.Config.BatchTimeout {
		b.flushUnlocked()
	}
}

func (b *BufferedFlush[T]) flushUnlocked() {
	if len(b.Buffer.Events) == 0 {
		return
	}

	if err := b.Config.FlusherFunc(b.Buffer.Events); err != nil {
		return
	}

	b.Buffer.Events = make([]*T, 0, b.Config.BatchSize)
	b.Buffer.LastFlush = time.Now()
}

func (b *BufferedFlush[T]) FlushBuffer() {
	b.Buffer.Mutex.Lock()
	defer b.Buffer.Mutex.Unlock()
	b.flushUnlocked()
}

func (b *BufferedFlush[T]) HandlePeriodicFlush(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(b.Config.BatchTimeout)
	defer ticker.Stop()

	for range ticker.C {
		b.FlushBuffer()
	}
}

func NewBufferedFlush[T any](config BufferConfig[T]) *BufferedFlush[T] {
	return &BufferedFlush[T]{
		Config: config,
		Buffer: Buffer[T]{
			Events:    make([]*T, 0, config.BatchSize),
			Mutex:     sync.Mutex{},
			LastFlush: time.Now(),
		},
	}
}
