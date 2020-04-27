package store

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type NewStoreFunc func(path string) (SearchClient, error)


type Registration struct {
	Name        string // unique name
	Title       string // human-readable name
	FactoryFunc NewStoreFunc
}

var registry = make(map[string]*Registration)


func Register(reg *Registration) {
	if reg.Name == "" {
		zlog.Fatal("name cannot be blank")
	} else if _, ok := registry[reg.Name]; ok {
		zlog.Fatal("already registered", zap.String("name", reg.Name))
	}
	registry[reg.Name] = reg
}

func New(dsn string) (SearchClient, error) {
	chunks := strings.Split(dsn, ":")
	reg, found := registry[chunks[0]]
	if !found {
		return nil, fmt.Errorf("no such dmesh store registered %q", chunks[0])
	}
	return reg.FactoryFunc(dsn)
}

// ByName returns a registered store driver
func ByName(name string) *Registration {
	r, ok := registry[name]
	if !ok {
		return nil
	}
	return r
}
