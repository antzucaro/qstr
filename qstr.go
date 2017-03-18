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

// RGBColor is a color in the RGB space. R, G, and B are in the range [0.0, 255.0]
type RGBColor struct {
	// Red, Green, and Blue
	R, G, B float64
}

// HexToRGB converts a sequence of three hexadecimal characters into an RGBColor
func HexToRGB(r string, g string, b string) (c RGBColor) {

	red, _ := strconv.ParseInt(fmt.Sprintf("%s%s", r, r), 16, 0)
	green, _ := strconv.ParseInt(fmt.Sprintf("%s%s", g, g), 16, 0)
	blue, _ := strconv.ParseInt(fmt.Sprintf("%s%s", b, b), 16, 0)

	return RGBColor{float64(red), float64(green), float64(blue)}
}

// SpanStr converts an RGBColor into a string representing an
// HTML span with inline coloring
func (c *RGBColor) SpanStr() string {
	return fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", int(c.R), int(c.G), int(c.B))
}

// HSL converts an RGBColor into an HSLColor
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

	h = math.Mod((h / 6.0), 1.0)
	if h < 0.0 {
		h = h + 1.0
	}

	return HSLColor{h, s, l}
}

// Capped returns an RGB color that is trimmed to have a lightness
// value between floor and ceiling, where floor < ceiling and both
// floor and ceiling are of the range [0.0, 1.0]
func (c *RGBColor) Capped(floor float64, ceiling float64) (r RGBColor) {
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

// RGB converts an HSLColor to an RGBColor
//
// Adapted from http://code.google.com/p/gorilla/source/browse/color/hsl.go,
// Ported from http://goo.gl/Vg1h9
func (c *HSLColor) RGB() (r RGBColor) {
	var fR, fG, fB float64
	if c.S == 0 {
		fR, fG, fB = c.L, c.L, c.L
	} else {
		var q float64
		if c.L < 0.5 {
			q = c.L * (1 + c.S)
		} else {
			q = c.L + c.S - c.S*c.L
		}
		p := 2*c.L - q
		fR = hueToRGB(p, q, c.H+1.0/3)
		fG = hueToRGB(p, q, c.H)
		fB = hueToRGB(p, q, c.H-1.0/3)
	}
	r.R = (fR * 255) + 0.5
	r.G = (fG * 255) + 0.5
	r.B = (fB * 255) + 0.5
	return
}

// hueToRGB is a helper function for HSLToRGB.
// Adapted from http://code.google.com/p/gorilla/source/browse/color/hsl.go,
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
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
		c = c.Capped(0.5, 1.0)
		r = strings.Replace(r, v[0], c.SpanStr(), 1)
	}

	// add the appropriate amount of closing spans
	for i := 0; i < (len(matchedDecStrings) + len(matchedHexStrings)); i++ {
		r = fmt.Sprintf("%s%s", r, "</span>")
	}

	return template.HTML(r)
}
