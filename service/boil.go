package service

import (
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
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
	return err
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

func BoilArrToBoilVoArr(boilArr []model.Boil) []model.BoilVo {
	boilVoArr := []model.BoilVo{}
	for _, boil := range boilArr {
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.CommentCount, _ = CountBoilComment(boil.ID)
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
