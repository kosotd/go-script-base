package timer

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"fmt"
	"time"
)

var times = sync.Map{}

func Start(label string) {
	if label == "" {
		label = "Timer"
	}

	times.Store(label, time.Now())
}

func Stop(label string) time.Duration {
	if label == "" {
		label = "Timer"
	}

	startTime, ok := times.Load(label)

	if !ok {
		panic(fmt.Sprintf("call stop before start for label: %s", label))
	}

	elapsed := time.Since(startTime.(time.Time))
	log.Info(fmt.Sprintf("%.3f seconds elapsed for %s", elapsed.Seconds(), label))

	return elapsed
}
