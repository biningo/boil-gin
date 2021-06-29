package v1

import (
	"github.com/biningo/boil-gin/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

/**
*@Author lyer
*@Date 4/15/21 15:34
*@Describe
**/

func TagList(c *gin.Context) {
	tags, err := service.GetAllTags()
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": tags})
}

func TagCreate(c *gin.Context) {
	title := c.Param("title")
	err := service.InsertTag(title)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "插入成功"})
}

func TagDelete(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Param("tid"))
	err := service.DeleteTagById(tid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}

func TagCountBoil(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Param("tid"))
	count, _ := service.CountBoilByTag(tid)
	c.JSON(200, gin.H{"data": count})
}
