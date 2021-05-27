package stream

const (
	stageNone        = 0
	stageStateless   = 1
	stageStateful    = 2
	stageNonShortcut = 3
	stageShortcut    = 4
)

type stage struct {
	prepare   func()
	action    func(T, int) (R, bool)
	stageFlag int
}

type stageMachine struct {
	s      *stream
	stages [][]stage
}

func (m *stageMachine) run() {
	prod := make([]T, len(m.s.data))
	copy(prod, m.s.data)

	for _, s := range m.stages {
		switch s[0].stageFlag {
		case stageNone:
		case stageStateless:
			// stateless has multiple actions, which can be connected in series
			// each action receive an element and return the product
			ii := 0
			skip := false
			for i, e := range prod {
				for _, ss := range s {
					e, skip = ss.action(e, i)
					if e == nil {
						break
					}
				}
				if e != nil {
					prod[ii] = e
					ii++
				}
				if skip {
					break
				}
			}
			prod = prod[:ii]
		case stageStateful:
			// stateful only has one action, which receives a slice and returns a slice of products
			res, _ := s[0].action(prod, len(prod))
			prod = res.([]T)
		case stageNonShortcut:
			// non-shortcut only has one action, which receives a slice and returns a slice of products
			s[0].action(prod, len(prod))
		case stageShortcut:
			// shortcut only has one action, which receives an element and consumes it and will break when can be done
			for i, e := range prod {
				_, skip := s[0].action(e, i)
				if skip {
					break
				}
			}
		}
	}
	return
}
