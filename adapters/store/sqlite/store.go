package sqlite

import (
	"database/sql"
	"encoding/json"

	"brain/core/object"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func Open(path string) *Store {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS objects (
		id BLOB PRIMARY KEY,
		type TEXT,
		author BLOB,
		body BLOB,
		signature BLOB
	);
	CREATE TABLE IF NOT EXISTS links (
		parent BLOB,
		child BLOB
	);
	`)
	if err != nil {
		panic(err)
	}
	return &Store{db}
}

func (s *Store) Put(o *object.Object) {
	body, err := json.Marshal(o.Body)
	if err != nil {
		panic(err)
	}
	_, err = s.db.Exec(`INSERT OR IGNORE INTO objects VALUES (?,?,?,?,?)`,
		o.ID[:], o.Type, o.Author, body, o.Signature)
	if err != nil {
		panic(err)
	}

	for _, p := range o.Parents {
		_, err = s.db.Exec(`INSERT INTO links VALUES (?,?)`, p[:], o.ID[:])
		if err != nil {
			panic(err)
		}
	}
}

func (s *Store) All() []*object.Object {
	rows, err := s.db.Query(`SELECT id,type,author,body,signature FROM objects`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var out []*object.Object
	for rows.Next() {
		o := object.Object{}
		var body []byte
		var id []byte
		var author []byte
		err = rows.Scan(&id, &o.Type, &author, &body, &o.Signature)
		if err != nil {
			panic(err)
		}
		copy(o.ID[:], id)
		o.Author = author
		o.Body = body
		out = append(out, &o)
	}
	return out
}

type Link struct {
	Parent object.ID
	Child  object.ID
}

func (s *Store) Links() []Link {
	rows, err := s.db.Query(`SELECT parent, child FROM links`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var out []Link
	for rows.Next() {
		var l Link
		var p, c []byte
		err = rows.Scan(&p, &c)
		if err != nil {
			panic(err)
		}
		copy(l.Parent[:], p)
		copy(l.Child[:], c)
		out = append(out, l)
	}
	return out
}
