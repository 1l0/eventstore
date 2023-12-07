package lmdb

import (
	"context"
	"encoding/hex"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/nbd-wtf/go-nostr"
)

func (b *LMDBBackend) DeleteEvent(ctx context.Context, evt *nostr.Event) error {
	err := b.lmdbEnv.Update(func(txn *lmdb.Txn) error {
		id, _ := hex.DecodeString(evt.ID)
		idx, err := txn.Get(b.indexId, id)
		if operr, ok := err.(*lmdb.OpError); ok && operr.Errno == lmdb.NotFound {
			// we already do not have this
			return nil
		}
		if err != nil {
			return err
		}

		// calculate all index keys we have for this event and delete them
		for _, k := range b.getIndexKeysForEvent(evt) {
			if err := txn.Del(k.dbi, k.key, idx); err != nil {
				return err
			}
		}

		// delete the raw event
		return txn.Del(b.rawEventStore, idx, nil)
	})
	if err != nil {
		return err
	}

	return nil
}
