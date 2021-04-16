package service

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/16/21 21:38
*@Describe
**/


func GetBoilById(bid int) model.Boil{
	db:=global.G_DB
	result, _ := db.Query("select id,tag_id,user_id,create_time,content from boil_boil where id=?", bid)
	boil := model.Boil{}
	defer result.Close()
	if result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
	}
	return boil
}