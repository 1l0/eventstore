package pebble

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/cockroachdb/pebble"
	"github.com/fiatjaf/eventstore"
	"github.com/nbd-wtf/go-nostr"
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

var _ eventstore.Store = (*PebbleBackend)(nil)

type PebbleBackend struct {
	Path     string
	MaxLimit int

	// Experimental
	SkipIndexingTag func(event *nostr.Event, tagName string, tagValue string) bool
	// Experimental
	IndexLongerTag func(event *nostr.Event, tagName string, tagValue string) bool

	*pebble.DB

	serial atomic.Uint32
}

func (b *PebbleBackend) Init() error {
	db, err := pebble.Open(b.Path, nil)
	if err != nil {
		return err
	}
	b.DB = db

	if err := b.runMigrations(); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	if b.MaxLimit == 0 {
		b.MaxLimit = 500
	}

	it, err := b.NewIter(&pebble.IterOptions{
		LowerBound: []byte{0},
		UpperBound: []byte{1},
	})
	if err != nil {
		return err
	}
	it.Last()
	it.SeekLT([]byte{1})
	if it.Valid() {
		if k := it.Key(); k != nil {
			idx := k[1:]
			serial := binary.BigEndian.Uint32(idx)
			b.serial.Store(serial)
		} else {
			return fmt.Errorf("error initializing serial: %w", err)
		}
	}
	if err := it.Close(); err != nil {
		return err
	}

	return nil
}

func (b *PebbleBackend) Close() {
	b.DB.Close()
}

func (b *PebbleBackend) Serial() []byte {
	next := b.serial.Add(1)
	vb := make([]byte, 5)
	vb[0] = rawEventStorePrefix
	binary.BigEndian.PutUint32(vb[1:], next)
	return vb
}

func (b *PebbleBackend) ValidForPrefix(it *pebble.Iterator, prefix []byte) bool {
	return it.Valid() && bytes.HasPrefix(it.Key(), prefix)
}

func (b *PebbleBackend) UpperBound(lowerBound []byte) []byte {
	end := make([]byte, len(lowerBound))
	copy(end, lowerBound)
	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}
	return nil // no upper-bound
}

func (b *PebbleBackend) PrefixIterOptions(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: b.UpperBound(prefix),
	}
}
