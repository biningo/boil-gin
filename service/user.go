package service

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
)

/**
*@Author lyer
*@Date 4/16/21 23:10
*@Describe
**/

func GetUserById(uid int) (model.User, error) {
	db := global.G_DB
	result, err := db.Query("select username,bio,avatar_id from boil_user where id=?", uid)
	if err != nil {
		return model.User{}, err
	}
	defer result.Close()
	user := model.User{}
	result.Next()
	result.Scan(&user.UserName, &user.Bio, &user.AvatarID)
	return user,nil
}
