package proxy

import (
	"context"
)

type StartServ interface {
	WithContext(context.Context)

	Start() error
	Stop()
	Serve() error
}
