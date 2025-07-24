package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"neutron/models"
	"neutron/services/strutil"
)

type GalaxyConfigStore struct {
	cache      *cache.Cache
	galaxyUrl  string
	env        string
	svc        string
	httpClient *http.Client
	project    string
	app        string
}

func NewGalaxyConfigStore(galaxyUrl string, project, app, env, svc string) GalaxyConfigStore {
	return GalaxyConfigStore{
		galaxyUrl: galaxyUrl,
		cache:     cache.New(time.Second*30, time.Second*60),
		project:   project,
		app:       app,
		env:       env,
		svc:       svc,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type GalaxyConfigData struct {
	Project string `json:"project"`
	App     string `json:"app"`
	Env     string `json:"env"`
	Svc     string `json:"svc"`
	Scope   string `json:"scope"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

func (c GalaxyConfigStore) GetValue(key string) (configValue any, getError error) {
	key = strings.TrimSpace(key)

	nameList := strings.Split(key, ".")
	var scope, name string
	if len(nameList) == 1 {
		scope = "svc" // default scope is "svc"
		name = key
	} else if len(nameList) == 2 {
		scope = nameList[0]
		name = nameList[1]
	} else {
		return nil, fmt.Errorf("invalid key format: %s, expected format is 'scope.name' or name", key)
	}

	if !strutil.IsValidName(scope) || !strutil.IsValidName(name) {
		return nil, fmt.Errorf("配置项[%s]格式不正确", key)
	}
	getUrl := fmt.Sprintf("%s/config?project=%s&app=%s&env=%s&svc=%s&scope=%s&name=%s",
		c.galaxyUrl, c.project, c.app, c.env, c.svc, scope, name)

	cacheValue, found := c.cache.Get(getUrl)
	if found {
		return cacheValue, nil
	}

	newRequest, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return false, fmt.Errorf("http.NewRequest: %w", err)
	}

	res, err := c.httpClient.Do(newRequest)
	if err != nil {
		return false, fmt.Errorf("client.Do: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			getError = fmt.Errorf("关闭Body失败: %w", err)
		}
	}(res.Body)

	respResult := &models.NECommonResult{
		Data: &GalaxyConfigData{},
	}
	derr := json.NewDecoder(res.Body).Decode(respResult)
	if derr != nil {
		return false, fmt.Errorf("json.NewDecoder: %w", derr)
	}
	if respResult.Code != models.NECodeOk || respResult.Data == nil {
		return nil, fmt.Errorf("获取配置失败: %s", respResult.Message)
	}
	configData, ok := respResult.Data.(*GalaxyConfigData)
	if !ok {
		return nil, fmt.Errorf("invalid response data type, expected GalaxyConfigData")
	}
	c.cache.Set(getUrl, configData.Value, cache.DefaultExpiration)

	return configData.Value, nil
}

func (c GalaxyConfigStore) GetString(key string) (string, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return "", err
	}
	if stringValue, ok := value.(string); ok {
		return stringValue, nil
	}
	return "", nil
}

func (c GalaxyConfigStore) GetBool(key string) (bool, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return false, err
	}
	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}
	return false, nil
}

func (c GalaxyConfigStore) MustGetString(key string) string {
	value, err := c.GetString(key)
	if err != nil {
		logrus.Fatalf("配置项[%s]不存在2", key)
	}
	return value
}

func (c GalaxyConfigStore) GetInt64(key string) (int64, error) {
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
	return 0, fmt.Errorf("配置项[%s]不存在或格式有误2", key)
}
