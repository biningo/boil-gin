package model

import "time"

/**
*@Author lyer
*@Date 4/15/21 15:21
*@Describe
**/

type Comment struct {
	ID int
	BoilID int
	UserID int
	CreateTime time.Time
	Content string
}