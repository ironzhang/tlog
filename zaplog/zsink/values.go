package zsink

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type values url.Values

func (v values) Get(key string) (string, bool) {
	if v == nil {
		return "", false
	}
	vs := v[key]
	if len(vs) == 0 {
		return "", false
	}
	return vs[0], true
}

func (v values) GetInt(key string) (int, bool, error) {
	s, ok := v.Get(key)
	if !ok {
		return 0, false, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, err
	}
	return n, true, nil
}

func (v values) GetDuration(key string) (time.Duration, bool, error) {
	s, ok := v.Get(key)
	if !ok {
		return 0, false, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false, err
	}
	return d, true, nil
}

func (v values) GetSize(key string) (int, bool, error) {
	s, ok := v.Get(key)
	if !ok {
		return 0, false, nil
	}
	sz, err := parseSize(s)
	if err != nil {
		return 0, false, err
	}
	return sz, true, nil
}

var unitMap = map[string]int{
	"":   1,
	"B":  1,
	"K":  1024,
	"KB": 1024,
	"M":  1024 * 1024,
	"MB": 1024 * 1024,
	"G":  1024 * 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
}

func parseSize(s string) (int, error) {
	sizestr, unitstr := splitSizeUnit(s)
	unit, ok := unitMap[unitstr]
	if !ok {
		return 0, fmt.Errorf("unknown unit %s in size %s", unitstr, s)
	}
	size, err := strconv.Atoi(sizestr)
	if err != nil {
		return 0, err
	}
	return size * unit, nil
}

func splitSizeUnit(s string) (size string, unit string) {
	for n := len(s) - 1; n >= 0; n-- {
		if s[n] >= '0' && s[n] <= '9' {
			return s[:n+1], strings.ToUpper(s[n+1:])
		}
	}
	return "", s
}
