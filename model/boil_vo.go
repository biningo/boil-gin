package model

/**
*@Author lyer
*@Date 4/15/21 15:49
*@Describe
**/

type BoilPublishVo struct {
	Content string `json:"content"`
	TagID   int    `json:"tagId"`
	UserID  int    `json:"userId"`
}

type BoilVo struct {
	ID           int    `json:"id"`
	CreateTime   string `json:"createTime"`
	Content      string `json:"content"`
	LikeCount    int    `json:"likeCount"`
	CommentCount int    `json:"commentCount"`
	TagID        int    `json:"tagId"`
	TagTitle     string `json:"tagTitle"`
	UserID       int    `json:"userId"`
	UserName     string `json:"username"`
	UserBio      string `json:"userBio"`
	UserAvatarId int    `json:"userAvatarId"`
}
