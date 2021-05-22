package appLifeManage

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"syscall"
)

type SigintHandler struct {
	cancel func()
}

func (si *SigintHandler) Run() {
	si.cancel()
}

func Run() {

	ctx, cancel := context.WithCancel(context.Background())
	g, c := errgroup.WithContext(ctx)

	s := NewSignal(make(chan os.Signal, 1))
	s.AddSignalHandler(syscall.SIGINT, &SigintHandler{
		cancel: cancel,
	})

	s.RegisterSignals()
	g.Go(func() error {
		for {
			select {
			case sig := <-s.GetSignal():
				s.Handle(sig)
			case <-c.Done():
				return c.Err()
			}

		}
	})
	for _, s := range []Server{ServerBiz(":9834"), ServerDebug(":9835")} {
		s := s
		g.Go(func() error {
			return s.Start()
		})
		g.Go(func() error {
			<-c.Done()
			return s.Stop()
		})
	}
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}

}

func ServerBiz(addr string) Server {
	h := NewHttpHandler()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am business server"))
	})
	return NewHttpServer(addr, h)
}

func ServerDebug(addr string) Server {
	h := NewHttpHandler()
	h.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("I am debug server"))
	})
	return NewHttpServer(addr, h)
}
