package runtime

import (
	"unsafe"

	"github.com/arnodel/golua/runtime/internal/luagc"
)

// Table implements a Lua table.
type Table struct {
	// This is where the implementation details are.
	*mixedTable

	meta *Table
}

// NewTable returns a new Table.
func NewTable() *Table {
	return &Table{mixedTable: &mixedTable{}}
}

// Metatable returns the table's metatable.
func (t *Table) Metatable() *Table {
	return t.meta
}

// SetMetatable sets the table's metatable.
func (t *Table) SetMetatable(m *Table) {
	t.meta = m
}

var _ luagc.Value = (*Table)(nil)

func (t *Table) Key() luagc.Key {
	return unsafe.Pointer(t.mixedTable)
}

func (t *Table) Clone() luagc.Value {
	clone := new(Table)
	*clone = *t
	return clone
}

// Get returns t[k].
func (t *Table) Get(k Value) Value {
	return t.get(k)
}

// Set implements t[k] = v (doesn't check if k is nil).
func (t *Table) Set(k, v Value) uint64 {
	if v.IsNil() {
		t.mixedTable.remove(k)
		return 0
	}
	t.mixedTable.insert(k, v)
	return 16
}

// Reset implements t[k] = v only if t[k] was already non-nil.
func (t *Table) Reset(k, v Value) (wasSet bool) {
	if v.IsNil() {
		return t.mixedTable.remove(k)
	}
	return t.mixedTable.reset(k, v)
}

// Len returns a length for t (see lua docs for details).
func (t *Table) Len() int64 {
	return int64(t.mixedTable.len())
}

// Next returns the key-value pair that comes after k in the table t.
//   - If k is NilValue, the first key-value pair in the table t is returned.
//   - If k is the last key in the table t, a pair of NilValues is returned.
//   - If the table t is empty, the returned key-value pair is always a pair of NilValues, regardless of k.
//   - In all cases, ok is true if and only if k is either NilValue or a key present in the table t.
func (t *Table) Next(k Value) (next Value, val Value, ok bool) {
	return t.mixedTable.next(k)
}
