package conf

type WrongPasswordCounter struct {
	Captcha int
	Account []struct {
		Statistics    int
		TriggerTimes  int
		ForbiddenTime int
	}
	Ip []struct {
		Statistics    int
		TriggerTimes  int
		ForbiddenTime int
	}
}
