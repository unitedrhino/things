package schema

import "github.com/i-Things/things/shared/errors"

func CheckModify(oldT *Model, newT *Model) error {
	for _, p := range newT.Property {
		if oldP, ok := oldT.Property[p.Identifier]; ok {
			//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
			if !CheckDefine(&oldP.Define, &p.Define) {
				return errors.Parameter.WithMsgf("不支持类型修改,只支持新增或删除,标识符:%v", p.Identifier)
			}
		}
	}
	for _, e := range newT.Event {
		if oldE, ok := oldT.Event[e.Identifier]; ok {
			//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
			for _, p := range e.Param {
				if oldP, ok := oldE.Param[p.Identifier]; ok {
					//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
					if !CheckDefine(&oldP.Define, &p.Define) {
						return errors.Parameter.WithMsgf("不支持类型修改,只支持新增或删除,标识符:%v", p.Identifier)
					}
				}
			}
		}
	}
	return nil
}

func CheckDefine(oldDef *Define, newDef *Define) bool {
	if oldDef == nil || newDef == nil { //新增删除是支持的
		return true
	}
	if oldDef.Type != newDef.Type {
		return false
	}
	switch oldDef.Type {
	case DataTypeStruct:
		for _, s := range newDef.Spec {
			if olds, ok := oldDef.Spec[s.Identifier]; ok {
				//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
				if !CheckDefine(&olds.DataType, &s.DataType) {
					return false
				}
			}
		}
	case DataTypeArray:
		return CheckDefine(oldDef.ArrayInfo, newDef.ArrayInfo)
	}
	return true
}
