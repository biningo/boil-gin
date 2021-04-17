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
*@Date 4/15/21 14:57
*@Describe
**/
func CommentPublish(c *gin.Context) {
	uid := c.GetHeader("userId")
	bid := c.Param("bid")
	comment := model.Comment{}
	c.ShouldBindJSON(&comment)
	comment.UserID, _ = strconv.Atoi(uid)
	comment.BoilID, _ = strconv.Atoi(bid)
	comment.CreateTime = time.Now()
	err := service.InsertComment(comment)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "评论成功"})
}

func CommentBoilList(c *gin.Context) {
	bid := c.Param("bid")
	commentArr, err := service.GetComments("boil_id=?", bid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	commentVoArr, err := service.CommentArrToCommentVoArr(commentArr)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": commentVoArr})
}

func CommentDelete(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("cid"))
	err := service.DeleteCommentById(cid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}
