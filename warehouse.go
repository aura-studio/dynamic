package dynamic

import (
	"errors"
	"log"
)

type Warehouse struct {
	Local  *Local
	Remote Remote
}

var warehouse = NewWarehouse()

func NewWarehouse() *Warehouse {
	return &Warehouse{}
}

func (w *Warehouse) Init(localPath, remotePath string) {
	w.Local = NewLocal(localPath)
	w.Remote = NewRemote(remotePath)
}

func (w *Warehouse) Load(name string) (any, error) {
	log.Printf("[dynamic] warehouse loading package: %s", name)

	if w.Local == nil {
		return nil, errors.New("dynamic: warehouse package not exists")
	}

	if !w.Local.Exists(name) {
		if w.Remote == nil {
			return nil, errors.New("dynamic: warehouse package not exists")
		}

		if err := w.Remote.Sync(name); err != nil {
			return nil, err
		}

		if !w.Local.Exists(name) {
			return nil, errors.New("dynamic: warehouse package not exists")
		}
	}

	pkg, err := w.Local.Load(name)
	if err != nil {
		log.Printf("[dynamic] warehouse load package %s failed: %v", name, err)
		return nil, err
	}

	log.Printf("[dynamic] warehouse load package: %s success", name)
	return pkg, nil
}
