package store

import "brain/core/object"

type Store interface {
	Put(obj *object.Object) error
	Get(id object.ID) (*object.Object, error)
	Has(id object.ID) (bool, error)
}
