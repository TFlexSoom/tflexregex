package parse

type visitor interface {
	anchor() visitor
	dollar() visitor
	char(byte) visitor
	unicode(rune) visitor
	class([]rune) visitor
	dot() visitor
	modifier(min uint, max uint) visitor
	openParenthesis() visitor
	closeParenthesis() visitor
	union() visitor
}
