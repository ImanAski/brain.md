package sqlite

import (
	"brain/core/object"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"

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

func (s *Store) Get(id object.ID) (*object.Object, error) {
	row := s.db.QueryRow(`SELECT id,type,author,body,signature FROM objects WHERE id = ?`, id[:])
	o := object.Object{}
	var body []byte
	var sid []byte
	var author []byte
	err := row.Scan(&sid, &o.Type, &author, &body, &o.Signature)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	copy(o.ID[:], sid)
	o.Author = author
	o.Body = body
	return &o, nil
}

func (s *Store) Has(id object.ID) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM objects WHERE id = ?`, id[:]).Scan(&count)
	return count > 0, err
}

func (s *Store) ResolveID(prefix string) (object.ID, error) {
	if len(prefix) == 64 { // Full ID
		var id object.ID
		raw, err := hex.DecodeString(prefix)
		if err != nil {
			return id, err
		}
		copy(id[:], raw)
		return id, nil
	}

	rows, err := s.db.Query(`SELECT id FROM objects WHERE hex(id) LIKE ?`, prefix+"%")
	if err != nil {
		return object.ID{}, err
	}
	defer rows.Close()

	var ids []object.ID
	for rows.Next() {
		var raw []byte
		var id object.ID
		if err := rows.Scan(&raw); err != nil {
			return id, err
		}
		copy(id[:], raw)
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return object.ID{}, fmt.Errorf("no object found with prefix %s", prefix)
	}
	if len(ids) > 1 {
		return object.ID{}, fmt.Errorf("ambiguous prefix %s: matched %d objects", prefix, len(ids))
	}

	return ids[0], nil
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

func (s *Store) GetByType(typ string) []*object.Object {
	if typ == "" {
		return s.All()
	}
	rows, err := s.db.Query(`SELECT id,type,author,body,signature FROM objects WHERE type = ?`, typ)
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

func (s *Store) Heads() []*object.Object {
	rows, err := s.db.Query(`
		SELECT id, type, author, body, signature 
		FROM objects 
		WHERE id NOT IN (SELECT parent FROM links)
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var out []*object.Object
	for rows.Next() {
		o := object.Object{}
		var body, id, author []byte
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
