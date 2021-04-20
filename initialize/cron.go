package initialize

import (
	"context"
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/service"
	"strconv"
	"strings"
	"time"
)

/**
*@Author lyer
*@Date 4/20/21 11:11
*@Describe
**/

func InitRedisToMySqlCron(duration time.Duration) {
	go func() {
		for {
			select {
			case <-time.Tick(duration):
				redisCli := global.RedisClient
				redisCli.Del(context.Background(), "boil_like_count")
				keys, err := redisCli.Keys(context.Background(), fmt.Sprintf("*_like_boils")).Result()
				if err != nil {
					continue
				}
				for _, key := range keys {
					arr := strings.Split(key, "_")
					if len(arr) < 1 {
						continue
					}
					suid := strings.Split(arr[0], ":")[1]
					uid, _ := strconv.Atoi(suid)
					bids, err := redisCli.SMembers(context.Background(), key).Result()
					if err != nil {
						continue
					}
					for _, sbid := range bids {
						bid, _ := strconv.Atoi(sbid)
						service.InsertUserLikeBoil(uid, bid)
					}
					redisCli.Del(context.Background(), key)
				}
			}
		}
	}()
}
