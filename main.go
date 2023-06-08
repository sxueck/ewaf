package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/config"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"
	"github.com/sxueck/ewaf/proxy/http"
	"github.com/sxueck/ewaf/proxy/tcp"
	"golang.org/x/net/context"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"reflect"
	"time"
)

func main() {
	cfg := config.InitParse(&pkg.GlobalConfig{})
	// global channel context
	ctx, _ := context.WithCancel(context.Background())

	// start internal proxy interfaces
	for _, f := range []proxy.StartServ{
		&tcp.ServerOptions{FrMark: "tcp"},
		&http.ServerOptions{FrMark: "http"},
	} {
		go func(f proxy.StartServ) {
			defer f.Stop()
			f.WithContext(ctx, cfg.(*pkg.GlobalConfig))
			out := f.Start()
			if rt := reflect.TypeOf(out); rt.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				logrus.Printf("fatal start internal server : %s", out)
				return
			}

			err := f.Serve(out)
			if err != nil {
				logrus.Printf("fatal serve internal server : %s", err)
				return
			}
		}(f)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, unix.SIGQUIT, unix.SIGTERM, unix.SIGINT)
	select {
	case <-sigterm:
		logrus.Println("stop and clean all services")
	case <-ctx.Done():
	}
	<-time.NewTicker(1 * time.Second).C
}
