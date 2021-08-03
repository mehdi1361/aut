package router

import (
	"aut/common"
	"github.com/patrickmn/go-cache"
	"time"
)

var AppCache *common.AppCache

func init() {
	AppCache = &common.AppCache{
		Client: cache.New(10*time.Minute, 10*time.Minute),
	}
}
