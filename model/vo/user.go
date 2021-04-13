package vo

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
}