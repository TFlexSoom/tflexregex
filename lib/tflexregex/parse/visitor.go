package parse

import (
	"fmt"
	"strconv"

	"github.com/tflexsoom/tflexregex/lib/tflexregex/progression"
)

type visitor struct {
	progressionRoot    progression.Progression
	progressionCurrent progression.Progression
	buffer             []byte
}

const maxBuffer = 32

func newVisitor() visitor {
	p := progression.NewProgression()

	return visitor{
		progressionRoot:    p,
		progressionCurrent: p,
		buffer:             make([]byte, 0, maxBuffer),
	}
}

func (v *visitor) addToBuffer(b byte) {
	(*v).buffer = append((*v).buffer, b)
}

func (v *visitor) addManyToBuffer(bs []byte) {
	(*v).buffer = append((*v).buffer, bs...)
}

func (v *visitor) pump(nt nodeType) {
	switch nt {
	case NODE_TYPE_NONE:
		panic("node type none passed to visit")
	case NODE_TYPE_ANCHOR:
		v.anchor()
		break
	case NODE_TYPE_DOLLAR:
		v.dollar()
		break
	case NODE_TYPE_CHAR:
		v.char()
		break
	case NODE_TYPE_UNICODE:
		v.unicode()
		break
	case NODE_TYPE_CLASS:
		v.class()
		break
	case NODE_TYPE_DOT:
		v.dot()
		break
	case NODE_TYPE_MODIFIER:
		v.modifier()
		break
	case NODE_TYPE_PAR_BEGIN:
		v.beginParenthesis()
		break
	case NODE_TYPE_PAR_END:
		v.endParenthesis()
		break
	case NODE_TYPE_UNION:
		v.union()
		break
	default:
		panic("unknown node type passed to visit (possibly unimplemented)")
	}

}

func (v visitor) validateLength() {
	if len(v.buffer) == 0 {
		panic("empty buffer")
	}
}

func (v *visitor) clearBuffer() {
	clear((*v).buffer)
}

func (v *visitor) dropNext(next byte) {
	if v.buffer[0] != next {
		panic("worng character in buffer")
	}

	v.clearBuffer()
}

func (v *visitor) anchor() {
	v.validateLength()
	v.clearBuffer()
	(*v).progressionRoot.Anchored()
}

func (v *visitor) dollar() {
	v.validateLength()
	v.clearBuffer()
	(*v).progressionRoot.Dollared()
}

func (v *visitor) char() {
	v.validateLength()
	(*v).progressionCurrent.AddCharFilter(v.buffer[0])
	v.clearBuffer()
}

func (v *visitor) unicode() {
	currentIndex := 0
	current := v.buffer[currentIndex]
	r := int32(0)
	if current >= byte(0b1111_0000) {
		r += int32(current) << 24
		currentIndex += 1
		current = v.buffer[currentIndex]
	}

	if current >= byte(0b1110_0000) {
		r += int32(current) << 16
		currentIndex += 1
		current = v.buffer[currentIndex]
	}

	if current >= byte(0b1100_0000) {
		r += int32(current) << 8
		r += int32(v.buffer[currentIndex+1])
		(*v).progressionCurrent.AddRuneFilter(rune(r))
		v.clearBuffer()
		return
	}

	if current >= byte(0b1000_0000) {
		panic("invalid unicode character detected")
	}

	v.char()
}

func (v *visitor) class() {
	v.validateLength()
	class := (*v).buffer[:]
	v.clearBuffer()
	length := len(class)

	if class[0] != '[' {
		panic("unexpected first byte for class")
	} else if class[length-1] != ']' {
		panic("unexpected last byte for class")
	}

	runeList := make([]rune, 0, 32)
	for i := 1; i < length-1; i++ {
		currentIndex := i
		current := class[currentIndex]
		r := int32(0)
		if current >= byte(0b1111_0000) {
			r += int32(current) << 24
			currentIndex += 1
			current = class[currentIndex]
		}

		if current >= byte(0b1110_0000) {
			r += int32(current) << 16
			currentIndex += 1
			current = class[currentIndex]
		}

		if current >= byte(0b1100_0000) {
			r += int32(current) << 8
			r += int32(class[currentIndex+1])
			runeList = append(runeList, rune(r))
			i = currentIndex + 1
			continue
		}

		runeList = append(runeList, rune(current))
	}

	(*v).progressionCurrent.AddRuneListFilter(runeList)
}

func (v *visitor) dot() {
	v.validateLength()
	v.clearBuffer()
	(*v).progressionCurrent.AddDotFilter()
}

func (v *visitor) rangeModifier(modifier []byte) {

}

func (v *visitor) modifier() {
	v.validateLength()
	modifier := (*v).buffer[:]
	clear((*v).buffer)
	length := len(modifier)

	if length != 1 {
		(*v).rangeModifier(modifier)
	}

	// TODO Split on comma and strconv both sides
	res, err := strconv.Atoi(string(v.buffer))
	if err != nil {
		panic(fmt.Sprintf("could not convert decimal: %v"))
	}

	(*v).progressionCurrent.AddModifier(uint(res), uint(res))
}

func (v *visitor) beginParenthesis() {
	(*v).progressionCurrent = (*v).progressionCurrent.Group()
	v.clearBuffer()
}

func (v *visitor) endParenthesis() {
	(*v).progressionCurrent = (*v).progressionCurrent.Degroup()
	v.clearBuffer()
}

func (v *visitor) union() {
	(*v).progressionCurrent = (*v).progressionCurrent.Union()
	v.clearBuffer()
}
