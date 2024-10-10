package pebble

import (
	"context"
	"encoding/hex"

	"github.com/cockroachdb/pebble"
	"github.com/nbd-wtf/go-nostr"
)

func (b *PebbleBackend) DeleteEvent(ctx context.Context, evt *nostr.Event) error {
	idx := make([]byte, 1, 5)
	idx[0] = rawEventStorePrefix

	// query event by id to get its idx
	idPrefix8, _ := hex.DecodeString(evt.ID[0 : 8*2])
	prefix := make([]byte, 1+8)
	prefix[0] = indexIdPrefix
	copy(prefix[1:], idPrefix8)

	it, err := b.NewIter(nil)
	if err != nil {
		return err
	}
	it.SeekGE(prefix)
	if b.ValidForPrefix(it, prefix) {
		idx = append(idx, it.Key()[1+8:]...)
	}
	it.Close()

	// if no idx was found, end here, this event doesn't exist
	if len(idx) == 1 {
		return nil
	}

	// calculate all index keys we have for this event and delete them
	for k := range b.getIndexKeysForEvent(evt, idx[1:]) {
		if err := b.Delete(k, pebble.Sync); err != nil {
			return err
		}
	}

	// delete the raw event
	b.Delete(idx, pebble.Sync)

	return nil
}
