package templateModel

import "github.com/i-Things/things/shared/errors"

func CheckModify(oldT *Template, newT *Template) error {
	for _, p := range newT.Property {
		if oldP, ok := oldT.Property[p.ID]; ok {
			//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
			if !CheckDefine(&oldP.Define, &p.Define) {
				return errors.Parameter.WithMsgf("不支持类型修改,只支持新增或删除,标识符:%v", p.ID)
			}
		}
	}
	for _, e := range newT.Event {
		if oldE, ok := oldT.Event[e.ID]; ok {
			//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
			for _, p := range e.Param {
				if oldP, ok := oldE.Param[p.ID]; ok {
					//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
					if !CheckDefine(&oldP.Define, &p.Define) {
						return errors.Parameter.WithMsgf("不支持类型修改,只支持新增或删除,标识符:%v", p.ID)
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
	case STRUCT:
		for _, s := range newDef.Spec {
			if olds, ok := oldDef.Spec[s.ID]; ok {
				//需要判断类型是否相同,如果不相同不可以修改,只能删除新增
				if !CheckDefine(&olds.DataType, &s.DataType) {
					return false
				}
			}
		}
	case ARRAY:
		return CheckDefine(oldDef.ArrayInfo, newDef.ArrayInfo)
	}
	return true
}
