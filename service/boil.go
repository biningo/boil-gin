package service

import (
	"context"
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
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
	for result.Next() {
		err := result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		if err != nil {
			return []model.Boil{}, err
		}
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
		boilVo.UserIsFollow, _ = IsUserFollow(loginUserId, boil.UserID)
		boilVoArr = append(boilVoArr, boilVo)
	}
	return boilVoArr
}

//like
func BoilUserLike(bid int, uid int) error {
	redisCli := global.RedisClient
	redisCli.SAdd(context.Background(), fmt.Sprintf("user::%d::like_boils", uid), bid)
	redisCli.SAdd(context.Background(), fmt.Sprintf("boil::%d::like_users", bid), uid)
	return nil
}
func BoilUserUnLike(bid, uid int) error {
	redisCli := global.RedisClient
	redisCli.SRem(context.Background(), fmt.Sprintf("boil::%d::like_users", bid), 0, uid)
	if exist := redisCli.SRem(context.Background(), fmt.Sprintf("user::%d::like_boils", uid), 0, bid).Val(); exist == 0 {
		return DeleteUserLikeBoil(uid, bid)
	}
	return nil
}

func CountBoilLike(bid int) (int, error) {
	redisCli := global.RedisClient
	db := global.G_DB
	//redis
	countRedis := redisCli.SCard(context.Background(), fmt.Sprintf("boil::%d::like_users", bid)).Val()
	//mysql
	r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=?", bid)
	if err != nil {
		return 0, err
	}
	countDB := 0
	r.Next()
	err = r.Scan(&countDB)
	if err != nil {
		return 0, err
	}
	r.Close()
	return int(countRedis) + countDB, nil
}
func CountUserLikeBoil(uid int) (int, error) {
	redisCli := global.RedisClient
	db := global.G_DB

	countRedis := redisCli.SCard(context.Background(), fmt.Sprintf("user::%d::like_boils", uid)).Val()
	r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE user_id=?", uid)
	if err != nil {
		return 0, err
	}
	r.Next()
	countDB := 0
	err = r.Scan(&countDB)
	if err != nil {
		return 0, err
	}
	r.Close()
	return int(countRedis) + countDB, err
}

func BoilUserIsLike(bid, uid int) bool {
	redisCli := global.RedisClient
	result := redisCli.SIsMember(context.Background(), fmt.Sprintf("user::%d::like_boils", uid), bid).Val()
	if result {
		return result
	}
	//if redis does not exist then query mysql
	db := global.G_DB
	count := 0
	r, err := db.Query("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=? AND user_id=?", bid, uid)
	if err != nil {
		return false
	}
	r.Next()
	err = r.Scan(&count)
	if err != nil {
		return false
	}
	r.Close()
	return count > 0
}

func ClearBoilUserLike(bid int) error {
	redisCli := global.RedisClient
	userIds := redisCli.SMembers(context.Background(), fmt.Sprintf("boil::%d::like_users", bid)).Val()
	redisCli.Del(context.Background(), fmt.Sprintf("boil::%d::like_users", bid))
	for _, uid := range userIds {
		redisCli.SRem(context.Background(), fmt.Sprintf("user::%s::like_boils", uid), bid)
	}
	return nil
}

func BoilListUserLike(uid int) ([]model.Boil, error) {
	redisCli := global.RedisClient
	bids := redisCli.SMembers(context.Background(), fmt.Sprintf("user::%s::like_boils", uid)).Val()
	db := global.G_DB
	r, err := db.Query("SELECT boil_id FROM boil_user_like_boil WHERE user_id=?", uid)
	if err != nil {
		return nil, err
	}
	for r.Next() {
		bid := 0
		err := r.Scan(&bid)
		if err != nil {
			return []model.Boil{}, err
		}
		bids = append(bids, strconv.Itoa(bid))
	}
	r.Close()
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

//like end

func CountUserBoil(uid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_boil WHERE user_id=?", uid)
	if err != nil {
		return
	}
	result.Next()
	err = result.Scan(&count)
	if err != nil {
		return 0, err
	}
	result.Close()
	return
}
