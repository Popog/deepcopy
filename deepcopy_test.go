package deepcopy

import (
	. "reflect"
	"testing"
)

type Basic struct {
	X int
	Y float32
}

type NotBasic Basic

type DeepEqualTest struct {
	a  interface{}
	eq bool
}

type Node *Node

func Addr(n Node, depth int) Node {
	for i := 0; i < depth; i++ {
		n2 := n
		n = &n2
	}
	return n
}

func DeepCompare(a, b Node) bool {
	for a != nil && b != nil {
		if a == b {
			return true
		}
		a, b = *a, *b
	}
	return false
}

// Simple functions for DeepEqual tests.
var (
	fn1 func() // nil.
)

var (
	node1 Node = new(Node)
	node2 Node = Addr(node1, 1)
	node3 Node = Addr(node1, 4)
	node4 Node = Addr(node1, 4)
)

func init() {
	*node1 = node3
}

var deepEqualTests = []interface{}{
	// Equalities
	nil,
	1,
	int32(1),
	0.5,
	float32(0.5),
	"hello",
	make([]int, 10),
	&[3]int{1, 2, 3},
	Basic{1, 0.5},
	error(nil),
	map[int]string{1: "one", 2: "two"},
	fn1,

	// Nil vs empty: not the same.
	[]int{},
	[]int(nil),
	map[int]int{},
	map[int]int(nil),

	// Mismatched types
	[]int{1, 2, 3},
	[3]int{1, 2, 3},
	&[3]interface{}{1, 2, 4},
	&[3]interface{}{1, 2, "s"},
	NotBasic{1, 0.5},
	map[uint]string{1: "one", 2: "two"},

	// Pointers
	new(int),
	new(float32),
	new(map[uint]string),
	new([]map[uint]string),
	new(interface{}),
	node1,
	node2,
	node3,
	node4,
}

func TestDeepCopy(t *testing.T) {
	for _, test := range deepEqualTests {
		//t.Logf("DeepCopy(%v)", test)
		if r := DeepCopy(test); !DeepEqual(test, r) {
			t.Errorf("DeepCopy(%v) = %v %v", test, r, TypeOf(r))
		}
	}
}

func BenchmarkDeepCopyNode(b *testing.B) {
	b.StopTimer()
	v := Addr(nil, 100)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(v)
	}
}
