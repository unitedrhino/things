package schema

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func productSchemaToApi(v *dm.ProductSchema) types.ProductSchema {
	return types.ProductSchema{
		CreatedTime: v.CreatedTime, //创建时间 只读
		ProductID:   v.ProductID,   //产品id 只读
		Schema:      v.Schema,      //数据模板
	}
}
