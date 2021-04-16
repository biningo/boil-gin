package v1

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
	"github.com/gin-gonic/gin"
	"log"
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
	db := global.G_DB

	boil := model.Boil{Content: boilPublish.Content, TagID: boilPublish.TagID, UserID: boilPublish.UserID}
	boil.CreateTime = time.Now()
	_, err := db.Exec("insert into boil_boil(tag_id,user_id,content,create_time) value(?,?,?,?)",
		boil.TagID, boil.UserID, boil.Content, boil.CreateTime)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "插入成功"})
}

func BoilDelete(c *gin.Context) {
	bid := c.Param("bid")
	db := global.G_DB
	_, err := db.Exec("delete from boil_boil where id=?", bid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}

func BoilListTag(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Param("tid"))
	db := global.G_DB
	result, err := db.Query("select id,tag_id,user_id,create_time,content from boil_boil where tag_id=? order by create_time desc", tid)
	data := []model.BoilVo{}
	boil := model.Boil{}
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer result.Close()
	for result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID

		result2, err := db.Query("select title from boil_tag where id=?", boil.TagID)
		if err != nil {
			c.JSON(500, err)
			return
		}
		result2.Next()
		result2.Scan(&boilVo.TagTitle)
		result2.Close()

		result2, err = db.Query("select username,bio,avatar_id from boil_user where id=?", boil.UserID)
		if err != nil {
			c.JSON(500, "3")
			return
		}
		result2.Next()
		result2.Scan(&boilVo.UserName, &boilVo.UserBio, &boilVo.UserAvatarId)
		result2.Close()
		data = append(data, boilVo)
	}
	c.JSON(200, gin.H{"data": data})
}

func BoilListUser(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	db := global.G_DB
	result, err := db.Query("select id,tag_id,user_id,create_time,content from boil_boil where user_id=? order by create_time desc", uid)
	data := []model.BoilVo{}
	boil := model.Boil{}
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer result.Close()
	for result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID

		result2, err := db.Query("select title from boil_tag where id=?", boil.TagID)
		if err != nil {
			c.JSON(500, err)
			return
		}
		result2.Next()
		result2.Scan(&boilVo.TagTitle)
		result2.Close()

		result2, err = db.Query("select username,bio,avatar_id from boil_user where id=?", boil.UserID)
		if err != nil {
			c.JSON(500, "3")
			return
		}
		result2.Next()
		result2.Scan(&boilVo.UserName, &boilVo.UserBio, &boilVo.UserAvatarId)
		result2.Close()
		data = append(data, boilVo)
	}
	c.JSON(200, gin.H{"data": data})
}

func BoilAll(c *gin.Context) {
	db := global.G_DB
	result, err := db.Query("select id,tag_id,user_id,create_time,content from boil_boil order by create_time desc")
	data := []model.BoilVo{}
	boil := model.Boil{}
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer result.Close()
	for result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID

		result2, err := db.Query("select title from boil_tag where id=?", boil.TagID)
		if err != nil {
			c.JSON(500, err)
			return
		}
		result2.Next()
		result2.Scan(&boilVo.TagTitle)
		result2.Close()

		result2, err = db.Query("select username,bio,avatar_id from boil_user where id=?", boil.UserID)
		if err != nil {
			c.JSON(500, "3")
			return
		}
		result2.Next()
		result2.Scan(&boilVo.UserName, &boilVo.UserBio, &boilVo.UserAvatarId)
		result2.Close()
		data = append(data, boilVo)
	}
	c.JSON(200, gin.H{"data": data})
}

func GetBoil(c *gin.Context) {
	bid := c.Param("bid")
	db := global.G_DB
	result, _ := db.Query("select id,tag_id,user_id,create_time,content from boil_boil where id=?", bid)
	boil := model.Boil{}

	defer result.Close()
	if result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID

		result2, err := db.Query("select title from boil_tag where id=?", boil.TagID)
		if err != nil {
			c.JSON(500, err)
			return
		}
		result2.Next()
		result2.Scan(&boilVo.TagTitle)
		result2.Close()

		result2, err = db.Query("select username,bio,avatar_id from boil_user where id=?", boil.UserID)
		if err != nil {
			c.JSON(500, "3")
			return
		}
		result2.Next()
		result2.Scan(&boilVo.UserName, &boilVo.UserBio, &boilVo.UserAvatarId)
		result2.Close()
		c.JSON(200, gin.H{"data": boilVo})
		return
	}
	c.JSON(404, gin.H{"msg": "boil不存在"})
}

func BoilListUserLike(c *gin.Context) {
	c.JSON(200, gin.H{"data": []model.BoilVo{}})
}

func BoilListUserComment(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Param("uid"))
	bids := []string{}
	bid := "0"
	db := global.G_DB
	result, err := db.Query("select boil_id from boil_comment where user_id=? order by create_time desc", uid)
	for result.Next() {
		result.Scan(&bid)
		bids = append(bids, bid)
	}
	data := []model.BoilVo{}
	if len(bids)==0{
		c.JSON(200,gin.H{"data":data})
		return
	}
	inStr:="("+strings.Join(bids,",")+")"
	log.Println(inStr)
	result, err = db.Query("select id,tag_id,user_id,create_time,content from boil_boil where id in "+inStr)
	boil := model.Boil{}
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer result.Close()
	for result.Next() {
		result.Scan(&boil.ID, &boil.TagID, &boil.UserID, &boil.CreateTime, &boil.Content)
		boilVo := model.BoilVo{}
		boilVo.ID = boil.ID
		boilVo.TagID = boil.TagID
		boilVo.Content = boil.Content
		boilVo.CreateTime = boil.CreateTime.Format("2006-01-02 15:04:05")
		boilVo.UserID = boil.UserID

		result2, err := db.Query("select title from boil_tag where id=?", boil.TagID)
		if err != nil {
			c.JSON(500, err)
			return
		}
		result2.Next()
		result2.Scan(&boilVo.TagTitle)
		result2.Close()

		result2, err = db.Query("select username,bio,avatar_id from boil_user where id=?", boil.UserID)
		if err != nil {
			c.JSON(500, "3")
			return
		}
		result2.Next()
		result2.Scan(&boilVo.UserName, &boilVo.UserBio, &boilVo.UserAvatarId)
		result2.Close()
		data = append(data, boilVo)
	}
	c.JSON(200, gin.H{"data": data})
}
