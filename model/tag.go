package model

/**
*@Author lyer
*@Date 4/15/21 15:31
*@Describe
**/

type Tag struct {
	ID    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}
