package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"neutron/services/datastore"
	"neutron/services/strutil"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type PgConfigStore struct {
	cache   *cache.Cache
	pgUrl   string
	env     string
	svc     string
	project string
	app     string
}

const pgDbName = "configdb"

func NewPgConfigStore(pgUrl string, project, app, env, svc string) (*PgConfigStore, error) {
	err := datastore.InitFor(pgDbName, pgUrl)
	if err != nil {
		return nil, fmt.Errorf("初始化配置数据库失败: %w", err)
	}
	pgStore := &PgConfigStore{
		pgUrl:   pgUrl,
		cache:   cache.New(time.Second*30, time.Second*60),
		project: project,
		app:     app,
		env:     env,
		svc:     svc,
	}
	logrus.Infof("配置数据库连接成功")
	return pgStore, nil
}

func (c *PgConfigStore) GetValue(key string) (configValue any, getError error) {
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
	baseSqlText := ` select c.content 
from galaxy.configuration c `
	whereText := ` where c.name = :name `

	if !strutil.IsValidName(c.svc) || !strutil.IsValidName(c.env) || !strutil.IsValidName(name) ||
		!strutil.IsValidName(c.project) || !strutil.IsValidName(c.app) || !strutil.IsValidName(scope) {
		return "", fmt.Errorf("invalid parameter format")
	}
	if c.project != "huable" && c.project != "weable" && c.project != "calieo" {
		return "", fmt.Errorf("invalid project 2")
	}
	sqlParams := map[string]interface{}{
		"environment": c.env,
		"name":        name,
	}
	//if scope == "project" || scope == "app" || scope == "svc" {
	//baseSqlText += ` join projects p on c.project = p.uid `
	//whereText += ` and p.name = :project `
	//sqlParams["project"] = c.project
	//if scope == "app" || scope == "svc" {
	//baseSqlText += ` join applications a on c.application = a.uid `
	//whereText += ` and a.name = :application `
	//sqlParams["application"] = c.app
	//if scope == "svc" {
	//	baseSqlText += ` join services s on c.service = s.uid `
	//	whereText += ` and s.name = :service `
	//	sqlParams["service"] = c.svc
	//}
	//}
	//} else {
	//	return "", fmt.Errorf("invalid scope")
	//}
	fullSqlText := fmt.Sprintf("%s%s%s", baseSqlText, whereText, " limit 1;")

	sqlResults := map[string]interface{}{}

	rows, err := datastore.NamedQueryFor(pgDbName, fullSqlText, sqlParams)
	if err != nil {
		return "", fmt.Errorf("NamedQuery: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logrus.Warnf("rows.Close: %v", closeErr)
		}
	}()
	for rows.Next() {
		if err = sqlx.MapScan(rows, sqlResults); err != nil {
			return "", fmt.Errorf("StructScan: %w", err)
		}
		if len(sqlResults) == 0 {
			return "", ErrConfigNotFound
		}
		tableMap := datastore.MapToDataRow(sqlResults)
		contentValue := tableMap.GetNullString("content")
		if !contentValue.Valid {
			return "", ErrConfigNotFound
		}
		return contentValue.String, nil
	}
	return "", ErrConfigNotFound
}

func (c *PgConfigStore) GetString(key string) (string, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return "", err
	}
	if stringValue, ok := value.(string); ok {
		return stringValue, nil
	}
	return "", nil
}

func (c *PgConfigStore) GetBool(key string) (bool, error) {
	value, err := c.GetValue(key)
	if err != nil {
		return false, err
	}
	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}
	return false, nil
}

func (c *PgConfigStore) MustGetString(key string) string {
	value, err := c.GetString(key)
	if err != nil {
		logrus.Fatalf("配置项[%s]不存在2", key)
	}
	return value
}

func (c *PgConfigStore) GetInt64(key string) (int64, error) {
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
