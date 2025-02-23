package main

import (
	"fmt"
	"time"

	"github.com/ahmad-abuziad/transaction-processor/internal/data"
	"go.uber.org/zap"
)

const (
	BatchSize  = 100
	BatchTime  = 5 * time.Second
	NumWorkers = 5
)

var (
	txnsChan    = make(chan data.SalesTransaction, 1000)
	refreshChan = make(chan int, 200)
	quitChan    = make(chan struct{})
)

func (app *application) aggregateTransactions() {
	batch := []data.SalesTransaction{}
	ticker := time.NewTicker(BatchTime)

	for {
		select {
		case txn, ok := <-txnsChan:
			if !ok {
				if len(batch) > 0 {
					app.logger.Info("batching due to closing channel", zap.Int("batch_size", len(batch)))
					app.batchInsertWithRetry(batch)
				}
				return
			}

			batch = append(batch, txn)
			if len(batch) >= BatchSize {
				app.logger.Info("batching due to batch size limit", zap.Int("batch_size", len(batch)))
				app.batchInsertWithRetry(batch)
				refreshChan <- len(batch)
				batch = []data.SalesTransaction{}
			}

		case <-ticker.C:
			if len(batch) > 0 {
				app.logger.Info("batching due to time limit", zap.Int("batch_size", len(batch)))
				app.batchInsertWithRetry(batch)
				refreshChan <- len(batch)
				batch = []data.SalesTransaction{}
			}
		}
	}
}

func (app *application) batchInsertWithRetry(txns []data.SalesTransaction) {
	// uncomment if you want to simulate retries
	// if rand.Intn(100) == 1 {
	// 	txns[len(txns)-1].QuantitySold = -1 // will fail due database constraint
	// }

	retries := 3
	delay := 10 * time.Second
	var err error
	for i := 0; i < retries; i++ {
		err = app.models.SalesTransactions.InsertBatch(txns)
		if err == nil {
			return // Success
		}
		app.logger.Error("retrying...", zap.Int("attempt", i+1), zap.String("error", err.Error()), zap.Duration("delay", delay))
		time.Sleep(delay)
	}

	app.logger.Error("insert batch failed after retries", zap.Any("transactions", txns))
}

func (app *application) startWorkers() {
	app.startAggregateTransactionsWorkers()
	app.startRefreshTopSellingWorker()
}

func (app *application) stopWorkers() {
	close(txnsChan)
	close(quitChan)
	app.wg.Wait()
}

func (app *application) startAggregateTransactionsWorkers() {
	for i := 1; i <= NumWorkers; i++ {
		app.worker(app.aggregateTransactions)
	}
}

const (
	RefreshEveryXTransaction = 1000
	TxnsProcessedTime        = 1 * time.Second
)

func (app *application) startRefreshTopSellingWorker() {
	app.worker(func() {
		app.refreshTopSellingProducts() // first load
		var txnsNum = 0
		for {
			select {
			case txnNum := <-refreshChan:
				txnsNum = txnsNum + txnNum
				if txnsNum >= RefreshEveryXTransaction {
					app.refreshTopSellingProducts()
					txnsNum = 0
				}
			case <-quitChan:
				return
			}
		}
	})
}

func (app *application) refreshTopSellingProducts() {
	topSellingProducts, err := app.models.SalesTransactions.GetTopSellingProducts(10)
	if err != nil {
		app.logger.Error("Failed to read top selling products cache", zap.String("error", err.Error()))
	}

	err = app.cache.SalesCache.SetTopSellingProducts(topSellingProducts)
	if err != nil {
		app.logger.Error("Failed to set top selling products cache", zap.String("error", err.Error()))
	}
}

func (app *application) worker(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}
