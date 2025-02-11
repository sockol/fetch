package main

import (
	"testing"
)

func TestReceiptsGetPoints(t *testing.T) {
	NewReceipt := func(
		_retailer string,
		_total string,
		_date string,
		_time string,
	) Receipt {
		total := "1.1"
		if _total != "" {
			total = _total
		}

		date := "2000-01-02"
		if _date != "" {
			date = _date
		}

		tm := "00:00"
		if _time != "" {
			tm = _time
		}
		return Receipt{
			Retailer:     _retailer,
			PurchaseDate: date,
			PurchaseTime: tm,
			Total:        total,
			Items:        []ReceiptItem{},
		}
	}
	NewReceiptWithItems := func(items ...ReceiptItem) Receipt {
		r := NewReceipt("", "", "", "")
		r.Items = items
		return r
	}
	NewReceiptItem := func() ReceiptItem {
		return ReceiptItem{
			ShortDescription: "",
			Price:            "0",
		}
	}
	tests := []struct {
		description    string
		receipt        Receipt
		expectedPoints int
		expectError    bool
	}{
		{
			description:    "Empty base case, will fail",
			receipt:        NewReceipt("", "", "", ""),
			expectedPoints: 0,
			expectError:    true,
		},
		{
			description:    "Bad total format",
			receipt:        NewReceipt("", "abc", "", ""),
			expectedPoints: 0,
			expectError:    true,
		},
		{
			description: "Bad item price format",
			receipt: NewReceiptWithItems(
				ReceiptItem{"", "abc"},
			),
			expectedPoints: 0,
			expectError:    true,
		},
		{
			description:    "Bad date format",
			receipt:        NewReceipt("", "", "10-10-10", ""),
			expectedPoints: 0,
			expectError:    true,
		},
		{
			description:    "Bad time format",
			receipt:        NewReceipt("", "", "", "10:10:10"),
			expectedPoints: 0,
			expectError:    true,
		},
		{
			description:    "Basic receipt, expect no points",
			receipt:        NewReceipt("", "", "", ""),
			expectedPoints: 0,
		},
		{
			description:    "0 alphanumeric characters",
			receipt:        NewReceipt("$ %", "", "", ""),
			expectedPoints: 0,
		},
		{
			description:    "4 alphanumeric characters",
			receipt:        NewReceipt("ab34$ %", "", "", ""),
			expectedPoints: 4,
		},
		{
			description:    "total is multiple of .25",
			receipt:        NewReceipt("", "1.25", "", ""),
			expectedPoints: 25,
		},
		{
			description:    "total is multiple of .25 and round",
			receipt:        NewReceipt("", "1", "", ""),
			expectedPoints: 75,
		},
		{
			description:    "no items",
			receipt:        NewReceiptWithItems(),
			expectedPoints: 0,
		},
		{
			description: "1 item, expect no extra points",
			receipt: NewReceiptWithItems(
				NewReceiptItem(),
			),
			expectedPoints: 0,
		},
		{
			description: "2 items, expect +5 extra points",
			receipt: NewReceiptWithItems(
				NewReceiptItem(),
				NewReceiptItem(),
			),
			expectedPoints: 5,
		},
		{
			description: "4 items, expect +5 * 2 extra points",
			receipt: NewReceiptWithItems(
				NewReceiptItem(),
				NewReceiptItem(),
				NewReceiptItem(),
				NewReceiptItem(),
			),
			expectedPoints: 10,
		},
		{
			description: "1 item, no points for trimmed description",
			receipt: NewReceiptWithItems(
				ReceiptItem{" ab ", "0"},
			),
			expectedPoints: 0,
		},
		{
			description: "1 item, yes points for trimmed description",
			receipt: NewReceiptWithItems(
				ReceiptItem{" abc ", "10"},
			),
			expectedPoints: 2,
		},
		{
			description: "1 item, yes points for trimmed description, rounded up",
			receipt: NewReceiptWithItems(
				ReceiptItem{" abc ", "11"},
			),
			expectedPoints: 3,
		},
		{
			description: "2 items, yes points for trimmed description, rounded up for both",
			receipt: NewReceiptWithItems(
				ReceiptItem{" abc ", "11"},
				ReceiptItem{" def ", "11"},
			),
			expectedPoints: 5 + 6,
		},
		{
			description:    "odd purcase date",
			receipt:        NewReceipt("", "", "2006-01-01", ""),
			expectedPoints: 6,
		},
		{
			description:    "extra points for time on left boundary",
			receipt:        NewReceipt("", "", "", "14:00"),
			expectedPoints: 10,
		},
		{
			description:    "no points for time on right boundary",
			receipt:        NewReceipt("", "", "", "16:00"),
			expectedPoints: 0,
		},
		{
			description:    "extra points for time inbetween left and right boundary",
			receipt:        NewReceipt("", "", "", "15:00"),
			expectedPoints: 10,
		},
		{
			description: "full test #1",
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []ReceiptItem{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					}, {
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					}, {
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					}, {
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					}, {
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
				Total: "35.35",
			},
			expectedPoints: 28,
		},
		{
			description: "full test #2",
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []ReceiptItem{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					}, {
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
				Total: "9.00",
			},
			expectedPoints: 109,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {

			points, err := tc.receipt.GetPoints()
			if err != nil {
				t.Log(err)
				if !tc.expectError {
					t.Fatalf("expected error")
				}
			}

			if points != tc.expectedPoints {
				t.Fatalf("Wrong points. expected: %v, got: %v", tc.expectedPoints, points)
			}
		})
	}
}
