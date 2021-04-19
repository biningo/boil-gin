package v1

import (
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

/**
*@Author lyer
*@Date 4/13/21 15:37
*@Describe
**/

func Login(c *gin.Context) {
	userLoginVo := model.UserLoginVo{}
	c.ShouldBindJSON(&userLoginVo)
	count, _ := service.CountUserByName(userLoginVo.UserName)
	if count != 1 {
		c.JSON(200, gin.H{"msg": "用户名或则密码不正确"})
		return
	}
	user, err := service.GetUserByName(userLoginVo.UserName)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	if !service.CheckUserToken(userLoginVo, user.Salt, user.PassWord) {
		c.JSON(403, gin.H{"msg": "用户名或密码不正确"})
		return
	}
	tokenString, err := service.CreateUserToken(user)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	userInfo := model.UserInfo{ID: user.ID, UserName: user.UserName, Bio: user.Bio, AvatarID: user.AvatarID}
	c.JSON(200, gin.H{"data": userInfo, "token": tokenString})
}

func Logout(c *gin.Context) {
	//TODO jwt token clear
	c.JSON(200, gin.H{"msg": "Logout Successful!"})
}

func Registry(c *gin.Context) {
	userRegistryVo := model.UserRegistryVo{}
	c.ShouldBindJSON(&userRegistryVo)
	count, _ := service.CountUserByName(userRegistryVo.UserName)
	if count != 0 {
		c.JSON(200, gin.H{"msg": "用户名存在!!"})
		return
	}
	user := service.UserRegistryVoToUser(userRegistryVo)
	err := service.InsertUser(user)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	tokenString, err := service.CreateUserToken(user)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	service.CleanUser(&user)
	c.JSON(200, gin.H{"data": user, "token": tokenString, "msg": "注册成功"})
}

func UserStatus(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	userStatusVo := model.UserStatusVo{}
	userStatusVo.UserID = uid
	userStatusVo.UserBoilCount, _ = service.CountUserBoil(uid)
	userStatusVo.CommentCount, _ = service.CountUserCommentBoil(uid)
	userStatusVo.LikeBoilCount, _ = service.CountUserLikeBoil(uid)
	c.JSON(200, gin.H{"data": userStatusVo})
}

func UpdateUserBio(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bioVo := struct {
		Bio string `json:"bio"`
	}{}
	c.ShouldBindJSON(&bioVo)
	err := service.UpdateUserBio(uid, bioVo.Bio)
	if err != nil {
		c.JSON(500, gin.H{"msg": "更新错误"})
		return
	}
	c.JSON(200, gin.H{"msg": "更新成功!"})
}
