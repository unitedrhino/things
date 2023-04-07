//api网关接口代理模块-apisvr
package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/i-Things/things/src/apisvr/apidirect"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/core/proc" //开启pprof采集 https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA
	"time"
)

func main() {
	logx.DisableStat()
	//Test()
	apiCtx := apidirect.NewApi(apidirect.ApiCtx{})
	apiCtx.Server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", apiCtx.Svc.Config.Host, apiCtx.Svc.Config.Port)
	apiCtx.Server.Start()
	defer apiCtx.Server.Stop()
}

func Test() {
	sc := gocron.NewScheduler(time.Local)
	job, err := sc.Tag("cron 1s test").CronWithSeconds("* * * * * ?").Do(func() {
		logx.Infof("hello world")
	})
	fmt.Println(job, err)
	job, err = sc.Tag("every 2s test").Every("2s").Do(func() {
		//	logx.Infof("hello world")
	})
	fmt.Println(job, err)
	sc.StartAsync()
	//fmt.Println(s)
	jobs, err := sc.FindJobsByTag("cron 1s test")
	fmt.Println(jobs, err)
}
