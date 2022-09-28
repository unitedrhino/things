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

func NewSchemaTsl(schemaStr []byte) (*Model, error) {
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
			d.Spec[p.Identifier] = p
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
		e.Param[p.Identifier] = p
	}
	return e
}

func (a *Action) init() *Action {
	a.In = make(map[string]*Param, len(a.Input)+1)
	a.Out = make(map[string]*Param, len(a.Output)+1)
	for i := 0; i < len(a.Input); i++ {
		p := &a.Input[i]
		a.In[p.Identifier] = p
		p.init()
	}
	for i := 0; i < len(a.Output); i++ {
		p := &a.Output[i]
		a.Out[p.Identifier] = p
		p.init()
	}
	return a
}

func (p *Property) init() *Property {
	p.Define.init()
	return p
}

func (m *Model) init() *Model {
	m.Property = make(map[string]*Property, len(m.Properties)+1)
	m.Event = make(map[string]*Event, len(m.Events)+1)
	m.Action = make(map[string]*Action, len(m.Actions)+1)
	for i := 0; i < len(m.Properties); i++ {
		p := &m.Properties[i]
		m.Property[p.Identifier] = p
		p.init()
	}
	for i := 0; i < len(m.Events); i++ {
		p := &m.Events[i]
		m.Event[p.Identifier] = p
		p.init()
	}
	for i := 0; i < len(m.Actions); i++ {
		p := &m.Actions[i]
		m.Action[p.Identifier] = p
		p.init()
	}
	return m
}
