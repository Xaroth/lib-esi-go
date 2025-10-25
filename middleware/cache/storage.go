package cache

import (
	"database/sql"
	"net/url"
	"strings"

	"github.com/bartventer/httpcache/store"
	"github.com/bartventer/httpcache/store/driver"
	"github.com/bartventer/httpcache/store/expapi"
)

type storage struct {
	db        *sql.DB
	tableName string
	encode    func(value []byte) []byte
	decode    func(value []byte) ([]byte, error)
}

var _ driver.Conn = &storage{}
var _ expapi.KeyLister = &storage{}

const (
	CacheDriverName          = "lib-esi-go"
	DefaultTableName         = "cache"
	DefaultDriver            = "sqlite"
	DefaultEnableCompression = true
)

func noopEncode(value []byte) []byte {
	return value
}

func noopDecode(value []byte) ([]byte, error) {
	return value, nil
}

func init() {
	store.Register(CacheDriverName, driver.DriverFunc(func(u *url.URL) (driver.Conn, error) {
		opts := make([]Option, 0)

		switch {
		case u.Host == ":memory:":
			opts = append(opts, WithMemoryStore())
		case u.Host == ".":
			opts = append(opts, WithPath(strings.TrimLeft(u.Path, "/")))
		default:
			opts = append(opts, WithPath(u.Path))
		}

		query := u.Query()

		if val := query.Get("driver"); val != "" {
			opts = append(opts, WithDriver(val))
		}
		if val := query.Get("table"); val != "" {
			opts = append(opts, WithTableName(val))
		}
		if query.Has("compress") {
			val := query.Get("compress")
			opts = append(opts, WithCompression(val != "false"))
		}

		return NewStorage(opts...)
	}))
}

func NewStorage(opts ...Option) (*storage, error) {
	config := NewConfig(opts...)

	db, err := sql.Open(config.driver, config.dsn)
	if err != nil {
		return nil, err
	}

	storage := &storage{
		db:        db,
		tableName: config.tableName,
	}
	if config.enableCompression {
		storage.encode = ZstdEncode
		storage.decode = ZstdDecode
	} else {
		storage.encode = noopEncode
		storage.decode = noopDecode
	}

	if err := storage.initialize(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *storage) initialize() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + s.tableName + ` (
			key TEXT PRIMARY KEY,
			value BLOB
		)
	`)
	return err
}

func (s *storage) Get(key string) ([]byte, error) {
	row := s.db.QueryRow(`SELECT value FROM `+s.tableName+` WHERE key = ?`, key)
	var value []byte
	err := row.Scan(&value)
	if err != nil {
		return nil, err
	}
	return s.decode(value)
}

func (s *storage) Set(key string, value []byte) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO `+s.tableName+` (key, value) VALUES (?, ?)`, key, s.encode(value))
	return err
}

func (s *storage) Delete(key string) error {
	_, err := s.db.Exec(`DELETE FROM `+s.tableName+` WHERE key = ?`, key)
	return err
}

func (s *storage) Keys(prefix string) ([]string, error) {
	rows, err := s.db.Query(`SELECT key FROM `+s.tableName+` WHERE key LIKE ?`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		err := rows.Scan(&key)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *storage) Close() error {
	return s.db.Close()
}
