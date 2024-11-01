package core

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/knadh/goyesql/v2"
	_ "github.com/lib/pq"

	"github.com/sounishnath003/url-shortner-service-golang/internal/bloom"
	"github.com/sounishnath003/url-shortner-service-golang/internal/utils"
)

// InitCore helps to initialize all the necessary configuration
// to run the backend services.
//
// This handles all the heavy workloads of initiatizes the dependencies into
// single construct reducing the error prone and checks all params before the WP.
func InitCore() *Core {

	co := &Core{
		Version:     "0.0.1",
		Port:        utils.GetEnv("PORT", 3000).(int),
		JwtSecret:   utils.GetEnv("JWT_SECRET", "24C3bebd3c22f155e57b6d426d2e7dfEZX2@1564#!").(string),
		dbType:      utils.GetEnv("DB_DRIVER", "postgres").(string),
		dsn:         utils.GetEnv("DSN", "postgres://root:root@127.0.0.1:5432/postgres?sslmode=disable").(string),
		BloomFilter: bloom.NewBloomFilter(10000000),

		RedisClientAddr: utils.GetEnv(
			"REDIS_CLIENT_ADDR",
			"localhost:6379",
		).(string),
		Lo: slog.Default(),
	}

	// Attach the db
	db, err := co.initDatabase()
	if err != nil {
		co.Lo.Error("Error initializing database", "error", err)
		panic(err)
	}

	co.db = db

	// Attach the redis client.
	rdb, err := co.initRedisConf()
	if err != nil {
		co.Lo.Error("Error initializing redis client", "error", err)
		panic(err)
	}
	co.rdb = rdb

	stmts, err := co.prepareSQLQueryStmts()
	if err != nil {
		co.Lo.Error("Error preparing the sql statements", "error", err)
		panic(err)
	}

	co.QueryStmts = stmts

	// Load the bloom filter with the shortUrl alias.
	// Runs in a separate go routine.
	go co.PreloadBloomFilter()
	go co.CacheShortOriginalUrls()

	return co
}

// Core struct holds up all the configuration required.
// It helps the application to run smoothly without fail.
type Core struct {
	Port            int
	Version         string
	JwtSecret       string
	QueryStmts      *UrlShorterServiceQueries
	Lo              *slog.Logger
	BloomFilter     *bloom.BloomFilter
	RedisClientAddr string

	dbType string
	dsn    string
	db     *sql.DB
	rdb    *redis.Client
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

func (co *Core) initRedisConf() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     co.RedisClientAddr,
		Password: "",
		DB:       0,
	})

	return rdb, nil
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

// PreloadBloomFilter helps to preload the bloom filter with all the existing short urls.
// This is done to improve the performance of the CustomAliasAvailabilityHandler.
//
// caller must run it in separate go routine. As the alias will be huge distributed.
func (co *Core) PreloadBloomFilter() error {
	rows, err := co.QueryStmts.GetAllShortUrlAliasQuery.Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		var shortUrl string
		err = rows.Scan(&shortUrl)
		if err != nil {
			return err
		}
		co.Lo.Info("added to preloading bloom filter", "shortUrl", shortUrl)
		co.BloomFilter.Add(shortUrl)
	}
	return nil
}

// CacheShortOriginalUrls helps to cache the short url and original url mappings into redis.
// This is done to improve the performance of the GetOriginalUrlHandler.
// This is an entire blocking infinite loop.
//
// caller must run it in separate go routine. As the alias will be huge distributed.
func (co *Core) CacheShortOriginalUrls() error {

	for {
		rows, err := co.QueryStmts.MostActiveHitsQuery.Query()
		if err != nil {
			co.Lo.Info("an error occured", "error", err)
			return err
		}

		for rows.Next() {
			var originalUrl string
			var shortUrl string

			err = rows.Scan(&originalUrl, &shortUrl)
			if err != nil {
				co.Lo.Info("an error occured", "error", err)
				return err
			}

			// Add in redis cache for 1 Hour eviction
			err = co.rdb.Set(context.Background(), shortUrl, originalUrl, 1*time.Hour).Err()

			co.Lo.Info("added to cache", "originalUrl", originalUrl, "shortUrl", shortUrl)
		}
		time.Sleep(1 * time.Minute)
	}
}

func (co *Core) FindOriginalUrlFromCache(shortUrl string) (string, error) {
	return co.rdb.Get(context.Background(), shortUrl).Result()
}

// CreateNewShortUrl helps to add a shortURL for the user.
// Execute and write the data using the database trasactions.
func (co *Core) CreateNewShortUrlAsTxn(OriginalUrl, shortUrl string, expiryDate time.Time, userID int) error {
	// Transaction init.
	tx, err := co.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	// Rollback aborts the transaction.
	defer tx.Rollback()

	// Insert the short url into the url_mappings table.
	_, err = tx.Exec("INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ($1, $2, $3, $4) RETURNING id", OriginalUrl, shortUrl, expiryDate, userID)
	if err != nil {
		return err
	}

	var shortUrlID int
	// Check if the short url is already present in the database.
	err = tx.QueryRow(
		"SELECT id FROM url_mappings WHERE original_url = $1 AND short_url = $2 AND user_id = $3",
		OriginalUrl,
		shortUrl,
		userID,
	).Scan(&shortUrlID)
	if err != nil {
		return err
	}

	// Insert the shortURL ID and userID into the users_url_mappings table.
	_, err = tx.Exec(
		"INSERT INTO users_url_mappings (UrlID, UserID) VALUES ($1, $2)",
		shortUrlID,
		userID,
	)

	if err != nil {
		return err
	}

	// Add the shortUrl into the bloom filter
	co.BloomFilter.Add(shortUrl)

	return tx.Commit()
}
