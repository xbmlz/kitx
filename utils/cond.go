package utils

func Or[T any](funcs ...func(T) bool) func(T) bool {
	return func(input T) bool {
		for _, f := range funcs {
			if f(input) {
				return true
			}
		}
		return false
	}
}

// Returns whenTrue if b is true; otherwise, returns whenFalse. IfElse should only be used when branches are either
// constant or precomputed as both branches will be evaluated regardless as to the value of b.
func IfElse[T any](b bool, whenTrue T, whenFalse T) T {
	if b {
		return whenTrue
	}
	return whenFalse
}

// Returns value if value is not the zero value of T; Otherwise, returns defaultValue. OrElse should only be used when
// defaultValue is constant or precomputed as its argument will be evaluated regardless as to the content of value.
func OrElse[T comparable](value T, defaultValue T) T {
	if value != *new(T) {
		return value
	}
	return defaultValue
}

// Returns `a` if `a` is not `nil`; Otherwise, returns `b`. Coalesce is roughly analogous to `??` in JS, except that it
// non-shortcutting, so it is advised to only use a constant or precomputed value for `b`
func Coalesce[T *U, U any](a T, b T) T {
	if a == nil {
		return b
	} else {
		return a
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Memoize[T any](create func() T) func() T {
	var value T
	return func() T {
		if create != nil {
			value = create()
			create = nil
		}
		return value
	}
}
