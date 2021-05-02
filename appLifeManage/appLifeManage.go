package appLifeManage

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

type SigintHandler struct {
	stop chan struct{}
}

func (si *SigintHandler) Run() {
	si.stop <- struct{}{}
}

type SighupHandler struct {
	exitChan chan struct{}
	stopChan chan struct{}
}

func (sh *SighupHandler) Run() {
	close(sh.stopChan)
	wg.Wait()
	sh.exitChan <- struct{}{}
}

func Run() {
	exitChan := make(chan struct{})
	stop := make(chan struct{})
	go ServerBiz("127.0.0.1:8888", stop)
	go ServerDebug("127.0.0.1:6666", stop)
	s := NewSignal()
	s.Add(syscall.SIGINT, &SigintHandler{
		stop: stop,
	})

	s.Add(syscall.SIGHUP, &SighupHandler{
		exitChan: exitChan,
		stopChan: stop,
	})
	s.Reg()
	go s.Handle()
	<-exitChan

}

func ServerBiz(addr string, stop chan struct{}) {
	h := http.NewServeMux()
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am business server"))
	})

	wg.Add(1)
	err := serverNew(addr, h, stop)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("business server exit")
		wg.Done()
	} else {
		fmt.Println(err)
	}

}

func ServerDebug(addr string, stop chan struct{}) {
	h := http.DefaultServeMux
	h.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("I am debug server"))
	})

	wg.Add(1)
	err := serverNew(addr, h, stop)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("debug exit")
		wg.Done()
	} else {
		fmt.Println(err)
	}
}

func serverNew(addr string, handler http.Handler, stop chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()

}
