package opencv

import "fmt"

func validateColorConversion(src ColorModel, code ColorConversionCode) error {
	// Strict mode validates model/code compatibility for the currently exposed
	// conversion code range (0..11). Constants are aliased by OpenCV semantics,
	// so validation is done on numeric code + source model class.
	if src == Unknown {
		return nil
	}

	c := int(code)
	if c < 0 || c > 11 {
		return fmt.Errorf("cvtColor strict: unsupported conversion code %d", c)
	}

	// Gray source must use Gray->color only.
	if src == Gray {
		if c == 8 || c == 9 {
			return nil
		}
		return fmt.Errorf("cvtColor strict: source model Gray incompatible with conversion code %d", c)
	}

	// Non-Gray source must not use Gray->* conversions.
	if src != Gray && (c == 8 || c == 9) {
		return fmt.Errorf("cvtColor strict: non-Gray source model %v incompatible with Gray->* conversion code %d", src, c)
	}

	// BGR/RGB 3-channel sources are incompatible with 4->3 conversions.
	if (src == BGR || src == RGB) && (c == 1 || c == 3 || c == 5 || c == 10 || c == 11) {
		// c==10/11 are *2Gray from 4ch in strict mapping.
		return fmt.Errorf("cvtColor strict: source model %v incompatible with conversion code %d", src, c)
	}

	// RGBA 4-channel source incompatible with 3->4 additions.
	if src == RGBA && (c == 0 || c == 2 || c == 9) {
		return fmt.Errorf("cvtColor strict: source model RGBA incompatible with conversion code %d", c)
	}

	// 3-channel source incompatible with 4-channel grayscale codes.
	if (src == BGR || src == RGB) && (c == 10 || c == 11) {
		return fmt.Errorf("cvtColor strict: source model %v incompatible with 4ch->gray conversion code %d", src, c)
	}

	// RGBA source incompatible with 3-channel grayscale codes.
	if src == RGBA && (c == 6 || c == 7) {
		return fmt.Errorf("cvtColor strict: source model RGBA incompatible with 3ch->gray conversion code %d", c)
	}
	return nil
}

func outputModelForCode(src ColorModel, code ColorConversionCode) (ColorModel, bool) {
	c := int(code)
	if c < 0 || c > 11 {
		return Unknown, false
	}

	// Gray outputs
	if c >= 6 && c <= 11 {
		return Gray, true
	}

	// Gray->color
	if c == 8 {
		return BGR, true
	}
	if c == 9 {
		return RGBA, true
	}

	// 3<->4 channel toggles; with no BGRA enum, use RGBA for 4ch class.
	if c >= 0 && c <= 3 {
		if src == RGBA {
			return BGR, true
		}
		return RGBA, true
	}

	// Byte-order permutations preserve channel count/model class.
	if c == 4 || c == 5 {
		return src, true
	}

	return Unknown, false
}
