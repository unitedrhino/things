package dto

import (
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dcsvr/dc"
	"github.com/go-things/things/src/webapi/internal/types"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

func GrouInfoToApi(v *dc.GroupInfo) *types.GroupInfo {
	return &types.GroupInfo{
		GroupID:     v.GroupID,     //组id
		Name:        v.Name,        //组名
		Uid:         v.Uid,         //管理员用户id
		CreatedTime: v.CreatedTime, //创建时间 只读
	}
}

func GroupMemberToApi(v *dc.GroupMember) *types.GroupMember {
	return &types.GroupMember{
		GroupID:     v.GroupID,     //组id
		MemberID:    v.MemberID,    //成员id
		MemberType:  v.MemberType,  //成员类型:1:设备 2:用户
		CreatedTime: v.CreatedTime, //创建时间 只读
	}
}

func GetGroupInfoReqToRpc(req *types.GetGroupInfoReq) (*dc.GetGroupInfoReq, error) {
	dcReq := &dc.GetGroupInfoReq{
		GroupID: req.GroupID,
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dcReq.Page = &dc.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	}
	return dcReq, nil
}

func GetGroupInfoRespToApi(resp *dc.GetGroupInfoResp) (*types.GetGroupInfoResp, error) {
	gis := make([]*types.GroupInfo, 0, len(resp.Info))
	for _, v := range resp.Info {
		gi := GrouInfoToApi(v)
		gis = append(gis, gi)
	}
	return &types.GetGroupInfoResp{
		Total: resp.Total,
		Info:  gis,
		Num:   int64(len(gis)),
	}, nil
}

func GetGroupMemberReqToRpc(req *types.GetGroupMemberReq) (*dc.GetGroupMemberReq, error) {
	dcReq := &dc.GetGroupMemberReq{
		GroupID:    req.GroupID,
		MemberID:   req.MemberID,
		MemberType: req.MemberType,
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dcReq.Page = &dc.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	}
	return dcReq, nil
}

func GetGroupMemberRespToApi(resp *dc.GetGroupMemberResp) (*types.GetGroupMemberResp, error) {
	gis := make([]*types.GroupMember, 0, len(resp.Info))
	for _, v := range resp.Info {
		gi := GroupMemberToApi(v)
		gis = append(gis, gi)
	}
	return &types.GetGroupMemberResp{
		Total: resp.Total,
		Info:  gis,
		Num:   int64(len(gis)),
	}, nil
}

func ManageGroupInfoReqToRpc(req *types.ManageGroupInfoReq) (*dc.ManageGroupInfoReq, error) {
	dcReq := &dc.ManageGroupInfoReq{
		Opt: req.Opt,
		Info: &dc.GroupInfo{
			GroupID: req.Info.GroupID, //组id
			Name:    req.Info.Name,    //组名
			Uid:     req.Info.Uid,     //管理员用户id
		},
	}
	return dcReq, nil
}

func ManageGroupMemberReqToRpc(req *types.ManageGroupMemberReq) (*dc.ManageGroupMemberReq, error) {
	dcReq := &dc.ManageGroupMemberReq{
		Opt: req.Opt,
		Info: &dc.GroupMember{
			GroupID:    req.Info.GroupID,    //组id
			MemberID:   req.Info.MemberID,   //成员id
			MemberType: req.Info.MemberType, //成员类型:1:设备 2:用户
		},
	}
	return dcReq, nil
}
