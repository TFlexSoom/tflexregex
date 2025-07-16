package regex

const FULL_BITS = 0 - 1

type recursiveTreeFlags byte

type recursiveMonad struct {
	source Ir
	stack  []*recursiveTree
	group  IrGroup
}

type recursiveTree struct {
	children []recursiveTree
	and      rune
	valid    rune
	min      uint
	max      uint
	union    bool
}

type RecursiveTreeRoot struct {
	tree     recursiveTree
	anchored bool
	dollared bool
}

type RecursiveTreeGroup struct {
	trees []RecursiveTreeRoot
}

var recursiveTreeConsumers = map[IrGroupCode]func(recursiveMonad) recursiveMonad{
	IR_GROUP_CODE_NO_CODE:       noCodeConsumer,
	IR_GROUP_CODE_SUBGROUP_STUB: subgroupConsumer,
	IR_GROUP_CODE_CLASS_STUB:    classConsumer,
	IR_GROUP_CODE_DOT:           dotConsumer,
	IR_GROUP_CODE_UNION:         unionConsumer,
}

func top(stack []*recursiveTree) *recursiveTree {
	length := len(stack)
	return stack[length-1]
}

func recursiveTreeFromSingle(ir Ir) RecursiveTreeRoot {
	monad := recursiveMonad{
		source: ir,
		stack:  make([]*recursiveTree, 0, 16),
		group:  IrGroup{},
	}
	monad.stack = append(monad.stack, &recursiveTree{
		children: make([]recursiveTree, 0, 32),
		and:      0,
		valid:    0,
		min:      1,
		max:      1,
		union:    false,
	})

	for _, v := range ir.Groups {
		consumer := recursiveTreeConsumers[v]
	}

	return RecursiveTreeRoot{
		tree:     (*monad.stack[0]),
		anchored: ir.Anchored,
		dollared: ir.Dollared,
	}
}

func noCodeConsumer(monad recursiveMonad) recursiveMonad {
	root := top(monad.stack)

	(*root).children = append(root.children, recursiveTree{
		children: []recursiveTree{},
		and:      rune(FULL_BITS),
		valid:    rune(monad.group.RuneOrIndex),
		min:      monad.group.Min,
		max:      monad.group.Max,
		union:    false,
	})

	return monad
}

// TODO
func subgroupConsumer(monad recursiveMonad) recursiveMonad {
	root := top(monad.stack)

	(*root).children = append(root.children, recursiveTree{
		children: []recursiveTree{},
		and:      rune(0),
		valid:    rune(0),
		min:      monad.group.Min,
		max:      monad.group.Max,
		union:    false,
	})

	return monad
}

// TODO
func classConsumer(monad recursiveMonad) recursiveMonad {
	root := top(monad.stack)

	(*root).children = append(root.children, recursiveTree{
		children: []recursiveTree{},
		and:      rune(0),
		valid:    rune(0),
		min:      monad.group.Min,
		max:      monad.group.Max,
		union:    false,
	})

	return monad
}

func dotConsumer(monad recursiveMonad) recursiveMonad {
	root := top(monad.stack)

	(*root).children = append(root.children, recursiveTree{
		children: []recursiveTree{},
		and:      rune(0),
		valid:    rune(0),
		min:      monad.group.Min,
		max:      monad.group.Max,
		union:    false,
	})

	return monad
}

// TODO
func unionConsumer(monad recursiveMonad) recursiveMonad {
	root := top(monad.stack)

	(*root).children = append(root.children, recursiveTree{
		children: []recursiveTree{},
		and:      rune(0),
		valid:    rune(0),
		min:      monad.group.Min,
		max:      monad.group.Max,
		union:    false,
	})

	return monad
}

func RecursiveTreeFromSingle(ir Ir) Regex {
	return recursiveTreeFromSingle(ir)
}

func RecursiveTreeFromGroup(irs []Ir) RegexGroup {
	length := len(irs)
	trees := make([]RecursiveTreeRoot, 0, length)

	for _, v := range irs {
		trees = append(trees, recursiveTreeFromSingle(v))
	}

	return RecursiveTreeGroup{
		trees: trees,
	}
}

func (r RecursiveTreeRoot) Matches([]byte) bool {
	return false
}

func (r RecursiveTreeRoot) MatchesUnicode(string) bool {
	return false
}

func (r RecursiveTreeGroup) First([]byte) int {
	return -1
}

func (r RecursiveTreeGroup) All([]byte) []int {
	return []int{}
}
