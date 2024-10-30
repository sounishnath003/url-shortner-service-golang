package core

import "database/sql"

// UrlShorterServiceQueries helps to prepare SQL statements
// to be executed and required by the backend service.
type UrlShorterServiceQueries struct {
	GetUserByEmail        *sql.Stmt `query:"GetUserByEmail"`
	CreateShortUrlQuery   *sql.Stmt `query:"CreateShortUrlQuery"`
	GetShortUrlQuery      *sql.Stmt `query:"GetShortUrlQuery"`
	GetIncrementalIDQuery *sql.Stmt `query:"GetIncrementalIDQuery"`
}
