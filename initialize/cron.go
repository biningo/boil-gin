package initialize

import (
	"context"
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
				keys := redisCli.Keys(context.Background(), "%::like_boils").Val()
				for _, key := range keys {
					uid := strings.Split(key, "::")[1]
					bids := redisCli.SMembers(context.Background(), key).Val()
					for _, bid := range bids {
						userId, _ := strconv.Atoi(uid)
						boilId, _ := strconv.Atoi(bid)
						service.InsertUserLikeBoil(userId, boilId)
					}
				}
				//service.ClearBoilUserLike()
			}
		}
	}()
}
