package schemaDataRepo

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d *DeviceDataRepo) DeleteProduct(ctx context.Context, t *schema.Model, productID string) error {
	tableList := d.GetStableNameList(t, productID)
	for _, v := range tableList {
		sql := fmt.Sprintf("drop stable if exists %s;", v)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}
	return nil
}

func (d *DeviceDataRepo) InitProduct(ctx context.Context, t *schema.Model, productID string) error {
	if t != nil {
		for _, p := range t.Property {
			err := d.createPropertyStable(ctx, p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.createPropertyStable product_id:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return errors.Database.AddDetail(err)
			}
		}
	}
	{
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s "+
			"(`ts` timestamp,`event_id` BINARY(50),`event_type` BINARY(20), `param` BINARY(5000)) "+
			"TAGS (`product_id` BINARY(50),`device_name` BINARY(50));",
			d.GetEventStableName())
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}

	return nil
}

func (d *DeviceDataRepo) CreateProperty(ctx context.Context, p *schema.Property, productID string) error {
	err := d.createPropertyStable(ctx, p, productID)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.createPropertyStable product_id:%v,properties:%v,err:%v",
			utils.FuncName(), productID, p, err)
		return errors.Database.AddDetail(err)
	}
	return nil
}
func (d *DeviceDataRepo) DeleteProperty(ctx context.Context, productID string, identifier string) error {
	sql := fmt.Sprintf("drop stable if exists %s;", d.GetPropertyStableName(productID, identifier))
	if _, err := d.t.ExecContext(ctx, sql); err != nil {
		return errors.Database.AddDetail(err)
	}
	return nil
}

func (d *DeviceDataRepo) UpdateProperty(
	ctx context.Context,
	oldP *schema.Property,
	newP *schema.Property,
	productID string) error {
	if newP.Define.Type != oldP.Define.Type {
		return errors.Parameter.AddMsg("类型不能修改,可以删除后新增")
	}
	if newP.Define.Type != schema.DataTypeStruct {
		//不需要修改数据库
		return nil
	}
	for _, newS := range newP.Define.Spec {
		if _, ok := oldP.Define.Spec[newS.Identifier]; ok {
			//如果老的物模型有这个字段则不处理
			delete(oldP.Define.Spec, newS.Identifier)
		} else {
			//新增
			sql := fmt.Sprintf("ALTER STABLE %s ADD COLUMN `%s` %s; ",
				d.GetPropertyStableName(productID, newP.Identifier), newS.Identifier, stores.GetTdType(newS.DataType))
			if _, err := d.t.ExecContext(ctx, sql); err != nil {
				return errors.Database.AddDetail(err)
			}
		}
	}
	for _, oldS := range oldP.Define.Spec {
		//这里是需要删除的字段
		sql := fmt.Sprintf("ALTER STABLE %s DROP COLUMN `%s`; ",
			d.GetPropertyStableName(productID, newP.Identifier), oldS.Identifier)
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}
	return nil
}

func (d *DeviceDataRepo) UpdateProduct(
	ctx context.Context, oldT *schema.Model, newt *schema.Model, productID string) error {
	for _, p := range newt.Property {
		if oldP, ok := oldT.Property[p.Identifier]; ok {
			//这里需要走修改流程
			err := d.UpdateProperty(ctx, oldP, p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.UpdateProperty product_id:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return errors.Database.AddDetail(err)
			}
		} else { //新增流程
			err := d.createPropertyStable(ctx, p, productID)
			if err != nil {
				logx.WithContext(ctx).Errorf("%s.createPropertyStable productID:%v,properties:%v,err:%v",
					utils.FuncName(), productID, p, err)
				return errors.Database.AddDetail(err)
			}
		}
		//已经修改过的需要删除了,新增修改完了的就是需要删除的
		delete(oldT.Property, p.Identifier)
	}
	//处理删除的属性
	for _, p := range oldT.Property {
		sql := fmt.Sprintf("drop stable if exists %s;", d.GetPropertyStableName(productID, p.Identifier))
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			logx.WithContext(ctx).Errorf("%s drop table product_id:%v,properties:%v,err:%v",
				utils.FuncName(), productID, p, err)
			return errors.Database.AddDetail(err)
		}
	}
	return nil
}

func (d *DeviceDataRepo) createPropertyStable(
	ctx context.Context, p *schema.Property, productID string) error {
	var sql string
	if p.Define.Type != schema.DataTypeStruct {
		sql = fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp,`param` %s)"+
			" TAGS (`product_id` BINARY(50),`device_name` BINARY(50),`"+PropertyType+"` BINARY(50));",
			d.GetPropertyStableName(productID, p.Identifier), stores.GetTdType(p.Define))
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	} else {
		sql := fmt.Sprintf("CREATE STABLE IF NOT EXISTS %s (`ts` timestamp, %s)"+
			" TAGS (`product_id` BINARY(50),`device_name` BINARY(50),`"+PropertyType+"` BINARY(50));",
			d.GetPropertyStableName(productID, p.Identifier), d.GetSpecsCreateColumn(p.Define.Specs))
		if _, err := d.t.ExecContext(ctx, sql); err != nil {
			return errors.Database.AddDetail(err)
		}
	}
	return nil
}
