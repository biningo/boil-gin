package model

/**
*@Author lyer
*@Date 4/13/21 18:10
*@Describe
**/

type User struct {
	ID       int
	UserName string
	PassWord string
	AvatarID int //[1,5]
	Salt     string
}
