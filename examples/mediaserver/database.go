package main

import (
	"fmt"
	"sync"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

////////////////////////////////////////////////////////////////////////////////
// SCHEMA

const (
	tableMediaItem = "id INTEGER NOT NULL PRIMARY KEY,fullpath TEXT,rootpath TEXT,relpath TEXT"
)

var schema = map[string]string{
	"media_item": tableMediaItem,
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Database struct {
	db           *sql.DB    `json:"-"`
	lock         sync.Mutex `json:"-"`
	Hostname     string     `json:"hostname,omitempty"`
	LastModified time.Time  `json:"last_modified,omitempty"`
}

type MediaItem struct {
	ID       uint `json:"id"`
	FullPath string `json:"-"`
	RootPath string `json:"-"`
	RelPath  string `json:"path"`
}

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func NewDatabase(hostname string) (*Database, error) {
	var err error

	// create this
	this := new(Database)
	this.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	this.Hostname = hostname
	this.LastModified = time.Now()

	// create table schema
	for k, v := range schema {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", k, v)
		_, err = this.db.Exec(sql)
		if err != nil {
			this.db.Close()
			return nil, err
		}
	}

	return this, nil
}

func (this *Database) Terminate() {
	if this.db != nil {
		this.db.Close()
	}
}

////////////////////////////////////////////////////////////////////////////////
// INSERT ITEM INTO DATABASE

func (this *Database) Insert(item MediaItem) error {

	// mutex lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// create atomic transaction
	tx, err := this.db.Begin()
	if err != nil {
		return err
	}

	// insert media_item
	statement, err := tx.Prepare("INSERT INTO media_item (fullpath,rootpath,relpath) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(item.FullPath, item.RootPath, item.RelPath)
	if err != nil {
		return err
	}

	// commit transaction
	tx.Commit()

	// update last modified date
	this.LastModified = time.Now()

	// return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// QUERY DATABASE

func (this *Database) Query() ([]MediaItem, error) {
	rows, err := this.db.Query("SELECT id,relpath FROM media_item ORDER BY id")
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	// make empty result object
	r := make([]MediaItem,0)

	// append media items
	for rows.Next() {
		var m MediaItem
		rows.Scan(&m.ID,&m.RelPath)
		r = append(r,m)
	}

	return r,nil
}

/*

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
}

*/
