package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"kaffein/pkg"
	"sync"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type setPgsql struct {
	addr string
	port int
	db   string
	user string
	pass string
	sslm bool

	conn *sql.DB
}

var (
	instance *setPgsql
	once     sync.Once
	_        pkg.Server = (*setPgsql)(nil)
)

type Option func(*setPgsql)

func Host(s string) Option { return func(c *setPgsql) { c.addr = s } }
func Port(p int) Option    { return func(c *setPgsql) { c.port = p } }
func DB(s string) Option   { return func(c *setPgsql) { c.db = s } }
func User(s string) Option { return func(c *setPgsql) { c.user = s } }
func Pass(s string) Option { return func(c *setPgsql) { c.pass = s } }
func SSLMode(enable bool) Option {
	return func(c *setPgsql) { c.sslm = enable }
}

func NewDatabasePGSQL(opts ...Option) *setPgsql {
	once.Do(func() {
		c := &setPgsql{}
		for _, opt := range opts {
			opt(c)
		}
		instance = c
	})
	return instance
}

func (c *setPgsql) Start(ctx context.Context) error {
	if c.addr == "" || c.user == "" || c.db == "" {
		log.Error().
			Msg("incomplete database configuration")
		return nil
	}

	sslMode := "disable"
	if c.sslm {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.addr, c.port, c.user, c.pass, c.db, sslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to open connection")
		return err
	}

	if err = db.PingContext(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("failed to ping database")
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	c.conn = db
	log.Info().
		Str("host", c.addr).
		Int("port", c.port).
		Str("db", c.db).
		Msg("connected to database")

	return nil
}

func (c *setPgsql) Stop(ctx context.Context) error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		log.Warn().
			Msg("connection closed")
		return nil
	}
	return nil
}

func GetInstance() *setPgsql {
	if instance == nil {
		log.Panic().
			Msg("instance not initialized")
	}
	return instance
}
