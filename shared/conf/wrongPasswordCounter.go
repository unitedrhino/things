package conf

type WrongPasswordCounter struct {
	Captcha int `json:",default=5"`
	Account []struct {
		Statistics    int `json:",default=1440"`
		TriggerTimes  int `json:",default=10"`
		ForbiddenTime int `json:",default=10"`
	}
	Ip []struct {
		Statistics    int `json:",default=1440"`
		TriggerTimes  int `json:",default=200"`
		ForbiddenTime int `json:",default=60"`
	}
}
