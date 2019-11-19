package common

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"hash"
	"io"
)

//打印元素
func DD(obj ...interface{}) {
	fmt.Printf("%v\n", obj)
}

//md5
func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//sha1
func Sha1(str []byte) string {
	h := sha1.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

//计算文件hash
func FileHash(reader io.Reader, tp string) string {
	var result []byte
	var h hash.Hash
	if tp == "md5" {
		h = md5.New()
	} else {
		h = sha1.New()
	}
	if _, err := io.Copy(h, reader); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(result))
}

//生成uuid
func GetUuid() string {
	var u uuid.UUID
	var err error
	for i := 0; i < 3; i++ {
		u, err = uuid.NewUUID()
		if err == nil {
			return u.String()
		}
	}
	return ""
}

//生成uuid v4
func GetUuidV4() string {
	var u uuid.UUID
	var err error
	for i := 0; i < 3; i++ {
		u, err = uuid.NewRandom()
		if err == nil {
			return u.String()
		}
	}
	return ""
}
