package core

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/knadh/goyesql/v2"
	_ "github.com/lib/pq"

	"github.com/sounishnath003/url-shortner-service-golang/cmd/utils"
)

// InitCore helps to initialize all the necessary configuration
// to run the backend services.
func InitCore() *Core {

	co := &Core{
		Version:   "0.0.1",
		Port:      utils.GetEnv("PORT", 3000).(int),
		JwtSecret: utils.GetEnv("JWT_SECRET", "24C3bebd3c22f155e57b6d426d2e7dfEZX2@1564#!").(string),
		dbType:    utils.GetEnv("DB_DRIVER", "postgres").(string),
		dsn:       utils.GetEnv("DSN", "postgres://root:root@127.0.0.1:5432/postgres?sslmode=disable").(string),
		Lo:        slog.Default(),
	}

	// Attach the db
	db, err := co.initDatabase()
	if err != nil {
		co.Lo.Error("Error initializing database", "error", err)
		panic(err)
	}

	co.db = db

	co.parseSQLQueries()

	return co
}

// Core struct holds up all the configuration required.
// It helps the application to run smoothly without fail.
type Core struct {
	Port      int
	Version   string
	JwtSecret string
	Lo        *slog.Logger

	dbType string
	dsn    string
	db     *sql.DB
}

// initDatabase helps to instantiate a database connection.
// Throws *db, error. caller must handle the error incase.
func (co *Core) initDatabase() (*sql.DB, error) {
	// Open a database connection.
	db, err := sql.Open(co.dbType, co.dsn)
	if err != nil {
		return nil, err
	}
	// Set the connections defaults.
	db.SetMaxIdleConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(100 * time.Second)

	// Perform a ping check.
	err = db.Ping()
	if err != nil {
		co.Lo.Error("error pinging database", "error", err)
		return nil, err
	}
	co.Lo.Info("database connection has been established successfully.")
	return db, nil
}

func (co *Core) parseSQLQueries() {
	queries := goyesql.MustParseFile("queries.sql")

	for name, query := range queries {
		co.Lo.Info("query", "name", name, "query", query)
	}
}
