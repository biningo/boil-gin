package v1

import (
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/service"
	"github.com/gin-gonic/gin"
	"log"
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
	userInfo, err := service.GetUserInfoById(user.ID, 0)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
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

func UserInfo(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	userInfoVo, err := service.GetUserInfoById(uid, c.GetInt("loginUserId"))
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": userInfoVo})
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

func UserFollow(c *gin.Context) {
	followerId := c.GetInt("loginUserId")
	uid, _ := strconv.Atoi(c.Param("uid"))
	err := service.InsertUserFollow(followerId, uid)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "successful"})
}

func UserUnFollow(c *gin.Context) {
	followerId := c.GetInt("loginUserId")
	uid, _ := strconv.Atoi(c.Param("uid"))
	err := service.DeleteUserFollow(followerId, uid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "successful"})
}

func ListUserFollower(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	userInfoArr, err := service.GetUserFollower(uid, c.GetInt("loginUserId"))
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": userInfoArr})
}
func ListUserFollowing(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	userInfoArr, err := service.GetUserFollowing(uid, c.GetInt("loginUserId"))
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": userInfoArr})
}

func ListRecommendUser(c *gin.Context) {
	loginUserId := c.GetInt("loginUserId")
	userInfoArr, err := service.GetRecommendUsers(loginUserId)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": userInfoArr})
}

func ListUser(c *gin.Context) {
	loginUserId := c.GetInt("loginUserId")
	userInfoArr, err := service.GetAllUser(loginUserId)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": userInfoArr})
}
