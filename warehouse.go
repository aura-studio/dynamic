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

func (w *Warehouse) Init(localPath string, remotePath string) *Warehouse {
	log.Println("dynamic: Warehouse Init", localPath, remotePath)
	return &Warehouse{
		Local:  NewLocal(localPath),
		Remote: NewRemote(remotePath),
	}
}

func (w *Warehouse) Load(name string) (any, error) {
	if w.Local == nil {
		return nil, errors.New("dynamic: warehouse plugin not exists")
	}

	if !w.Local.Exists(name) {
		if w.Remote == nil {
			return nil, errors.New("dynamic: warehouse plugin not exists")
		}

		if err := w.Remote.Sync(name); err != nil {
			return nil, err
		}

		if !w.Local.Exists(name) {
			return nil, errors.New("dynamic: warehouse plugin not exists")
		}
	}

	return w.Local.PluginLoad(name)
}
