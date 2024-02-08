package lipgloss

import "os"

// Environ represents environment variables.
type Environ interface {
	Getenv(string) string
	LookupEnv(string) (string, bool)
	Environ() []string
}

// OsEnviron is an implementation of Environ that uses the os package.
type OsEnviron struct{}

// Getenv returns the value of the environment variable.
func (OsEnviron) Getenv(key string) string {
	return os.Getenv(key)
}

// LookupEnv retrieves the value of the environment variable.
func (OsEnviron) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Environ returns a copy of strings representing the environment.
func (OsEnviron) Environ() []string {
	return os.Environ()
}
