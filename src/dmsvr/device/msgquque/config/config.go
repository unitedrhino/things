package config

type Config struct {
	Brokers		[]string	//kafka的节点
	Group 		string		//kafka的分组
}