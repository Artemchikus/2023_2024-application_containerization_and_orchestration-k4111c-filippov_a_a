package metrix

import (
	"context"
	"sync"
	"time"
)

type (
	ResponseTime struct {
		count    float32
		syncData *syncData
	}
)

func mustResponseTime(ctx context.Context) *ResponseTime {
	respT := &ResponseTime{
		syncData: &syncData{
			mx: sync.RWMutex{},
		}}

	respT.startTicker(ctx)

	return respT
}

func (respT *ResponseTime) startTicker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case <-ticker.C:
				respT.count++
				if respT.count == 900000 {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *ResponseTime) Set(mSec int64) {
	fMSec := float32(mSec)

	r.syncData.mx.Lock()
	defer r.syncData.mx.Unlock()

	switch {
	case r.count < 6000:
		r.syncData.oneMinData = (r.syncData.oneMinData*r.count + fMSec) / (r.count + 1)
		r.syncData.fiveMinData = r.syncData.oneMinData
		r.syncData.fifteenMinData = r.syncData.oneMinData
	case r.count < 30000:
		r.syncData.oneMinData = (r.syncData.oneMinData*6000 + fMSec) / 6001
		r.syncData.fiveMinData = (r.syncData.oneMinData*r.count + fMSec) / (r.count + 1)
		r.syncData.fifteenMinData = r.syncData.fiveMinData
	case r.count < 90000:
		r.syncData.oneMinData = (r.syncData.oneMinData*6000 + fMSec) / 6001
		r.syncData.fiveMinData = (r.syncData.fiveMinData*30000 + fMSec) / 30001
		r.syncData.fifteenMinData = (r.syncData.fifteenMinData*r.count + fMSec) / (r.count + 1)
	default:
		r.syncData.oneMinData = (r.syncData.oneMinData*6000 + fMSec) / 6001
		r.syncData.fiveMinData = (r.syncData.fiveMinData*30000 + fMSec) / 30001
		r.syncData.fifteenMinData = (r.syncData.fifteenMinData*90000 + fMSec) / 90001
	}
}

func (r *ResponseTime) Get() *RespTimeData {
	return &RespTimeData{
		OneMinAvg:     r.syncData.oneMinData,
		FiveMinAvg:    r.syncData.fiveMinData,
		FifteenMinAvg: r.syncData.fifteenMinData,
	}
}
