package config

// Config defines the flags used as structure for internal usage.
type Config struct {
	Datadir      string
	Port         string
	CrossCompile bool
	Branch       string
}
