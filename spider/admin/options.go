package admin

import "github.com/doubletrey/crawlab-core/interfaces"

type Option func(svc interfaces.SpiderAdminService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.SpiderAdminService) {
		svc.SetConfigPath(path)
	}
}
