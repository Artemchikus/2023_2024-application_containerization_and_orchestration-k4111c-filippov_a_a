package metrix

import (
	"context"
	"sync"
	"time"
)

type (
	AvgRequests struct {
		count    float32
		syncData *syncData
	}
)

func mustAvgRequests(ctx context.Context) *AvgRequests {
	avgR := &AvgRequests{
		syncData: &syncData{
			mx: sync.RWMutex{},
		}}

	avgR.startTicker(ctx)

	return avgR
}

func (avgR *AvgRequests) startTicker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case <-ticker.C:
				avgR.count++
				if avgR.count == 90000 {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *AvgRequests) Set(number int64) {
	fNumber := float32(number)

	a.syncData.mx.Lock()
	defer a.syncData.mx.Unlock()

	switch {
	case a.count < 6000:
		a.syncData.oneMinData = (a.syncData.oneMinData*a.count + fNumber) / (a.count + 1)
		a.syncData.fiveMinData = a.syncData.oneMinData
		a.syncData.fifteenMinData = a.syncData.oneMinData
	case a.count < 30000:
		a.syncData.oneMinData = (a.syncData.oneMinData*6000 + fNumber) / 6001
		a.syncData.fiveMinData = (a.syncData.oneMinData*a.count + fNumber) / (a.count + 1)
		a.syncData.fifteenMinData = a.syncData.fiveMinData
	case a.count < 90000:
		a.syncData.oneMinData = (a.syncData.oneMinData*6000 + fNumber) / 6001
		a.syncData.fiveMinData = (a.syncData.fiveMinData*30000 + fNumber) / 30001
		a.syncData.fifteenMinData = (a.syncData.fifteenMinData*a.count + fNumber) / (a.count + 1)
	default:
		a.syncData.oneMinData = (a.syncData.oneMinData*6000 + fNumber) / 6001
		a.syncData.fiveMinData = (a.syncData.fiveMinData*30000 + fNumber) / 30001
		a.syncData.fifteenMinData = (a.syncData.fifteenMinData*90000 + fNumber) / 90001
	}
}

func (a *AvgRequests) Get() *AvgReqData {
	return &AvgReqData{
		OneMinAvg:     a.syncData.oneMinData,
		FiveMinAvg:    a.syncData.fiveMinData,
		FifteenMinAvg: a.syncData.fifteenMinData,
	}
}
