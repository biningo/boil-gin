package config

/**
*@Author lyer
*@Date 2/20/21 15:30
*@Describe
**/

type MySql struct {
	Host      string `json:"host" yaml:"host"`
	Port      string `json:"port" yaml:"port"`
	User      string `json:"user" yaml:"user"`
	Password  string `json:"password" yaml:"password"`
	DB        string `json:"db" yaml:"db"`
	Collation string `json:"collation" yaml:"collation"`
}
