package config

import "errors"

var ErrConfigNotFound = errors.New("config not found")

type IConfigStore interface {
	GetValue(key string) (any, error)
	GetString(key string) (string, error)
	GetBool(key string) (bool, error)
	MustGetString(key string) string
	GetInt64(key string) (int64, error)
}
