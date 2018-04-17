package env

import (
	"os"

	"github.com/frozzare/go/cast"
	"github.com/joho/godotenv"
)

// Get retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// If value is empty and a default value is passed it will be returned instead.
func Get(key string, def ...string) string {
	v := os.Getenv(key)

	if v == "" && len(def) > 0 {
		return def[0]
	}

	return v
}

// Load .env files using godotenv package.
func Load(files ...string) error {
	return godotenv.Overload(files...)
}

// Set sets the value of the environment variable named by the key.
// It returns an error, if any.
func Set(key string, value interface{}) error {
	return os.Setenv(key, cast.MustString(value))
}
