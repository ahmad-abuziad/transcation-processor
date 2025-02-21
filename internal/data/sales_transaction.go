package data

import "time"

type SalesTransaction struct {
	TenantID     int64
	BranchID     int64
	ProductID    int64
	QuantitySold int
	PricePerUnit int64
	Created      time.Time
}
