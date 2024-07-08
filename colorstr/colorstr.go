package colorstr

// the color format information comes from https://en.wikipedia.org/wiki/ANSI_escape_code

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var color = map[string]string{
	"BlackFg":         "30",
	"RedFg":           "31",
	"GreenFg":         "32",
	"YellowFg":        "33",
	"BlueFg":          "34",
	"MagentaFg":       "35",
	"CyanFg":          "36",
	"WhiteFg":         "37",
	"BrightBlackFg":   "90",
	"BrightRedFg":     "91",
	"BrightGreenFg":   "92",
	"BrightYellowFg":  "93",
	"BrightBlueFg":    "94",
	"BrightMagentaFg": "95",
	"BrightCyanFg":    "96",
	"BrightWhiteFg":   "97",
	"BlackBg":         "40",
	"RedBg":           "41",
	"GreenBg":         "42",
	"YellowBg":        "43",
	"BlueBg":          "44",
	"MagentaBg":       "45",
	"CyanBg":          "46",
	"WhiteBg":         "47",
	"BrightBlackBg":   "100",
	"BrightRedBg":     "101",
	"BrightGreenBg":   "102",
	"BrightYellowBg":  "103",
	"BrightBlueBg":    "104",
	"BrightMagentaBg": "105",
	"BrightCyanBg":    "106",
	"BrightWhiteBg":   "107",
	"End":             "0",
}

// can receive more than one color, foreground color and background color
// cause color format depends on the color number, so the order is not important
// \033[30;45m is the same as \033[45;30m
func Colorize(text string, colorname ...string) string {
	var colortext string

	colorcode1, exists := color[colorname[0]]
	if !exists {
		io.WriteString(os.Stdout, "color name is not right\n")
		os.Exit(1)
	}
	colortext = fmt.Sprintf("\033[%sm%s\033[0m", colorcode1, text)

	if len(colorname) == 2 {
		colorcode2, exists := color[colorname[1]]
		if !exists {
			io.WriteString(os.Stdout, "color name is not right\n")
			os.Exit(1)
		}

		colortext = fmt.Sprintf("\033[%s;%sm%s\033[0m", colorcode1, colorcode2, text)
	}

	return colortext
}

func ColorizeRgbFg(rgb, text string) string {
	if len(rgb) != 7 || !strings.HasPrefix(rgb, "#") {
		io.WriteString(os.Stdout, "rgb color not right\n")
		os.Exit(1)
	}

	r, g, b := rgb[1:3], rgb[3:5], rgb[5:7]
	numr, _ := strconv.ParseUint(r, 16, 8)
	numg, _ := strconv.ParseUint(g, 16, 8)
	numb, _ := strconv.ParseUint(b, 16, 8)

	colorizeText := fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", numr, numg, numb, text)
	return colorizeText
}

func ColorizeRgbBg(rgb, text string) string {
	if len(rgb) != 7 || !strings.HasPrefix(rgb, "#") {
		io.WriteString(os.Stdout, "rgb color not right\n")
		os.Exit(1)
	}

	r, g, b := rgb[1:3], rgb[3:5], rgb[5:7]
	numr, _ := strconv.ParseUint(r, 16, 8)
	numg, _ := strconv.ParseUint(g, 16, 8)
	numb, _ := strconv.ParseUint(b, 16, 8)

	colorizeText := fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", numr, numg, numb, text)
	return colorizeText
}

// can receive forground color and background color but the order in important
// rgb1 = foreground color, rgb2 = background color
func ColorizeRgb(rgb1, rgb2, text string) string {
	if len(rgb1) != 7 || len(rgb2) != 7 || !strings.HasPrefix(rgb1, "#") || !strings.HasPrefix(rgb2, "#") {
		io.WriteString(os.Stdout, "rgb color not right\n")
		os.Exit(1)
	}

	r1, g1, b1 := rgb1[1:3], rgb1[3:5], rgb1[5:7]
	numr1, _ := strconv.ParseUint(r1, 16, 8)
	numg1, _ := strconv.ParseUint(g1, 16, 8)
	numb1, _ := strconv.ParseUint(b1, 16, 8)

	r2, g2, b2 := rgb2[1:3], rgb2[3:5], rgb2[5:7]
	numr2, _ := strconv.ParseUint(r2, 16, 8)
	numg2, _ := strconv.ParseUint(g2, 16, 8)
	numb2, _ := strconv.ParseUint(b2, 16, 8)

	colorizeText := fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s\033[0m", numr1, numg1, numb1, numr2, numg2, numb2, text)
	return colorizeText
}
