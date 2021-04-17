package service

import (
	"crypto/sha256"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/middleware"
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/**
*@Author lyer
*@Date 4/16/21 23:10
*@Describe
**/

func GetUserById(uid int) (model.User, error) {
	db := global.G_DB
	result, err := db.Query("SELECT username,bio,avatar_id FROM boil_user WHERE id=?", uid)
	if err != nil {
		return model.User{}, err
	}
	defer result.Close()
	user := model.User{}
	result.Next()
	result.Scan(&user.UserName, &user.Bio, &user.AvatarID)
	result.Close()
	return user, nil
}

func GetUserByName(username string) (user model.User, err error) {
	db := global.G_DB
	exec, _ := db.Prepare("SELECT id,username,password,bio,avatar_id,salt FROM boil_user WHERE username=?")
	result, err := exec.Query(username)
	if err != nil {
		return
	}
	result.Next()
	result.Scan(&user.ID, &user.UserName, &user.PassWord, &user.Bio, &user.AvatarID, &user.Salt)
	return
}

func UpdateUserBio(uid int, bio string) error {
	db := global.G_DB
	_, err := db.Exec("update boil_user set bio=? where id=?", bio, uid)
	return err
}

func InsertUser(user model.User) error {
	db := global.G_DB
	_, err := db.Exec(
		"insert into boil_user(username,password,avatar_id,salt,bio) value(?,?,?,?,?)",
		user.UserName, user.PassWord, user.AvatarID, user.Salt, user.Bio)
	return err
}

func CountUserByName(username string) (count int, err error) {
	db := global.G_DB
	result, err := db.Query("select count(*) from boil_user where username=?", username)
	if err != nil {
		return 0, nil
	}
	result.Next()
	result.Scan(&count)
	return
}

func UserRegistryVoToUser(userRegistryVo model.UserRegistryVo) (user model.User) {
	salt := utils.Rand5Str()
	pwd256Token := sha256.Sum256([]byte(userRegistryVo.PassWord + salt))
	if userRegistryVo.AvatarID <= 0 || userRegistryVo.AvatarID > 5 {
		userRegistryVo.AvatarID = utils.RandInt(1, 5)
	}
	if userRegistryVo.Bio == "" {
		userRegistryVo.Bio = "这个家伙很懒,什么也没留下"
	}
	user.UserName = userRegistryVo.UserName
	user.PassWord = string(pwd256Token[:])
	user.Salt = salt
	user.Bio = userRegistryVo.Bio
	user.AvatarID = userRegistryVo.AvatarID
	return
}

func CreateUserToken(user model.User) (string, error) {
	tokenString, err := middleware.NewJwt().CreateToken(&middleware.CustomClaims{
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "boil",
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(global.G_CONFIG.Jwt.TokenTime)).Unix(),
		},
	})
	return tokenString, err
}

func CheckUserToken(userLoginVo model.UserLoginVo, salt string, password string) bool {
	pwd256Bytes := sha256.Sum256([]byte(userLoginVo.PassWord + salt))
	userLoginPwdToken := string(pwd256Bytes[:])
	return userLoginPwdToken == password
}

func CleanUser(user *model.User) {
	user.PassWord = ""
	user.Salt = ""
}
