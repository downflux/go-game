package fsm

import (
	fcpb "github.com/downflux/game/fsm/api/constants_go_proto"
)

type State string

type Transition struct {
	From        State
	To          State
	VirtualOnly bool
}

type FSM struct {
	fsmType     fcpb.FSMType
	transitions map[State]map[State]bool
}

func New(transitions []Transition, fsmType fcpb.FSMType) *FSM {
	fsm := &FSM{
		transitions: map[State]map[State]bool{},
		fsmType:     fsmType,
	}

	for _, t := range transitions {
		if fsm.transitions[t.From] == nil {
			fsm.transitions[t.From] = map[State]bool{}
		}
		fsm.transitions[t.From][t.To] = t.VirtualOnly
	}

	return fsm
}

func (f *FSM) Type() fcpb.FSMType { return f.fsmType }

func (f *FSM) Exists(from State, to State) (bool, bool) {
	if _, found := f.transitions[from]; !found {
		return false, false
	}
	virtualOnly, found := f.transitions[from][to]
	return found, virtualOnly
}
