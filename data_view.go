package puzzle

import (
	"encoding/binary"
	"math"
)

var order = binary.LittleEndian

// DataView provides an interface to read and write different primitives from
// a byte array in little endian order. You must ensure that the underlying
// byte array has a size large enough when writing to it.
type DataView []byte

// Uint8 reads a uint8 (single byte) from the given offset.
func (view DataView) Uint8(offset int) uint8 {
	return view[offset]
}

// Uint16 reads a uint16 from the given offset.
func (view DataView) Uint16(offset int) uint16 {
	return order.Uint16(view[offset : offset+2])
}

// Uint32 reads a uint32 from the given offset.
func (view DataView) Uint32(offset int) uint32 {
	return order.Uint32(view[offset : offset+4])
}

// Float32 reads a float32 from the given offset.
func (view DataView) Float32(offset int) float32 {
	bits := order.Uint32(view[offset : offset+4])
	return math.Float32frombits(bits)
}

// PutUint8 writes a uint8 (single byte) at the given offset.
func (view DataView) PutUint8(val uint8, offset int) {
	view[offset] = val
}

// PutUint16 writes a uint16 at the given offset.
func (view DataView) PutUint16(val uint16, offset int) {
	order.PutUint16(view[offset:offset+2], val)
}

// PutUint32 writes a uint32 at the given offset.
func (view DataView) PutUint32(val uint32, offset int) {
	order.PutUint32(view[offset:offset+4], val)
}

// PutFloat32 writes a float32 at the given offset.
func (view DataView) PutFloat32(val float32, offset int) {
	bits := math.Float32bits(val)
	order.PutUint32(view[offset:offset+4], bits)
}

// ReadString reads a length-prefixed string starting at the given offset.
// Returns the string read and the number of bytes read in total (length of
// string plus 2 bytes for length prefix).
func (view DataView) ReadString(offset int) (string, int) {
	length := view.Uint16(offset)
	start := uint16(offset + 2)
	strBytes := view[start : start+length]
	return string(strBytes), int(2 + length)
}

// PutRawBytes copies the bytes from data at the given offset.
func (view DataView) PutRawBytes(data []byte, offset int) {
	copy(view[offset:], data)
}
