package model

/**
*@Author lyer
*@Date 4/13/21 17:31
*@Describe
**/

type UserLoginVo struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type UserRegistryVo struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	AvatarID int `json:"avatarId"`
	Bio string `json:"bio"`
}

type UserInfo struct {
	ID int `json:"id"`
	UserName string `json:"username"`
	Bio string `json:"bio"`
	AvatarID int `json:"avatarId"`
}