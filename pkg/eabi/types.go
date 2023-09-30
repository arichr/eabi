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
	// TODO (arichr): This almost repeats the behaviour of marhsalInt64() and should be refactored.
	arrayLen := v.Len()
	var cursor byte = byte(arrayLen & 0b00000111)

	i := 4

	// 1100             ?         000
	// Data type  Finalizing bit  Data
	//
	// Check if we can fit the array length to 3 bits
	if arrayLen>>i&intChunkBitmask == 0 {
		*buf = append(*buf, cursor|0b1100_1_000)
	} else {
		*buf = append(*buf, cursor|0b1100_0_000)
	}

	for i < 64 {
		cursor = byte(arrayLen >> i & intChunkBitmask)
		if cursor == 0 {
			// Mark the last written byte as final
			(*buf)[len(*buf)-1] |= 0b10000000
			break
		}

		*buf = append(*buf, cursor)
		i += 7
	}

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
	// fmt.Printf("Encoding %v...\n", v) // FIXME (arichr): Remove

	if v.IsNil() {
		*buf = append(*buf, 0x08)
		// fmt.Printf("buf: %v\n", buf) // FIXME (arichr): Remove
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

	// fmt.Printf("buf: %v\n", buf) // FIXME (arichr): Remove
	return nil
}

const intChunkBitmask = 0b01111111

func marhsalInt64(v reflect.Value, buf *[]byte) {
	intval := v.Elem().Int()
	var cursor byte = byte(intval & 0b00000111)

	i := 3

	// 0100             ?         000
	// Data type  Finalizing bit  Data
	//
	// Check if we can fit the integer to 3 bits
	if intval>>i&intChunkBitmask == 0 {
		*buf = append(*buf, cursor|0b0100_1_000)
		return
	} else {
		*buf = append(*buf, cursor|0b0100_0_000)
	}

	for i < 64 {
		cursor = byte(intval >> i & intChunkBitmask)
		if cursor == 0 {
			// Mark the last written byte as final
			(*buf)[len(*buf)-1] |= 0b10000000
			break
		}

		*buf = append(*buf, cursor)
		i += 7
	}
}
