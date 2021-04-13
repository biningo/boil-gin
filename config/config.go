package config

/**
*@Author lyer
*@Date 2/20/21 15:15
*@Describe
**/

type Config struct {
	Server Server `json:"server" yaml:"server"`
	MySql  MySql  `json:"mysql" yaml:"mysql"`
	Jwt Jwt `json:"jwt" yaml:"jwt"`
}
