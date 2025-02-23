package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

func main() {

	type tenant struct {
		id       int64
		branches []int64
		products []int64
	}

	tenants := []tenant{
		{
			id:       1,
			branches: []int64{1, 2, 3, 4, 5},
			products: []int64{1, 2, 3, 4, 5},
		},
		{
			id:       2,
			branches: []int64{6, 7, 8, 9, 10},
			products: []int64{6, 7, 8, 9, 10},
		},
		{
			id:       3,
			branches: []int64{11, 12, 13, 14, 15},
			products: []int64{11, 12, 13, 14, 15},
		},
	}

	var wg sync.WaitGroup

	for _, tenant := range tenants {
		for _, branch := range tenant.branches {
			wg.Add(1)
			go func() {
				for {
					//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
					txn := salesTransaction{
						TenantID:     int64(tenant.id),
						BranchID:     int64(branch),
						ProductID:    int64(tenant.products[rand.Intn(5)]),
						QuantitySold: rand.Intn(20),
						PricePerUnit: decimal.NewFromFloat(rand.Float64() * 100).StringFixed(2),
						Timestamp:    time.Now(),
					}
					request(txn)
				}
			}()
		}
	}

	wg.Wait()
}

type salesTransaction struct {
	TenantID     int64     `json:"tenantID"`
	BranchID     int64     `json:"branchID"`
	ProductID    int64     `json:"productID"`
	QuantitySold int       `json:"quantitySold"`
	PricePerUnit string    `json:"pricePerUnit"`
	Timestamp    time.Time `json:"timestamp"`
}

func request(txn salesTransaction) {
	jsonData, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost/tenant/%v/branch/%v/sales-transaction", txn.TenantID, txn.BranchID)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err.Error())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	fmt.Println(string(body))
}
