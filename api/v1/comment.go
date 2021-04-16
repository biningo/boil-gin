package v1

import (
	"github.com/biningo/boil-gin/global"
	"github.com/biningo/boil-gin/model"
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
	db := global.G_DB
	_, err := db.Exec("insert into boil_comment(user_id,boil_id,content,create_time) value(?,?,?,?)", comment.UserID, comment.BoilID, comment.Content, comment.CreateTime)
	if err != nil {
		print(err.Error())
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "评论成功"})
}

func CommentBoilList(c *gin.Context) {

	bid := c.Param("bid")
	db := global.G_DB
	result, err := db.Query("select id,user_id,boil_id,content,create_time from boil_comment where boil_id=?", bid)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	commentVo := model.CommentVo{}
	data := []model.CommentVo{}
	for result.Next() {
		createTime:=time.Now()
		result.Scan(&commentVo.ID, &commentVo.UserID, &commentVo.BoilId, &commentVo.Content, &createTime)
		commentVo.CreateTime = createTime.Format("2006-01-02 15:04:05")
		result2, err := db.Query("select username,bio,avatar_id from boil_user where id=?", commentVo.UserID)
		if err != nil {
			c.JSON(500, gin.H{"msg": err.Error()})
			return
		}
		result2.Next()
		result2.Scan(&commentVo.UserName, &commentVo.UserBio, &commentVo.UserAvatarId)
		result2.Close()
		data = append(data, commentVo)
	}
	result.Close()
	c.JSON(200, gin.H{"data": data})
}


func CommentDelete(c *gin.Context){
	cid:=c.Param("cid")
	db:=global.G_DB
	db.Exec("delete from boil_comment where id=?",cid)
	c.JSON(200,gin.H{"msg":"删除成功"})
}