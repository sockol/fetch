package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestDB(t *testing.T) {
	NewReceipt := func(retrailer string) *Receipt {
		return &Receipt{
			Retailer: retrailer,
		}
	}
	NewReceiptWithID := func(r *Receipt, id uuid.UUID) *ReceiptWithID {
		return &ReceiptWithID{r, id}
	}
	tests := []struct {
		description      string
		receipts         []*Receipt
		expectedReceipts []*ReceiptWithID
	}{
		{
			description:      "Empty base case",
			receipts:         nil,
			expectedReceipts: nil,
		},
		{
			description: "1 item",
			receipts: []*Receipt{
				NewReceipt("1"),
			},
			expectedReceipts: []*ReceiptWithID{
				NewReceiptWithID(NewReceipt("1"), uuid.New()),
			},
		},
		{
			description: "2 items",
			receipts: []*Receipt{
				NewReceipt("1"),
				NewReceipt("2"),
			},
			expectedReceipts: []*ReceiptWithID{
				NewReceiptWithID(NewReceipt("1"), uuid.New()),
				NewReceiptWithID(NewReceipt("2"), uuid.New()),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			db := NewDB()
			for _, r := range tc.receipts {
				db.Add(r)
			}
			if len(tc.expectedReceipts) != len(db.receipts) {
				t.Fatalf("Invalid recepits length. expected: %v, got: %v", len(tc.expectedReceipts), len(db.receipts))
			}
			for i, r := range tc.expectedReceipts {
				if r.Retailer != db.receipts[i].Retailer {
					t.Fatalf("expected: %v, got: %v", r, db.receipts[i])
				}

				e := db.Get(db.receipts[i].ID)
				if e.Retailer != db.receipts[i].Retailer {
					t.Fatalf("expected: %v, got: %v", e, db.receipts[i])
				}
			}
		})
	}
}
