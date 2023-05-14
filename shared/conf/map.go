package conf

type MapConf struct {
	Mode      string `json:",default=baidu"`
	AccessKey string
}

type WrongPasswordCounter struct {
	Captcha int `json:",default=3"`
	Account []struct {
		Level int `yaml:"level"`
		Year  int `yaml:"year,omitempty"`
		Month int `yaml:"month,omitempty"`
		Day   int `yaml:"day,omitempty"`
		Times int `default:"42"`
	} `yaml:"Account"`
	Ip []struct {
		Level int `yaml:"level"`
		Year  int `yaml:"year,omitempty"`
		Month int `yaml:"month,omitempty"`
		Day   int `yaml:"day,omitempty"`
		Times int `yaml:"times"`
	} `yaml:"Ip"`
}
