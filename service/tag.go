package service

import (
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/16/21 23:08
*@Describe
**/

func GetTagTitleById(tid int) (string, error) {
	db := global.G_DB
	result, err := db.Query("SELECT title FROM boil_tag WHERE id=?", tid)
	if err != nil {
		return "", err
	}
	tagTitle := ""
	result.Next()
	if err := result.Scan(&tagTitle); err != nil {
		return "", nil
	}
	result.Close()
	return tagTitle, nil
}

func GetTags(querySql string, args ...interface{}) (tags []model.TagVo, err error) {
	db := global.G_DB
	tags = []model.TagVo{}
	exec, err := db.Prepare(fmt.Sprintf("SELECT id,title FROM boil_tag WHERE %s", querySql))
	if err != nil {
		return
	}
	result, err := exec.Query(args...)
	if err != nil {
		return
	}
	for result.Next() {
		tagVo := model.TagVo{}
		if err := result.Scan(&tagVo.ID, &tagVo.Title); err != nil {
			return []model.TagVo{}, err
		}
		tagVo.BoilCount, _ = CountBoilByTag(tagVo.ID)
		tags = append(tags, tagVo)
	}
	result.Close()
	return
}

func InsertTag(tagTitle string) error {
	db := global.G_DB
	_, err := db.Exec("INSERT INTO boil_tag(title) value(?)", tagTitle)
	return err
}

func DeleteTagById(tid int) error {
	db := global.G_DB
	_, err := db.Exec("DELETE FROM boil_tag WHERE id=?", tid)
	return err
}

func CountBoilByTag(tid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_boil WHERE tag_id=?", tid)
	if err != nil {
		return 0, err
	}
	result.Next()
	if err := result.Scan(&count); err != nil {
		return 0, err
	}
	result.Close()
	return
}
