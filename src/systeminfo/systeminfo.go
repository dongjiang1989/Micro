package systeminfo

import (
	"os/user"
	"runtime"

	log "github.com/cihub/seelog"
)

type Status struct{}

type Info struct {
	State       STATE  `json:"state"`
	DiskUsage   uint64 `json:"diskUsage"`
	DiskTotal   uint64 `json:"diskTotal"`
	MemoryUsage uint64 `json:"memoryUsage"`
	MemoryTotal uint64 `json:"memoryTotal"`
	HomeDir     string `json:"homeDir"`
	Uname       string `json:"uname"`
	NetIn       uint64 `json:"netIn"`
	NetOut      uint64 `json:"netOut"`
}

type memory struct {
	Usage uint64 `json:"diskUsage"`
	Total uint64 `json:"diskTotal"`
}

type disk struct {
	Usage uint64 `json:"diskUsage"`
	Total uint64 `json:"diskTotal"`
}

type net struct {
	In  uint64 `json:"netIn"`
	Out uint64 `json:"netOut"`
}

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Error("err: ", err)
		return ""
	}
	return usr.HomeDir
}

func New() (*Info, error) {
	mem, err := memoryStats()
	if err != nil {
		return nil, err
	}

	disk, err := diskStats()
	if err != nil {
		return nil, err
	}

	return &Info{
		State:       RUNNING,
		Uname:       runtime.GOOS,
		HomeDir:     HomeDir(),
		DiskUsage:   disk.Usage,
		DiskTotal:   disk.Total,
		MemoryTotal: mem.Total,
		MemoryUsage: mem.Usage,
	}, nil
}
