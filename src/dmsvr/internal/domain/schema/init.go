package schema

import (
	"encoding/json"
)

var (
	defaultSchema = Model{
		Version:    "",
		Properties: []Property{},
		Events:     []Event{},
		Actions:    []Action{},
		Property:   map[string]*Property{},
		Event:      map[string]*Event{},
		Action:     map[string]*Action{},
	}
	DefaultSchema string
)

func init() {
	t, _ := json.Marshal(defaultSchema)
	DefaultSchema = string(t)
}

func NewSchema(schemaStr []byte) (*Model, error) {
	schema := Model{}
	//如果没有需要返回默认值
	if len(schemaStr) == 0 {
		return &schema, nil
	}
	err := json.Unmarshal(schemaStr, &schema)
	if err != nil {
		return nil, err
	}
	schema.init()
	return &schema, nil
}

func (d *Define) init() *Define {
	if d.Specs != nil {
		d.Spec = make(map[string]*Spec, len(d.Specs)+1)
		for i := 0; i < len(d.Specs); i++ {
			p := &d.Specs[i]
			d.Spec[p.ID] = p
		}
	}
	if d.ArrayInfo != nil {
		d.ArrayInfo.init()
	}
	return d
}

func (p *Param) init() *Param {
	p.Define.init()
	return p
}

func (e *Event) init() *Event {
	e.Param = make(map[string]*Param, len(e.Params)+1)
	for i := 0; i < len(e.Params); i++ {
		p := &e.Params[i]
		p.init()
		e.Param[p.ID] = p
	}
	return e
}

func (a *Action) init() *Action {
	a.In = make(map[string]*Param, len(a.Input)+1)
	a.Out = make(map[string]*Param, len(a.Output)+1)
	for i := 0; i < len(a.Input); i++ {
		p := &a.Input[i]
		a.In[p.ID] = p
		p.init()
	}
	for i := 0; i < len(a.Output); i++ {
		p := &a.Output[i]
		a.Out[p.ID] = p
		p.init()
	}
	return a
}

func (p *Property) init() *Property {
	p.Define.init()
	return p
}

func (t *Model) init() *Model {
	t.Property = make(map[string]*Property, len(t.Properties)+1)
	t.Event = make(map[string]*Event, len(t.Events)+1)
	t.Action = make(map[string]*Action, len(t.Actions)+1)
	for i := 0; i < len(t.Properties); i++ {
		p := &t.Properties[i]
		t.Property[p.ID] = p
		p.init()
	}
	for i := 0; i < len(t.Events); i++ {
		p := &t.Events[i]
		t.Event[p.ID] = p
		p.init()
	}
	for i := 0; i < len(t.Actions); i++ {
		p := &t.Actions[i]
		t.Action[p.ID] = p
		p.init()
	}
	return t
}
