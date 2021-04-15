package v1

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"github.com/gin-gonic/gin"
)

/**
*@Author lyer
*@Date 4/15/21 15:34
*@Describe
**/

func TagList(c *gin.Context) {
	db := global.G_DB
	result, err := db.Query("select id,title from boil_tag")
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	data := []model.Tag{}
	tag := model.Tag{}
	for result.Next() {
		result.Scan(&tag.ID, &tag.Title)
		data = append(data, tag)
	}
	c.JSON(200, gin.H{"data": data})
}
func TagCreate(c *gin.Context) {
	title := c.Param("title")
	db := global.G_DB
	_, err := db.Exec("insert into boil_tag(title) value(?)", title)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "插入成功"})
}
func TagDelete(c *gin.Context) {
	tid := c.Param("tid")
	db := global.G_DB
	if _, err := db.Exec("delete from boil_tag where id=?", tid); err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}
func TagBoilCount(c *gin.Context) {
	tid := c.Param("tid")
	db := global.G_DB
	count := 0
	result, err := db.Query("select count(*) from boil_boil where tag_id=?", tid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error(), "data": 0})
		return
	}
	result.Next()
	result.Scan(&count)
	c.JSON(200, gin.H{"data": count})
}
