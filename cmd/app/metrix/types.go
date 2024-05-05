package metrix

import (
	"sync"
)

type (
	AllMetrixData struct {
		AvgReq     *AvgReqData     `json:"avg_req"`
		RespTime   *RespTimeData   `json:"resp_time"`
		CurrentReq *CurrentReqData `json:"current_req"`
		DBMetrix   *DbMetrixData   `json:"db_metrix"`
	}
	syncData struct {
		mx             sync.RWMutex
		oneMinData     float32
		fiveMinData    float32
		fifteenMinData float32
	}
	AvgReqData struct {
		OneMinAvg     float32 `json:"one_min_avg_req"`
		FiveMinAvg    float32 `json:"five_min_avg_req"`
		FifteenMinAvg float32 `json:"fifteen_min_avg_req"`
	}
	CurrentReqData struct {
		CurrentReq int64 `json:"current_req_num"`
	}
	RespTimeData struct {
		OneMinAvg     float32 `json:"one_min_avg_resp_time"`
		FiveMinAvg    float32 `json:"five_min_avg_resp_time"`
		FifteenMinAvg float32 `json:"fifteen_min_avg_resp_time"`
	}
	DbMetrixData struct {

		// AcquireCount returns the cumulative count of successful acquires from the pool.
		AcquireCount int64 `json:"acquire_count"`

		// AcquireDuration returns the total duration of all successful acquires from
		// the pool.
		AcquireDuration int64 `json:"acquire_duration"`

		// AcquiredConns returns the number of currently acquired connections in the pool.
		AcquiredConns int32 `json:"acquired_conns"`

		// CanceledAcquireCount returns the cumulative count of acquires from the pool
		// that were canceled by a context.
		CanceledAcquireCount int64 `json:"canceled_acquire_count"`

		// ConstructingConns returns the number of conns with construction in progress in
		// the pool.
		ConstructingConns int32 `json:"constructing_conns"`

		// EmptyAcquireCount returns the cumulative count of successful acquires from the pool
		// that waited for a resource to be released or constructed because the pool was
		// empty.
		EmptyAcquireCount int64 `json:"empty_acquire_count"`

		// IdleConns returns the number of currently idle conns in the pool.
		IdleConns int32 `json:"idle_conns"`

		// MaxConns returns the maximum size of the pool.
		MaxConns int32 `json:"max_conns"`

		// TotalConns returns the total number of resources currently in the pool.
		// The value is the sum of ConstructingConns, AcquiredConns, and
		// IdleConns.
		TotalConns int32 `json:"total_conns"`

		// NewConnsCount returns the cumulative count of new connections opened.
		NewConnsCount int64 `json:"new_conns_count"`

		// MaxLifetimeDestroyCount returns the cumulative count of connections destroyed
		// because they exceeded MaxConnLifetime.
		MaxLifetimeDestroyCount int64 `json:"max_lifetime_destroy_count"`

		// MaxIdleDestroyCount returns the cumulative count of connections destroyed because
		// they exceeded MaxConnIdleTime.
		MaxIdleDestroyCount int64 `json:"max_idle_destroy_count"`
	}
)
