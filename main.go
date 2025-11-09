package main

import (
	"kaffein/pkg/server/app/_http"
	"os"

	"kaffein/config"
	"kaffein/helper"
	"kaffein/pkg"
	"kaffein/pkg/adapter/pgsql"
	"kaffein/pkg/adapter/rdb"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func init() {
	if err := config.SetupApp("config.json"); err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}
}

func main() {
	pgsqlx := pgsql.NewDatabasePGSQL(
		pgsql.Host(helper.GetString("postgres.host")),
		pgsql.Port(helper.GetInt("postgres.port")),
		pgsql.User(helper.GetString("postgres.user")),
		pgsql.Pass(helper.GetString("postgres.password")),
		pgsql.DB(helper.GetString("postgres.dbname")),
		pgsql.SSLMode(helper.GetBool("postgres.sslmode")),
	)

	rdbs := rdb.NewDatabaseRedis(
		rdb.Host(helper.GetString("redis.host")),
		rdb.Port(helper.GetInt("redis.port")),
		rdb.Pass(helper.GetString("redis.password")),
		rdb.DB(helper.GetInt("redis.db")),
	)

	httpServer := _http.NewServer(
		_http.WithHost(helper.GetString("server.app_host")),
		_http.WithPort(helper.GetInt("server.app_port")),
		_http.WithEngine(gin.Default()),
	)

	server := pkg.NewServer(
		pkg.WithName(helper.GetString("server.app_name")),
		pkg.WithVersion(helper.GetString("server.app_version")),
		pkg.WithAdapter(pgsqlx, rdbs, httpServer),
	)

	if err := server.Run(); err != nil {
		log.Error().
			Err(err).
			Msg("server exited with error")
		os.Exit(1)
	}

	log.Info().
		Msg("server exited gracefully ✅")
}
