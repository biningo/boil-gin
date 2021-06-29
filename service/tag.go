package service

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/16/21 23:08
*@Describe
**/

func GetTagTitleById(tid int) (string, error) {
	var tag model.Tag
	if err := global.G_DB.Get(&tag, "SELECT title FROM boil_tag WHERE id=?", tid); err != nil {
		return "", err
	}
	return tag.Title, nil
}

func InsertTag(tagTitle string) error {
	_, err := global.G_DB.Exec("INSERT INTO boil_tag(title) value(?)", tagTitle)
	return err
}

func DeleteTagById(tid int) error {
	_, err := global.G_DB.Exec("DELETE FROM boil_tag WHERE id=?", tid)
	return err
}

func CountBoilByTag(tid int) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_boil WHERE tag_id=?", tid)
	err = r.Scan(&count)
	return
}

func GetAllTags() (tags []model.Tag, err error) {
	err = global.G_DB.Select(&tags, "SELECT id,title FROM boil_tag")
	return
}
