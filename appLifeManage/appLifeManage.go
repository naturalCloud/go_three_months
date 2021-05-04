package appLifeManage

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"syscall"
)

type SigintHandler struct {
	stop func()
}

func (si *SigintHandler) Run() {
	si.stop()
}

type SighupHandler struct {
	exitChan chan struct{}
	stopChan chan struct{}
}

func (sh *SighupHandler) Run() {
	close(sh.stopChan)
	sh.exitChan <- struct{}{}
}

func Run() {

	ctx, cancel := context.WithCancel(context.Background())
	g, c := errgroup.WithContext(ctx)

	exitChan := make(chan struct{})
	stop := make(chan struct{})

	s := NewSignal()
	s.Add(syscall.SIGINT, &SigintHandler{
		stop: cancel,
	})
	s.Add(syscall.SIGHUP, &SighupHandler{
		exitChan: exitChan,
		stopChan: stop,
	})
	s.Reg([]os.Signal{syscall.SIGHUP, syscall.SIGINT})
	go s.Handle()
	for _, s := range []Server{ServerBiz(":9834"), ServerDebug(":8789")} {
		s := s
		g.Go(func() error {
			return s.Start()
		})
		g.Go(func() error {
			<-c.Done()
			return s.Stop()
		})
	}
	err := g.Wait()
	fmt.Println(err)

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
