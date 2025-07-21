package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type FileConfigStore map[string]interface{}

func (c FileConfigStore) GetValue(key string) (any, error) {
	key = strings.TrimSpace(key)

	// when the key is in the format of scope.name, extract scope and name, local file ignores scope
	nameList := strings.Split(key, ".")
	var name string
	if len(nameList) == 1 {
		name = key
	} else if len(nameList) == 2 {
		name = nameList[1]
	}

	return c[name], nil
}

func (c FileConfigStore) GetString(key string) (string, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return "", err
	}
	if stringValue, ok := value.(string); ok {
		return stringValue, nil
	}
	return "", nil
}

func (c FileConfigStore) GetBool(key string) (bool, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return false, err
	}
	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}
	return false, nil
}

func (c FileConfigStore) MustGetString(key string) string {
	value, err := c.GetString(key)
	if err != nil {
		logrus.Fatalf("配置项[%s]不存在3", key)
	}
	return value
}

func (c FileConfigStore) GetInt64(key string) (int64, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case string:
		stringValue, ok := value.(string)
		if !ok {
			return 0, nil
		}
		intValue, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	}
	return 0, fmt.Errorf("配置项[%s]不存在或格式有误", key)
}

func ParseConfigContent(fileContent string) (FileConfigStore, error) {
	configMap := make(map[string]any)

	err := yaml.Unmarshal([]byte(fileContent), &configMap)
	if err != nil {
		return nil, fmt.Errorf("解析配置内容出错: %w", err)
	}

	var cmdEnv []string

	osEnv := os.Environ()
	cmdEnv = append(cmdEnv, osEnv...)
	for _, e := range cmdEnv {
		index := strings.Index(e, "=")
		if index > 0 {
			configMap[e[:index]] = e[index+1:]
		}
	}

	var model FileConfigStore = configMap

	var filePrefix = "include://"
	var contentPrefix = "content://"
	for key, value := range configMap {
		stringValue, ok := value.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(stringValue, filePrefix) {
			valueData, err := os.ReadFile(value.(string)[len(filePrefix):])
			if err != nil {
				return nil, fmt.Errorf("读取配置文件出错: %w", err)
			}
			value = string(valueData)
			configMap[key] = value
		} else if strings.HasPrefix(value.(string), contentPrefix) {
			value = value.(string)[len(contentPrefix):]
			configMap[key] = value
		}
	}

	return model, nil
}

func ParseConfigFile(filePath string) (FileConfigStore, error) {
	var model FileConfigStore

	if _, err := os.Stat(filePath); err == nil {
		configData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件出错: %w", err)
		}
		model, err = ParseConfigContent(string(configData))
		if err != nil {
			return nil, fmt.Errorf("解析配置文件出错: %w", err)
		}
	}

	return model, nil
}
