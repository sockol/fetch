package main

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`

	Items []ReceiptItem `json:"items"`
}

type ReceiptItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// GetPoints returns how many points should be awarded to a receipt according to the rules listed below.
func (r *Receipt) GetPoints() (int, error) {
	points := float64(0)
	// - One point for every alphanumeric character in the retailer name.
	{
		for _, s := range r.Retailer {
			if isAlphanumeric(s) {
				points++
			}
		}
	}
	// - 50 points if the total is a round dollar amount with no cents.
	{
		total, err := strconv.ParseFloat(r.Total, 64)
		if err != nil {
			return 0, err
		}
		if isRoundNumber(total) {
			points += 50
		}
	}
	// - 25 points if the total is a multiple of 0.25.
	{
		total, err := strconv.ParseFloat(r.Total, 64)
		if err != nil {
			return 0, err
		}
		if isRoundNumber(total / 0.25) {
			points += 25
		}
	}
	// - 5 points for every two items on the receipt.
	{
		i := len(r.Items) / 2
		points += float64(i) * 5
	}
	// - If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	{
		for _, item := range r.Items {

			desc := strings.TrimSpace(item.ShortDescription)
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			if isRoundNumber(float64(len(desc)) / 3) {
				price *= 0.2
				if !isRoundNumber(price) {
					price = float64(int(price) + 1)
				}
				points += price
			}
		}
	}
	// - 6 points if the day in the purchase date is odd.
	{
		t, err := time.Parse("2006-01-02", r.PurchaseDate)
		if err != nil {
			return 0, err
		}
		isOdd := t.Day()%2 == 1
		if isOdd {
			points += 6
		}
	}
	// - 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	{
		t, err := time.Parse("15:04:05", r.PurchaseTime+":00")
		if err != nil {
			return 0, err
		}
		purchaseDate := time.Date(2000, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
		startDate := time.Date(2000, 1, 1, 14, 0, 0, 0, time.UTC)
		endDate := time.Date(2000, 1, 1, 16, 0, 0, 0, time.UTC)
		if startDate.Unix() <= purchaseDate.Unix() && purchaseDate.Unix() < endDate.Unix() {
			points += 10
		}

	}
	return int(points), nil
}

func isAlphanumeric(s rune) bool {
	return unicode.IsLetter(s) || unicode.IsDigit(s)
}

func isRoundNumber(num float64) bool {
	return math.Abs(num-math.Round(num)) < 1e-10
}
