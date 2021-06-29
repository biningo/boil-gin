package service

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/17/21 09:58
*@Describe
**/

func GetCommentBidsByUid(uid int) (bids []string, err error) {
	if err = global.G_DB.Select(&bids, "SELECT boil_id FROM boil_comment WHERE user_id=? ORDER BY create_time DESC", uid); err != nil {
		return []string{}, err
	}
	return
}

func InsertComment(comment model.Comment) error {
	_, err := global.G_DB.NamedExec("INSERT INTO boil_comment(user_id,boil_id,content,create_time) VALUE(:user_id,:boil_id,:content,:create_time)", comment)
	return err
}

func DeleteCommentById(cid int) error {
	_, err := global.G_DB.Exec("DELETE FROM boil_comment WHERE id=?", cid)
	return err
}

func CommentArrToCommentVoArr(commentArr []model.Comment) (commentVoArr []model.CommentVo, err error) {
	for _, comment := range commentArr {
		var commentVo model.CommentVo
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
	return
}

func CountUserCommentBoil(uid int) (count int, err error) {
	r := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_comment WHERE user_id=?", uid)
	err = r.Scan(&count)
	return
}

func CountBoilComment(bid int) (count int, err error) {
	result := global.G_DB.QueryRowx("SELECT COUNT(*) FROM boil_comment WHERE boil_id=?", bid)
	err = result.Scan(&count)
	return
}

func GetCommentsByBoilId(bid string) (comments []model.Comment, err error) {
	err = global.G_DB.Select(&comments, "SELECT id,user_id,boil_id,content,create_time FROM boil_comment WHERE boil_id=? ORDER BY create_time DESC", bid)
	return
}
