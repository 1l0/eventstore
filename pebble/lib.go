package pebble

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"github.com/nbd-wtf/go-nostr"
	// "github.com/fiatjaf/eventstore"
)

const (
	dbVersionKey          byte = 255
	rawEventStorePrefix   byte = 0
	indexCreatedAtPrefix  byte = 1
	indexIdPrefix         byte = 2
	indexKindPrefix       byte = 3
	indexPubkeyPrefix     byte = 4
	indexPubkeyKindPrefix byte = 5
	indexTagPrefix        byte = 6
	indexTag32Prefix      byte = 7
	indexTagAddrPrefix    byte = 8
)

// var _ eventstore.Store = (*PebbleBackend)(nil)

type PebbleBackend struct {
	Path     string
	MaxLimit int

	// Experimental
	SkipIndexingTag func(event *nostr.Event, tagName string, tagValue string) bool
	// Experimental
	IndexLongerTag func(event *nostr.Event, tagName string, tagValue string) bool

	*pebble.DB

	lastId atomic.Uint32
}

func (b *PebbleBackend) Init() error {
	db, err := pebble.Open(b.Path, nil)
	if err != nil {
		return err
	}
	b.DB = db

	itr, err := b.NewIter(nil)
	if err != nil {
		return err
	}
	if itr.Last() {
		if k := itr.Key(); k != nil {
			b.lastId.Store(binary.BigEndian.Uint32(k))
		}
	}
	if err := itr.Close(); err != nil {
		return err
	}

	return b.runMigrations()
}

func (b *PebbleBackend) Close() {
	b.DB.Close()
}

func (b *PebbleBackend) Serial() []byte {
	v := b.lastId.Add(1)
	vb := make([]byte, 4)
	binary.BigEndian.PutUint32(vb[:], uint32(v))
	return vb
}
