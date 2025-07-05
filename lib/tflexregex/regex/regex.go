package regex

type Regex interface {
	Matches([]byte) bool
	MatchesUnicode(string) bool
}

type RegexGroup interface {
	First([]byte) int
	All([]byte) []int
}
