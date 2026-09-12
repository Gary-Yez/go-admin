package monitor

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/metrics"
	"strings"
	"time"
)

type Collector struct {
	base            Snapshot
	process         *process.Process
	previousCPU     *cpu.TimesStat
	previousProcess float64
	previousTime    time.Time
	previousGC      uint64
	previousPause   float64
	sampled         bool
	diskPath        string
}

func NewCollector(name string) *Collector {
	host, _ := os.Hostname()
	if name == "" {
		name = host
	}
	started := time.Now().UTC()
	p, _ := process.NewProcess(int32(os.Getpid()))
	path, _ := os.Getwd()
	version := "开发构建"
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				version = setting.Value
				if len(version) > 12 {
					version = version[:12]
				}
			}
		}
	}
	return &Collector{base: Snapshot{ID: rand.Text(), Name: strings.TrimSpace(name), Hostname: host, PID: os.Getpid(), OS: runtime.GOOS, Arch: runtime.GOARCH, Version: version, StartedAt: started, CPUCount: runtime.NumCPU(), GoVersion: runtime.Version()}, process: p, diskPath: path}
}

func (c *Collector) Collect(ctx context.Context) Snapshot {
	row := c.base
	now := time.Now().UTC()
	row.SampledAt = now
	row.Uptime = now.Sub(row.StartedAt).Seconds()
	row.GOMAXPROCS = runtime.GOMAXPROCS(0)
	row.Notices = []string{"机器资源采用操作系统可见范围；容器中不代表容器配额。磁盘统计运行目录所在卷。", "进程 CPU 以单核 100% 计，多核使用时可超过 100%。"}
	if values, err := cpu.TimesWithContext(ctx, false); err == nil && len(values) > 0 {
		value := values[0]
		if c.previousCPU != nil {
			busy := func(t cpu.TimesStat) float64 { return t.User + t.System + t.Nice + t.Irq + t.Softirq + t.Steal }
			before, after := busy(*c.previousCPU), busy(value)
			total := after - before + value.Idle - c.previousCPU.Idle + value.Iowait - c.previousCPU.Iowait
			if total > 0 && after >= before {
				percent := min(100.0, max(0.0, (after-before)/total*100))
				row.CPU = &percent
			}
		}
		c.previousCPU = &value
	} else {
		row.Notices = append(row.Notices, "机器 CPU 暂不可采集")
	}
	if value, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		row.MemoryTotal = &value.Total
		row.MemoryUsed = &value.Used
		row.MemoryPercent = &value.UsedPercent
	} else {
		row.Notices = append(row.Notices, "机器内存暂不可采集")
	}
	if value, err := disk.UsageWithContext(ctx, c.diskPath); err == nil {
		row.DiskTotal = &value.Total
		row.DiskUsed = &value.Used
		row.DiskPercent = &value.UsedPercent
	} else {
		row.Notices = append(row.Notices, "磁盘暂不可采集")
	}
	if c.process != nil {
		if value, err := c.process.MemoryInfoWithContext(ctx); err == nil {
			row.RSS = &value.RSS
		} else {
			row.Notices = append(row.Notices, "进程内存暂不可采集")
		}
		if value, err := c.process.TimesWithContext(ctx); err == nil {
			total := value.User + value.System
			if !c.previousTime.IsZero() && total >= c.previousProcess {
				percent := (total - c.previousProcess) / now.Sub(c.previousTime).Seconds() * 100
				row.ProcessCPU = &percent
			}
			c.previousTime = now
			c.previousProcess = total
		} else {
			row.Notices = append(row.Notices, "进程 CPU 暂不可采集")
		}
	} else {
		row.Notices = append(row.Notices, "进程资源暂不可采集")
	}
	samples := []metrics.Sample{{Name: "/sched/goroutines:goroutines"}, {Name: "/memory/classes/heap/objects:bytes"}, {Name: "/gc/heap/objects:objects"}, {Name: "/gc/cycles/total:gc-cycles"}}
	metrics.Read(samples)
	row.Goroutines = samples[0].Value.Uint64()
	row.HeapBytes = samples[1].Value.Uint64()
	row.HeapObjects = samples[2].Value.Uint64()
	row.GCCycles = samples[3].Value.Uint64()
	// 用累计暂停纳秒的差值展示采样间隔内 GC 暂停总时长。
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	pause := float64(stats.PauseTotalNs) / float64(time.Second)
	if c.sampled {
		delta := row.GCCycles - c.previousGC
		row.GCDelta = &delta
		elapsed := pause - c.previousPause
		row.GCPauseSeconds = &elapsed
	} else {
		row.Notices = append(row.Notices, fmt.Sprintf("每 %.0f 秒采样，CPU 和 GC 增量等待第二次采样", Interval.Seconds()))
	}
	c.previousGC = row.GCCycles
	c.previousPause = pause
	c.sampled = true
	return row
}
