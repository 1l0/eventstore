package sqlite

import (
	"github.com/jmoiron/sqlx"
)

type SQLiteBackend struct {
	*sqlx.DB
	DatabaseURL       string
	QueryLimit        int
	QueryIDsLimit     int
	QueryAuthorsLimit int
	QueryKindsLimit   int
	QueryTagsLimit    int
}

func (b *SQLiteBackend) Close() {
	b.DB.Close()
}
