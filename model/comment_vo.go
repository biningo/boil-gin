package model

/**
*@Author lyer
*@Date 4/16/21 20:37
*@Describe
**/
type CommentVo struct {
	ID           int    `json:"id"`
	BoilId       int    `json:"boilId"`
	CreateTime   string `json:"createTime"`
	Content      string `json:"content"`
	UserID       int    `json:"userId"`
	UserName     string `json:"username"`
	UserBio      string `json:"userBio"`
	UserAvatarId int    `json:"userAvatarId"`
}
