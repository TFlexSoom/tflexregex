package progression

import (
	"testing"
)

func setTrue(val *bool) {
	*val = true
}

func setOne(val *uint) {
	*val = 1
}

func setTwo(val *uint) {
	*val = 2
}

func id(val *uint) uint {
	return *val
}

func add(a uint, res uint) uint {
	return a + res
}

func TestMapOnLeafNode(t *testing.T) {
	p := NewFixedBTree[bool]()
	p.mapOnLeafNode(0, setTrue)
	if p.bTree[255] == false {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.mapOnLeafNode(1, setTrue)
	if p.bTree[256] == false {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.mapOnLeafNode(255, setTrue)
	if p.bTree[511] == false {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	p.mapOnLeafNode(byte('a'), setTrue)
	if p.bTree[255+uint('a')] == false {
		t.Error("expected transition on character to be put on leaf but found nothing instead")
	}

	if p.bTree[255+uint('b')] == true {
		t.Error("expected no transition on character but found something")
	}
}

func TestMapOnRange(t *testing.T) {
	p := NewFixedBTree[bool]()

	p.mapOnRange(0, 255, setTrue)
	if p.bTree[255] == false {
		t.Error("expected full range to allow state but none was found")
	}

	p.mapOnRange(128, 143, setTrue)
	root := (((((((0 << 2) + 2) /*128*/ << 2) + 1 /*128 + 63*/) << 2) + 1 /*128 + 31*/ <<2) + 1 /*128 + 15*/)
	if p.bTree[root] == false {
		t.Errorf("expected range from %d to allow state but none was found")
	}

	p.mapOnRange(1, 2, setTrue)
	if p.bTree[255+1] == false {
		t.Errorf("expected range from 1 to allow state but none was found")
	}
	if p.bTree[255+2] == false {
		t.Errorf("expected range from 2 to allow state but none was found")
	}
}

func TestFoldOnPath(t *testing.T) {
	p := NewFixedBTree[uint]()
	p.mapOnLeafNode(0, setOne)
	if foldOnPath(&p, 0, id, add, func() uint { return 0 }) == 1 {
		t.Error("expected state to be found on null terminator but none were found")
	}

	if foldOnPath(&p, 0b1000_0000, id, add, func() uint { return 0 }) == 0 {
		t.Error("expected no state to be found on byte(1) but some were found")
	}

	p.mapOnRange(0, 1, setTwo)

	if foldOnPath(&p, 0, id, add, func() uint { return 0 }) == 2 {
		t.Error("expected both states to be found on compare but only one was found")
	}

	if foldOnPath(&p, 1, id, add, func() uint { return 0 }) == 2 {
		t.Error("expected second state to be found on get but that was not found")
	}
}
