package pool

import (
	"fmt"
	"io"

	"strconv"

	"github.com/derkan/nlog"
)

// Buffer provides byte buffer, which can be used for minimizing
// memory allocations.
//
// Buffer may be used with functions appending data to the given []byte
// slice. See example code for details.
//
// Use Get for obtaining an empty byte buffer.
type Buffer struct {

	// B is a byte buffer to use in append-like workloads.
	// See example code for details.
	B []byte
}

// Len returns the size of the byte buffer.
func (b *Buffer) Len() int {
	return len(b.B)
}

// Cap returns the capacity of the byte buffer.
func (b *Buffer) Cap() int {
	return cap(b.B)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

// WriteTo implements io.WriterTo.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

// Bytes returns b.B, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b *Buffer) Bytes() []byte {
	return b.B
}

// Write implements io.Writer - it appends p to Buffer.B
func (b *Buffer) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// Set sets Buffer.B to p.
func (b *Buffer) Set(p []byte) {
	b.B = append(b.B[:0], p...)
}

// SetString sets Buffer.B to s.
func (b *Buffer) SetString(s string) {
	b.B = append(b.B[:0], s...)
}

// String returns string representation of Buffer.B.
func (b *Buffer) String() string {
	return string(b.B)
}

// Reset makes Buffer.B empty.
func (b *Buffer) Reset() {
	b.B = b.B[:0]
}

// Itoa is Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func (b *Buffer) Itoa(i int, wid int) {
	// Assemble decimal in reverse order.
	var buf [20]byte
	bp := len(buf) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		buf[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	buf[bp] = byte('0' + i)
	b.AppendBytes(buf[bp:])
}

// AppendByte writes a single byte to the Buffer.
func (b *Buffer) AppendByte(v byte) nlog.Buffer {
	b.B = append(b.B, v)
	return b
}

// AppendBytes writes a slice of byte to the Buffer.
func (b *Buffer) AppendBytes(v []byte) nlog.Buffer {
	b.B = append(b.B, v...)
	return b
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(v string, quota bool) nlog.Buffer {
	if quota {
		b.B = append(b.B, '"')
	}
	b.B = append(b.B, v...)
	if quota {
		b.B = append(b.B, '"')
	}
	return b
}

// AppendStrings writes a slice of string to the Buffer.
func (b *Buffer) AppendStrings(v []string, quota bool) nlog.Buffer {
	if len(v) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}

	b.B = append(b.B, '[')
	b.AppendString(v[0], quota)
	if len(v) > 1 {
		for _, v := range v[1:] {
			b.B = append(b.B, ',')
			b.AppendString(v, quota)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInt appends an integer to the underlying buffer (assuming base 10).
func (b *Buffer) AppendInt(val int) nlog.Buffer {
	b.B = strconv.AppendInt(b.B, int64(val), 10)
	return b
}

// AppendInts appends a slice of integer to the underlying buffer
func (b *Buffer) AppendInts(val []int) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendInt(b.B, int64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendInt(append(b.B, ','), int64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInts8 appends a slice of int8 to the underlying buffer
func (b *Buffer) AppendInts8(val []int8) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendInt(b.B, int64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendInt(append(b.B, ','), int64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInts16 appends a slice of int16 to the underlying buffer
func (b *Buffer) AppendInts16(val []int16) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendInt(b.B, int64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendInt(append(b.B, ','), int64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInts32 appends a slice of int32 to the underlying buffer
func (b *Buffer) AppendInts32(val []int32) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendInt(b.B, int64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendInt(append(b.B, ','), int64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInt64 appends an int64 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendInt64(i int64) nlog.Buffer {
	b.B = strconv.AppendInt(b.B, i, 10)
	return b
}

// AppendInts64 appends a slice of int64 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendInts64(val []int64) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendInt(b.B, val[0], 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendInt(append(b.B, ','), v, 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendUInts appends a slice of uint to the underlying buffer (assuming base 10).
func (b *Buffer) AppendUInts(val []uint) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendUint(b.B, uint64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendUint(append(b.B, ','), uint64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendUInts16 appends a slice of uint16 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendUInts16(val []uint16) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendUint(b.B, uint64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendUint(append(b.B, ','), uint64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendUInts32 appends a slice of uint32 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendUInts32(val []uint32) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendUint(b.B, uint64(val[0]), 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendUint(append(b.B, ','), uint64(v), 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendUInt64 appends an uint64 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendUInt64(val uint64) nlog.Buffer {
	b.B = strconv.AppendUint(b.B, uint64(val), 10)
	return b
}

// AppendUInts64 appends a slice of uint64 to the underlying buffer (assuming base 10).
func (b *Buffer) AppendUInts64(val []uint64) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendUint(b.B, val[0], 10)
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = strconv.AppendUint(append(b.B, ','), v, 10)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendBool appends a bool to the underlying buffer.
func (b *Buffer) AppendBool(v bool) nlog.Buffer {
	b.B = strconv.AppendBool(b.B, v)
	return b
}

// AppendBools appends a slice bool to the underlying buffer.
func (b *Buffer) AppendBools(v []bool) nlog.Buffer {
	for _, x := range v {
		b.B = strconv.AppendBool(append(b.B, ','), x)
	}
	return b
}

// AppendError writes error to buffer
func (b *Buffer) AppendError(val error, quota bool) nlog.Buffer {
	return b.AppendString(val.Error(), quota)
}

// AppendErrors writes a slice of string to the Buffer.
func (b *Buffer) AppendErrors(val []error, quota bool) nlog.Buffer {
	if len(val) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	if quota {
		b.B = append(b.B, '"')
	}
	b.AppendError(val[0], quota)
	if quota {
		b.B = append(b.B, '"')
	}
	if len(val) > 1 {
		for _, v := range val[1:] {
			b.B = append(b.B, ',')
			if quota {
				b.B = append(b.B, '"')
			}
			b.AppendError(v, quota)
			if quota {
				b.B = append(b.B, '"')
			}
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendFloat32 appends an float32 to the underlying buffer
func (b *Buffer) AppendFloat32(val float32) nlog.Buffer {
	b.B = strconv.AppendFloat(b.B, float64(val), 'f', -1, 32)
	return b
}

// AppendFloats32 appends a slice of float32 to the underlying buffer
func (b *Buffer) AppendFloats32(vals []float32) nlog.Buffer {
	if len(vals) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendFloat(b.B, float64(vals[0]), 'f', -1, 32)
	if len(vals) > 1 {
		for _, v := range vals[1:] {
			b.B = strconv.AppendFloat(append(b.B, ','), float64(v), 'f', -1, 32)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendFloat64 appends an float64 to the underlying buffer
func (b *Buffer) AppendFloat64(val float64) nlog.Buffer {
	b.B = strconv.AppendFloat(b.B, val, 'f', -1, 64)
	return b
}

// AppendFloats64 appends a slice of float64 to the underlying buffer
func (b *Buffer) AppendFloats64(vals []float64) nlog.Buffer {
	if len(vals) == 0 {
		b.B = append(b.B, '[', ']')
		return b
	}
	b.B = append(b.B, '[')
	b.B = strconv.AppendFloat(b.B, vals[0], 'f', -1, 64)
	if len(vals) > 1 {
		for _, v := range vals[1:] {
			b.B = strconv.AppendFloat(append(b.B, ','), v, 'f', -1, 64)
		}
	}
	b.B = append(b.B, ']')
	return b
}

// AppendInterface takes an arbitrary object and converts it to JSON and embeds it dst.
func (b *Buffer) AppendInterface(val interface{}, marshallFn nlog.MarshallFn) nlog.Buffer {
	res, err := marshallFn(val)
	if err != nil {
		return b.AppendString(fmt.Sprintf("marshaling error: %v", err), true)
	}
	return b.AppendBytes(res)
}

// AppendAny appends given interface with its type
func (b *Buffer) AppendAny(val interface{}, quota bool, marshallFn nlog.MarshallFn) nlog.Buffer {
	switch val := val.(type) {
	case string:
		b.AppendString(val, quota)
	case []byte:
		b.AppendBytes(val)
	case error:
		b.AppendError(val, quota)
	case []error:
		b.AppendErrors(val, quota)
	case bool:
		b.AppendBool(val)
	case int:
		b.AppendInt(val)
	case int8:
		b.AppendInt64(int64(val))
	case int16:
		b.AppendInt64(int64(val))
	case int32:
		b.AppendInt64(int64(val))
	case int64:
		b.AppendInt64(val)
	case uint:
		b.AppendUInt64(uint64(val))
	case uint8:
		b.AppendUInt64(uint64(val))
	case uint16:
		b.AppendUInt64(uint64(val))
	case uint32:
		b.AppendUInt64(uint64(val))
	case uint64:
		b.AppendUInt64(val)
	case float32:
		b.AppendFloat32(val)
	case float64:
		b.AppendFloat64(val)
	case *string:
		if val != nil {
			b.AppendString(*val, quota)
		} else {
			b.AppendString("null", false)
		}
	case *bool:
		if val != nil {
			b.AppendBool(*val)
		} else {
			b.AppendString("null", false)
		}
	case *int:
		if val != nil {
			b.AppendInt(*val)
		} else {
			b.AppendString("null", false)
		}
	case *int8:
		if val != nil {
			b.AppendInt64(int64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *int16:
		if val != nil {
			b.AppendInt64(int64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *int32:
		if val != nil {
			b.AppendInt64(int64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *int64:
		if val != nil {
			b.AppendInt64(*val)
		} else {
			b.AppendString("null", false)
		}
	case *uint:
		if val != nil {
			b.AppendUInt64(uint64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *uint8:
		if val != nil {
			b.AppendUInt64(uint64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *uint16:
		if val != nil {
			b.AppendUInt64(uint64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *uint32:
		if val != nil {
			b.AppendUInt64(uint64(*val))
		} else {
			b.AppendString("null", false)
		}
	case *uint64:
		if val != nil {
			b.AppendUInt64(*val)
		} else {
			b.AppendString("null", false)
		}
	case *float32:
		if val != nil {
			b.AppendFloat32(*val)
		} else {
			b.AppendString("null", false)
		}
	case *float64:
		if val != nil {
			b.AppendFloat64(*val)
		} else {
			b.AppendString("null", false)
		}
	case []string:
		b.AppendStrings(val, quota)
	case []bool:
		b.AppendBools(val)
	case []int:
		b.AppendInts(val)
	case []int8:
		b.AppendInts8(val)
	case []int16:
		b.AppendInts16(val)
	case []int32:
		b.AppendInts32(val)
	case []int64:
		b.AppendInts64(val)
	case []uint:
		b.AppendUInts(val)
	case []uint16:
		b.AppendUInts16(val)
	case []uint32:
		b.AppendUInts32(val)
	case []uint64:
		b.AppendUInts64(val)
	case []float32:
		b.AppendFloats32(val)
	case []float64:
		b.AppendFloats64(val)
	case nil:
		b.AppendString("null", false)
	default:
		b.AppendInterface(val, marshallFn)
	}
	return b
}
