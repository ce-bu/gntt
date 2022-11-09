package gntt_optional

type Optional[T any] struct {
	value *T
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		value: &value,
	}
}

func (opt *Optional[T]) HasValue() bool {
	return opt.value != nil
}

func (opt *Optional[T]) Get() T {
	return *(opt.value)
}

func (opt *Optional[T]) Set(value T) {
	opt.value = &value
}
