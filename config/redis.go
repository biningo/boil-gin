package config

/**
*@Author lyer
*@Date 4/19/21 10:11
*@Describe
**/

type Redis struct {
	Addr     string `json:"addr" yaml:"addr"`
	DB       int    `json:"db" yaml:"db"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}
