package deepcopy

import . "reflect"

type visit struct {
	val uintptr
	typ Type
}

// Returns a deep copy of the value passed in, copying all elements and exported fields. Note that Chan and Func values are are copied as shallow copies. This result of this function may not compare as equal with reflect.DeepEqual, as Func comparison will return false if both are not nil. Pointer hierarchies are preserved.
func DeepCopy(i interface{}) interface{} {
	if i == nil {
		return nil
	}
	return DeepCopyValue(ValueOf(i)).Interface()
}

// Returns a deep copy of the data contained in the reflect.Value passed in, copying all elements and exported fields. Note that Chan and Func values are are copied as shallow copies. This result of this function may not compare as equal with reflect.DeepEqual, as Func comparison will return false if both are not nil. Pointer hierarchies are preserved.
func DeepCopyValue(val Value) Value {
	return deepCopyValue(val, make(map[visit]Value))
}

func deepCopyValue(val Value, visited map[visit]Value) Value {
	switch typ := val.Type(); typ.Kind() {

	// Just return everything that can be shallow copied
	case Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128, Chan, Func, String:
		return val

	// Deal with the types which can contain references
	case Array, Map, Ptr, Slice, Struct:
		// Allocate the return value
		new_val := New(typ).Elem()
		switch val.Kind() {
		case Slice:
			if !val.IsNil() {
				new_val.Set(MakeSlice(typ, val.Len(), val.Cap()))
			}
		case Map:
			if !val.IsNil() {
				new_val.Set(MakeMap(typ))
			}
		}

		// Store the new value for reference recreation
		if val.CanAddr() {
			v := visit{val.UnsafeAddr(), typ}
			if v, ok := visited[v]; ok {
				return v
			}
			visited[v] = new_val
		}

		// Copy the elements
		switch val.Kind() {
		case Array, Slice:
			for i := 0; i < val.Len(); i++ {
				new_val.Index(i).Set(deepCopyValue(val.Index(i), visited))
			}
		case Map:
			for _, k := range val.MapKeys() {
				new_val.SetMapIndex(k, deepCopyValue(val.MapIndex(k), visited))
			}
		case Ptr:
			if !val.IsNil() {
				new_val.Set(deepCopyValue(val.Elem(), visited).Addr())
			}
		case Struct:
			for i, n := 0, val.NumField(); i < n; i++ {
				if new_val.Field(i).CanSet() {
					new_val.Field(i).Set(deepCopyValue(val.Field(i), visited))
				}
			}
		}
		return new_val

	// Deal with the types which can contain references
	case Interface:
		new_val := New(typ).Elem()
		if !val.IsNil() {
			new_val.Set(deepCopyValue(val.Elem(), visited))
		}
		return new_val
	}

	panic("deepcopy: DeepCopyValue of invalid type")
}
