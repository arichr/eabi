package eabi

// 	0: I8
// 	1: U8
// 	2: I16
// 	3: U16
// 	4: I64
// 	5: U64
// 	6: F64
// 	7: Bytes
// 	8: Null
// 	10: Map
// 	12: Array

import (
	"fmt"
	"reflect"
)

type Marshaler interface {
	MarshalEabi() ([]byte, error)
}

func marshalArray(v reflect.Value, buf *[]byte) error {
	*buf = append(*buf, 0x0C) // TODO: Specify array length

	for i := 0; i < v.Len(); i++ {
		if err := marshalToBuffer(v.Index(i), buf); err != nil {
			return err
		}
	}
	return nil
}

func marshalMap(v reflect.Value, buf *[]byte) error {
	panic("Not implemented: marshalMap")
}

func marshalElement(v reflect.Value, buf *[]byte) error {
	fmt.Printf("Encoding %v...\n", v) // FIXME: Remove

	if v.IsNil() {
		*buf = append(*buf, 0x08)
		fmt.Printf("buf: %v\n", buf) // FIXME: Remove
		return nil
	}

	kind := v.Elem().Kind()
	switch kind {
	case reflect.Int, reflect.Int64:
		marhsalInt64(v, buf)
	case reflect.Uint, reflect.Uint64:
		panic("Not implemented: marshalElement: reflect.Uint64")
	case reflect.Int16:
		panic("Not implemented: marshalElement: reflect.Int16")
	case reflect.Uint16:
		panic("Not implemented: marshalElement: reflect.Uint16")
	case reflect.Int8:
		panic("Not implemented: marshalElement: reflect.Int8")
	case reflect.Uint8:
		panic("Not implemented: marshalElement: reflect.Uint8")
	case reflect.Float64:
		panic("Not implemented: marshalElement: reflect.Float64")
	case reflect.Array, reflect.Slice, reflect.String:
		if kind == reflect.String {
			panic("Not implemented: marshalElement: reflect.String")
		}
		if v.Elem().Index(0).Kind() == reflect.Uint8 {
			panic("Not implemented: marshalElement: [?]byte")
		}
		panic("Not implemented: marshalElement: [?]any")
	default:
		marshalerType := reflect.TypeOf((*Marshaler)(nil)).Elem()
		if v.Elem().Type().Implements(marshalerType) {
			panic("Not implemented: marshalElement: Marshaler support")
		}
		return fmt.Errorf("unknown type: %v", kind)
	}

	fmt.Printf("buf: %v\n", buf) // FIXME: Remove
	return nil
}

// Splits the binary representation of int64 into 7-bit chunks.
// Each chunk is preceded by the "finalising bit".
//
// Think of it as a simple VarInt implementation.
func marhsalInt64(v reflect.Value, buf *[]byte) {
	var cursor byte
	*buf = append(*buf, 0x04)
	intval := v.Elem().Int()

	for i := 0; i < 64; i += 7 {
		cursor = byte(intval >> i & 0b01111111)
		if cursor == 0 {
			(*buf)[len(*buf)-1] |= 0b10000000
			break
		}
		*buf = append(*buf, cursor)
	}
}
