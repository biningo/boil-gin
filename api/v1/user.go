package v1

import (
	"crypto/sha256"
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/model/vo"
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
	userLoginVo := vo.UserLoginVo{}
	if err := c.ShouldBindJSON(&userLoginVo); err != nil {
		c.JSON(400, gin.H{"msg": "输入信息不正确"})
		return
	}

	db := global.G_DB
	result, err := db.Query("select username,password,salt,avatar_id from boil_user where username=?", userLoginVo.UserName)
	if err != nil {
		c.JSON(500, gin.H{"msg": "服务器错误"})
		return
	}
	if !result.Next() {
		c.JSON(200, gin.H{"msg": "用户不存在"})
		return
	}
	userInfo := model.User{}
	if err = result.Scan(&userInfo.UserName, &userInfo.PassWord, &userInfo.Salt, &userInfo.AvatarID); err != nil {
		c.JSON(500, gin.H{"msg": "服务器错误"})
		return
	}
	pwdToken := string(sha256.New().Sum([]byte(userLoginVo.PassWord + userInfo.Salt)))
	if pwdToken != userInfo.PassWord {
		c.JSON(403, gin.H{"msg": "密码错误"})
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
	c.JSON(200, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "Logout Successful!"})
}

func Registry(c *gin.Context) {
	user := vo.UserRegistryVo{}
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
	pwd256Token := sha256.New().Sum([]byte(user.PassWord + salt))
	if user.AvatarID <= 0 || user.AvatarID > 5 {
		user.AvatarID = utils.RandInt(1, 5)
	}
	_, err = db.Exec("insert into boil_user(username,password,avatar_id,salt) value(?,?,?,?)", user.UserName, pwd256Token, user.AvatarID, salt)
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
	c.JSON(200, gin.H{"token": tokenString, "msg": "注册成功"})
}
