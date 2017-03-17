package qstr

import "testing"

func TestHSL(t *testing.T) {
	rgbcolors := []RGBColor{
		{0, 0, 0},
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
		{0, 255, 255},
		{255, 255, 0},
		{255, 0, 255},
		{255, 255, 255},
	}

	hslcolors := []HSLColor{
		{0.0, 0.0, 0.0},
		{0, 1, 0.5},
		{0.3333333333333333, 1, 0.5},
		{0.6666666666666666, 1, 0.5},
		{0.5, 1, 0.5},
		{0.16666666666666666, 1, 0.5},
		{0.8333333333333334, 1, 0.5},
		{0, 0, 1},
	}

	for i, v := range rgbcolors {
		expected := hslcolors[i]
		received := v.HSL()

		if received.H != expected.H ||
			received.S != expected.S ||
			received.L != expected.L {
			t.Errorf("Incorrect HSL translation for RGB color %v. Expected: %v, Got: %v.", v, expected, received)
		}
	}
}

func TestCappedLow(t *testing.T) {
	c := RGBColor{255, 255, 255}
	cbar := c.Capped(0.0, 0.5)
	h := cbar.HSL()

	// some flexibility needed due to the rounding to int
	if h.L > 0.51 {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected: <= 50.0, Got: %v.", c, h.L)
	}
}

func TestCappedMid(t *testing.T) {
	c := RGBColor{127, 127, 127}
	cbar := c.Capped(0.25, 0.75)
	h := cbar.HSL()

	// some flexibility needed due to the rounding to int
	// but essentially it must remain in the same range as before
	if h.L < 0.49 || h.L > 0.51 {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected around 50.0, Got: %v.", c, h.L)
	}
}

func TestCappedHigh(t *testing.T) {
	c := RGBColor{0, 0, 0}
	cbar := c.Capped(0.5, 1)
	h := cbar.HSL()

	// some flexibility needed due to the rounding to int
	if h.L < 0.5 {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected: >= 50.0, Got: %v.", c, h.L)
	}
}
