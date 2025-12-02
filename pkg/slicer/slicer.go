package slicer

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func AddIfNotFound[S []E, E any](s S, f func(E) bool, e E) S {
	result := s
	if !ContainsF(s, f) {
		result = append(result, e)
	}
	return result
}

// Group groups the elements of a collection by a specified key.
func Group[S ~[]E, E any, K comparable](s S, f func(E) K) map[K]S {
	result := make(map[K]S)

	for _, v := range s {
		k := f(v)
		result[k] = append(result[k], v)
	}

	return result
}

// ContainsF checks if the given function returns true for any element in the slice.
func ContainsF[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func ContainsStr(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Includes[S ~[]E, E any](s S, e E) bool {
	for _, v := range s {
		if reflect.DeepEqual(v, e) {
			return true
		}
	}
	return false
}

func Map[S ~[]E, E any, T any](s S, f func(E) T) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}
	return result
}

func Each[S ~[]E, E any](s S, f func(E)) {
	for _, v := range s {
		f(v)
	}
}

func Filter[S ~[]E, E any](s S, f func(E) bool) []E {
	result := make([]E, 0, len(s))
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func First[S ~[]E, E any](s S) E {
	if len(s) == 0 {
		return *new(E)
	}
	return s[0]
}

func Last[S ~[]E, E any](s S) E {
	if len(s) == 0 {
		return *new(E)
	}
	return s[len(s)-1]
}

func Find[S ~[]E, E any](s S, f func(E) bool) (int, E) {
	for i, v := range s {
		if f(v) {
			return i, v
		}
	}
	return -1, *new(E)
}

func FindAndDo[S ~[]E, E any](s S, f func(E) bool, do func(int, E)) {
	for i, v := range s {
		if f(v) {
			do(i, v)
		}
	}
}

func FindMaxByTimeF[S ~[]E, E comparable](s S, f func(E) time.Time) E {
	if len(s) == 0 {
		return *new(E)
	}

	result := s[0]
	for _, v := range s {
		if f(v).After(f(result)) {
			result = v
		}
	}
	return result
}

func FindMinByTimeF[S ~[]E, E comparable](s S, f func(E) time.Time) E {
	if len(s) == 0 {
		return *new(E)
	}

	result := s[0]
	for _, v := range s {
		if f(v).Before(f(result)) {
			result = v
		}
	}
	return result
}

func Merge[S ~[]E, S2 ~[]E, E any](s1 S, s2 S2) []E {
	result := make([]E, 0, len(s1)+len(s2))
	result = append(result, s1...)
	result = append(result, s2...)
	return result
}

func MergeF[S ~[]E, S2 ~[]E2, E any, E2 any](s1 S, s2 S2, f func(E2) E) []E {
	result := make([]E, 0, len(s1)+len(s2))
	for _, v := range s1 {
		result = append(result, v)
	}
	for _, v := range s2 {
		result = append(result, f(v))
	}
	return result
}

func Unique[S ~[]E, E comparable](s S) []E {
	result := make([]E, 0, len(s))
	m := make(map[E]struct{})
	for _, v := range s {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func UniqueF[S ~[]E, E any, K comparable](s S, f func(E) K) []E {
	result := make([]E, 0, len(s))
	m := make(map[K]struct{})
	for _, v := range s {
		k := f(v)
		if _, ok := m[k]; !ok {
			m[k] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func CleanEmptyStrings(s []string) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func Order[S ~[]E, E any](s S, f func(E, E) bool) {
	sort.Slice(s, func(i, j int) bool {
		return f(s[i], s[j])
	})
}

func Reduce[S ~[]E, E any, V any](s S, f func(V, E) V, initialValue V) V {
	result := initialValue
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

func SumF[S ~[]E, E any, T constraints.Signed | constraints.Unsigned | ~float32 | ~float64](s S, f func(E) T) T {
	result := T(0)
	for _, v := range s {
		result += f(v)
	}
	return result
}

func Increment[S ~[]E, E any](s S, f func(E) bool) int {
	result := 0
	for _, v := range s {
		if f(v) {
			result++
		}
	}
	return result
}

func Pluck[S ~[]E, E any, K comparable](s S, f func(E) K) []K {
	result := make([]K, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}
	return result
}

func PluckByField[S ~[]E, E any](s S, fieldName string) []interface{} {
	result := make([]interface{}, 0, len(s))
	for _, v := range s {
		value, ok := GetFieldValue(v, fieldName)
		if ok {
			result = append(result, value)
		}
	}
	return result
}

func AToStrings[S ~[]E, E any](s S) []string {
	result := make([]string, 0, len(s))
	for _, v := range s {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return result
}

func AToInterfaces[S ~[]E, E any](s S) []interface{} {
	result := make([]interface{}, 0, len(s))
	for _, v := range s {
		result = append(result, v)
	}
	return result
}

func LastIndexOf[E any](s []E) int {
	return len(s) - 1
}

func IsLastIndex[E any](i int, s []E) bool {
	return i == LastIndexOf(s)
}

func Join[E any](s []E, sep string) string {
	strs := []string{}

	for _, v := range s {
		strs = append(strs, fmt.Sprintf("%v", v))
	}

	return strings.Join(strs, sep)
}

// GetFieldValue gets the value of a struct field by its name
func GetFieldValue(s interface{}, fieldName string) (interface{}, bool) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if t := v.Type(); t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == fieldName {
				return v.Field(i).Interface(), true
			}
		}
	}
	return nil, false
}

// GetFieldValueRecursive retrieves the value of a struct field, potentially nested within multiple levels of structs.
// It handles pointers and returns the found value or nil if not found.
func GetFieldValueRecursive(s interface{}, fieldNames []string) (interface{}, bool) {
	// Each through each field name in the provided path
	tmp := s
	for _, fieldName := range fieldNames {
		if v, ok := GetFieldValue(tmp, fieldName); ok {
			tmp = v
		} else {
			return nil, false
		}
	}

	// Unreachable code if the path is properly structured
	return tmp, true
}
