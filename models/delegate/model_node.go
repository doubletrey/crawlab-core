package delegate

import (
	"github.com/doubletrey/crawlab-core/constants"
	"github.com/doubletrey/crawlab-core/interfaces"
	"time"
)

type ModelNodeDelegate struct {
	n interfaces.Node
	interfaces.ModelDelegate
}

func (d *ModelNodeDelegate) UpdateStatus(active bool, activeTs *time.Time, status string) (err error) {
	d.n.SetActive(active)
	if activeTs != nil {
		d.n.SetActiveTs(*activeTs)
	}
	d.n.SetStatus(status)
	return d.Save()
}

func (d *ModelNodeDelegate) UpdateStatusOnline() (err error) {
	now := time.Now()
	return d.UpdateStatus(true, &now, constants.NodeStatusOnline)
}

func (d *ModelNodeDelegate) UpdateStatusOffline() (err error) {
	return d.UpdateStatus(false, nil, constants.NodeStatusOffline)
}

func NewModelNodeDelegate(n interfaces.Node) interfaces.ModelNodeDelegate {
	return &ModelNodeDelegate{
		n:             n,
		ModelDelegate: NewModelDelegate(n),
	}
}
