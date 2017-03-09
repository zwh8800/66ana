package util

type empty struct{}
type Semaphore chan empty

func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}

func (s Semaphore) P(n int) {
	for i := 0; i < n; i++ {
		s <- empty{}
	}
}

func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

/* mutexes */

func (s Semaphore) Lock() {
	s.P(1)
}

func (s Semaphore) Unlock() {
	s.V(1)
}

/* signal-wait */

func (s Semaphore) Signal() {
	s.V(1)
}

func (s Semaphore) Wait(n int) {
	s.P(n)
}
