package qstr

import (
	"fmt"
	"math"
	"testing"
)

func TestHSL(t *testing.T) {
	rgbColors := []RGBColor{
		{0, 0, 0},
		{0, 0, 1},
		{0, 1, 0},
		{0, 1, 1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, 1},
	}

	hslcolors := []HSLColor{
		{0, 0, 0},
		{0.666666666667, 1, 0.5},
		{0.333333333333, 1, 0.5},
		{0.5, 1, 0.5},
		{0.0, 1, 0.5},
		{0.833333333333, 1, 0.5},
		{0.0, 0.0, 1.0},
	}

	// if the diff goes beyond this value, the test will fail
	tolerance := 0.005

	for i, v := range rgbColors {
		expected := hslcolors[i]
		received := v.HSL()

		hDiff := math.Abs(expected.H - received.H)
		sDiff := math.Abs(expected.S - received.S)
		lDiff := math.Abs(expected.L - received.L)
		if hDiff > tolerance || sDiff > tolerance || lDiff > tolerance {
			t.Errorf("Incorrect HSL translation for RGB color %v. Expected: %v, Got: %v.", v, expected, received)
		}
	}
}

func TestCappedLow(t *testing.T) {
	c := RGBColor{1, 1, 1}
	cbar := c.CapLightness(0.0, 0.5)
	h := cbar.HSL()

	// some flexibility needed due to the rounding to int
	if h.L > 0.5 {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected: <= 0.5, Got: %f.", c, h.L)
	}
}

func TestCappedHigh(t *testing.T) {
	c := RGBColor{0, 0, 0}
	cbar := c.CapLightness(0.5, 1)
	h := cbar.HSL()

	// some flexibility needed due to the rounding to int
	if h.L < 0.5 {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected: >= 50.0, Got: %v.", c, h.L)
	}
}

func TestCappedInvalid(t *testing.T) {
	c := RGBColor{0, 0, 0}

	// this floor value is invalid, so c should not be modified
	cbar := c.CapLightness(-1, 1)

	if cbar != c {
		t.Errorf("Incorrect HSL cap for RGB color %v. Expected the same value, but got: %v.", c, cbar)
	}
}

func TestCappedNoChange(t *testing.T) {
	c := RGBColor{0, 0, 0}

	// this floor value is invalid, so c should not be modified
	cbar := c.CapLightness(0, 1)

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
		{"0", "0", "0"},
		{"0", "0", "F"},
		{"0", "F", "0"},
		{"F", "0", "0"},
		{"F", "0", "F"},
		{"F", "F", "F"},
	}

	expectedList := []RGBColor{
		RGBColor{0, 0, 0},
		RGBColor{0, 0, 1},
		RGBColor{0, 1, 0},
		RGBColor{1, 0, 0},
		RGBColor{1, 0, 1},
		RGBColor{1, 1, 1},
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
	expected := fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", 255, 127, 0)
	color := RGBColor{1, 0.5, 0}
	received := color.SpanStr()

	if received != expected {
		t.Errorf("Incorrect SpanStr value returned. Expected: %v, Got: %v.", expected, received)
	}
}

func TestDecode(t *testing.T) {
	input := QStr("abcdefgh")
	expected := QStr("abcd😊😊😊efgh")

	decodeMap := map[rune]rune{'': '😊'}

	decoded := input.Decode(decodeMap)
	if decoded != expected {
		t.Errorf("Incorrect decoding. Expected: %v, Got: %v.", expected, decoded)
	}
}
