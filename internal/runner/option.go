package runner

import (
	"log"

	"github.com/frozzare/max/internal/backend"
	"github.com/frozzare/max/internal/config"
)

// Option configures a runtime option.
type Option func(*Runner)

// Config returns an option configured with a config value.
func Config(config *config.Config) Option {
	return func(r *Runner) {
		r.config = config
	}
}

// Engine returns an option configured with a runner engine.
func Engine(engine backend.Engine) Option {
	return func(r *Runner) {
		r.engine = engine
	}
}

// Log returns an option configured with a log value.
func Log(log *log.Logger) Option {
	return func(r *Runner) {
		r.log = log
	}
}

// Once returns an option configured with a once value.
func Once(once bool) Option {
	return func(r *Runner) {
		r.once = once
	}
}

// Quiet returns an option configured with a quiet value.
func Quiet(quiet bool) Option {
	return func(r *Runner) {
		r.quiet = quiet
	}
}

// Verbose returns an option configured with a verbose value.
func Verbose(verbose bool) Option {
	return func(r *Runner) {
		r.verbose = verbose
	}
}
