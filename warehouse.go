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
	log.Printf("[dynamic] warehouse load plugin: %s", name)

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

	plugin, err := w.Local.PluginLoad(name)
	if err != nil {
		log.Printf("[dynamic] warehouse load plugin %s failed: %v", name, err)
		return nil, err
	}

	log.Printf("[dynamic] warehouse load plugin %s success", name)
	return plugin, nil
}
