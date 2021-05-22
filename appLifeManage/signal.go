package appLifeManage

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type SignalHandle interface {
	Run()
}

type Signal struct {
	handler map[os.Signal]SignalHandle
	lock    sync.Mutex
	sigChan chan os.Signal
}

func NewSignal(sigChan chan os.Signal) *Signal {
	return &Signal{
		handler: make(map[os.Signal]SignalHandle),
		sigChan: sigChan,
	}
}

func (s *Signal) AddSignalHandler(sigFlag os.Signal, handle SignalHandle) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.handler[sigFlag] = handle
}

func (s *Signal) Remove(sigFlag os.Signal) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.handler[sigFlag]; ok {
		delete(s.handler, sigFlag)
	}
}
func (s *Signal) RegisterSignals() {
	signal.Notify(s.sigChan, s.Signals()...)
}

func (s *Signal) Handle(signal os.Signal) {
	if f, ok := s.handler[signal]; ok {
		f.Run()
	} else {
		fmt.Printf("%d 没有对应处理函数", signal)
	}

}

func (s *Signal) GetSignal() chan os.Signal {
	return s.sigChan
}

func (s *Signal) Signals() []os.Signal {
	var sis []os.Signal
	for si, _ := range s.handler {
		sis = append(sis, si)
	}
	if len(sis) <= 0 {
		panic("must add handler")
	}
	return sis
}
