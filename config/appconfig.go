package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	v2 "neutron/config/v2"

	"github.com/sirupsen/logrus"
)

var appConfigStore v2.IConfigStore

func InitAppConfig(configUrl string, project, app, env, svc string) error {
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     false,
		TimestampFormat: time.RFC3339,
	})

	if strings.HasPrefix(configUrl, "galaxy://") {
		galaxyUrl := strings.Replace(configUrl, "galaxy://", "http://", 1)
		appConfigStore = v2.NewGalaxyConfigStore(galaxyUrl, project, app, env, svc)
	} else {
		fileStore, err := v2.ParseConfigFile(configUrl)
		if err != nil {
			return fmt.Errorf("配置文件解析失败: %w", err)
		}
		appConfigStore = fileStore
	}

	if Debug() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Println("日志级别:", logrus.GetLevel())
	return nil
}

func GetConfiguration(key interface{}) (interface{}, bool) {
	if key, ok := key.(string); ok {
		if value, err := appConfigStore.GetValue(key); err == nil {
			return value, true
		}
	}
	return nil, false
}

func GetConfigurationString(key interface{}) (string, bool) {
	if key, ok := key.(string); ok {
		if value, err := appConfigStore.GetString(key); err == nil {
			return value, true
		}
	}
	return "", false
}

func MustGetConfigurationString(key interface{}) string {
	value, ok := GetConfigurationString(key)
	if !ok {
		logrus.Fatalf("配置项[%s]不存在1", key)
	}
	return value
}

func GetConfigurationInt64(key interface{}) (int64, bool) {
	if key, ok := key.(string); ok {
		if value, err := appConfigStore.GetInt64(key); err == nil {
			return value, true
		}
	}
	return 0, false
}

func GetConfigOrDefaultInt64(key interface{}, defaultValue int64) int64 {
	value, ok := GetConfigurationInt64(key)
	if !ok {
		return defaultValue
	}
	return value
}

func Debug() bool {
	// check if environment variable is set
	if os.Getenv("DEBUG") == "true" {
		return true
	}
	if os.Getenv("MODE") == "DEBUG" {
		return true
	}
	//mode, ok = GetConfiguration("RUN_MODE")
	//if ok && mode == "development" {
	//	return true
	//}
	return false
}

func GetEnvName() string {
	if Debug() {
		return "development"
	}
	return "production"
}
