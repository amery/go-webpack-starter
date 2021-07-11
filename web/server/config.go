package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
)

type Reloader interface {
	Reload() error
}

type ServerConfig struct {
	Port            uint16        `default:"8080"`
	PIDFile         string        `default:"/tmp/tableflip.pid"`
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"5s"`
	IdleTimeout     time.Duration `default:"30s"`
	GracefulTimeout time.Duration `default:"60s"`
}

func (sc *ServerConfig) ListenAndServe(router http.Handler) error {
	// Graceful restart mode
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: sc.PIDFile,
	})
	if err != nil {
		return err
	}
	defer upg.Stop()

	// prepare server
	if router == nil {
		router = http.NotFoundHandler()
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", sc.Port),
		Handler:      router,
		ReadTimeout:  sc.ReadTimeout,
		WriteTimeout: sc.WriteTimeout,
		IdleTimeout:  sc.IdleTimeout,
	}

	// listen service port
	l, err := upg.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	log.Print(fmt.Sprintf("Listening %s", l.Addr()))

	// watch signals
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

		for signum := range sig {

			switch signum {
			case syscall.SIGHUP:
				// attempt to reload config on SIGHUP if supported
				// or an upgrade if isn't
				if r, ok := router.(Reloader); ok {
					if err := r.Reload(); err != nil {
						log.Println("Reload failed:", err)
					}
				} else if err := upg.Upgrade(); err != nil {
					log.Println("Upgrade failed:", err)
				}
			case syscall.SIGUSR2:
				// attempt to upgrade on SIGUSR2
				if err := upg.Upgrade(); err != nil {
					log.Println("Upgrade failed:", err)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				// terminate on SIGINT or SIGTERM
				log.Println("Terminating...")
				upg.Stop()
			}
		}
	}()

	// close router before exiting
	if r, ok := router.(io.Closer); ok {
		defer r.Close()
	}

	// start servicing
	go func() {
		err := s.Serve(l)
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// notify being ready for service
	if err = upg.Ready(); err != nil {
		return err
	}
	<-upg.Exit()

	if sc.GracefulTimeout > 0 {
		// graceful shutdown timeout
		time.AfterFunc(sc.GracefulTimeout, func() {
			log.Println("Graceful shutdown timed out")
			os.Exit(1)
		})
	}

	// Wait for connections to drain.
	s.Shutdown(context.Background())

	return nil
}
