package service

import (
	"context"
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"github.com/jmoiron/sqlx"
)

/**
*@Author lyer
*@Date 4/16/21 21:38
*@Describe
**/

func InsertBoil(boil model.Boil) error {
	_, err := global.G_DB.Exec(
		"INSERT INTO boil_boil(tag_id,user_id,content,create_time) VALUE(:tag_id,:user_id,:content,:create_time)",
		boil,
	)
	return err
}

func DeleteBoilById(bid int) error {
	if _, err := global.G_DB.Exec("DELETE FROM boil_boil WHERE id=?", bid); err != nil {
		return err
	}
	return ClearBoilUserLike(bid) //clear redis
}


func BoilArrToBoilVoArr(boilArr []model.Boil, loginUserId int) []model.BoilVo {
	var boilVoArr []model.BoilVo
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
	countRedis := global.RedisClient.SCard(context.Background(), fmt.Sprintf("boil::%d::like_users", bid)).Val()
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=?", bid)
	countDB := 0
	if err := r.Scan(&countDB); err != nil {
		return 0, err
	}
	return int(countRedis) + countDB, nil
}
func CountUserLikeBoil(uid int) (int, error) {
	redisCli := global.RedisClient
	db := global.G_DB
	countRedis := redisCli.SCard(context.Background(), fmt.Sprintf("user::%d::like_boils", uid)).Val()
	r := db.QueryRowx("SELECT COUNT(*) FROM boil_user_like_boil WHERE user_id=?", uid)
	countDB := 0
	if err := r.Scan(&countDB); err != nil {
		return 0, err
	}
	return int(countRedis) + countDB, nil
}

func BoilUserIsLike(bid, uid int) bool {
	if result := global.RedisClient.SIsMember(context.Background(), fmt.Sprintf("user::%d::like_boils", uid), bid).Val(); result {
		return result
	}
	//if redis does not exist then query mysql
	count := 0
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_user_like_boil WHERE boil_id=? AND user_id=?", bid, uid)
	if err := r.Scan(&count); err != nil {
		return false
	}
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

func BoilListUserLike(uid int) (boils []model.Boil, err error) {
	redisCli := global.RedisClient
	bids := redisCli.SMembers(context.Background(), fmt.Sprintf("user::%s::like_boils", uid)).Val()
	if err = global.G_DB.Select(&bids, "SELECT boil_id FROM boil_user_like_boil WHERE user_id=?", uid); err != nil {
		return []model.Boil{}, err
	}
	if len(bids) == 0 {
		return []model.Boil{}, nil
	}
	boilsSql, args, err := sqlx.In(
		"SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE id in (?) ORDER BY create_time DESC", bids)
	if err = global.G_DB.Select(&boils, boilsSql, args); err != nil {
		return []model.Boil{}, err
	}
	return
}

func InsertUserLikeBoil(uid, bid int) error {
	_, err := global.G_DB.Exec("INSERT INTO boil_user_like_boil(user_id,boil_id) VALUE (?,?)", uid, bid)
	return err
}
func DeleteUserLikeBoil(uid, bid int) error {
	_, err := global.G_DB.Exec("DELETE FROM boil_user_like_boil WHERE user_id=? AND boil_id=?", uid, bid)
	return err
}

//like end

func CountUserBoil(uid int) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_boil WHERE user_id=?", uid)
	err = r.Scan(&count)
	return
}

func GetBoilsByTagId(tid string) (boils []model.Boil, err error) {
	err = global.G_DB.Select(&boils, "SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE tag_id=? ORDER BY create_time DESC", tid)
	return
}
func GetBoilsByUserId(uid string) (boils []model.Boil, err error) {
	err = global.G_DB.Select(&boils, "SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE user_id=? ORDER BY create_time DESC", uid)
	return
}

func GetBoilsByUserIds(userIds []string) (boils []model.Boil, err error) {
	querySql, args, err := sqlx.In("SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE user_id in (?) ORDER BY create_time DESC", userIds)
	if err != nil {
		return
	}
	err = global.G_DB.Select(&boils, querySql, args)
	return
}

func GetAllBoil() (boils []model.Boil, err error) {
	err = global.G_DB.Select(&boils, "SELECT id,tag_id,user_id,create_time,content FROM boil_boil ORDER BY create_time DESC")
	return
}

func GetBoilById(bid string) (boil model.Boil, err error) {
	err = global.G_DB.Get(&boil, "SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE id=? ORDER BY create_time DESC", bid)
	return
}

func GetBoilsByIds(bids []string) (boils []model.Boil, err error) {
	querySql, args, err := sqlx.In("SELECT id,tag_id,user_id,create_time,content FROM boil_boil WHERE id in (?) ORDER BY create_time DESC", bids)
	if err != nil {
		return
	}
	err = global.G_DB.Select(&boils, querySql, args)
	return
}
