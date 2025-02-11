package main

import (
	"github.com/google/uuid"
)

// DB is storing the whole object and building the points on the read path, because space is cheap
// and maybe the point generating algo will change.
type DB struct {
	receipts []*ReceiptWithID
}

type ReceiptWithID struct {
	*Receipt
	ID uuid.UUID
}

func NewDB() DB {
	return DB{
		receipts: []*ReceiptWithID{},
	}
}

func (db *DB) Add(r *Receipt) uuid.UUID {
	// O(1) insert, because no need to optimize.
	id := uuid.New()
	db.receipts = append(db.receipts, &ReceiptWithID{r, id})
	return id
}

func (db *DB) Get(id uuid.UUID) *Receipt {
	// Sequential scan, because no need to optimize.
	for _, r := range db.receipts {
		if r.ID == id {
			return r.Receipt
		}
	}
	return nil
}
