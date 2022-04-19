package deviceDataRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) DropProduct(ctx context.Context, t *deviceTemplate.Template, productID string) error {
	tableList := getStableNameList(t, productID)
	for _, v := range tableList {
		sql := fmt.Sprintf("drop stable if exists %s;", v)
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) InitProduct(ctx context.Context, t *deviceTemplate.Template, productID string) error {
	if t != nil {
		for _, p := range t.Properties {
			err := d.createPropertyStable(ctx, &p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s|createPropertyStable|prodecutID:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return err
			}
		}
	}
	{
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`event_id` BINARY(50),`event_type` BINARY(20), `param` BINARY(5000)) "+
			"TAGS (device_name BINARY(50));",
			getEventStableName(productID))
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}

	return nil
}

func (d *DeviceDataRepo) ModifyProperty(
	oldP *deviceTemplate.Property,
	newP *deviceTemplate.Property,
	productID string) error {
	if newP.Define.Type != deviceTemplate.STRUCT {
		//不需要修改数据库
		return nil
	}
	for _, newS := range newP.Define.Spec {
		if _, ok := oldP.Define.Spec[newS.ID]; ok {
			//如果老的物模型有这个字段则不处理
			delete(oldP.Define.Spec, newS.ID)
		} else {
			//新增
			sql := fmt.Sprintf("ALTER STABLE %s ADD COLUMN %s %s; ",
				getPropertyStableName(productID, newP.ID), newS.ID, getTdType(newS.DataType))
			if _, err := d.t.Exec(sql); err != nil {
				return err
			}
		}
	}
	for _, oldS := range oldP.Define.Spec {
		//这里是需要删除的字段
		sql := fmt.Sprintf("ALTER STABLE %s DROP COLUMN %s; ",
			getPropertyStableName(productID, newP.ID), oldS.ID)
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) ModifyProduct(
	ctx context.Context, oldT *deviceTemplate.Template, newt *deviceTemplate.Template, productID string) error {
	for _, p := range newt.Property {
		if oldP, ok := oldT.Property[p.ID]; ok {
			//这里需要走修改流程
			err := d.ModifyProperty(oldP, p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s|ModifyProperty|prodecutID:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return err
			}
		} else { //新增流程
			err := d.createPropertyStable(ctx, p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s|createPropertyStable|prodecutID:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return err
			}
		}
		//已经修改过的需要删除了,新增修改完了的就是需要删除的
		delete(oldT.Property, p.ID)
	}
	//处理删除的属性
	for _, p := range oldT.Property {
		sql := fmt.Sprintf("drop stable if exists %s;", getPropertyStableName(productID, p.ID))
		if _, err := d.t.Exec(sql); err != nil {
			logx.WithContext(ctx).Errorf("%s|drop table|prodecutID:%v,properties:%v,err:%v",
				utils.FuncName(), productID, p, err)
			return err
		}
	}
	return nil
}

func (d *DeviceDataRepo) createPropertyStable(
	ctx context.Context, p *deviceTemplate.Property, productID string) error {
	var sql string
	if p.Define.Type != deviceTemplate.STRUCT {
		sql = fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp,`param` %s)"+
			" TAGS (`device_name` BINARY(50),`"+PROPERTY_TYPE+"` BINARY(50));",
			getPropertyStableName(productID, p.ID), getTdType(p.Define))
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	} else {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp, %s)"+
			" TAGS (`device_name` BINARY(50),`"+PROPERTY_TYPE+"` BINARY(50));",
			getPropertyStableName(productID, p.ID), getSpecsColumn(p.Define.Specs))
		if _, err := d.t.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}
