package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Color represents an RGBA color with values 0-255 for RGB and 0-1 for alpha.
type Color struct {
	R, G, B uint8
	A       float64
}

// ColorConvertRequest is the JSON body for color convert endpoint.
type ColorConvertRequest struct {
	Value string `json:"value"`
}

// ColorConvertResponse contains all color format representations.
type ColorConvertResponse struct {
	Hex   string `json:"hex"`
	RGB   string `json:"rgb"`
	HSL   string `json:"hsl"`
	OKLCH string `json:"oklch"`
}

// ContrastRequest is the JSON body for contrast check endpoint.
type ContrastRequest struct {
	Color1 string `json:"color1"`
	Color2 string `json:"color2"`
}

// ContrastResponse contains contrast ratio and WCAG compliance info.
type ContrastResponse struct {
	Ratio         float64 `json:"ratio"`
	RatioText     string  `json:"ratioText"`
	AANormal      bool    `json:"aaNormal"`
	AALarge       bool    `json:"aaLarge"`
	AAANormal     bool    `json:"aaaNormal"`
	AAALarge      bool    `json:"aaaLarge"`
	Color1Hex     string  `json:"color1Hex"`
	Color2Hex     string  `json:"color2Hex"`
	Color1Lighter bool    `json:"color1Lighter"`
}

// BlindnessRequest is the JSON body for color blindness simulation.
type BlindnessRequest struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// BlindnessResponse contains the simulated color.
type BlindnessResponse struct {
	Original  string `json:"original"`
	Simulated string `json:"simulated"`
	Type      string `json:"type"`
}

// PaletteRequest is the JSON body for palette extraction.
type PaletteRequest struct {
	Image string `json:"image"`
	Count int    `json:"count"`
}

// PaletteResponse contains extracted colors.
type PaletteResponse struct {
	Colors []string `json:"colors"`
}

// ShadesRequest is the JSON body for shade/tint generation.
type ShadesRequest struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// ShadesResponse contains the generated shades and tints.
type ShadesResponse struct {
	Shades []string `json:"shades"`
	Tints  []string `json:"tints"`
	Base   string   `json:"base"`
}

// HarmoniesRequest is the JSON body for color harmonies.
type HarmoniesRequest struct {
	Value string `json:"value"`
}

// HarmoniesResponse contains color harmony values.
type HarmoniesResponse struct {
	Base          string   `json:"base"`
	Complementary string   `json:"complementary"`
	Triadic       []string `json:"triadic"`
	Analogous     []string `json:"analogous"`
	SplitComp     []string `json:"splitComplementary"`
	Tetradic      []string `json:"tetradic"`
}

// Named CSS colors (subset of common colors)
var namedColors = map[string]string{
	"black":         "#000000",
	"white":         "#ffffff",
	"red":           "#ff0000",
	"green":         "#008000",
	"blue":          "#0000ff",
	"yellow":        "#ffff00",
	"cyan":          "#00ffff",
	"magenta":       "#ff00ff",
	"silver":        "#c0c0c0",
	"gray":          "#808080",
	"grey":          "#808080",
	"maroon":        "#800000",
	"olive":         "#808000",
	"lime":          "#00ff00",
	"aqua":          "#00ffff",
	"teal":          "#008080",
	"navy":          "#000080",
	"fuchsia":       "#ff00ff",
	"purple":        "#800080",
	"orange":        "#ffa500",
	"pink":          "#ffc0cb",
	"brown":         "#a52a2a",
	"coral":         "#ff7f50",
	"gold":          "#ffd700",
	"indigo":        "#4b0082",
	"ivory":         "#fffff0",
	"khaki":         "#f0e68c",
	"lavender":      "#e6e6fa",
	"salmon":        "#fa8072",
	"turquoise":     "#40e0d0",
	"violet":        "#ee82ee",
	"rebeccapurple": "#663399",
	"aliceblue":     "#f0f8ff",
	"antiquewhite":  "#faebd7",
	"beige":         "#f5f5dc",
	"bisque":        "#ffe4c4",
	"blanchedalmond": "#ffebcd",
	"blueviolet":    "#8a2be2",
	"burlywood":     "#deb887",
	"cadetblue":     "#5f9ea0",
	"chartreuse":    "#7fff00",
	"chocolate":     "#d2691e",
	"cornflowerblue": "#6495ed",
	"cornsilk":      "#fff8dc",
	"crimson":       "#dc143c",
	"darkblue":      "#00008b",
	"darkcyan":      "#008b8b",
	"darkgoldenrod": "#b8860b",
	"darkgray":      "#a9a9a9",
	"darkgreen":     "#006400",
	"darkkhaki":     "#bdb76b",
	"darkmagenta":   "#8b008b",
	"darkolivegreen": "#556b2f",
	"darkorange":    "#ff8c00",
	"darkorchid":    "#9932cc",
	"darkred":       "#8b0000",
	"darksalmon":    "#e9967a",
	"darkseagreen":  "#8fbc8f",
	"darkslateblue": "#483d8b",
	"darkslategray": "#2f4f4f",
	"darkturquoise": "#00ced1",
	"darkviolet":    "#9400d3",
	"deeppink":      "#ff1493",
	"deepskyblue":   "#00bfff",
	"dimgray":       "#696969",
	"dodgerblue":    "#1e90ff",
	"firebrick":     "#b22222",
	"floralwhite":   "#fffaf0",
	"forestgreen":   "#228b22",
	"gainsboro":     "#dcdcdc",
	"ghostwhite":    "#f8f8ff",
	"greenyellow":   "#adff2f",
	"honeydew":      "#f0fff0",
	"hotpink":       "#ff69b4",
	"indianred":     "#cd5c5c",
	"lawngreen":     "#7cfc00",
	"lemonchiffon":  "#fffacd",
	"lightblue":     "#add8e6",
	"lightcoral":    "#f08080",
	"lightcyan":     "#e0ffff",
	"lightgray":     "#d3d3d3",
	"lightgreen":    "#90ee90",
	"lightpink":     "#ffb6c1",
	"lightsalmon":   "#ffa07a",
	"lightseagreen": "#20b2aa",
	"lightskyblue":  "#87cefa",
	"lightslategray": "#778899",
	"lightsteelblue": "#b0c4de",
	"lightyellow":   "#ffffe0",
	"limegreen":     "#32cd32",
	"linen":         "#faf0e6",
	"mediumaquamarine": "#66cdaa",
	"mediumblue":    "#0000cd",
	"mediumorchid":  "#ba55d3",
	"mediumpurple":  "#9370db",
	"mediumseagreen": "#3cb371",
	"mediumslateblue": "#7b68ee",
	"mediumspringgreen": "#00fa9a",
	"mediumturquoise": "#48d1cc",
	"mediumvioletred": "#c71585",
	"midnightblue":  "#191970",
	"mintcream":     "#f5fffa",
	"mistyrose":     "#ffe4e1",
	"moccasin":      "#ffe4b5",
	"navajowhite":   "#ffdead",
	"oldlace":       "#fdf5e6",
	"olivedrab":     "#6b8e23",
	"orangered":     "#ff4500",
	"orchid":        "#da70d6",
	"palegoldenrod": "#eee8aa",
	"palegreen":     "#98fb98",
	"paleturquoise": "#afeeee",
	"palevioletred": "#db7093",
	"papayawhip":    "#ffefd5",
	"peachpuff":     "#ffdab9",
	"peru":          "#cd853f",
	"plum":          "#dda0dd",
	"powderblue":    "#b0e0e6",
	"rosybrown":     "#bc8f8f",
	"royalblue":     "#4169e1",
	"saddlebrown":   "#8b4513",
	"sandybrown":    "#f4a460",
	"seagreen":      "#2e8b57",
	"seashell":      "#fff5ee",
	"sienna":        "#a0522d",
	"skyblue":       "#87ceeb",
	"slateblue":     "#6a5acd",
	"slategray":     "#708090",
	"snow":          "#fffafa",
	"springgreen":   "#00ff7f",
	"steelblue":     "#4682b4",
	"tan":           "#d2b48c",
	"thistle":       "#d8bfd8",
	"tomato":        "#ff6347",
	"wheat":         "#f5deb3",
	"whitesmoke":    "#f5f5f5",
	"yellowgreen":   "#9acd32",
}

// parseColor parses a color string in various formats and returns a Color.
func parseColor(s string) (Color, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// Check named colors first
	if hex, ok := namedColors[s]; ok {
		s = hex
	}

	// Hex format
	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}

	// RGB/RGBA format
	if strings.HasPrefix(s, "rgb") {
		return parseRGB(s)
	}

	// HSL/HSLA format
	if strings.HasPrefix(s, "hsl") {
		return parseHSL(s)
	}

	// OKLCH format
	if strings.HasPrefix(s, "oklch") {
		return parseOKLCH(s)
	}

	// Try parsing as hex without #
	if len(s) == 6 || len(s) == 3 {
		return parseHex("#" + s)
	}

	return Color{}, fmt.Errorf("unsupported color format: %s", s)
}

func parseHex(s string) (Color, error) {
	s = strings.TrimPrefix(s, "#")
	var r, g, b uint8
	var a float64 = 1.0

	switch len(s) {
	case 3:
		r64, _ := strconv.ParseUint(string(s[0])+string(s[0]), 16, 8)
		g64, _ := strconv.ParseUint(string(s[1])+string(s[1]), 16, 8)
		b64, _ := strconv.ParseUint(string(s[2])+string(s[2]), 16, 8)
		r, g, b = uint8(r64), uint8(g64), uint8(b64)
	case 6:
		r64, _ := strconv.ParseUint(s[0:2], 16, 8)
		g64, _ := strconv.ParseUint(s[2:4], 16, 8)
		b64, _ := strconv.ParseUint(s[4:6], 16, 8)
		r, g, b = uint8(r64), uint8(g64), uint8(b64)
	case 8:
		r64, _ := strconv.ParseUint(s[0:2], 16, 8)
		g64, _ := strconv.ParseUint(s[2:4], 16, 8)
		b64, _ := strconv.ParseUint(s[4:6], 16, 8)
		a64, _ := strconv.ParseUint(s[6:8], 16, 8)
		r, g, b = uint8(r64), uint8(g64), uint8(b64)
		a = float64(a64) / 255.0
	default:
		return Color{}, fmt.Errorf("invalid hex color: #%s", s)
	}
	return Color{R: r, G: g, B: b, A: a}, nil
}

var rgbRe = regexp.MustCompile(`rgba?\s*\(\s*([\d.]+)\s*[,\s]\s*([\d.]+)\s*[,\s]\s*([\d.]+)\s*(?:[,/]\s*([\d.]+))?\s*\)`)

func parseRGB(s string) (Color, error) {
	matches := rgbRe.FindStringSubmatch(s)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid rgb format: %s", s)
	}
	r, _ := strconv.ParseFloat(matches[1], 64)
	g, _ := strconv.ParseFloat(matches[2], 64)
	b, _ := strconv.ParseFloat(matches[3], 64)
	a := 1.0
	if matches[4] != "" {
		a, _ = strconv.ParseFloat(matches[4], 64)
		if a > 1 {
			a = a / 255.0
		}
	}
	return Color{R: uint8(clamp(r, 0, 255)), G: uint8(clamp(g, 0, 255)), B: uint8(clamp(b, 0, 255)), A: a}, nil
}

var hslRe = regexp.MustCompile(`hsla?\s*\(\s*([\d.]+)\s*[,\s]\s*([\d.]+)%?\s*[,\s]\s*([\d.]+)%?\s*(?:[,/]\s*([\d.]+))?\s*\)`)

func parseHSL(s string) (Color, error) {
	matches := hslRe.FindStringSubmatch(s)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid hsl format: %s", s)
	}
	h, _ := strconv.ParseFloat(matches[1], 64)
	sat, _ := strconv.ParseFloat(matches[2], 64)
	l, _ := strconv.ParseFloat(matches[3], 64)
	a := 1.0
	if matches[4] != "" {
		a, _ = strconv.ParseFloat(matches[4], 64)
	}
	r, g, b := hslToRGB(h, sat/100, l/100)
	return Color{R: r, G: g, B: b, A: a}, nil
}

var oklchRe = regexp.MustCompile(`oklch\s*\(\s*([\d.]+)%?\s+([.\d]+)\s+([.\d]+)\s*(?:[,/]\s*([\d.]+))?\s*\)`)

func parseOKLCH(s string) (Color, error) {
	matches := oklchRe.FindStringSubmatch(s)
	if matches == nil {
		return Color{}, fmt.Errorf("invalid oklch format: %s", s)
	}
	L, _ := strconv.ParseFloat(matches[1], 64)
	C, _ := strconv.ParseFloat(matches[2], 64)
	H, _ := strconv.ParseFloat(matches[3], 64)
	a := 1.0
	if matches[4] != "" {
		a, _ = strconv.ParseFloat(matches[4], 64)
	}
	if L > 1 {
		L = L / 100
	}
	r, g, b := oklchToRGB(L, C, H)
	return Color{R: r, G: g, B: b, A: a}, nil
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// hslToRGB converts HSL to RGB
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8(math.Round((r + m) * 255)),
		uint8(math.Round((g + m) * 255)),
		uint8(math.Round((b + m) * 255))
}

// rgbToHSL converts RGB to HSL
func rgbToHSL(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l := (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	d := max - min
	var s float64
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	var h float64
	switch max {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/d + 2
	case bf:
		h = (rf-gf)/d + 4
	}
	h *= 60

	return h, s, l
}

// oklchToRGB converts OKLCH to RGB
func oklchToRGB(L, C, H float64) (uint8, uint8, uint8) {
	hRad := H * math.Pi / 180
	a := C * math.Cos(hRad)
	b := C * math.Sin(hRad)

	l_ := L + 0.3963377774*a + 0.2158037573*b
	m_ := L - 0.1055613458*a - 0.0638541728*b
	s_ := L - 0.0894841775*a - 1.2914855480*b

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	rLin := +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	gLin := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bLin := -0.0041960863*l - 0.7034186147*m + 1.7076147010*s

	rLin = clamp(rLin, 0, 1)
	gLin = clamp(gLin, 0, 1)
	bLin = clamp(bLin, 0, 1)

	toSRGB := func(v float64) uint8 {
		if v <= 0.0031308 {
			return uint8(math.Round(v * 12.92 * 255))
		}
		return uint8(math.Round((1.055*math.Pow(v, 1/2.4) - 0.055) * 255))
	}

	return toSRGB(rLin), toSRGB(gLin), toSRGB(bLin)
}

// rgbToOKLCH converts RGB to OKLCH
func rgbToOKLCH(r, g, b uint8) (float64, float64, float64) {
	toLinear := func(v uint8) float64 {
		vf := float64(v) / 255
		if vf <= 0.04045 {
			return vf / 12.92
		}
		return math.Pow((vf+0.055)/1.055, 2.4)
	}

	rLin := toLinear(r)
	gLin := toLinear(g)
	bLin := toLinear(b)

	l := 0.4122214708*rLin + 0.5363325363*gLin + 0.0514459929*bLin
	m := 0.2119034982*rLin + 0.6806995451*gLin + 0.1073969566*bLin
	s := 0.0883024619*rLin + 0.2817188376*gLin + 0.6299787005*bLin

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	L := 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a := 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	bVal := 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_

	C := math.Sqrt(a*a + bVal*bVal)
	H := math.Atan2(bVal, a) * 180 / math.Pi
	if H < 0 {
		H += 360
	}

	return L, C, H
}

// toHex formats Color as hex string
func (c Color) toHex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// toRGB formats Color as rgb() string
func (c Color) toRGB() string {
	if c.A < 1 {
		return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", c.R, c.G, c.B, c.A)
	}
	return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

// toHSL formats Color as hsl() string
func (c Color) toHSL() string {
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	if c.A < 1 {
		return fmt.Sprintf("hsla(%.0f, %.0f%%, %.0f%%, %.2f)", h, s*100, l*100, c.A)
	}
	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s*100, l*100)
}

// toOKLCH formats Color as oklch() string
func (c Color) toOKLCH() string {
	L, C, H := rgbToOKLCH(c.R, c.G, c.B)
	if c.A < 1 {
		return fmt.Sprintf("oklch(%.1f%% %.3f %.1f / %.2f)", L*100, C, H, c.A)
	}
	return fmt.Sprintf("oklch(%.1f%% %.3f %.1f)", L*100, C, H)
}

// ColorConvert parses a color and returns it in all formats.
func ColorConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ColorConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	color, err := parseColor(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, ColorConvertResponse{
		Hex:   color.toHex(),
		RGB:   color.toRGB(),
		HSL:   color.toHSL(),
		OKLCH: color.toOKLCH(),
	})
}

// relativeLuminance calculates WCAG relative luminance
func relativeLuminance(c Color) float64 {
	toLinear := func(v uint8) float64 {
		vf := float64(v) / 255
		if vf <= 0.03928 {
			return vf / 12.92
		}
		return math.Pow((vf+0.055)/1.055, 2.4)
	}

	r := toLinear(c.R)
	g := toLinear(c.G)
	b := toLinear(c.B)

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// ContrastCheck calculates WCAG contrast ratio between two colors.
func ContrastCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ContrastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	c1, err := parseColor(req.Color1)
	if err != nil {
		http.Error(w, "invalid color1: "+err.Error(), http.StatusBadRequest)
		return
	}
	c2, err := parseColor(req.Color2)
	if err != nil {
		http.Error(w, "invalid color2: "+err.Error(), http.StatusBadRequest)
		return
	}

	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)

	lighter := l1
	darker := l2
	c1Lighter := true
	if l2 > l1 {
		lighter = l2
		darker = l1
		c1Lighter = false
	}

	ratio := (lighter + 0.05) / (darker + 0.05)

	writeJSON(w, ContrastResponse{
		Ratio:         math.Round(ratio*100) / 100,
		RatioText:     fmt.Sprintf("%.2f:1", ratio),
		AANormal:      ratio >= 4.5,
		AALarge:       ratio >= 3.0,
		AAANormal:     ratio >= 7.0,
		AAALarge:      ratio >= 4.5,
		Color1Hex:     c1.toHex(),
		Color2Hex:     c2.toHex(),
		Color1Lighter: c1Lighter,
	})
}

// Color blindness simulation matrices
type matrix [3][3]float64

var blindnessMatrices = map[string]matrix{
	"protanopia": {
		{0.567, 0.433, 0},
		{0.558, 0.442, 0},
		{0, 0.242, 0.758},
	},
	"deuteranopia": {
		{0.625, 0.375, 0},
		{0.7, 0.3, 0},
		{0, 0.3, 0.7},
	},
	"tritanopia": {
		{0.95, 0.05, 0},
		{0, 0.433, 0.567},
		{0, 0.475, 0.525},
	},
	"achromatopsia": {
		{0.299, 0.587, 0.114},
		{0.299, 0.587, 0.114},
		{0.299, 0.587, 0.114},
	},
}

// ColorBlindness simulates how a color appears with color vision deficiency.
func ColorBlindness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req BlindnessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	color, err := parseColor(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blindType := strings.ToLower(req.Type)
	if blindType == "" {
		blindType = "deuteranopia"
	}

	mat, ok := blindnessMatrices[blindType]
	if !ok {
		http.Error(w, "unknown blindness type: "+req.Type, http.StatusBadRequest)
		return
	}

	rf := float64(color.R)
	gf := float64(color.G)
	bf := float64(color.B)

	newR := mat[0][0]*rf + mat[0][1]*gf + mat[0][2]*bf
	newG := mat[1][0]*rf + mat[1][1]*gf + mat[1][2]*bf
	newB := mat[2][0]*rf + mat[2][1]*gf + mat[2][2]*bf

	simulated := Color{
		R: uint8(clamp(newR, 0, 255)),
		G: uint8(clamp(newG, 0, 255)),
		B: uint8(clamp(newB, 0, 255)),
		A: color.A,
	}

	writeJSON(w, BlindnessResponse{
		Original:  color.toHex(),
		Simulated: simulated.toHex(),
		Type:      blindType,
	})
}

// ExtractPalette extracts dominant colors from a base64-encoded image.
func ExtractPalette(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req PaletteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	count := req.Count
	if count <= 0 {
		count = 5
	}
	if count > 10 {
		count = 10
	}

	imgData := req.Image
	if idx := strings.Index(imgData, ","); idx > 0 {
		imgData = imgData[idx+1:]
	}

	decoded, err := base64.StdEncoding.DecodeString(imgData)
	if err != nil {
		http.Error(w, "invalid base64 image", http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(strings.NewReader(string(decoded)))
	if err != nil {
		http.Error(w, "failed to decode image: "+err.Error(), http.StatusBadRequest)
		return
	}

	colors := extractDominantColors(img, count)
	hexColors := make([]string, len(colors))
	for i, c := range colors {
		hexColors[i] = c.toHex()
	}

	writeJSON(w, PaletteResponse{Colors: hexColors})
}

// extractDominantColors uses a simple color quantization approach
func extractDominantColors(img image.Image, count int) []Color {
	bounds := img.Bounds()
	colorCounts := make(map[uint32]int)

	step := 1
	if bounds.Dx()*bounds.Dy() > 10000 {
		step = int(math.Sqrt(float64(bounds.Dx()*bounds.Dy()) / 10000))
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			r, g, b, a := img.At(x, y).RGBA()
			if a < 128*256 {
				continue
			}
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			r8 = (r8 / 32) * 32
			g8 = (g8 / 32) * 32
			b8 = (b8 / 32) * 32
			key := uint32(r8)<<16 | uint32(g8)<<8 | uint32(b8)
			colorCounts[key]++
		}
	}

	type colorCount struct {
		color uint32
		count int
	}
	var sorted []colorCount
	for c, cnt := range colorCounts {
		sorted = append(sorted, colorCount{c, cnt})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	var result []Color
	for i := 0; i < count && i < len(sorted); i++ {
		key := sorted[i].color
		result = append(result, Color{
			R: uint8((key >> 16) & 0xFF),
			G: uint8((key >> 8) & 0xFF),
			B: uint8(key & 0xFF),
			A: 1.0,
		})
	}

	if len(result) == 0 {
		result = append(result, Color{R: 128, G: 128, B: 128, A: 1.0})
	}

	return result
}

// GenerateShades generates tints (lighter) and shades (darker) of a color.
func GenerateShades(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ShadesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	color, err := parseColor(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count := req.Count
	if count <= 0 {
		count = 5
	}
	if count > 10 {
		count = 10
	}

	h, s, l := rgbToHSL(color.R, color.G, color.B)

	shades := make([]string, count)
	tints := make([]string, count)

	for i := 0; i < count; i++ {
		shadeFactor := float64(i+1) / float64(count+1)
		shadeL := l * (1 - shadeFactor)
		sr, sg, sb := hslToRGB(h, s, shadeL)
		shades[i] = fmt.Sprintf("#%02x%02x%02x", sr, sg, sb)

		tintFactor := float64(i+1) / float64(count+1)
		tintL := l + (1-l)*tintFactor
		tr, tg, tb := hslToRGB(h, s, tintL)
		tints[i] = fmt.Sprintf("#%02x%02x%02x", tr, tg, tb)
	}

	writeJSON(w, ShadesResponse{
		Base:   color.toHex(),
		Shades: shades,
		Tints:  tints,
	})
}

// GenerateHarmonies generates color harmonies based on color theory.
func GenerateHarmonies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req HarmoniesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	color, err := parseColor(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h, s, l := rgbToHSL(color.R, color.G, color.B)

	hueShift := func(hue, shift float64) float64 {
		h := hue + shift
		for h >= 360 {
			h -= 360
		}
		for h < 0 {
			h += 360
		}
		return h
	}

	toHex := func(hue float64) string {
		r, g, b := hslToRGB(hue, s, l)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	complementary := toHex(hueShift(h, 180))

	triadic := []string{
		toHex(hueShift(h, 120)),
		toHex(hueShift(h, 240)),
	}

	analogous := []string{
		toHex(hueShift(h, -30)),
		toHex(hueShift(h, 30)),
	}

	splitComp := []string{
		toHex(hueShift(h, 150)),
		toHex(hueShift(h, 210)),
	}

	tetradic := []string{
		toHex(hueShift(h, 90)),
		toHex(hueShift(h, 180)),
		toHex(hueShift(h, 270)),
	}

	writeJSON(w, HarmoniesResponse{
		Base:          color.toHex(),
		Complementary: complementary,
		Triadic:       triadic,
		Analogous:     analogous,
		SplitComp:     splitComp,
		Tetradic:      tetradic,
	})
}
