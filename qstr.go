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

// RGBColor is a color in the RGB space. R, G, and B are in the range [0, 255]
type RGBColor struct {
	// Red
	R int

	// Green
	G int

	// Blue
	B int
}

// HexToRGB converts a sequence of three hexadecimal characters into an RGBColor
func HexToRGB(r string, g string, b string) (c RGBColor) {

	red, _ := strconv.ParseInt(fmt.Sprintf("%s%s", r, r), 16, 0)
	green, _ := strconv.ParseInt(fmt.Sprintf("%s%s", g, g), 16, 0)
	blue, _ := strconv.ParseInt(fmt.Sprintf("%s%s", b, b), 16, 0)

	return RGBColor{int(red), int(green), int(blue)}
}

// SpanStr converts an RGBColor into a string representing an
// HTML span with inline coloring
func (c *RGBColor) SpanStr() string {
	return fmt.Sprintf("<span style=\"color:rgb(%d,%d,%d)\">", c.R, c.G, c.B)
}

// HSL converts an RGBColor to a HSLColor
//
// Adapted from http://code.google.com/p/gorilla/source/browse/color/hsl.go,
// which in turn was ported from http://goo.gl/Vg1h9.
func (c *RGBColor) HSL() (h HSLColor) {
	fR := float64(c.R) / 255
	fG := float64(c.G) / 255
	fB := float64(c.B) / 255

	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)

	h.L = (max + min) / 2

	if max == min {
		// Achromatic.
		h.H, h.S = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if h.L > 0.5 {
			h.S = d / (2.0 - max - min)
		} else {
			h.S = d / (max + min)
		}
		switch max {
		case fR:
			h.H = (fG - fB) / d
			if fG < fB {
				h.H += 6
			}
		case fG:
			h.H = (fB-fR)/d + 2
		case fB:
			h.H = (fR-fG)/d + 4
		}
		h.H /= 6
	}
	return
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
	// Hue
	H float64

	// Saturation
	S float64

	// Lightness
	L float64
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
	r.R = int((fR * 255) + 0.5)
	r.G = int((fG * 255) + 0.5)
	r.B = int((fB * 255) + 0.5)
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

var allColors = regexp.MustCompile(`\^(\d|x[\dA-Fa-f]{3})`)
var decColors = regexp.MustCompile(`\^(\d)`)
var hexColors = regexp.MustCompile(`\^x([\dA-Fa-f])([\dA-Fa-f])([\dA-Fa-f])`)

type QStr string

func (s *QStr) Stripped() string {
	return allColors.ReplaceAllString(string(*s), "")
}

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
