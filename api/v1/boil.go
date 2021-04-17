package v1

import (
	"github.com/biningo/boil-gin/model"
	"github.com/biningo/boil-gin/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

/**
*@Author lyer
*@Date 4/15/21 13:58
*@Describe
**/

func BoilPublish(c *gin.Context) {
	boilPublish := model.BoilPublishVo{}
	c.ShouldBindJSON(&boilPublish)
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
	tid, _ := strconv.Atoi(c.Param("tid"))
	boilArr, err := service.GetBoils("tag_id=?", tid)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr)
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilListByUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	boilArr, err := service.GetBoils("user_id=?", uid)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr)
	c.JSON(200, gin.H{"data": boilVoArr})
}

func BoilAll(c *gin.Context) {
	boilArr, err := service.GetBoils("1=1")
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr)
	c.JSON(200, gin.H{"data": boilVoArr})
}

func GetBoilById(c *gin.Context) {
	bid, _ := strconv.Atoi(c.Param("bid"))
	boilArr, err := service.GetBoils("id=?", bid)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr)
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
	boilArr, err := service.GetBoils("id in " + "(" + strings.Join(bids, ",") + ")")
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	boilVoArr := service.BoilArrToBoilVoArr(boilArr)
	c.JSON(200, gin.H{"data": boilVoArr})
}
