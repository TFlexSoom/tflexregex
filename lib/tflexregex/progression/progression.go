package progression

type Progression interface {
	Clear()
	Anchored()
	Dollared()
	AddCharFilter(byte)
	AddRuneFilter(rune)
	AddRuneListFilter([]rune)
	AddDotFilter()
	AddModifier(uint, uint)
	Group() Progression
	Degroup() Progression
	Union() Progression
}

func NewProgression() Progression {
	return nil
}

// TODO IMPLEMENT
