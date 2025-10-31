package log

type Config struct {
	Level      string
	Format     string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Option func(*Config)

func WithLevel(level string) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

func WithFormat(format string) Option {
	return func(cfg *Config) {
		cfg.Format = format
	}
}

func WithFilename(filename string) Option {
	return func(cfg *Config) {
		cfg.Filename = filename
	}
}

func WithMaxSize(maxSize int) Option {
	return func(cfg *Config) {
		cfg.MaxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(cfg *Config) {
		cfg.MaxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) Option {
	return func(cfg *Config) {
		cfg.MaxAge = maxAge
	}
}
