package tflexregex

import (
	"slices"
	"testing"
)

func TestTransitionOnCharacter(t *testing.T) {
	state := uint(1)
	p := NewProgression()
	p.TransitionOnCharacter(0, state)
	if p.binaryTreeOfSets[255] == nil || p.binaryTreeOfSets[255][state] != true {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.TransitionOnCharacter(0, state+1)
	if p.binaryTreeOfSets[255] == nil || p.binaryTreeOfSets[255][state+1] != true {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.TransitionOnCharacter(255, state)
	if p.binaryTreeOfSets[511] == nil || p.binaryTreeOfSets[511][state] != true {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.TransitionOnCharacter(byte('a'), state)
	if p.binaryTreeOfSets[255+uint('a')] == nil || p.binaryTreeOfSets[255+uint('a')][state] != true {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	if p.binaryTreeOfSets[510] != nil {
		t.Error("expected no transition on character but found something")
	}
}

func TestTransitionOnRange(t *testing.T) {
	state_0 := uint(1)
	state_1 := uint(2)
	state_2 := uint(3)
	p := NewProgression()

	p.TransitionOnRange(0, 255, state_0)
	if p.binaryTreeOfSets[0] == nil || p.binaryTreeOfSets[0][state_0] != true {
		t.Error("expected full range to allow state but none was found")
	}

	p.TransitionOnRange(128, 143, state_1)
	root := (((((((0 << 2) + 2) /*128*/ << 2) + 1 /*128 + 63*/) << 2) + 1 /*128 + 31*/ <<2) + 1 /*128 + 15*/)
	if p.binaryTreeOfSets[root] == nil || p.binaryTreeOfSets[root][state_1] != true {
		t.Errorf("expected range from %d to allow state but none was found")
	}

	p.TransitionOnRange(1, 2, state_2)
	if p.binaryTreeOfSets[255+1] == nil || p.binaryTreeOfSets[255+1][state_2] != true {
		t.Errorf("expected range from 1 to allow state but none was found")
	}
	if p.binaryTreeOfSets[255+2] == nil || p.binaryTreeOfSets[255+2][state_2] != true {
		t.Errorf("expected range from 2 to allow state but none was found")
	}
}

func TestGetTransition(t *testing.T) {
	state := uint(1)
	p := NewProgression()
	p.TransitionOnCharacter(0, state)
	if slices.Compare(p.GetTransitions(0), []uint{state}) == 0 {
		t.Error("expected state to be found on null terminator but none were found")
	}

	if slices.Compare(p.GetTransitions(1), []uint{}) == 0 {
		t.Error("expected no state to be found on byte(1) but some were found")
	}

	p.TransitionOnRange(0, 1, state+1)

	if slices.Compare(p.GetTransitions(0), []uint{state, state + 1}) == 0 {
		t.Error("expected both states to be found on compare but only one was found")
	}

	if slices.Compare(p.GetTransitions(1), []uint{state + 1}) == 0 {
		t.Error("expected second state to be found on get but that was not found")
	}
}
