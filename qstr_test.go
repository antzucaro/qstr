package qstr

import (
	"fmt"
	"math"
	"testing"
)

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
		{1, 0, 0},
	}

	hslcolors := []HSLColor{
		{0.0, 0.0, 0.0},
		{0.0, -1.00790513834, 127.5},
		{0.333333333333, -1.00790513834, 127.5},
		{0.666666666667, -1.00790513834, 127.5},
		{0.5, -1.00790513834, 127.5},
		{0.166666666667, -1.00790513834, 127.5},
		{0.833333333333, -1.00790513834, 127.5},
		{0.0, 0.0, 255.0},
		{0.0, 1.0, 0.5},
	}

	// if the diff goes beyond this value, the test will fail
	tolerance := 0.005

	for i, v := range rgbcolors {
		expected := hslcolors[i]
		received := v.HSL()

		hDiff := math.Abs(expected.H - received.H)
		sDiff := math.Abs(expected.S - received.S)
		lDiff := math.Abs(expected.L - received.L)
		if hDiff > tolerance || sDiff > tolerance || lDiff > tolerance {
			t.Errorf("Incorrect HSLv2 translation for RGB color %v. Expected: %v, Got: %v.", v, expected, received)
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

func TestCappedInvalid(t *testing.T) {
	c := RGBColor{0, 0, 0}

	// this floor value is invalid, so c should not be modified
	cbar := c.Capped(-1, 1)

	if cbar != c {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected the same value, but got: %v.", c, cbar)
	}
}

func TestStrippedQStr(t *testing.T) {
	nicks := []string{
		"Anti^x444body",
		"^x444Antibody",
		"Antibody^x444",
		"Anti^7body",
		"^7Antibody",
		"Antibody^7",
	}

	expected := "Antibody"
	for _, nick := range nicks {
		nickQ := QStr(nick)
		received := nickQ.Stripped()

		if received != expected {
			t.Errorf("Incorrect stripping applied to %v. Expected: %v, Got: %v.", nick, expected, received)
		}
	}
}

func TestHexToRGB(t *testing.T) {
	var hexRGBList = []struct {
		R string
		G string
		B string
	}{
		{"A", "A", "A"},
		{"0", "0", "0"},
		{"4", "4", "4"},
		{"F", "F", "F"},
	}

	expectedList := []RGBColor{
		RGBColor{170, 170, 170},
		RGBColor{0, 0, 0},
		RGBColor{68, 68, 68},
		RGBColor{255, 255, 255},
	}

	for i, input := range hexRGBList {
		received := HexToRGB(input.R, input.G, input.B)
		expected := expectedList[i]

		if received != expected {
			t.Errorf("Incorrect HexToRGB value returned. Expected: %v, Got: %v.", expected, received)
		}
	}
}

func TestSpanStr(t *testing.T) {
	expected := fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", 1, 2, 3)
	color := RGBColor{1, 2, 3}
	received := color.SpanStr()

	if received != expected {
		t.Errorf("Incorrect SpanStr value returned. Expected: %v, Got: %v.", expected, received)
	}
}
