package utils

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"time"
)

func SingletonRun(ctx context.Context, store kv.Store, singletonKey string, f func(ctx2 context.Context)) {
	for true { //定时任务为单例执行模式,有效期15秒,如果服务挂了,其他服务每隔10秒检测到就抢到执行
		ok, err := store.SetnxExCtx(ctx, "singleton:"+singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Store.SetnxExCtx singletonKey:%v err:%v", FuncName(), singletonKey, err)
			time.Sleep(time.Second * 10)
			continue
		}
		if ok { //抢到锁了
			break
		}
		logx.WithContext(ctx).Infof("SingletonRun not get  singletonKey:%v", singletonKey)
		//没抢到锁,10秒钟后继续
		time.Sleep(time.Second * 10)
	}
	logx.WithContext(ctx).Infof("SingletonRun start running singletonKey:%v", singletonKey)
	//抢到锁需要维系锁
	//每隔10秒刷新锁,如果服务挂了,锁才能退出
	Go(ctx, func() {
		defer Recover(ctx)
		ticker := time.NewTicker(time.Second * 10)
		for range ticker.C {
			store.SetexCtx(ctx, singletonKey, time.Now().Format("2006-01-02 15:04:05.999"), 15)
		}
	})
	f(ctx)
}
