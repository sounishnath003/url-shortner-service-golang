package core

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/knadh/goyesql/v2"
	_ "github.com/lib/pq"

	"github.com/sounishnath003/url-shortner-service-golang/internal/utils"
)

// InitCore helps to initialize all the necessary configuration
// to run the backend services.
//
// This handles all the heavy workloads of initiatizes the dependencies into
// single construct reducing the error prone and checks all params before the WP.
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

	stmts, err := co.prepareSQLQueryStmts()
	if err != nil {
		co.Lo.Error("Error preparing the sql statements", "error", err)
		panic(err)
	}

	co.QueryStmts = stmts

	return co
}

// Core struct holds up all the configuration required.
// It helps the application to run smoothly without fail.
type Core struct {
	Port       int
	Version    string
	JwtSecret  string
	QueryStmts *UrlShorterServiceQueries
	Lo         *slog.Logger

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
	db.SetConnMaxLifetime(10 * time.Second)
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

// prepareSQLQueryStmts helps to prepare the raw sql queries.
// Which handles and parses the .sql file and return the required SQL
// statements by the backend service to run.
//
// It pre-loads the query, no bottleneck to only assign the args to be executed.
func (co *Core) prepareSQLQueryStmts() (*UrlShorterServiceQueries, error) {
	queries := goyesql.MustParseFile("queries.sql")
	var queryStmts UrlShorterServiceQueries
	// prepares a given set of Queries and assigns the resulting *sql.Stmt statements to the fields
	err := goyesql.ScanToStruct(&queryStmts, queries, co.db)
	return &queryStmts, err
}
