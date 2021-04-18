package service

import (
	"fmt"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/17/21 09:58
*@Describe
**/

func GetCommentBidsByUid(uid int) ([]string, error) {
	db := global.G_DB
	bids := []string{}
	bid := "0"
	result, err := db.Query("SELECT boil_id FROM boil_comment WHERE user_id=? ORDER BY create_time DESC", uid)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		result.Scan(&bid)
		bids = append(bids, bid)
	}
	result.Close()
	return bids, nil
}

func InsertComment(comment model.Comment) error {
	db := global.G_DB
	_, err := db.Exec("INSERT INTO boil_comment(user_id,boil_id,content,create_time) VALUE(?,?,?,?)", comment.UserID, comment.BoilID, comment.Content, comment.CreateTime)
	return err
}

func DeleteCommentById(cid int) error {
	db := global.G_DB
	_, err := db.Exec("DELETE FROM boil_comment WHERE id=?", cid)
	return err
}

func GetComments(querySql string, args ...interface{}) ([]model.Comment, error) {
	db := global.G_DB
	strSql := fmt.Sprintf("SELECT id,user_id,boil_id,content,create_time FROM boil_comment WHERE %s ORDER BY create_time DESC", querySql)
	result, err := db.Query(strSql, args...)
	if err != nil {
		return nil, err
	}
	comments := []model.Comment{}
	comment := model.Comment{}
	for result.Next() {
		result.Scan(&comment.ID, &comment.UserID, &comment.BoilID, &comment.Content, &comment.CreateTime)
		comments = append(comments, comment)
	}
	result.Close()
	return comments, nil
}

func CommentArrToCommentVoArr(commentArr []model.Comment) ([]model.CommentVo, error) {
	commentVoArr := []model.CommentVo{}
	for _, comment := range commentArr {
		commentVo := model.CommentVo{}
		commentVo.CreateTime = comment.CreateTime.Format("2006-01-02 15:04:05")
		commentVo.ID = comment.ID
		commentVo.BoilId = comment.BoilID
		commentVo.Content = comment.Content
		commentVo.UserID = comment.UserID
		user, _ := GetUserById(comment.UserID)
		commentVo.UserName = user.UserName
		commentVo.UserBio = user.Bio
		commentVo.UserAvatarId = user.AvatarID
		commentVoArr = append(commentVoArr, commentVo)
	}
	return commentVoArr, nil
}

func CountUserCommentBoil(uid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_comment WHERE user_id=?", uid)
	if err != nil {
		return
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}

func CountBoilComment(bid int) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("SELECT COUNT(*) FROM boil_comment WHERE boil_id=?", bid)
	if err != nil {
		return
	}
	result.Next()
	result.Scan(&count)
	result.Close()
	return
}
