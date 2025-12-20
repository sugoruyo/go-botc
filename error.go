package botc

import "fmt"

type ConversionError[T any] struct {
	key   string
	value T
}

func (e *ConversionError[T]) Error() string {
	return fmt.Sprintf("error converting %s: %v", e.key, e.value)
}

func NewConversionError[T any](key string, value T) *ConversionError[T] {
	return &ConversionError[T]{
		key:   key,
		value: value,
	}
}

type IllegalValueForEnumError[T any] struct {
	key   string
	val   T
	valid []T
}

func (e *IllegalValueForEnumError[T]) Error() string {
	return fmt.Sprintf("value %v was invalid for %s, must be one of %v", e.val, e.key, e.valid)
}

func NewIllegalValueForEnumError[T any](key string, val T, valid []T) *IllegalValueForEnumError[T] {
	return &IllegalValueForEnumError[T]{
		key:   key,
		val:   val,
		valid: valid,
	}
}

type RequiredFieldMissingError struct {
	key string
}

func (e *RequiredFieldMissingError) Error() string {
	return fmt.Sprintf("required field %s missing", e.key)
}

func NewRequiredFieldMissingError(key string) *RequiredFieldMissingError {
	return &RequiredFieldMissingError{
		key: key,
	}
}
