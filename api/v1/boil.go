package v1

import (
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

/**
*@Author lyer
*@Date 4/15/21 13:58
*@Describe
**/

func BoilPublish(c *gin.Context) {
	boilPublish := model.BoilPublishVo{}
	if err := c.ShouldBindJSON(&boilPublish); err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	boil := model.Boil{Content: boilPublish.Content, TagID: boilPublish.TagID, UserID: boilPublish.UserID}
	boil.CreateTime = time.Now()
	err := service.InsertBoil(boil)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "插入成功"})
}

func BoilDelete(c *gin.Context) {
	bid, _ := strconv.Atoi(c.Param("bid"))
	err := service.DeleteBoilById(bid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}

func BoilListByTag(c *gin.Context) {
	boilArr, err := service.GetBoilsByTagId(c.Param("tid"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilListByUser(c *gin.Context) {
	boilArr, err := service.GetBoilsByUserId(c.Param("uid"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilListByFollowing(c *gin.Context) {
	loginUserId := c.GetInt("loginUserId")
	followingIds, err := service.GetUserFollowingIds(loginUserId)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	if len(followingIds) == 0 {
		c.JSON(200, gin.H{"data": []model.BoilVo{}})
		return
	}
	boilArr, err := service.GetBoilsByUserIds(followingIds)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilListUserLike(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	boilArr, err := service.BoilListUserLike(uid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilAll(c *gin.Context) {
	boilArr, err := service.GetAllBoil()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func GetBoilById(c *gin.Context) {
	boil, err := service.GetBoilById(c.Param("bid"))
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr([]model.Boil{boil}, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr[0]})
}

func BoilListUserComment(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bids, err := service.GetCommentBidsByUid(uid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	if len(bids) == 0 {
		c.JSON(200, gin.H{"data": []model.BoilVo{}})
		return
	}
	boilArr, err := service.GetBoilsByIds(bids)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr, c.GetInt("loginUserId"))
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilUserLike(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bid, _ := strconv.Atoi(c.Param("bid"))
	if err := service.BoilUserLike(bid, uid); err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "successful"})
}
func BoilUserUnLike(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bid, _ := strconv.Atoi(c.Param("bid"))
	if err := service.BoilUserUnLike(bid, uid); err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "successful"})
}

func BoilUserIsLike(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bid, _ := strconv.Atoi(c.Param("bid"))
	c.JSON(200, gin.H{"data": service.BoilUserIsLike(bid, uid)})
}
