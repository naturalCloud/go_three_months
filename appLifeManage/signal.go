package appLifeManage

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SignalHandle interface {
	Run()
}

type Signal struct {
	handler map[os.Signal]SignalHandle
	lock    sync.Mutex
	sigChan chan os.Signal
}

func NewSignal() *Signal {
	return &Signal{handler: make(map[os.Signal]SignalHandle), sigChan: make(chan os.Signal)}
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
func (s *Signal) Reg() *Signal {
	signal.Notify(s.sigChan, syscall.SIGHUP, syscall.SIGINT)
	return s
}

func (s *Signal) Handle() {
	for value := range s.sigChan {
		if f, ok := s.handler[value]; ok {
			f.Run()
		} else {
			fmt.Println(value)
		}
	}
}
