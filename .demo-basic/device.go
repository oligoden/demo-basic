package demobasic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/oligoden/chassis/device"
	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
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

func (d Device) Read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		v := d.NewView(w)

		rgx, _ := regexp.Compile("^/basics(/?|/([a-zA-Z0-9]+))$")
		matches := rgx.FindStringSubmatch(r.URL.Path)

		if len(matches) == 0 {
			v.Error(m)
			return
		}

		if matches[2] != "" {
			m.ReadUC(matches[2])
			v.JSON(m)
			return
		}

		e := NewList()
		m.Data(e)
		m.Read()

		v.JSON(m)
	})
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
	m.NewData = func() data.Operator { return NewRecord() }
	m.Data(NewRecord())
	return m
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

func (v View) JSON(m model.Operator) {
	if m.Err() != nil {
		v.Error(m)
		return
	}

	out, err := json.Marshal(m.Data())
	if err != nil {
		log.Println(err)
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	v.Response.Header().Set("Content-Type", "application/json")
	n, err := v.Response.Write(out)
	if err != nil {
		log.Println("an error occured writing to response", err)
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}
	fmt.Printf("\n--< written %d bytes\n--< %s\n", n, string(out))
}
