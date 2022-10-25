package apihandlers

import (
	"errors"
	"sync"

	"github.com/ivankuchin/timecard.ru-api/logs"
)

var errKeyNotFound = errors.New("key not found")

var cs = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key, val string) error {
	cs.Lock()
	cs.m[key] = val
	cs.Unlock()

	return nil
}

func Get(key, tID string) (string, error) {
	cs.RLock()
	val, err := cs.m[key]
	cs.RUnlock()

	if !err {
		logs.Sugar.Infow("key("+key+") not found in captcha store", "traceID", tID)
		return "", nil
	}

	return val, nil
}

func Delete(key, tID string) error {
	val, err := Get(key, tID)

	if err != nil {
		logs.Sugar.Errorw(err.Error(), "traceID", tID)
		return err
	}

	if val == "" {
		logs.Sugar.Infow("key("+key+") not found in captcha store", "traceID", tID)
		return nil
	}

	cs.Lock()
	delete(cs.m, key)
	cs.Unlock()

	return nil
}
