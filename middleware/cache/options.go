package cache

type config struct {
	driver            string
	dsn               string
	tableName         string
	enableCompression bool
}

type Option func(*config)

func NewConfig(opts ...Option) *config {
	c := &config{
		driver:            DefaultDriver,
		tableName:         DefaultTableName,
		enableCompression: DefaultEnableCompression,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithDriver(driver string) Option {
	return func(c *config) {
		c.driver = driver
	}
}

func WithDsn(dsn string) Option {
	return func(c *config) {
		c.dsn = dsn
	}
}

func WithPath(path string) Option {
	return WithDsn(path)
}

func WithMemoryStore() Option {
	return WithDsn(":memory:")
}

func WithTableName(tableName string) Option {
	return func(c *config) {
		c.tableName = tableName
	}
}

func WithCompression(enable bool) Option {
	return func(c *config) {
		c.enableCompression = enable
	}
}
