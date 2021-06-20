package config

type Config struct {
	Mysql struct {
		DataSource string
	}
	Brokers		[]string	//kafka的节点
	Group 		string		//kafka的分组
}