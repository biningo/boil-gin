package model

import "time"

/**
*@Author lyer
*@Date 4/15/21 15:21
*@Describe
**/

type Boil struct {
	ID         int       `db:"id"`
	TagID      int       `db:"tag_id"`
	UserID     int       `db:"user_id"`
	CreateTime time.Time `db:"create_time"`
	Content    string    `db:"content"`
}
