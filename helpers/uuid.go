package helpers

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewUuid() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func MustUuid() string {
	id, err := uuid.NewUUID()
	if err != nil {
		logrus.Fatalln("初始化uuid出错", err)
		return ""
	}
	return id.String()
}

// StringToMD5UUID converts a string to a MD5 hash and then to a UUID
// Dreprecated: Use `MustUuid` instead.
// 把md5换成uuid格式，有可能会产生非法的uuid，所以该函数弃用
//func StringToMD5UUID(input string) (string, error) {
//	// Compute MD5 hash
//	hash := md5.Sum([]byte(input))
//
//	// Convert hash to UUID
//	uuid, err := uuid.FromBytes(hash[:])
//	if err != nil {
//		return "", err
//	}
//
//	return uuid.String(), nil
//}

func IsUuid(stringValue string) bool {
	_, err := uuid.Parse(stringValue)
	return err == nil
}

func EmptyUuid() string {
	return "00000000-0000-0000-0000-000000000000"
}
