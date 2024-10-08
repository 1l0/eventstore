package pebble

import (
	"context"
	"encoding/hex"

	"github.com/cockroachdb/pebble"
	"github.com/fiatjaf/eventstore"
	bin "github.com/fiatjaf/eventstore/internal/binary"
	"github.com/nbd-wtf/go-nostr"
)

func (b *PebbleBackend) SaveEvent(ctx context.Context, evt *nostr.Event) error {
	// batch := b.NewBatch()
	// defer batch.Close()

	// query event by id to ensure we don't save duplicates
	id, _ := hex.DecodeString(evt.ID)
	prefix := make([]byte, 1+8)
	prefix[0] = indexIdPrefix
	copy(prefix[1:], id)
	it, err := b.NewIter(nil)
	if err != nil {
		return err
	}
	defer it.Close()
	it.SeekGE(prefix)
	if b.ValidForPrefix(it, prefix) {
		// event exists
		return eventstore.ErrDupEvent
	}

	// encode to binary
	bin, err := bin.Marshal(evt)
	if err != nil {
		return err
	}

	idx := b.Serial()
	// raw event store
	if err := b.Set(idx, bin, pebble.Sync); err != nil {
		return err
	}

	for k := range b.getIndexKeysForEvent(evt, idx[1:]) {
		if err := b.Set(k, nil, pebble.Sync); err != nil {
			return err
		}
	}

	// if err := batch.Commit(nil); err != nil {
	// 	return err
	// }

	return nil
}
