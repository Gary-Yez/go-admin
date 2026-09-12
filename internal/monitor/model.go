package monitor

import "time"

const Interval = 5 * time.Second
const OfflineAfter = 20 * time.Second
const Retention = 10 * time.Minute

// 可选指标使用指针，不可采集时输出 null，而不是伪装成零。
type Snapshot struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Hostname       string    `json:"hostname"`
	PID            int       `json:"pid"`
	OS             string    `json:"os"`
	Arch           string    `json:"arch"`
	Version        string    `json:"version"`
	StartedAt      time.Time `json:"started_at"`
	SampledAt      time.Time `json:"sampled_at"`
	LastSeen       time.Time `json:"last_seen"`
	Online         bool      `json:"online"`
	Uptime         float64   `json:"uptime"`
	CPUCount       int       `json:"cpu_count"`
	CPU            *float64  `json:"cpu"`
	MemoryTotal    *uint64   `json:"memory_total"`
	MemoryUsed     *uint64   `json:"memory_used"`
	MemoryPercent  *float64  `json:"memory_percent"`
	DiskTotal      *uint64   `json:"disk_total"`
	DiskUsed       *uint64   `json:"disk_used"`
	DiskPercent    *float64  `json:"disk_percent"`
	ProcessCPU     *float64  `json:"process_cpu"`
	RSS            *uint64   `json:"rss"`
	GoVersion      string    `json:"go_version"`
	Goroutines     uint64    `json:"goroutines"`
	HeapBytes      uint64    `json:"heap_bytes"`
	HeapObjects    uint64    `json:"heap_objects"`
	GCCycles       uint64    `json:"gc_cycles"`
	GCDelta        *uint64   `json:"gc_delta"`
	GCPauseSeconds *float64  `json:"gc_pause_seconds"`
	GOMAXPROCS     int       `json:"gomaxprocs"`
	Notices        []string  `json:"notices"`
}
