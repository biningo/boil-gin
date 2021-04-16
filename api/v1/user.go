package v1

import (
	"crypto/sha256"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)
import "github.com/biningo/boil-gin/middleware"

/**
*@Author lyer
*@Date 4/13/21 15:37
*@Describe
**/

func Login(c *gin.Context) {
	j := middleware.NewJwt()
	userLoginVo := model.UserLoginVo{}
	if err := c.ShouldBindJSON(&userLoginVo); err != nil {
		c.JSON(400, gin.H{"msg": "输入信息不正确"})
		return
	}
	db := global.G_DB
	result, err := db.Query("select id,username,password,salt,avatar_id,bio from boil_user where username=?", userLoginVo.UserName)
	if err != nil {
		c.JSON(500, gin.H{"msg": "服务器错误1"})
		return
	}
	if !result.Next() {
		c.JSON(403, gin.H{"msg": "用户名或者密码错误1"})
		return
	}
	user := model.User{}
	if err = result.Scan(&user.ID, &user.UserName, &user.PassWord, &user.Salt, &user.AvatarID, &user.Bio); err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	byteArr := sha256.Sum256([]byte(userLoginVo.PassWord + user.Salt))
	pwdToken := string(byteArr[:])
	if pwdToken != user.PassWord {
		c.JSON(403, gin.H{"msg": "用户名或密码错误2"})
		return
	}
	claims := &middleware.CustomClaims{
		UserName: userLoginVo.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(global.G_CONFIG.Jwt.TokenTime)).Unix(),
			Issuer:    "lyer",
		},
	}
	tokenString, _ := j.CreateToken(claims)

	userInfo := model.UserInfo{ID: user.ID, UserName: user.UserName, Bio: user.Bio, AvatarID: user.AvatarID}
	c.JSON(200, gin.H{"data": userInfo, "token": tokenString})
}

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "Logout Successful!"})
}

func Registry(c *gin.Context) {
	user := model.UserRegistryVo{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	db := global.G_DB

	result, _ := db.Query("select count(*) from boil_user where username=?", user.UserName)
	count := 0
	result.Next()
	if err := result.Scan(&count); err != nil {
		c.JSON(200, gin.H{"msg": err.Error()})
		return
	}
	if count != 0 {
		c.JSON(200, gin.H{"msg": "用户名存在!!"})
		return
	}

	salt := utils.Rand5Str()
	pwd256Token := sha256.Sum256([]byte(user.PassWord + salt))
	if user.AvatarID <= 0 || user.AvatarID > 5 {
		user.AvatarID = utils.RandInt(1, 5)
	}
	if user.Bio == "" {
		user.Bio = "这个家伙很懒,什么也没留下"
	}
	_, err = db.Exec("insert into boil_user(username,password,avatar_id,salt,bio) value(?,?,?,?,?)", user.UserName, string(pwd256Token[:]), user.AvatarID, salt, user.Bio)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	tokenString, err := middleware.NewJwt().CreateToken(&middleware.CustomClaims{
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "boil",
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(global.G_CONFIG.Jwt.TokenTime)).Unix(),
		},
	})
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	user.PassWord = ""
	c.JSON(200, gin.H{"data": user, "token": tokenString, "msg": "注册成功"})
}

func UserStatus(c *gin.Context) {
	uid := c.Param("uid")
	userStatusVo := model.UserStatusVo{}
	db := global.G_DB
	result, _ := db.Query("select count(*) from boil_boil where user_id=?", uid)
	result.Next()
	result.Scan(&userStatusVo.UserBoilCount)

	result, _ = db.Query("select count(*) from boil_comment where user_id=?", uid)
	result.Next()
	result.Scan(&userStatusVo.CommentCount)

	result.Close()
	c.JSON(200, gin.H{"data": userStatusVo})
}

func UserUpdateBio(c *gin.Context) {
	uid := c.Param("uid")
	bioVo := struct {
		Bio string `json:"bio"`
	}{}
	c.ShouldBindJSON(&bioVo)
	db := global.G_DB
	_, err := db.Exec("update boil_user set bio=? where id=?", bioVo.Bio, uid)
	if err != nil {
		c.JSON(500, gin.H{"msg": "更新错误"})
		return
	}
	c.JSON(200, gin.H{"msg": "更新成功!"})
}
