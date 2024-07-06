package sqlite

import (
	_ "embed"

	"github.com/fiatjaf/eventstore"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "modernc.org/sqlite"
)

const (
	queryLimit        = 100
	queryIDsLimit     = 500
	queryAuthorsLimit = 500
	queryKindsLimit   = 10
	queryTagsLimit    = 10
)

var _ eventstore.Store = (*SQLiteBackend)(nil)

//go:embed create_table.sql
var createTable string

//go:embed create_index_id.sql
var createIndexID string

//go:embed create_index_pubkey.sql
var createIndexPubkey string

//go:embed create_index_time.sql
var createIndexTime string

//go:embed create_index_kind.sql
var createIndexKind string

//go:embed create_index_kindtime.sql
var createIndexKindTime string

var ddls = []string{
	createTable,
	createIndexID,
	createIndexPubkey,
	createIndexTime,
	createIndexKind,
	createIndexKindTime,
}

func (b *SQLiteBackend) Init() error {
	db, err := sqlx.Connect("sqlite", b.DatabaseURL)
	if err != nil {
		return err
	}

	db.Mapper = reflectx.NewMapperFunc("json", sqlx.NameMapper)
	b.DB = db

	for _, ddl := range ddls {
		_, err = b.DB.Exec(ddl)
		if err != nil {
			return err
		}
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
