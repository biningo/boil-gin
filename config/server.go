package config

/**
*@Author lyer
*@Date 2/20/21 15:30
*@Describe
**/

type Server struct {
	Addr string `json:"addr" yaml:"addr"`
	Mode string `json:"mode" yaml:"mode"`
}
