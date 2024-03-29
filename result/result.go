// Package result implements the Result type.
package result

import (
	"github.com/onur1/data"
)

// Ok creates a result which never fails and returns a value of type A.
func Ok[A any](a A) data.Result[A] {
	return func() (A, error) {
		return a, nil
	}
}

// Error creates a result which always fails with an error.
func Error[A any](err error) data.Result[A] {
	return func() (a A, _ error) {
		return a, err
	}
}

// Zero creates a result which never fails and returns a zero value
// of the type that it is initialized with.
func Zero[A any]() (a A, _ error) {
	return
}

// Map creates a result by applying a function on a succeeding result.
func Map[A, B any](fa data.Result[A], f func(A) B) data.Result[B] {
	if a, err := fa(); err != nil {
		return Error[B](err)
	} else {
		return Ok(f(a))
	}
}

// MapError creates a result by applying a function on a failing result.
func MapError[A any](fa data.Result[A], f func(error) error) data.Result[A] {
	if _, err := fa(); err != nil {
		return Error[A](f(err))
	}
	return fa
}

// Ap creates a result by applying a function contained in the first result
// on the value contained in the second result.
func Ap[A, B any](fab data.Result[func(A) B], fa data.Result[A]) data.Result[B] {
	var (
		err error
		ab  func(A) B
	)

	if ab, err = fab(); err != nil {
		return Error[B](err)
	}

	var a A

	if a, err = fa(); err != nil {
		return Error[B](err)
	}

	return Ok(ab(a))
}

// Chain creates a result which combines two results in sequence, using the
// return value of one result to determine the next one.
func Chain[A, B any](ma data.Result[A], f func(A) data.Result[B]) data.Result[B] {
	if a, err := ma(); err != nil {
		return Error[B](err)
	} else {
		return f(a)
	}
}

// ChainFirst composes two results in sequence, using the return value of one result
// to determine the next one, keeping only the first result.
func ChainFirst[A, B any](ma data.Result[A], f func(A) data.Result[B]) data.Result[A] {
	return Chain(ma, func(a A) data.Result[A] {
		return Map(f(a), fst[A, B](a))
	})
}

// Bimap creates a result by mapping a pair of functions over an error or a value
// contained in a result.
func Bimap[A, B any](fa data.Result[A], f func(error) error, g func(A) B) data.Result[B] {
	if a, err := fa(); err != nil {
		return Error[B](f(err))
	} else {
		return Ok(g(a))
	}
}

// ApFirst creates a result by combining two effectful computations, keeping
// only the result of the first.
func ApFirst[A, B any](fa data.Result[A], fb data.Result[B]) data.Result[A] {
	return Ap(Map(fa, fst[A, B]), fb)
}

// ApSecond creates a result by combining two effectful computations, keeping
// only the result of the second.
func ApSecond[A, B any](fa data.Result[A], fb data.Result[B]) data.Result[B] {
	return Ap(Map(fa, snd[A, B]), fb)
}

// Fold takes two functions and a result and returns a value by applying
// one of the supplied functions to the inner value.
func Fold[A, B any](ma data.Result[A], onError func(error) B, onSuccess func(A) B) B {
	if a, err := ma(); err != nil {
		return onError(err)
	} else {
		return onSuccess(a)
	}
}

// GetOrElse creates a result which can be used to recover from a failing result
// with a new value.
func GetOrElse[A any](ma data.Result[A], onError func(error) A) A {
	if a, err := ma(); err != nil {
		return onError(err)
	} else {
		return a
	}
}

// OrElse creates a result which can be used to recover from a failing result
// by switching to a new result.
func OrElse[A any](ma data.Result[A], onError func(error) data.Result[A]) data.Result[A] {
	if _, err := ma(); err != nil {
		return onError(err)
	}
	return ma
}

// FilterOrElse creates a result which can be used to fail with an error unless
// a predicate holds on a succeeding result.
func FilterOrElse[A any](ma data.Result[A], predicate data.Predicate[A], onFalse func(A) error) data.Result[A] {
	return Chain(ma, func(a A) data.Result[A] {
		if predicate(a) {
			return Ok(a)
		} else {
			return Error[A](onFalse(a))
		}
	})
}

// Fork is like Fold but it doesn't have a return value.
func Fork[A any](ma data.Result[A], onError func(error), onSuccess func(A)) {
	if a, err := ma(); err != nil {
		onError(err)
	} else {
		onSuccess(a)
	}
}

// FromNilable creates a result from a nilable, returning the supplied error
// for nil values.
func FromNilable[A any](ma data.Nilable[A], onNil func() error) data.Result[A] {
	if ma == nil {
		return Error[A](onNil())
	}
	return Ok(*ma)
}

func fst[A, B any](a A) func(B) A {
	return func(B) A {
		return a
	}
}

func snd[A, B any](A) func(B) B {
	return func(b B) B {
		return b
	}
}
