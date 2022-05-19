package demobasic_test

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/oligoden/chassis/adapter"
	"github.com/oligoden/chassis/storage/gosql"
	demobasic "github.com/oligoden/demo-basic"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	uri := "user:pass@tcp(localhost:3320)/test?charset=utf8&parseTime=True&loc=Local"

	db := testDBDropTables(t, "mysql", uri)
	defer db.Close()

	if err := os.MkdirAll("testing", 0755); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll("testing")

	c := []byte(`<html></html>`)
	if err := ioutil.WriteFile("testing/index.html", c, 0644); err != nil {
		t.Error(err)
	}

	c = []byte(`.css{}`)
	if err := ioutil.WriteFile("testing/index.css", c, 0644); err != nil {
		t.Error(err)
	}

	store := gosql.New(uri)
	if store.Err() != nil {
		t.Fatal("could not connect to store")
	}

	mux := adapter.NewMux().
		SetURL("http://test.com").
		SetStaticDir("testing").
		SetStore("mysqldb", store).
		AddRPD("api.test.com").
		Compile(demobasic.MuxFunc())

	qs := []string{
		"INSERT INTO `routings` (`uc`, `domain`, `path`, `url`, `reset_cors`, `owner_id`, `perms`, `hash`) VALUES ('a', 'api.test.com', '/', 'api.test.com', false, 1, ':::', 'abc')",
	}
	testDBSetup(db, t, qs...)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html></html>").
		End()

	req = httptest.NewRequest(http.MethodGet, "/static/index.css", nil)
	req.Host = "test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body(".css{}").
		End()

	swaggerMock := apitest.NewMock().
		Get("http://api.test.com/").
		RespondWith().
		Body(`<html>swagger</html>`).
		Status(http.StatusOK).
		End()

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "api.test.com"
	expiration := time.Now().Add(time.Minute)
	cookie := &http.Cookie{Name: "session", Value: "xyz", Expires: expiration}
	req.AddCookie(cookie)

	apitest.New().
		Mocks(swaggerMock).
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html>swagger</html>").
		End()
}

func TestBasicCRUD(t *testing.T) {
	uri := "user:pass@tcp(localhost:3320)/test?charset=utf8&parseTime=True&loc=Local"

	db := testDBDropTables(t, "mysql", uri)
	defer db.Close()

	store := gosql.New(uri)
	if store.Err() != nil {
		t.Fatal("could not connect to store")
	}

	mux := adapter.NewMux().
		SetURL("http://test.com:8080").
		SetDomain("test.com").
		SetStore("mysqldb", store).
		AddRPD("localhost:8081").
		Compile(demobasic.MuxFunc())

	qs := []string{
		"INSERT INTO `basics` (`uc`, `name`, `owner_id`, `perms`, `hash`) VALUES ('a', 'Goat', 1, ':::', 'abc')",
	}
	testDBSetup(db, t, qs...)

	f := make(url.Values)
	f.Set("name", "Tree")
	req := httptest.NewRequest(http.MethodPost, "/basics", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = "api.test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Assert(
			jsonpath.Chain().
				Equal("name", "Tree").
				End(),
		).
		End()

	var id uint
	err := db.QueryRow("SELECT id from basics WHERE name='Tree'").Scan(&id)
	if err != nil {
		t.Error(err)
	}

	assert := assert.New(t)
	assert.Equal(uint(2), id)

	req = httptest.NewRequest(http.MethodGet, "/basics", nil)
	req.Host = "api.test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Assert(
			jsonpath.Chain().
				Contains(`$.*.name`, "Goat").
				Contains(`$.*.name`, "Tree").
				End(),
		).
		End()

	req = httptest.NewRequest(http.MethodGet, "/basics/a", nil)
	req.Host = "api.test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Assert(
			jsonpath.Equal("$.name", "Goat"),
		).
		End()

	f = make(url.Values)
	f.Set("name", "Fish")
	req = httptest.NewRequest(http.MethodPut, "/basics", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Host = "api.test.com"

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		End()

	var nm string
	err = db.QueryRow("SELECT id,name from basics WHERE name='Fish'").Scan(&id, &nm)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(uint(2), id)
	assert.Equal("Fish", nm)
}

func testDBDropTables(t *testing.T, dbt, uri string) *sql.DB {
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}

	db.Exec("DROP TABLE routings")
	db.Exec("DROP TABLE basics")
	return db
}

func testDBSetup(db *sql.DB, t *testing.T, qs ...string) {
	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatal(err)
		}
	}
}
