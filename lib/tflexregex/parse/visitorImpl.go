package parse

import (
	"github.com/tflexsoom/tflexregex/lib/tflexregex/regex"
)

type IrVisitor struct {
	ir              regex.Ir
	groupRowIndexes []uint
}

func NewIrVisitor() *IrVisitor {
	return &IrVisitor{
		ir: regex.Ir{
			Classes: make([][]rune, 0, 8),
			Groups:  make([][]regex.IrGroup, 0, 8),
		},
		groupRowIndexes: make([]uint, 1, 8),
	}
}

func (v *IrVisitor) Ir() regex.Ir {
	return (*v).ir
}

func (v *IrVisitor) row() *[]regex.IrGroup {
	length := len((*v).groupRowIndexes)
	groupRowIndex := (*v).groupRowIndexes[length-1]
	return &((*v).ir.Groups[groupRowIndex])
}

func (v *IrVisitor) anchor() visitor {
	(*v).ir.Anchored = true

	return v
}

func (v *IrVisitor) dollar() visitor {
	(*v).ir.Dollared = true

	return v
}

func (v *IrVisitor) char(b byte) visitor {
	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code:        regex.IR_GROUP_CODE_NO_CODE,
		RuneOrIndex: uint32(b),
	})

	return v
}

func (v *IrVisitor) unicode(r rune) visitor {
	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code:        regex.IR_GROUP_CODE_NO_CODE,
		RuneOrIndex: uint32(r),
	})

	return v
}

func (v *IrVisitor) class(rs []rune) visitor {
	(*v).ir.Classes = append((*v).ir.Classes, rs)

	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code:        regex.IR_GROUP_CODE_CLASS_STUB,
		RuneOrIndex: uint32(len((*v).ir.Classes) - 1),
	})

	return v
}

func (v *IrVisitor) dot() visitor {
	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code: regex.IR_GROUP_CODE_DOT,
	})

	return v
}

func (v *IrVisitor) modifier(min uint, max uint) visitor {
	row := v.row()
	length := len(*row)
	if length == 0 {
		panic("nothing to modify")
	}

	(*row)[length-1].Min = min
	(*row)[length-1].Max = max

	return v
}

func (v *IrVisitor) openParenthesis() visitor {
	(*v).ir.Groups = append((*v).ir.Groups, make([]regex.IrGroup, 0, 64))
	subIndex := uint(len((*v).ir.Groups) - 1)
	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code:        regex.IR_GROUP_CODE_SUBGROUP_STUB,
		RuneOrIndex: uint32(subIndex),
	})
	(*v).groupRowIndexes = append((*v).groupRowIndexes, subIndex)

	return v
}

func (v *IrVisitor) closeParenthesis() visitor {
	length := len((*v).groupRowIndexes)
	(*v).groupRowIndexes = (*v).groupRowIndexes[:length-1]

	return v
}

func (v *IrVisitor) union() visitor {
	row := v.row()
	*row = append(*row, regex.IrGroup{
		Code: regex.IR_GROUP_CODE_UNION,
	})

	return v
}
