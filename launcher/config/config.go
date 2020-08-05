package config

// Config is the struct wrapper for the launcher configurations
type Config struct {
	Password     string
	ExternalHost string
	Nodes        int
	Validators   int
	Source       bool
	Debug        bool
}
