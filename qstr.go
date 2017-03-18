package qstr

import (
	"fmt"
	"html"
	"html/template"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// RGBColor is a color in the RGB space. R, G, and B are in the range [0, 1]
type RGBColor struct {
	// Red, Green, and Blue
	R, G, B float64
}

func NewRGBColorFrom255(r, g, b float64) (RGBColor) {
	r = r/255.0
	g = g/255.0
	b = b/255.0

	return RGBColor{r, g, b}
}

// HexToRGB converts a sequence of three hexadecimal characters into an RGBColor
func HexToRGB(r string, g string, b string) (c RGBColor) {

	red, _ := strconv.ParseInt(fmt.Sprintf("%s%s", r, r), 16, 0)
	green, _ := strconv.ParseInt(fmt.Sprintf("%s%s", g, g), 16, 0)
	blue, _ := strconv.ParseInt(fmt.Sprintf("%s%s", b, b), 16, 0)

	return NewRGBColorFrom255(float64(red), float64(green), float64(blue))
}

// SpanStr converts an RGBColor into a string representing an
// HTML span with inline coloring
func (c *RGBColor) SpanStr() string {
	// convert to a [0, 255] range
	r255 := int(c.R*255.0)
	g255 := int(c.G*255.0)
	b255 := int(c.B*255.0)

	return fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", r255, g255, b255)
}

// HSL converts an RGBColor into an HSLColor. Ported from python's colorsys module.
func (c *RGBColor) HSL() HSLColor {
	maxC := math.Max(math.Max(c.R, c.G), c.B)
	minC := math.Min(math.Min(c.R, c.G), c.B)

	var h, l, s float64

	l = (minC + maxC) / 2.0
	if minC == maxC {
		return HSLColor{0.0, 0.0, l}
	}
	if l <= 0.5 {
		s = (maxC - minC) / (maxC + minC)
	} else {
		s = (maxC - minC) / (2.0 - maxC - minC)
	}
	rc := (maxC - c.R) / (maxC - minC)
	gc := (maxC - c.G) / (maxC - minC)
	bc := (maxC - c.B) / (maxC - minC)

	if c.R == maxC {
		h = bc - gc
	} else if c.G == maxC {
		h = 2.0 + rc - bc
	} else {
		h = 4.0 + gc - rc
	}

	h = math.Mod(h/6.0, 1.0)
	if h < 0.0 {
		h = h + 1.0
	}

	return HSLColor{h, s, l}
}

// CapLightness returns an RGB color that is trimmed to have a lightness
// value between floor and ceiling, where floor < ceiling and both
// floor and ceiling are between 0 and 1 where 0 is no
// light (black) and 1 is maximum light (white).
func (c *RGBColor) CapLightness(floor float64, ceiling float64) (r RGBColor) {
	// check invalid values
	if floor >= ceiling || floor < 0 || ceiling > 1 {
		return *c
	}

	h := c.HSL()
	if h.L < floor {
		h.L = floor
	} else if h.L > ceiling {
		h.L = ceiling
	} else {
		// no need to do any conversion, just return back what we had before
		return *c
	}
	return h.RGB()
}

// HSLColor is a color in the HSL space.
type HSLColor struct {
	// Hue, Saturation, and Lightness
	H, S, L float64
}

var ONE_THIRD = 1.0 / 3.0
var ONE_SIXTH = 1.0 / 6.0
var TWO_THIRD = 2.0 / 3.0

func v(m1 float64, m2 float64, hue float64) float64 {
	hue = math.Mod(hue, 1.0)
	if hue < 0.0 {
		hue = hue + 1.0
	}

	if hue < ONE_SIXTH {
		return m1 + (m2-m1)*hue*6.0
	}
	if hue < 0.5 {
		return m2
	}
	if hue < TWO_THIRD {
		return m1 + (m2-m1)*(TWO_THIRD-hue)*6.0
	}
	return m1
}

// RGB converts an HSLColor to an RGBColor. Ported from python's colorsys module.
func (c *HSLColor) RGB() RGBColor {
	if c.S == 0.0 {
		return RGBColor{c.L, c.L, c.L}
	}

	var m2 float64
	if c.L <= 0.5 {
		m2 = c.L * (1.0 + c.S)
	} else {
		m2 = c.L + c.S - (c.L * c.S)
	}

	m1 := 2.0*c.L - m2

	return RGBColor{
		R: v(m1, m2, c.H+ONE_THIRD),
		G: v(m1, m2, c.H),
		B: v(m1, m2, c.H-ONE_THIRD),
	}
}

// color codes of the form ^N
var decColors = regexp.MustCompile(`\^(\d)`)

// color codes of the form ^xNNN
var hexColors = regexp.MustCompile(`\^x([\dA-Fa-f])([\dA-Fa-f])([\dA-Fa-f])`)

// either of the above forms of color codes
var allColors = regexp.MustCompile(`\^(\d|x[\dA-Fa-f]{3})`)

// Type QStr is a Quake-style string with optional embedded color codes within
// it. The color codes can take a basic form of ^N, where N is in 0..9. These
// represent a basic color palette. The more expanded color code form is ^xNNN,
// where the Ns are hexadecimal characters. This form allows you to specify
// colors with greater precision.
type QStr string

// Stripped removes all of the color codes from string
func (s *QStr) Stripped() string {
	return allColors.ReplaceAllString(string(*s), "")
}

// HTML returns the HTML representation of the QStr. Color codes are converted
// into nested <span> elements with the appropriate color attached as inline
// CSS.
func (s *QStr) HTML() template.HTML {
	// color representation by key for the "^n" format, where n is 0-9
	var decimalSpans = map[string]string{
		"^0": "<span style='color:rgb(128,128,128)'>",
		"^1": "<span style='color:rgb(255,0,0)'>",
		"^2": "<span style='color:rgb(51,255,0)'>",
		"^3": "<span style='color:rgb(255,255,0)'>",
		"^4": "<span style='color:rgb(51,102,255)'>",
		"^5": "<span style='color:rgb(51,255,255)'>",
		"^6": "<span style='color:rgb(255,51,102)'>",
		"^7": "<span style='color:rgb(255,255,255)'>",
		"^8": "<span style='color:rgb(153,153,153)'>",
		"^9": "<span style='color:rgb(128,128,128)'>",
	}

	// cast once to the string representation 'r'
	r := string(*s)

	// remove HTMl special characters
	r = html.EscapeString(r)

	// substitute matches of the form ^n, with n in 0..9
	matchedDecStrings := decColors.FindAllStringSubmatch(r, -1)
	for _, v := range matchedDecStrings {
		r = strings.Replace(r, v[0], decimalSpans[v[0]], 1)
	}

	// substitute matches of the form ^xrgb
	// with r, g, and b being hexadecimal digits
	// also cap the lightness to be in the given range
	matchedHexStrings := hexColors.FindAllStringSubmatch(r, -1)
	for _, v := range matchedHexStrings {
		c := HexToRGB(v[1], v[2], v[3])
		c = c.CapLightness(0.5, 1.0)
		r = strings.Replace(r, v[0], c.SpanStr(), 1)
	}

	// add the appropriate amount of closing spans
	for i := 0; i < (len(matchedDecStrings) + len(matchedHexStrings)); i++ {
		r = fmt.Sprintf("%s%s", r, "</span>")
	}

	return template.HTML(r)
}
