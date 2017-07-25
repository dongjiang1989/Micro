// +build darwin freebsd linux netbsd

package systeminfo

// mac machine
import (
	"bytes"
	"encoding/binary"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
	"unsafe"

	log "github.com/cihub/seelog"
)

func diskStats() (*disk, error) {
	d := new(disk)
	stat := new(syscall.Statfs_t)

	if err := syscall.Statfs("/", stat); err != nil {
		return nil, err
	}

	bsize := stat.Bsize / 512

	d.Total = (uint64(stat.Blocks) * uint64(bsize)) >> 1
	free := (uint64(stat.Bfree) * uint64(bsize)) >> 1
	d.Usage = d.Total - free

	return d, nil
}

// generic Sysctl buffer unmarshalling
func sysctlByname(name string, data interface{}) (err error) {
	val, err := syscall.Sysctl(name)
	if err != nil {
		log.Error("err: ", err)
		return err
	}

	buf := []byte(val)

	switch v := data.(type) {
	case *uint64:
		*v = *(*uint64)(unsafe.Pointer(&buf[0]))
		return
	}

	bbuf := bytes.NewBuffer([]byte(val))
	return binary.Read(bbuf, binary.LittleEndian, data)
}

func vm_stat() (bytesFree, bytesInactive uint64, err error) {
	type memStat struct {
		regex *regexp.Regexp
		value uint64
		valid bool
	}

	stats := map[string]*memStat{
		"pageSize":      {regex: regexp.MustCompile("page size of (\\d+) bytes")},
		"pagesFree":     {regex: regexp.MustCompile("Pages free: *(\\d+).")},
		"pagesInactive": {regex: regexp.MustCompile("Pages inactive: *(\\d+).")},
	}

	cmd := exec.Command("vm_stat")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	// Parse lines in vm_stat output
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		for _, stat := range stats {
			match := stat.regex.FindSubmatch(line)
			if match == nil {
				continue
			}

			stat.value, err = strconv.ParseUint(string(match[1]), 10, 64)
			if err != nil {
				return 0, 0, err
			}
			stat.valid = true
		}
	}

	// Check every stat is found in output
	for _, stat := range stats {
		if !stat.valid {
			return 0, 0, errors.New("cannot parse vm_stat output")
		}
	}

	pageSize := stats["pageSize"].value
	return pageSize * stats["pagesFree"].value, pageSize * stats["pagesInactive"].value, nil
}

func memoryStats() (*memory, error) {
	m := new(memory)
	err := sysctlByname("hw.memsize", &m.Total)
	if err != nil {
		return nil, err
	}

	free, inactive, err := vm_stat()
	if err != nil {
		return nil, err
	}

	m.Usage = m.Total - free - inactive

	return m, nil
}
