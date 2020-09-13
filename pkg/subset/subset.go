package subset

import (
	"reflect"
	"unsafe"
)

func SubsetEqual(subset, superset interface{}) bool {
	t := newTreeWalk()
	return t.subsetValueEqual(reflect.ValueOf(subset), reflect.ValueOf(superset))
}

type treeWalk struct {
	matches int
	visited map[node]bool
}

func newTreeWalk() *treeWalk {
	t := treeWalk{}
	t.visited = make(map[node]bool)
	return &t
}

// track visited nodes to short circuit recursion
type node struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Tests for subset equality using reflected types.
func (t *treeWalk) subsetValueEqual(subset, superset reflect.Value) bool {
	// if subset side is undefined, then nothing to compare
	if !subset.IsValid() {
		return true
	}

	// sanity check, rest of code assume same type on both sides
	if !superset.IsValid() {
		return false
	}
	if subset.Type() != superset.Type() {
		return false
	}

	// short circuit references already seen
	isRef := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}
	if subset.CanAddr() && superset.CanAddr() && isRef(subset.Kind()) {
		n := node{
			unsafe.Pointer(subset.UnsafeAddr()),
			unsafe.Pointer(superset.UnsafeAddr()),
			subset.Type(),
		}
		if t.visited[n] {
			return true
		}
		t.visited[n] = true
	}

	// walk tree
	switch subset.Kind() {
	case reflect.Array, reflect.Slice:
		// recursive subset: superset may have extra elements at the end, only the subset members must match
		if superset.Len() < subset.Len() {
			return false
		}
		for i := 0; i < subset.Len(); i++ {
			if !t.subsetValueEqual(subset.Index(i), superset.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Interface, reflect.Ptr:
		return t.subsetValueEqual(subset.Elem(), superset.Elem())
	case reflect.Struct:
		for i, n := 0, subset.NumField(); i < n; i++ {
			if !t.subsetValueEqual(subset.Field(i), superset.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		for _, k := range subset.MapKeys() {
			if !t.subsetValueEqual(subset.MapIndex(k), superset.MapIndex(k)) {
				return false
			}
		}
		return true
	default:
		// Leaf node: if exported, non-default value, compare subset with superset.
		if subset.CanInterface() {
			// Ignore default values, like empty string, bool false, zero int... :-|
			// If needed, could be more selective in the kinds of zeros ignored.
			if subset.Interface() == reflect.Zero(subset.Type()).Interface() {
				return true
			}
			if subset.Interface() == superset.Interface() {
				t.matches++
				return true
			}
			return false
		}
		return true // Ignore non-exported opaque internal fields, like time.Time{}.
	}
}
