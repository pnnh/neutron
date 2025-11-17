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

	urlList := strings.Split(configUrl, ",")
	if len(urlList) > 2 {
		return fmt.Errorf("invalid config url: %s", configUrl)
	}
	if len(urlList) == 1 {
		store, err := configUrlToStore(urlList[0], project, app, env, svc)
		if err != nil {
			return err
		}
		appConfigStore = store

	} else if len(urlList) == 2 {
		originStore, err := configUrlToStore(urlList[0], project, app, env, svc)
		if err != nil {
			return err
		}
		overrideStore, err := configUrlToStore(urlList[1], project, app, env, svc)
		if err != nil {
			return err
		}
		appConfigStore = v2.NewOverrideConfigStore(originStore, overrideStore)
	} else {
		return fmt.Errorf("invalid config url: %s", configUrl)
	}

	if Debug() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}

func configUrlToStore(configUrl string, project, app, env, svc string) (v2.IConfigStore, error) {
	if strings.HasPrefix(configUrl, "galaxy://") {
		galaxyUrl := strings.Replace(configUrl, "galaxy://", "http://", 1)
		return v2.NewGalaxyConfigStore(galaxyUrl, project, app, env, svc), nil
	} else if strings.HasPrefix(configUrl, "pggo://") {
		pgUrl := strings.Replace(configUrl, "pggo://", "", 1)
		pgStore, err := v2.NewPgConfigStore(pgUrl, project, app, env, svc)
		if err != nil {
			return nil, fmt.Errorf("configUrlToStore NewPgConfigStore: %w", err)
		}
		return pgStore, nil
	} else if strings.HasPrefix(configUrl, "file:") {
		fileStore, err := v2.ParseConfigFile(configUrl)
		if err != nil {
			return nil, fmt.Errorf("configUrlToStore ParseConfigFile: %w", err)
		}
		return fileStore, nil
	} else {
		return nil, fmt.Errorf("unsupported config url: %s", configUrl)
	}
}
func GetConfiguration(key interface{}) (interface{}, bool) {
	if key, ok := key.(string); ok {
		if value, err := appConfigStore.GetValue(key); err == nil {
			return value, true
		} else {
			logrus.Errorf("GetConfiguration key=%s error: %v", key, err)
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
		logrus.Fatalf("MustGetConfigurationString GetConfigurationString not exists: %s", key)
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

// Debug returns true if the application is running in debug mode.
func Debug() bool {
	modeValue := os.Getenv("GXMODE")
	if modeValue != "" {
		return modeValue == "DEBUG" || modeValue == "debug"
	}
	return false
}

// Testing returns true if the application is running in testing mode.
func Testing() bool {
	modeValue := os.Getenv("GXMODE")
	if modeValue != "" {
		return modeValue == "TEST" || modeValue == "test"
	}
	return false
}

// Release returns true if the application is running in release mode.
func Release() bool {
	return !Debug() && !Testing()
}

func GetEnvName() string {
	envValue := os.Getenv("GXENV")
	if envValue != "" {
		return envValue
	}
	if Debug() {
		return "development"
	}
	return "production"
}
