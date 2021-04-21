package router

import (
	"github.com/patrickmn/go-cache"
	"login_service/common"
	"time"
)

var AppCache *common.AppCache

func init() {
	AppCache = &common.AppCache{
		Client: cache.New(10*time.Minute, 10*time.Minute),
	}
}
