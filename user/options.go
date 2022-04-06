package user

import "github.com/doubletrey/crawlab-core/interfaces"

type Option func(svc interfaces.UserService)

func WithJwtSecret(secret string) Option {
	return func(svc interfaces.UserService) {
		svc.SetJwtSecret(secret)
	}
}
