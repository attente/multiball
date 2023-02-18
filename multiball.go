package multiball

type Multiball struct {
	Finish chan<- int64

	next int64
	done []bool
	setNext func(next int64)
}

func NewMultiball(next int64, setNext func(next int64)) Multiball {
	finish := make(chan int64)

	m := Multiball{
		Finish: finish,
		next: next,
		done: nil,
		setNext: setNext,
	}

	go func() {
		for ball := range finish {
			if ball < m.next {
				continue
			}

			if ball - m.next >= int64(len(m.done)) {
				n := ball - m.next + 1 - int64(len(m.done))
				m.done = append(m.done, make([]bool, n)...)
			}

			m.done[ball - m.next] = true

			if ball == m.next {
				for len(m.done) != 0 && m.done[0] {
					m.next++
					m.done = m.done[1:]
				}

				m.setNext(m.next)
			}
		}
	}()

	return m
}

func (m *Multiball) Close() {
	close(m.Finish)
}
