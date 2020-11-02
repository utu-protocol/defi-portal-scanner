package collector

import (
	"time"

	"github.com/xujiajun/nutsdb"
)

type Store struct {
	db *nutsdb.DB
}

// OpenStore open the storage for writing
func OpenStore(storePath string) (db *Store, err error) {
	// Open the database located in the /tmp/nutsdb directory.
	// It will be created if it doesn't exist.
	store := new(Store)
	opt := nutsdb.DefaultOptions
	opt.Dir = storePath
	store.db, err = nutsdb.Open(opt)
	return
}

const (
	txBucket = "txs"
)

func (db *Store) Tx(from, to string, when time.Time) {
	//TODO: nohop
}

func (db *Store) Close() {
	db.db.Close()
}
