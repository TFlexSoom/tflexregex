package regex

type ThreeSatNFA struct {
}

func NFAFromSingle(ir Ir) Regex {
	return ThreeSatNFA{}
}

func NFAFromGroup(irs []Ir) RegexGroup {
	return ThreeSatNFA{}
}

func (r ThreeSatNFA) Matches([]byte) bool {
	return false
}

func (r ThreeSatNFA) MatchesUnicode(string) bool {
	return false
}

func (r ThreeSatNFA) First([]byte) int {
	return -1
}

func (r ThreeSatNFA) All([]byte) []int {
	return []int{}
}
