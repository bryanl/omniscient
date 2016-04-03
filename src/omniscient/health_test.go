package omniscient

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	var mu sync.Mutex
	healthStatus := true
	check := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return healthStatus
	}

	updateCheckInterval := func(h *Health) error {
		h.checkInterval = 1 * time.Millisecond
		return nil
	}

	h, err := NewHealth(
		HealthCheckOption(check),
		updateCheckInterval)
	assert.NoError(t, err)

	h.Start()
	defer h.Stop()

	assert.True(t, h.IsOK(), "health check status")

	mu.Lock()
	healthStatus = false
	mu.Unlock()

	<-h.statusUpdated
	assert.False(t, h.IsOK(), "health check status")
}

func TestHealthCannotStartIfStarted(t *testing.T) {
	updateCheckInterval := func(h *Health) error {
		h.checkInterval = 1 * time.Millisecond
		return nil
	}

	h, err := NewHealth(updateCheckInterval)
	assert.NoError(t, err)

	err = h.Start()
	assert.NoError(t, err)

	<-h.stateChange

	err = h.Start()
	assert.Error(t, err)
}

func TestHealthCannotStopIfNotStarted(t *testing.T) {
	updateCheckInterval := func(h *Health) error {
		h.checkInterval = 1 * time.Millisecond
		return nil
	}

	h, err := NewHealth(updateCheckInterval)
	assert.NoError(t, err)

	err = h.Stop()
	assert.Error(t, err)
}
