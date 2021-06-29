package model

/**
*@Author lyer
*@Date 4/13/21 18:10
*@Describe
**/

type User struct {
	ID       int    `db:"id"`
	UserName string `db:"username"`
	PassWord string `db:"password"`
	AvatarID int    `db:"avatar_id"` //[1,5]
	Salt     string `db:"salt"`
	Bio      string `db:"bio"`
}
