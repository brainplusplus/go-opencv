package opencv

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type matHandle = uint64

// Mat is an opaque OpenCV matrix owned by a Runtime.
// It intentionally mirrors OpenCV.js explicit lifecycle semantics: users must Close.
type Mat struct {
	runtime *Runtime
	handle  matHandle
	model   ColorModel
	closed  atomic.Bool
}

func (m *Mat) Close() error {
	if m == nil || m.runtime == nil || m.handle == 0 {
		return ErrInvalidMat
	}
	if m.closed.Swap(true) {
		return ErrClosed
	}
	return m.runtime.closeMat(m.handle)
}

// Delete is an OpenCV.js-compatible alias for Close.
func (m *Mat) Delete() error { return m.Close() }

func (m *Mat) Rows() (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	return m.runtime.backend.MatRows(m.runtime.ctx, m.handle)
}

func (m *Mat) Cols() (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	return m.runtime.backend.MatCols(m.runtime.ctx, m.handle)
}

func (m *Mat) Type() (MatType, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	typ, err := m.runtime.backend.MatType(m.runtime.ctx, m.handle)
	return MatType(typ), err
}

func (m *Mat) Clone() (*Mat, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	h, err := m.runtime.backend.MatClone(m.runtime.ctx, m.handle)
	if err != nil {
		return nil, err
	}
	return m.runtime.wrapMat(h), nil
}

func (m *Mat) Empty() (bool, error) {
	if err := m.validate(); err != nil {
		return false, err
	}
	return m.runtime.backend.MatEmpty(m.runtime.ctx, m.handle)
}

func (m *Mat) ElemSize() (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	return m.runtime.backend.MatElemSize(m.runtime.ctx, m.handle)
}

func (m *Mat) Step() (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	return m.runtime.backend.MatStep(m.runtime.ctx, m.handle)
}

func (m *Mat) Row(row int) (*Mat, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	h, err := m.runtime.backend.MatRow(m.runtime.ctx, m.handle, int32(row))
	if err != nil {
		return nil, err
	}
	return m.runtime.wrapMatWithModel(h, m.model), nil
}

func (m *Mat) Col(col int) (*Mat, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	h, err := m.runtime.backend.MatCol(m.runtime.ctx, m.handle, int32(col))
	if err != nil {
		return nil, err
	}
	return m.runtime.wrapMatWithModel(h, m.model), nil
}

func (m *Mat) Region(rect Rect) (*Mat, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	h, err := m.runtime.backend.MatRegion(m.runtime.ctx, m.handle, rect.X, rect.Y, rect.Width, rect.Height)
	if err != nil {
		return nil, err
	}
	return m.runtime.wrapMatWithModel(h, m.model), nil
}

func (m *Mat) Reshape(channels, rows int) (*Mat, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	h, err := m.runtime.backend.MatReshape(m.runtime.ctx, m.handle, int32(channels), int32(rows))
	if err != nil {
		return nil, err
	}
	return m.runtime.wrapMatWithModel(h, Unknown), nil
}

// Channels returns the number of channels in the Mat (e.g. 3 for BGR, 1 for Gray).
func (m *Mat) Channels() (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}
	return m.runtime.backend.MatChannels(m.runtime.ctx, m.handle)
}

// CopyTo copies Mat pixel data into dst. dst must have capacity >= rows*step.
// Returns the number of bytes copied.
func (m *Mat) CopyTo(dst []byte) (int, error) {
	if err := m.validate(); err != nil {
		return 0, err
	}

	rows, err := m.Rows()
	if err != nil {
		return 0, err
	}
	step, err := m.Step()
	if err != nil {
		return 0, err
	}

	total := rows * step
	if total == 0 {
		return 0, nil
	}
	if len(dst) < total {
		return 0, fmt.Errorf("copyto: dst buffer too small (%d < %d)", len(dst), total)
	}

	ptr, err := m.runtime.backend.MatDataPtr(m.runtime.ctx, m.handle)
	if err != nil {
		return 0, err
	}

	copy(dst, unsafe.Slice((*byte)(ptr), total))
	return total, nil
}

func (m *Mat) validate() error {
	if m == nil || m.runtime == nil || m.handle == 0 || m.closed.Load() {
		return ErrInvalidMat
	}
	return m.runtime.validateOpen()
}

// ColorModel returns tracked logical channel model metadata.
// Unknown means model cannot be trusted after ambiguous operations.
func (m *Mat) ColorModel() (ColorModel, error) {
	if err := m.validate(); err != nil {
		return Unknown, err
	}
	return m.model, nil
}

// IsColorKnown reports whether ColorModel is not Unknown.
func (m *Mat) IsColorKnown() (bool, error) {
	model, err := m.ColorModel()
	if err != nil {
		return false, err
	}
	return model != Unknown, nil
}

// Total returns the total number of array elements.
func (m *Mat) Total() int {
	if err := m.validate(); err != nil {
		return 0
	}
	n, err := m.runtime.backend.MatTotal(m.runtime.ctx, m.handle)
	if err != nil {
		return 0
	}
	return n
}
