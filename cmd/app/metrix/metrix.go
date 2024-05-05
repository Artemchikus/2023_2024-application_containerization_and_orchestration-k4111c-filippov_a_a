package metrix

import (
	"context"
	"find-ship/storage"
)

type (
	Metrix struct {
		AvgReq     *AvgRequests
		RespTime   *ResponseTime
		CurrentReq *CurrentRequests
		DBMetrix   *DBMetrix
	}
)

func New(context context.Context, storage storage.Storage) *Metrix {
	return &Metrix{
		AvgReq:     mustAvgRequests(context),
		RespTime:   mustResponseTime(context),
		CurrentReq: mustCurrentRequests(),
		DBMetrix:   mustDBMetrix(storage),
	}
}

func (m *Metrix) GetAll() AllMetrixData {
	return AllMetrixData{
		AvgReq:     m.AvgReq.Get(),
		RespTime:   m.RespTime.Get(),
		CurrentReq: m.CurrentReq.Get(),
		DBMetrix:   m.DBMetrix.Get(),
	}
}
