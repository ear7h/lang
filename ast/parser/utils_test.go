package parser_test

import (
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ear7h/lang/ast/parser"
)

func init() {
	defaultFi := parser.NewCursorString("", "").FileInfo()

	if reflect.DeepEqual(defaultFi, parser.FileInfo{}) {
		// the default file info should not be the zero
		// value. Firstly, it should be start on line 1
		// col 1. Secondly, a non-zero value as the
		// initial cursor FileInfo ensures that Parse
		// is properly initalizing the file info
		panic("default file info is zero value")
	}
}

func assertEq(t *testing.T, expect, got interface{}) {
	t.Helper()

	if expect ==  nil || got == nil {
		if expect != got {
			t.Fatalf("expected: %v (%[1]T)\ngot: %[2]v (%[2]T)", expect, got)
		}

		return
	}

	av := reflect.ValueOf(expect)
	bv := reflect.ValueOf(got)

	av.Type()
	bv.Type()

	if av.Type() != bv.Type() {
		t.Fatalf("expected: %v (%[1]T)\ngot: %[2]v (%[2]T)", expect, got)
	}

	if !astDeepValueEqual(av, bv, make(map[visit]bool), 0) {
		t.Fatalf("expected: %#v (%[1]T)\ngot: %#[2]v (%[2]T)", expect, got)
	}
}

func assertErrIs(t *testing.T, expect, got error) {
	t.Helper()

	if !errors.Is(expect, got) {
		t.Fatalf("expected: %v\ngot: %v", expect, got)
	}
}

// the following was mostly taken from then Go
// source tree, commit 872bbc

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// astDeepValueEqual works like reflect.DeepEqual, but with
func astDeepValueEqual(v1, v2 reflect.Value,
	visited map[visit]bool, depth int) bool {

	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	hard := func(v1, v2 reflect.Value) bool {
		switch v1.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			// Nil pointers cannot be cyclic. Avoid putting them in the visited map.
			return !v1.IsNil() && !v2.IsNil()
		}
		return false
	}

	if hard(v1, v2) {
		ptrval := func(v reflect.Value) unsafe.Pointer {
			switch v1.Kind() {
			case reflect.Interface:
				// internally, the reflect package
				// uses Value.ptr to get the pointer out
				// of an iface, but it's not exported
				// so we hack it here
				type iface struct {
					tab  unsafe.Pointer
					data unsafe.Pointer
				}

				ifacev := v.Interface()
				return (*iface)(unsafe.Pointer(&ifacev)).data
			default:
				return unsafe.Pointer(v.Pointer())
			}
		}
		addr1 := ptrval(v1)
		addr2 := ptrval(v2)
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !astDeepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}

		return true

	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !astDeepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true

	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return astDeepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)

	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return astDeepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)

	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {

			// ear7h modification, skip file info
			// in BaseNode. In the test suite the ast nodes
			// are better created with existing functions
			// rather than struct literals, ex:
			/*
				out: &ast.UnaryExpr{
					Op: '+',
					Operand: ast.MustParseString(
						&ast.NumberLiteral{},
						"123",
					),
				},
			*/
			if v1.Type().Name() == "BaseNode" &&
				v1.Type().Field(i).Name == "Fi" {
				continue
			}

			if !astDeepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1) {
				return false
			}
		}
		return true

	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !astDeepValueEqual(val1, val2, visited, depth+1) {
				return false
			}
		}
		return true

	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false

	default:
		// Normal equality suffices
		return v1.CanInterface() && v1.Interface() == v2.Interface()
	}
}
