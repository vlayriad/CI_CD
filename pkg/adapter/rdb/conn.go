package rdb

import (
	"context"
	"fmt"
	"kaffein/pkg"
	"net"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type setRedis struct {
	addr string
	port int
	db   int
	pass string
	conn *redis.Client
}

var (
	instance *setRedis
	once     sync.Once
	_        pkg.Server = (*setRedis)(nil)
)

type Option func(*setRedis)

func Host(s string) Option { return func(c *setRedis) { c.addr = s } }
func Port(p int) Option    { return func(c *setRedis) { c.port = p } }
func DB(d int) Option      { return func(c *setRedis) { c.db = d } }
func Pass(p string) Option { return func(c *setRedis) { c.pass = p } }

func NewDatabaseRedis(opts ...Option) *setRedis {
	once.Do(func() {
		c := &setRedis{}
		for _, opt := range opts {
			opt(c)
		}
		instance = c
	})
	return instance
}

func (s setRedis) Start(ctx context.Context) error {
	if s.addr == "" || s.port == 0 {
		log.Error().
			Msg("incomplete redis configuration")
		return nil
	}

	dsn := net.JoinHostPort(s.addr, strconv.Itoa(s.port))

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: s.pass,
		DB:       s.db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to redis")
		return nil
	}

	s.conn = client
	log.Info().
		Int("db", s.db).
		Str("host", s.addr).
		Int("port", s.port).
		Msg("connected to redis")

	return nil
}

func (s setRedis) Stop(ctx context.Context) error {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		log.Warn().
			Msg("redis connection closed")
		return nil
	}
	return nil
}

func GetInstance() *setRedis {
	if instance == nil {
		log.Panic().
			Msg("redis instance is not initialized")
	}
	return instance
}
