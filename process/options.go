package process

import (
	"github.com/doubletrey/crawlab-core/interfaces"
	"time"
)

type DaemonOption func(d interfaces.ProcessDaemon)

func WithDaemonMaxErrors(maxErrors int) DaemonOption {
	return func(d interfaces.ProcessDaemon) {
		d.SetMaxErrors(maxErrors)
	}
}

func WithExitTimeout(timeout time.Duration) DaemonOption {
	return func(d interfaces.ProcessDaemon) {

	}
}
