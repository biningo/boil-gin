package model

import "time"

/**
*@Author lyer
*@Date 4/15/21 15:21
*@Describe
**/
type UserMessage struct {
	ID         int
	UserID     int
	CreateTime time.Time
	Content    string
	BioID      int
	SrcUserID  int
}
