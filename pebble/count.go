package pebble

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/cockroachdb/pebble"
	bin "github.com/fiatjaf/eventstore/internal/binary"
	"github.com/nbd-wtf/go-nostr"
)

func (b *PebbleBackend) CountEvents(ctx context.Context, filter nostr.Filter) (int64, error) {
	var count int64 = 0

	queries, extraFilter, since, err := prepareQueries(filter)
	if err != nil {
		return 0, err
	}

	batch := b.NewBatch()
	defer batch.Close()

	// actually iterate
	for _, q := range queries {
		it, er := batch.NewIter(nil)
		if er != nil {
			return 0, err
		}
		defer it.Close()

		for it.SeekGE(q.startingPoint); b.ValidForPrefix(it, q.prefix); it.Prev() {
			select {
			case <-ctx.Done():
				return 0, fmt.Errorf("context canceled")
			default:
			}

			key := it.Key()

			idxOffset := len(key) - 4 // this is where the idx actually starts

			// "id" indexes don't contain a timestamp
			if !q.skipTimestamp {
				createdAt := binary.BigEndian.Uint32(key[idxOffset-4 : idxOffset])
				if createdAt < since {
					break
				}
			}

			idx := make([]byte, 5)
			idx[0] = rawEventStorePrefix
			copy(idx[1:], key[idxOffset:])

			if extraFilter == nil {
				count++
			} else {
				// fetch actual event
				value, closer, err := batch.Get(idx)
				if err != nil {
					return 0, err
				}
				defer closer.Close()

				evt := &nostr.Event{}
				if err := bin.Unmarshal(value, evt); err != nil {
					return 0, err
				}

				// check if this matches the other filters that were not part of the index
				if extraFilter.Matches(evt) {
					count++
				}
			}
		}
	}

	if err := batch.Commit(nil); err != nil {
		return 0, err
	}

	return count, nil
}

func (b *PebbleBackend) ValidForPrefix(it *pebble.Iterator, prefix []byte) bool {
	return it.Valid() && bytes.HasPrefix(it.Key(), prefix)
}
