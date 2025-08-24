package config

import (
	"fmt"
	"neutron/services/convert"

	"github.com/sirupsen/logrus"
)

type OverrideConfigStore struct {
	originStore IConfigStore
	maskStore   IConfigStore
}

func NewOverrideConfigStore(originStore, maskStore IConfigStore) OverrideConfigStore {
	return OverrideConfigStore{
		originStore: originStore,
		maskStore:   maskStore,
	}
}

func (c OverrideConfigStore) GetValue(key string) (any, error) {
	// try to get from maskStore first
	if c.maskStore != nil {
		value, err := c.maskStore.GetValue(key)
		if err != nil {
			return nil, fmt.Errorf("maskStore.GetValue: %w", err)
		}
		if value != nil {
			return value, nil
		}
	}
	// then get from originStore
	if c.originStore == nil {
		return nil, fmt.Errorf("MixedConfigStore.originStore is nil")
	}
	value, err := c.originStore.GetValue(key)
	if err != nil {
		return nil, fmt.Errorf("originStore.GetValue: %w", err)
	}
	if value == nil {
		return nil, ErrConfigNotFound
	} else {
		return value, nil
	}
}

func (c OverrideConfigStore) GetString(key string) (string, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return "", fmt.Errorf("MixedConfigStore GetValue: %w", err)
	}
	strValue, err := convert.ConvertString(value)
	if err != nil {
		return "", fmt.Errorf("MixedConfigStore ConvertString: %w", err)
	}
	return strValue, nil
}

func (c OverrideConfigStore) GetBool(key string) (bool, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return false, fmt.Errorf("MixedConfigStore GetBool: %w", err)
	}
	boolValue, err := convert.ConvertBool(value)
	if err != nil {
		return false, fmt.Errorf("MixedConfigStore ConvertBool: %w", err)
	}
	return boolValue, nil
}

func (c OverrideConfigStore) MustGetString(key string) string {
	value, err := c.GetString(key)
	if err != nil {
		logrus.Fatalf("配置项[%s]不存在: %v", key, err)
	}
	return value
}

func (c OverrideConfigStore) GetInt64(key string) (int64, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return 0, fmt.Errorf("MixedConfigStore GetInt64: %w", err)
	}
	intValue, err := convert.ToInt64(value)
	if err != nil {
		return 0, fmt.Errorf("MixedConfigStore ToInt64: %w", err)
	}
	return intValue, nil
}
