package fsm

type State string

type Transition struct {
	From State
	To   State
}

type FSM struct {
	transitions map[State]map[State]bool
}

func (f *FSM) Exists(from State, to State) bool {
	if _, found := f.transitions[from]; !found {
		return false
	}
	_, found := f.transitions[from][to]
	return found
}

func New(transitions []Transition) *FSM {
	fsm := &FSM{
		transitions: map[State]map[State]bool{},
	}

	for _, t := range transitions {
		if fsm.transitions[t.From] == nil {
			fsm.transitions[t.From] = map[State]bool{}
		}
		fsm.transitions[t.From][t.To] = true
	}

	return fsm
}
