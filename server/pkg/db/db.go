package db

import (
	"context"
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"server/pkg/config"
)

var db *sql.DB

func Setup() {
	c := config.GetConfig()
	d, err := sql.Open("sqlite", c.DbUrl)
	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}
	db = d
}
func GetDB() *sql.DB {
	return db
}

type Runner interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
