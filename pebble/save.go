package pebble

import (
	"context"
	"encoding/hex"

	"github.com/fiatjaf/eventstore"
	bin "github.com/fiatjaf/eventstore/internal/binary"
	"github.com/nbd-wtf/go-nostr"
)

func (b *PebbleBackend) SaveEvent(ctx context.Context, evt *nostr.Event) (err error) {
	batch := b.NewBatch()
	defer func() {
		err = batch.Close()
	}()

	// query event by id to ensure we don't save duplicates
	id, _ := hex.DecodeString(evt.ID)
	prefix := make([]byte, 1+8)
	prefix[0] = indexIdPrefix
	copy(prefix[1:], id)
	it, err := batch.NewIter(nil)
	if err != nil {
		return err
	}
	defer func() {
		err = it.Close()
	}()
	if it.SeekGE(prefix) {
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
	if err := batch.Set(idx, bin, nil); err != nil {
		return err
	}

	for _, k := range b.getIndexKeysForEvent(evt, idx[1:]) {
		if err := batch.Set(k, nil, nil); err != nil {
			return err
		}
	}

	if err := batch.Commit(nil); err != nil {
		return err
	}

	return nil
}
