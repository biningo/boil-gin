package service

import (
	"context"
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"log"
	"strconv"
	"strings"
)

/**
*@Author lyer
*@Date 4/16/21 21:38
*@Describe
**/

func InsertBoil(boil model.Boil) error {
	db := global.G_DB
	_, err := db.Exec("INSERT INTO boil_boil(tag_id,user_id,content,create_time) VALUE(?,?,?,?)",
		boil.TagID, boil.UserID, boil.Content, boil.CreateTime)
	return err
}

func DeleteBoilById(bid int) error {
	db := global.G_DB
	_, err := db.Exec("DELETE FROM boil_boil WHERE id=?", bid)
	if err != nil {
		return err
	}
	return ClearBoilUserLike(bid)
}

//Also can ById
func GetBoils(querySql string, args ...interface{}) ([]model.Boil, error) {
	boilArr := []model.Boil{}
	boil := model.Boil{}
	db := global.G_DB
	strSql := fmt.Sprintf("SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE %s ORDER BY create_time DESC", querySql)
	exec, err := db.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	result, err := exec.Query(args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilArr = append(boilArr, boil)
	}
	result.Close()
	return boilArr, nil
}

func BoilArrToBoilVoArr(boilArr []model.Boil, loginUserId int) []model.BoilVo {
	boilVoArr := []model.BoilVo{}
	for _, boil := range boilArr {
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.CommentCount, _ = CountBoilComment(boil.ID)
		boilVo.LikeCount, _ = CountBoilLike(boil.ID)
		boilVo.IsLike = BoilUserIsLike(boil.ID, loginUserId)
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID
		boilVo.TagTitle, _ = GetTagTitleById(boilVo.TagID)

		user, _ := GetUserById(boilVo.UserID)
		boilVo.UserName = user.UserName
		boilVo.UserBio = user.Bio
		boilVo.UserAvatarId = user.AvatarID
		boilVoArr = append(boilVoArr, boilVo)
	}
	return boilVoArr
}

//like
func BoilUserLike(bid int, uid int) error {
	redisCli := global.RedisClient
	_, err := redisCli.SAdd(context.Background(), fmt.Sprintf("user:%d_like_boils", uid), bid).Result()
	if err != nil {
		return InsertUserLikeBoil(uid, bid)
	}
	IncrBoilLikeCount(bid)
	return nil
}
func BoilUserUnLike(bid, uid int) error {
	redisCli := global.RedisClient
	r, err := redisCli.SRem(context.Background(), fmt.Sprintf("user:%d_like_boils", uid), 0, bid).Result()
	if r == 0 || err != nil {
		err := DeleteUserLikeBoil(uid, bid)
		if err != nil {
			return err
		}
	}
	if err == nil {
		DecrBoilLikeCount(bid)
	}
	return nil
}
func IncrBoilLikeCount(bid int) {
	redisCli := global.RedisClient
	redisCli.HIncrBy(context.Background(), "boil_like_count", strconv.Itoa(bid), 1)
}
func DecrBoilLikeCount(bid int) {
	redisCli := global.RedisClient
	count, _ := redisCli.HGet(context.Background(), "boil_like_count", strconv.Itoa(bid)).Result()
	if count == "0" || count == "" {
		return
	}
	redisCli.HIncrBy(context.Background(), "boil_like_count", strconv.Itoa(bid), -1)
}

func CountBoilLike(bid int) (int, error) {
	redisCli := global.RedisClient
	result, _ := redisCli.HGet(context.Background(), "boil_like_count", strconv.Itoa(bid)).Result()
	count, _ := strconv.Atoi(result)
	if count < 0 {
		count = 0
	}
	//mysql
	db := global.G_DB
	r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=?", bid)
	if err != nil {
		return 0, err
	}
	countDB := 0
	r.Next()
	r.Scan(&countDB)
	return count + countDB, nil
}
func CountUserLikeBoil(uid int) (int, error) {
	redisCli := global.RedisClient
	count, _ := redisCli.SCard(context.Background(), fmt.Sprintf("user:%d_like_boils", uid)).Result()
	db := global.G_DB
	r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE user_id=?", uid)
	if err != nil {
		return 0, err
	}
	r.Next()
	countDB := 0
	r.Scan(&countDB)
	return int(count) + countDB, err
}
func BoilUserIsLike(bid, uid int) bool {
	redisCli := global.RedisClient
	result, _ := redisCli.SIsMember(context.Background(), fmt.Sprintf("user:%d_like_boils", uid), bid).Result()
	//mysql
	if result == false {
		db := global.G_DB
		count := 0
		r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=? AND user_id=?", bid, uid)
		if err != nil {
			return false
		}
		r.Next()
		r.Scan(&count)
		if count > 0 {
			result = true
		}
	}
	return result
}
func ClearBoilUserLike(bid int) error {
	redisCli := global.RedisClient
	redisCli.HDel(context.Background(), "boil_like_count", strconv.Itoa(bid))
	keys, err := redisCli.Keys(context.Background(), fmt.Sprintf("*_like_boils")).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		redisCli.SRem(context.Background(), key, bid)
	}
	return nil
}

func BoilListUserLike(uid int) ([]model.Boil, error) {
	redisCli := global.RedisClient
	bids, _ := redisCli.SMembers(context.Background(), fmt.Sprintf("user:%d_like_boils", uid)).Result()
	db := global.G_DB
	r, err := db.Query("SELECT boil_id FROM boil_user_like_boil WHERE user_id=?", uid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for r.Next() {
		bid := 0
		r.Scan(&bid)
		bids = append(bids, strconv.Itoa(bid))
	}
	if len(bids) == 0 {
		return nil, nil
	}
	boilArr, err := GetBoils("id in " + "(" + strings.Join(bids, ",") + ")")
	if err != nil {
		return nil, err
	}
	return boilArr, nil
}

func InsertUserLikeBoil(uid, bid int) error {
	db := global.G_DB
	_, err := db.Exec("INSERT INTO boil_user_like_boil(user_id,boil_id) VALUE (?,?)", uid, bid)
	return err
}
func DeleteUserLikeBoil(uid, bid int) error {
	db := global.G_DB
	_, err := db.Exec("DELETE FROM boil_user_like_boil WHERE user_id=? AND boil_id=?", uid, bid)
	return err
}

//like

func CountUserBoil(uid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_boil WHERE user_id=?", uid)
	if err != nil {
		return
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}
