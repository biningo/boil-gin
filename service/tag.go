package service

import "github.com/biningo/boil-gin/global"

/**
*@Author lyer
*@Date 4/16/21 23:08
*@Describe
**/

func GetTagTitleById(tid int) (string, error) {
	db := global.G_DB
	result, err := db.Query("select title from boil_tag where id=?", tid)
	if err != nil {
		return "", err
	}
	tagTitle := ""
	result.Next()
	result.Scan(&tagTitle)
	return tagTitle, nil
}
