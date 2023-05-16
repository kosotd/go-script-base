package async

import "time"

func NewThrottle(count int, per time.Duration) *Throttle {
	return &Throttle{count, per, nil, nil, nil}
}

func (t *Throttle) Start() chan struct{} {
	t.counter = make(chan struct{}, t.count)
	t.ticker = time.NewTicker(t.per)
	t.done = make(chan struct{})
	t.fill()

	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.fill()
			case <-t.done:
				return
			}
		}
	}()

	return t.counter
}

func (t *Throttle) Stop() {
	defer close(t.counter)
	defer t.ticker.Stop()
	defer close(t.done)
}

func (t *Throttle) fill() {
	defer func() {
		recover()
	}()

	for i := 0; i < t.count; i++ {
		t.counter <- struct{}{}
	}
}

type Throttle struct {
	count   int
	per     time.Duration
	counter chan struct{}
	ticker  *time.Ticker
	done    chan struct{}
}
