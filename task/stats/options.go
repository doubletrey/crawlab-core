package stats

import "github.com/doubletrey/crawlab-core/interfaces"

type Option func(service interfaces.TaskStatsService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskStatsService) {
		svc.SetConfigPath(path)
	}
}
