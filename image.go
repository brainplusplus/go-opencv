package opencv

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var rgbaImagePool = sync.Pool{
	New: func() any { return &image.RGBA{} },
}

var byteBufPool = sync.Pool{
	New: func() any { return make([]byte, 0, 1<<20) },
}

func acquireRGBA(cols, rows int) *image.RGBA {
	rgba := rgbaImagePool.Get().(*image.RGBA)
	need := cols * rows * 4
	if cap(rgba.Pix) < need {
		rgba.Pix = make([]byte, need)
	} else {
		rgba.Pix = rgba.Pix[:need]
	}
	rgba.Stride = cols * 4
	rgba.Rect = image.Rect(0, 0, cols, rows)
	return rgba
}

func releaseRGBA(img image.Image) {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return
	}
	// keep capacity for reuse; zeroing not required for encoder correctness.
	rgbaImagePool.Put(rgba)
}

func acquireByteBuf(size int) []byte {
	b := byteBufPool.Get().([]byte)
	if cap(b) < size {
		return make([]byte, size)
	}
	return b[:size]
}

func releaseByteBuf(b []byte) {
	if b == nil {
		return
	}
	byteBufPool.Put(b[:0])
}

// IMRead decodes an image file (JPEG, PNG, GIF) and returns a *Mat.
// Pixel data is converted from Go's RGBA to OpenCV's native BGR channel ordering.
func (r *Runtime) IMRead(path string, model ...ColorModel) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opencv: open %q: %w", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("opencv: decode %q: %w", path, err)
	}

	mode := BGR
	if len(model) > 0 {
		mode = model[0]
	}
	return r.imageToMatWithModel(img, mode)
}

// IMReadBytes decodes an encoded image from memory (JPEG, PNG, GIF bytes) and returns a *Mat.
func (r *Runtime) IMReadBytes(data []byte, model ...ColorModel) (*Mat, error) {
	if err := r.validateOpen(); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("opencv: decode image bytes: %w", err)
	}

	mode := BGR
	if len(model) > 0 {
		mode = model[0]
	}
	return r.imageToMatWithModel(img, mode)
}

// imageToMat converts any image.Image to a BGR Mat (CV_8UC3) and uploads pixels to the backend.
func (r *Runtime) imageToMatWithModel(img image.Image, model ColorModel) (*Mat, error) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Convert to RGBA first (stdlib draw.Draw handles all image types)
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)

	var payload []byte
	var typ MatType

	switch model {
	case BGR:
		payload = make([]byte, h*w*3)
		for i := 0; i < h*w; i++ {
			off := i * 4
			payload[i*3+0] = rgba.Pix[off+2] // B
			payload[i*3+1] = rgba.Pix[off+1] // G
			payload[i*3+2] = rgba.Pix[off+0] // R
		}
		typ = CV8UC3
	case RGB:
		payload = make([]byte, h*w*3)
		for i := 0; i < h*w; i++ {
			off := i * 4
			payload[i*3+0] = rgba.Pix[off+0] // R
			payload[i*3+1] = rgba.Pix[off+1] // G
			payload[i*3+2] = rgba.Pix[off+2] // B
		}
		typ = CV8UC3
	case RGBA:
		payload = rgba.Pix
		typ = CV8UC4
	case Gray:
		payload = make([]byte, h*w)
		for i := 0; i < h*w; i++ {
			off := i * 4
			// BT.601 luma approximation
			rr := int(rgba.Pix[off+0])
			gg := int(rgba.Pix[off+1])
			bb := int(rgba.Pix[off+2])
			payload[i] = uint8((299*rr + 587*gg + 114*bb + 500) / 1000)
		}
		typ = CV8UC1
	default:
		return nil, fmt.Errorf("opencv: unsupported IMRead model %d", model)
	}

	handle, err := r.backend.NewMatFromData(r.ctx, payload, h, w, int32(typ))
	if err != nil {
		return nil, fmt.Errorf("opencv: upload pixels: %w", err)
	}

	return r.wrapMatWithModel(handle, model), nil
}

// IMWrite writes Mat pixel data to an image file.
// Format is determined by file extension: .png, .jpg/.jpeg, .gif.
// Currently supports Mat types CV_8UC1 (Gray), CV_8UC3 (BGR), CV_8UC4 (BGRA/RGBA).
func (r *Runtime) IMWrite(path string, m *Mat, model ...ColorModel) error {
	if err := r.validateOpen(); err != nil {
		return err
	}
	if err := m.validate(); err != nil {
		return err
	}

	writeModel := m.model
	if len(model) > 0 {
		writeModel = model[0]
	}
	if writeModel == Unknown {
		return fmt.Errorf("imwrite: color model unknown; pass explicit model override (BGR/RGB/RGBA/Gray)")
	}

	rows, err := m.Rows()
	if err != nil {
		return fmt.Errorf("save image: rows: %w", err)
	}
	cols, err := m.Cols()
	if err != nil {
		return fmt.Errorf("save image: cols: %w", err)
	}
	step, err := m.Step()
	if err != nil {
		return fmt.Errorf("save image: step: %w", err)
	}
	ch, err := m.Channels()
	if err != nil {
		return fmt.Errorf("save image: channels: %w", err)
	}

	// Read pixel data from Mat
	buf := acquireByteBuf(rows * step)
	defer releaseByteBuf(buf)
	if _, err := m.CopyTo(buf); err != nil {
		return fmt.Errorf("save image: copy: %w", err)
	}

	// Convert to Go image based on channel count
	var img image.Image
	defer releaseRGBA(img)
	switch writeModel {
	case Gray:
		gray := image.NewGray(image.Rect(0, 0, cols, rows))
		for y := 0; y < rows; y++ {
			copy(gray.Pix[y*gray.Stride:], buf[y*step:y*step+cols])
		}
		img = gray
	case RGB:
		rgba := acquireRGBA(cols, rows)
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				off := y*step + x*3
				poff := y*rgba.Stride + x*4
				rgba.Pix[poff+0] = buf[off+0]
				rgba.Pix[poff+1] = buf[off+1]
				rgba.Pix[poff+2] = buf[off+2]
				rgba.Pix[poff+3] = 255
			}
		}
		img = rgba
	case RGBA:
		rgba := acquireRGBA(cols, rows)
		if ch == 4 {
			for y := 0; y < rows; y++ {
				for x := 0; x < cols; x++ {
					off := y*step + x*4
					poff := y*rgba.Stride + x*4
					rgba.Pix[poff+0] = buf[off+0]
					rgba.Pix[poff+1] = buf[off+1]
					rgba.Pix[poff+2] = buf[off+2]
					rgba.Pix[poff+3] = buf[off+3]
				}
			}
		} else if ch == 3 {
			// Source is 3-channel; synthesize alpha=255.
			for y := 0; y < rows; y++ {
				for x := 0; x < cols; x++ {
					off := y*step + x*3
					poff := y*rgba.Stride + x*4
					rgba.Pix[poff+0] = buf[off+0]
					rgba.Pix[poff+1] = buf[off+1]
					rgba.Pix[poff+2] = buf[off+2]
					rgba.Pix[poff+3] = 255
				}
			}
		} else {
			return fmt.Errorf("imwrite: RGBA override requires 3 or 4 channels, got %d", ch)
		}
		img = rgba
	case BGR:
		fallthrough
	default:
		// Fallback to channel-count behavior if metadata unknown
		switch ch {
		case 1:
			gray := image.NewGray(image.Rect(0, 0, cols, rows))
			for y := 0; y < rows; y++ {
				copy(gray.Pix[y*gray.Stride:], buf[y*step:y*step+cols])
			}
			img = gray
		case 3:
			// Mat data is BGR — convert to RGB for Go
			rgba := acquireRGBA(cols, rows)
			for y := 0; y < rows; y++ {
				for x := 0; x < cols; x++ {
					off := y*step + x*3
					poff := y*rgba.Stride + x*4
					rgba.Pix[poff+0] = buf[off+2] // R
					rgba.Pix[poff+1] = buf[off+1] // G
					rgba.Pix[poff+2] = buf[off+0] // B
					rgba.Pix[poff+3] = 255        // A
				}
			}
			img = rgba
		case 4:
			// Mat data is BGRA/RGBA — swap R↔B
			rgba := acquireRGBA(cols, rows)
			for y := 0; y < rows; y++ {
				for x := 0; x < cols; x++ {
					off := y*step + x*4
					poff := y*rgba.Stride + x*4
					rgba.Pix[poff+0] = buf[off+2] // R
					rgba.Pix[poff+1] = buf[off+1] // G
					rgba.Pix[poff+2] = buf[off+0] // B
					rgba.Pix[poff+3] = buf[off+3] // A
				}
			}
			img = rgba
		default:
			return fmt.Errorf("save image: unsupported channel count %d", ch)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save image: create %q: %w", path, err)
	}
	defer f.Close()

	// Encode based on extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return png.Encode(f, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	default:
		return errors.New("save image: unsupported format " + ext + " (use .png or .jpg)")
	}
}
