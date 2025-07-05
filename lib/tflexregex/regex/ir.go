package regex

type Ir struct {
	Anchored bool
	Dollared bool
	Classes  [][]rune
	Groups   [][]IrGroup
}

type IrGroup struct {
	Code        IrGroupCode
	RuneOrIndex uint32
	Min         uint
	Max         uint // Max == 0 ? Infinite : Max
}

type IrGroupCode uint

const (
	IR_GROUP_CODE_NO_CODE IrGroupCode = iota
	IR_GROUP_CODE_SUBGROUP_STUB
	IR_GROUP_CODE_CLASS_STUB
	IR_GROUP_CODE_DOT
	IR_GROUP_CODE_UNION
)
