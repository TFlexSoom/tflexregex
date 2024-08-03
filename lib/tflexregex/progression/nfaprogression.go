package progression

type transition map[rune][]int

type nfaProgression struct {
	states    []transition
	iterators []int
}

const nullRune = rune(0)
const anyRune = rune(1)

const defaultStateCapacity = 256
const maxStateCapacity = 65536

const maxTransitionCapacity = 256

const defaultIteratorsCapacity = 256
const maxIteratorsCapacity = 65536

func newNfaProgression() nfaProgression {
	states := make([]transition, 1, defaultStateCapacity)
	states[0] = make(map[rune][]int, 2)
	states[0][nullRune] = []int{1}
	states[0][anyRune] = []int{0}

	return nfaProgression{
		states:    states,
		iterators: make([]int, 0, defaultIteratorsCapacity),
	}
}

func (n *nfaProgression) Clear() {
	(*n) = newNfaProgression()
}

func (n *nfaProgression) Anchored() {
	if len(n.states) != 0 {
		panic("anchored on unanchorable progression: nfa progression")
	}

	delete((*n).states[0], anyRune)
}

func (n *nfaProgression) Dollared() {
	(*n).states = append(n.states, make(map[rune][]int, 1))
}

func (n *nfaProgression) AddCharFilter(b byte) {
	length := len(n.states)
	transition := make(map[rune][]int, 1)
	transition[rune(b)] = []int{length + 1}
	(*n).states = append(n.states, transition)
}

func (n *nfaProgression) AddRuneFilter(r rune) {
	length := len(n.states)
	transition := make(map[rune][]int, 1)
	transition[r] = []int{length + 1}
	(*n).states = append(n.states, transition)
}

func (n *nfaProgression) AddRuneListFilter(lst []rune) {
	if len(lst) > maxTransitionCapacity {
		panic("rune list filter to large for performance: nfa progression")
	}

	length := len(n.states)
	transition := make(map[rune][]int, len(lst))
	for _, r := range lst {
		transition[r] = []int{length + 1}
	}
	(*n).states = append(n.states, transition)
}

func (n *nfaProgression) AddDotFilter() {
	length := len(n.states)
	transition := make(map[rune][]int, 1)
	transition[anyRune] = []int{length + 1}
	(*n).states = append(n.states, transition)
}

func (n *nfaProgression) AddModifier(uint, uint) {

}

func (n *nfaProgression) Group() Progression {

}

func (n *nfaProgression) Degroup() Progression {

}

func (n *nfaProgression) Union() Progression {

}
