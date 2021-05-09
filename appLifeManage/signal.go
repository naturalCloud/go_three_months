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

func (s *Signal) Add(sigFlag os.Signal, handle SignalHandle) {
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
func (s *Signal) Register(signals []os.Signal) {
	signal.Notify(s.sigChan, signals...)
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
