package omniscient

import (
	"errors"
	"sync"
	"time"
)

const (
	defaultCheckInterval = 5 * time.Second
)

// HealthCheckFn is a health check function that returns true for success.
type HealthCheckFn func() bool

// Health is a health checker.
type Health struct {
	checks        []HealthCheckFn
	checkInterval time.Duration
	checkQuit     chan struct{}

	// stateChange is used to to signal when the state has been changed. It's useful
	// for testing and why it is not exported.
	stateChange chan struct{}

	status        bool
	mu            sync.RWMutex
	statusUpdated chan bool
}

// HealthOption is a Healh configuration option.
type HealthOption func(*Health) error

// NewHealth builds an instance of Health.
func NewHealth(opts ...HealthOption) (*Health, error) {
	h := &Health{
		checkInterval: defaultCheckInterval,
		stateChange:   make(chan struct{}, 100000),
		statusUpdated: make(chan bool, 100000),
		status:        true,
	}

	for _, opt := range opts {
		err := opt(h)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

// HealthCheckOption adds a health check.
func HealthCheckOption(hc HealthCheckFn) HealthOption {
	return func(h *Health) error {
		h.checks = append(h.checks, hc)
		return nil
	}
}

// Start starts the health check verification loop.
func (h *Health) Start() error {
	if h.isStarted() {
		return errors.New("health check has already been started")
	}

	go func() {
		ticker := time.NewTicker(h.checkInterval)

		h.mu.Lock()
		h.checkQuit = make(chan struct{})
		h.mu.Unlock()
		h.stateChange <- struct{}{}

		for {
			select {
			case <-ticker.C:
				h.runChecks()
			case <-h.checkQuit:
				ticker.Stop()
				h.checkQuit = nil
			}
		}
	}()

	return nil
}

func (h *Health) isStarted() bool {
	return h.checkQuit != nil
}

// Stop stops the health check verification loop.
func (h *Health) Stop() error {
	if !h.isStarted() {
		return errors.New("health check has not previously been started")
	}

	close(h.checkQuit)
	h.stateChange <- struct{}{}
	return nil
}

// IsOK returns the current health status. True for OK, and false for not.
func (h *Health) IsOK() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.status
}

func (h *Health) runChecks() {
	h.mu.Lock()
	defer h.mu.Unlock()

	ok := true
	for _, check := range h.checks {
		if !check() {
			ok = false
		}
	}

	if ok != h.status {
		h.statusUpdated <- ok
		h.status = ok
	}
}
