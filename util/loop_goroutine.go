package util

import "time"

type LoopTask struct {
	f       func()
	c       chan bool
	timeout time.Duration
}

func NewLoopTask(f func(), timeout time.Duration) *LoopTask {
	loopTask := &LoopTask{
		f:       f,
		c:       make(chan bool),
		timeout: timeout,
	}
	go loopTask.run()
	return loopTask
}

func (l *LoopTask) Close() {
	if int64(l.timeout) == 0 {
		l.c <- true

	} else {
		select {
		case l.c <- true:
		case time.After(l.timeout):
		}
	}
	close(l.c)
}

func (l *LoopTask) run() {
	for {
		select {
		case <-l.c:
			return
		default:
		}
		l.f()
	}
}
