package core

import "database/sql"

// UrlShorterServiceQueries helps to prepare SQL statements
// to be executed and required by the backend service.
type UrlShorterServiceQueries struct {
	GetUserByEmail           *sql.Stmt `query:"GetUserByEmail"`
	CreateNewUser            *sql.Stmt `query:"CreateNewUser"`
	CreateShortUrlQuery      *sql.Stmt `query:"CreateShortUrlQuery"`
	GetShortUrlQuery         *sql.Stmt `query:"GetShortUrlQuery"`
	IncrUrlHitCountQuery     *sql.Stmt `query:"IncrUrlHitCountQuery"`
	GetIncrementalIDQuery    *sql.Stmt `query:"GetIncrementalIDQuery"`
	GetAllShortUrlAliasQuery *sql.Stmt `query:"GetAllShortUrlAliasQuery"`
	MostActiveHitsQuery      *sql.Stmt `query:"MostActiveHitsQuery"`
}
