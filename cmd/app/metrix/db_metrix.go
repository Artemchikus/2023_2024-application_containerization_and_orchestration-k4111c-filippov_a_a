package metrix

import "find-ship/storage"

type (
	DBMetrix struct {
		storage storage.Storage
	}
)

func mustDBMetrix(storage storage.Storage) *DBMetrix {
	return &DBMetrix{
		storage: storage,
	}
}

func (m *DBMetrix) Get() *DbMetrixData {
	stat := m.storage.Stat()

	return &DbMetrixData{
		AcquireCount:            stat.AcquireCount(),
		AcquireDuration:         stat.AcquireDuration().Microseconds(),
		AcquiredConns:           stat.AcquiredConns(),
		CanceledAcquireCount:    stat.CanceledAcquireCount(),
		ConstructingConns:       stat.ConstructingConns(),
		EmptyAcquireCount:       stat.EmptyAcquireCount(),
		IdleConns:               stat.IdleConns(),
		MaxConns:                stat.MaxConns(),
		TotalConns:              stat.TotalConns(),
		NewConnsCount:           stat.NewConnsCount(),
		MaxLifetimeDestroyCount: stat.MaxLifetimeDestroyCount(),
		MaxIdleDestroyCount:     stat.MaxIdleDestroyCount(),
	}
}
