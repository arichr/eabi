package eabi

import (
	"fmt"
	"reflect"
)

func marshalToBuffer(v reflect.Value, buf *[]byte) error {
	if !v.IsValid() {
		return fmt.Errorf("invalid value of %T", v)
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return marshalArray(v, buf)
	case reflect.Map:
		return marshalMap(v, buf)
	default:
		return marshalElement(v, buf)
	}
}

func marshalWithSize(v reflect.Value, bufSize int) ([]byte, error) {
	buf := make([]byte, 0, bufSize)

	if err := marshalToBuffer(v, &buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func Marshal(v any) ([]byte, error) {
	buf, err := marshalWithSize(reflect.ValueOf(v), 1024)
	if err != nil {
		return nil, fmt.Errorf("Marshal: %s", err)
	}
	return buf, nil
}

func MarshalToBuffer(v any, buf *[]byte) error {
	if err := marshalToBuffer(reflect.ValueOf(v), buf); err != nil {
		return fmt.Errorf("MarshalToBuffer: %s", err)
	}
	return nil
}
