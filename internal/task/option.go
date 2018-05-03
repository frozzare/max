package task

import "log"

// Option configures a runtime option.
type Option func(*Task)

// Log returns an option configured with a log value.
func Log(log *log.Logger) Option {
	return func(t *Task) {
		t.log = log
	}
}
