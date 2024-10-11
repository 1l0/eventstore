package sqlite3

import (
	_ "embed"

	"github.com/fiatjaf/eventstore"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	queryLimit        = 100
	queryIDsLimit     = 500
	queryAuthorsLimit = 500
	queryKindsLimit   = 10
	queryTagsLimit    = 10
)

var _ eventstore.Store = (*SQLite3Backend)(nil)

type SQLite3Backend struct {
	*sqlx.DB
	DatabaseURL       string
	QueryLimit        int
	QueryIDsLimit     int
	QueryAuthorsLimit int
	QueryKindsLimit   int
	QueryTagsLimit    int
}

//go:embed sql/init.sql
var sqlInit string

func (b *SQLite3Backend) Init() error {
	db, err := sqlx.Connect("sqlite3", b.DatabaseURL)
	if err != nil {
		return err
	}

	db.Mapper = reflectx.NewMapperFunc("json", sqlx.NameMapper)
	b.DB = db

	_, err = b.DB.Exec(sqlInit)
	if err != nil {
		return err
	}

	if b.QueryLimit == 0 {
		b.QueryLimit = queryLimit
	}
	if b.QueryIDsLimit == 0 {
		b.QueryIDsLimit = queryIDsLimit
	}
	if b.QueryAuthorsLimit == 0 {
		b.QueryAuthorsLimit = queryAuthorsLimit
	}
	if b.QueryKindsLimit == 0 {
		b.QueryKindsLimit = queryKindsLimit
	}
	if b.QueryTagsLimit == 0 {
		b.QueryTagsLimit = queryTagsLimit
	}
	return nil
}

func (b *SQLite3Backend) Close() {
	b.DB.Close()
}
