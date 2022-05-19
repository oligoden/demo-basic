package demobasic

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/adapter"
	"github.com/oligoden/chassis/storage/gosql"
	"github.com/oligoden/gateway/routing"
)

func setUser(user, session string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X_user") == "" {
			r.Header.Add("X_user", "1")
		}

		if r.Header.Get("X_session") == "" {
			r.Header.Add("X_session", "1")
		}
	})
}

func MuxFunc() func(*adapter.Mux) {
	return func(mux *adapter.Mux) {
		fp := filepath.Join(mux.StaticDir, "index.html")
		store := mux.Stores["mysqldb"]

		if store == nil {
			store = gosql.New(gosql.ConnURL())
			if store.Err() != nil {
				e := chassis.Mark("could not connect to store", store.Err())
				fmt.Println(chassis.ErrorTrace(e))
				os.Exit(1)
			}
		}

		d := NewDevice(store)
		store.Migrate(NewRecord())

		dRouting := routing.NewDevice(store, mux.RPDs...)
		store.Migrate(routing.NewRecord())

		mux.Handle("/").
			Core(adapter.ServeFile(fp)).
			SubDomain(dRouting.Check()).
			And(setUser("1", "1")).
			Notify().Entry()

		mux.Handle("/static/").
			Core(adapter.ServeFiles("/static/", mux.StaticDir)).
			SubDomain(dRouting.Check()).
			And(setUser("1", "1")).
			Entry()

		mux.Handle("/oas.json").
			Core(adapter.ServeFile("oas.json")).
			Notify().Entry()

		mux.Handle("/favicon.ico").
			Core(adapter.ServeFile(filepath.Join("static", "favicon.ico"))).
			Entry()

		mux.Handle("/basics").
			MethodNotAllowed().
			Put(d.Update()).
			Get(d.Read()).
			Post(d.Create()).
			//-{{if eq .Environment "development"}}
			//+Options("PUT").
			//+CORS().
			//+SubDomain(adapter.ServeFile(fp), "-api").
			//-{{end}}
			Notify("entering /basics").
			And(setUser("1", "1")).
			Entry()

		mux.Handle("/basics/").
			MethodNotAllowed().
			Put(d.Update()).
			Get(d.Read()).
			Post(d.Create()).
			//-{{if eq .Environment "development"}}
			//+Options("PUT").
			//+CORS().
			//+SubDomain(adapter.ServeFile(fp), "-api").
			//-{{end}}
			Notify("entering /basics/").
			And(setUser("1", "1")).
			Entry()
	}
}
