package utils

import "github.com/doubletrey/crawlab-core/constants"

func IsCancellable(status string) bool {
	switch status {
	case constants.TaskStatusPending,
		constants.TaskStatusRunning:
		return true
	default:
		return false
	}
}
