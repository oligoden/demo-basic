package demobasic

import (
	"net/http"

	"github.com/oligoden/chassis/device"
	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage/gosql"
)

type Device struct {
	device.Default
}

func NewDevice(s model.Connector) *Device {
	d := &Device{}
	nm := func(r *http.Request) model.Operator { return d.NewModel(r) }
	nv := func(w http.ResponseWriter) view.Operator { return NewView(w) }
	d.Default = device.NewDevice(nm, nv, s)
	return d
}

type Model struct {
	model.Default
}

func (d *Device) NewModel(r *http.Request) *Model {
	m := &Model{}
	m.Default = model.Default{}
	m.Request = r
	m.Store = d.Store
	m.BindUser()
	m.Data(NewRecord())
	return m
}

func (m *Model) NewData(t ...string) {
	if m.Err() != nil {
		return
	}

	e := NewList()
	m.Data(e)
}

func (m *Model) ReadUC(uc string) {
	if m.Err() != nil {
		return
	}

	c := m.Store.Connect(m.User())

	e := NewRecord()
	m.Data(e)
	where := gosql.NewWhere("uc=?", uc)
	c.AddModifiers(where)
	c.Read(e)
	if c.Err() != nil {
		m.Err(c.Err)
		return
	}

	err := e.Complete()
	if err != nil {
		m.Err(err)
		return
	}

	m.Hasher()
}

type View struct {
	view.Default
}

func NewView(w http.ResponseWriter) *View {
	v := &View{}
	v.Default = view.Default{}
	v.Response = w
	return v
}
