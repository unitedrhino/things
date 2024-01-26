package common

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/parnurzeal/gorequest"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type WeatherReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWeatherReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeatherReadLogic {
	return &WeatherReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var key = "b7de434f83c146e480d13ba6a565ce30"

type respType[t any] struct {
	Code string `json:"code"`
	Now  t      `json:"now"`
}

func (l *WeatherReadLogic) WeatherRead(req *types.WeatherReadReq) (resp *types.WeatherReadResp, err error) {
	var (
		weather respType[types.WeatherReadResp]
		air     respType[types.WeatherAir]
		greq    = gorequest.New().Retry(3, time.Second*2)
	)

	_, _, errs := greq.Get(fmt.Sprintf("https://devapi.qweather.com/v7/weather/now?location=%v,%v&key=%s",
		req.Position.Longitude, req.Position.Latitude, key)).EndStruct(&weather)
	if err != nil {
		return nil, errors.System.AddDetail(errs)
	}
	_, _, errs = greq.Get(fmt.Sprintf("https://devapi.qweather.com/v7/air/now?location=%v,%v&key=%s",
		req.Position.Longitude, req.Position.Latitude, key)).EndStruct(&air)
	if errs != nil {
		return nil, errors.System.AddDetail(errs)
	}
	weather.Now.Air = air.Now
	return &weather.Now, nil
}
