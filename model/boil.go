package model

import "time"

/**
*@Author lyer
*@Date 4/15/21 15:21
*@Describe
**/

type Boil struct {
	ID int
	TagID int
	UserID int
	CreateTime time.Time
	Content string
}