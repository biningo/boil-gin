package config

/**
*@Author lyer
*@Date 4/13/21 17:37
*@Describe
**/
type Jwt struct {
	Secret string `json:"secret" yaml:"secret"`
	TokenTime int `json:"tokentime" yaml:"tokentime"` //hour
}