package object

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"sort"

	"github.com/zeebo/blake3"
)

type ID [32]byte

type Object struct {
	ID        ID
	Author    ed25519.PublicKey
	Type      string
	Body      json.RawMessage
	Parents   []ID
	Signature []byte
}

type unsignedObject struct {
	Author  []byte
	Type    string
	Body    json.RawMessage
	Parents []ID
}

func (o *Object) canonical() []byte {

	u := unsignedObject{
		Author:  o.Author,
		Type:    o.Type,
		Body:    o.Body,
		Parents: append([]ID{}, o.Parents...),
	}

	sort.Slice(u.Parents, func(i, j int) bool {
		return bytes.Compare(u.Parents[i][:], u.Parents[j][:]) < 0
	})

	b, _ := json.Marshal(u)
	return b
}

func New(author ed25519.PublicKey, objType string, body json.RawMessage, parents []ID, privKey ed25519.PrivateKey) (*Object, error) {
	o := &Object{
		Author:  author,
		Type:    objType,
		Body:    body,
		Parents: parents,
	}

	cannon := o.canonical()
	o.ID = blake3.Sum256(cannon)
	o.Signature = ed25519.Sign(privKey, cannon)

	return o, nil
}

func (o *Object) Verify() bool {
	cannon := o.canonical()
	id := blake3.Sum256(cannon)
	if o.ID != id {
		return false
	}
	return ed25519.Verify(o.Author, o.ID[:], o.Signature)
}
