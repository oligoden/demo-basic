package demobasic

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/oligoden/chassis/device/model/data"
)

type Record struct {
	Name string `json:"name" form:"name"`
	data.Default
}

func NewRecord() *Record {
	r := &Record{}
	r.Default = data.Default{}
	r.Perms = ":::cr"
	return r
}

func (e Record) Prepare() error {
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9- ']{0,40}[a-zA-Z0-9]$", e.Name)
	if !match {
		return fmt.Errorf("bad request, name not valid")
	}

	return nil
}

func (Record) TableName() string {
	return "basics"
}

func (Record) Migrate(db *sql.DB) error {
	q := "CREATE TABLE IF NOT EXISTS `basics` (`name` varchar(255) DEFAULT NULL, `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `uc` varchar(255) DEFAULT NULL, `owner_id` int(10) unsigned DEFAULT NULL, `perms` varchar(255) DEFAULT NULL, `hash` varchar(255) DEFAULT NULL, PRIMARY KEY (`id`), UNIQUE KEY `uc` (`uc`) )"

	_, err := db.Exec(q)
	if err != nil {
		fmt.Println("basics migration error", err)
		return fmt.Errorf("doing migration: %w", err)
	}

	return nil
}

type List map[string]Record

func (List) TableName() string {
	return "basics"
}

func (e List) Permissions(p ...string) string {
	return ""
}

func (e List) Owner(o ...uint) uint {
	return 0
}

func (e List) Users(u ...uint) []uint {
	return []uint{}
}

func (e List) Groups(g ...uint) []uint {
	return []uint{}
}

func (List) IDValue(...uint) uint {
	return 0
}

func (e List) UniqueCode(uc ...string) string {
	return ""
}

func NewList() List {
	return List{}
}

func (List) Complete() error {
	return nil
}

func (List) Hasher() error {
	return nil
}

func (List) Prepare() error {
	return nil
}
