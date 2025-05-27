package productmanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCategoryMultiImportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCategoryMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCategoryMultiImportLogic {
	return &ProductCategoryMultiImportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCategoryMultiImportLogic) ProductCategoryMultiImport(in *dm.ProductCategoryImportReq) (*dm.ImportResp, error) {
	var pcs []*relationDB.DmProductCategory
	err := json.Unmarshal([]byte(in.Categories), &pcs)
	if err != nil {
		return nil, err
	}
	var resp = dm.ImportResp{Total: int64(len(pcs))}
	for _, node := range pcs {
		ss := node.Schemas
		node.Children = nil
		node.Schemas = nil
		err := relationDB.NewProductCategoryRepo(l.ctx).Insert(l.ctx, node)
		if err != nil {
			if errors.Cmp(err, errors.Duplicate) {
				resp.IgnoreCount++
				continue
			}
			l.Error(node, err)
			resp.ErrCount++
			continue
		}
		resp.SuccCount++
		if len(ss) > 0 {
			for _, v := range ss {
				v.ID = 0
				v.ProductCategoryID = node.ID
			}
			err = relationDB.NewProductCategorySchemaRepo(l.ctx).MultiInsert(l.ctx, ss)
			if err != nil {
				l.Error(node, err)
			}
		}
	}
	return &resp, nil
}

func (l *ProductCategoryMultiImportLogic) insertTree(parent *relationDB.DmProductCategory, node *relationDB.DmProductCategory, resp *dm.ImportResp) error {
	cs := node.Children
	ss := node.Schemas
	node.ID = 0
	node.ParentID = parent.ID
	node.Children = nil
	node.Schemas = nil
	err := stores.GetCommonConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewProductCategoryRepo(tx).Insert(l.ctx, node)
		if err != nil {
			return err
		}
		node.IDPath = fmt.Sprintf("%s%v-", parent.IDPath, node.ID)
		return relationDB.NewProductCategoryRepo(tx).Update(l.ctx, node)
	})
	if err != nil {
		l.Error(node, err)
		resp.ErrCount++
		resp.ErrCount += int64(len(cs))
		return err
	}
	resp.SuccCount++
	if len(ss) > 0 {
		for _, v := range ss {
			v.ID = 0
			v.ProductCategoryID = node.ID
		}
		err = relationDB.NewProductCategorySchemaRepo(l.ctx).MultiInsert(l.ctx, ss)
		if err != nil {
			l.Error(node, err)
		}
	}
	for _, v := range cs {
		l.insertTree(node, v, resp)
	}
	return nil
}

func fillTree(parent *relationDB.DmProductCategory, parentIDMap map[int64][]*relationDB.DmProductCategory) {
	if parent == nil {
		return
	}
	cs := parentIDMap[parent.ID]
	if len(cs) == 0 {
		return
	}
	parent.Children = cs
	for _, c := range cs {
		fillTree(c, parentIDMap)
	}
}
