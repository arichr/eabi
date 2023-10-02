package eabi

// Types:
// 0    (0b0000): Null
// 1    (0b0001): Integer
// 2    (0b0010): Float
// 3    (0b0011): String
// 4    (0b0100): Array
// 5    (0b0101): Map
// 6    (0b0110): Typedef
// 7    (0b0111): (unused)
//
// 8-15 (0b1xxx): User-defined types

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

	//     4            ?         000
	// Data type  Finalizing bit  Data
	//
	// Check if we can fit the array length to 3 bits
	if arrayLen>>i&intChunkBitmask == 0 {
		*buf = append(*buf, cursor|0b0100_1_000)
	} else {
		*buf = append(*buf, cursor|0b0100_0_000)
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
		*buf = append(*buf, 0)
		// fmt.Printf("buf: %v\n", buf) // FIXME (arichr): Remove
		return nil
	}

	kind := v.Elem().Kind()
	switch kind {
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		marhsalUint(v.Elem().Uint(), buf)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		// TODO (arichr): Minimize the size of encoded negative values
		marhsalUint(uint64(v.Elem().Int()), buf)
	case reflect.Float64:
		panic("Not implemented: marshalElement: reflect.Float64")
	case reflect.String:
		panic("Not implemented: marshalElement: reflect.String")
	case reflect.Array, reflect.Slice:
		if v.Elem().Index(0).Kind() == reflect.Uint8 {
			panic("Not implemented: marshalElement: [?]byte (aka string)")
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

func marhsalUint(value uint64, buf *[]byte) {
	i := 3
	var cursor byte = byte(value & 0b00000111)

	//     1            ?         000
	// Data type  Finalizing bit  Data
	//
	// Check if we can fit the integer to 3 bits
	if value>>i&intChunkBitmask == 0 {
		*buf = append(*buf, cursor|0b0001_1_000)
		return
	} else {
		*buf = append(*buf, cursor|0b0001_0_000)
	}

	for i < 64 {
		cursor = byte(value >> i & intChunkBitmask)
		if cursor == 0 {
			// Mark the last written byte as final
			(*buf)[len(*buf)-1] |= 0b10000000
			break
		}

		*buf = append(*buf, cursor)
		i += 7
	}
}
