package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"os/exec"
	"strings"
	"time"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func UuidGenerator() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(fmt.Sprintf("%s", out), "\n")
}

type AppCache struct {
	Client *cache.Cache
	Index  []string
}

func (r *AppCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, find := r.Find(key)
	if find {
		return fmt.Errorf("key %s already exist", key)
	}
	r.Client.Set(key, b, expiration)
	r.Index = append(r.Index, key)
	return nil
}

func (r *AppCache) Get(key string) ([]byte, error) {
	res, exist := r.Client.Get(key)
	if !exist {
		return nil, errors.New("")
	}

	resByte, ok := res.([]byte)
	if !ok {
		return nil, errors.New("format is not arr of bytes")
	}

	return resByte, nil
}

func (r *AppCache) Delete(key string) {
	r.Client.Delete(key)
}

func (r *AppCache) KeyContain(subStr string) []string {
	var result []string
	for _, item := range r.Index {
		if strings.Contains(item, subStr) {
			result = append(result, item)
		}
	}
	return result
}

func (r *AppCache) Find(val string) (int, bool) {
	for i, item := range r.Index {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func (r AppCache) KeyContainDelete(key string) {
	for _, val := range r.KeyContain(key) {
		r.Delete(val)
	}
}
