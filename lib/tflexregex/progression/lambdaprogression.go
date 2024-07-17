package progression

type filter func(byte) bool

type lambdaProgression struct {
	filters  []filter
	anchored bool
	dollared bool
}

const defaultCapacity = 256
const maxCapacity = 65536

func newLambdaProgression() lambdaProgression {
	return lambdaProgression{
		filters:  make([]filter, 0, defaultCapacity),
		anchored: false,
		dollared: false,
	}
}

func (l *lambdaProgression) Clear() {
	clear((*l).filters)
	(*l).anchored = false
	(*l).dollared = false
}

func (l *lambdaProgression) Anchored() {
	if l.anchored {
		panic("anchored called a second time: lambda progression")
	}

	(*l).anchored = true
}

func (l *lambdaProgression) Dollared() {
	if l.dollared {
		panic("dollared called a second time: lambda progression")
	}

	(*l).dollared = true
}

func (l *lambdaProgression) addFilter(filter filter) {
	if len(l.filters) >= maxCapacity {
		panic("filters have gone past max capacity")
	}

	(*l).filters = append((*l).filters, filter)
}

func (l *lambdaProgression) AddCharFilter(b byte) {
	(*l).addFilter(func(i byte) bool {
		return b == i
	})
}

func (l *lambdaProgression) AddRuneFilter(r rune) {
	(*l).addFilter(func(i byte) bool {
		return b == i
	})
}

func (l *lambdaProgression) AddRuneListFilter([]rune) {

}

func (l *lambdaProgression) AddDotFilter() {

}

func (l *lambdaProgression) AddModifier(uint, uint) {

}

func (l *lambdaProgression) Group() Progression {

}

func (l *lambdaProgression) Degroup() Progression {

}

func (l *lambdaProgression) Union() Progression {

}
