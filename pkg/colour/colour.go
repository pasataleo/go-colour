package colour

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	Reset         = "\033[0m"
	Bold          = "\033[1m"
	Italic        = "\033[3m"
	Underline     = "\033[4m"
	Strikethrough = "\033[9m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Default = "\033[39m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
	BgDefault = "\033[49m"

	Grey          = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	BgGrey          = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"
)

var (
	rgbPattern   = regexp.MustCompile(`\{rgb:(\d{1,3}),(\d{1,3}),(\d{1,3})\}`)
	bgRgbPattern = regexp.MustCompile(`\{bg-rgb:(\d{1,3}),(\d{1,3}),(\d{1,3})\}`)
)

// RGB returns the ANSI escape code for a 24-bit foreground colour.
func RGB(r, g, b uint8) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// BgRGB returns the ANSI escape code for a 24-bit background colour.
func BgRGB(r, g, b uint8) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

var (
	defaultColours = map[string]string{
		"{reset}":         Reset,
		"{bold}":          Bold,
		"{italic}":        Italic,
		"{underline}":     Underline,
		"{strikethrough}": Strikethrough,

		"{black}":   Black,
		"{red}":     Red,
		"{green}":   Green,
		"{yellow}":  Yellow,
		"{blue}":    Blue,
		"{magenta}": Magenta,
		"{cyan}":    Cyan,
		"{white}":   White,
		"{default}": Default,

		"{bg-black}":   BgBlack,
		"{bg-red}":     BgRed,
		"{bg-green}":   BgGreen,
		"{bg-yellow}":  BgYellow,
		"{bg-blue}":    BgBlue,
		"{bg-magenta}": BgMagenta,
		"{bg-cyan}":    BgCyan,
		"{bg-white}":   BgWhite,
		"{bg-default}": BgDefault,

		"{grey}":           Grey,
		"{bright-red}":     BrightRed,
		"{bright-green}":   BrightGreen,
		"{bright-yellow}":  BrightYellow,
		"{bright-blue}":    BrightBlue,
		"{bright-magenta}": BrightMagenta,
		"{bright-cyan}":    BrightCyan,
		"{bright-white}":   BrightWhite,

		"{bg-grey}":           BgGrey,
		"{bg-bright-red}":     BgBrightRed,
		"{bg-bright-green}":   BgBrightGreen,
		"{bg-bright-yellow}":  BgBrightYellow,
		"{bg-bright-blue}":    BgBrightBlue,
		"{bg-bright-magenta}": BgBrightMagenta,
		"{bg-bright-cyan}":    BgBrightCyan,
		"{bg-bright-white}":   BgBrightWhite,
	}
)

type Colour struct {
	Colours map[string]string
	Disable bool
}

type Opt func(Colour) Colour

func WithColour(colour, code string) Opt {
	return func(c Colour) Colour {
		c.Colours[colour] = code
		return c
	}
}

func WithColours(colours map[string]string) Opt {
	return func(c Colour) Colour {
		maps.Copy(c.Colours, colours)
		return c
	}
}

func Enable() Opt {
	return func(c Colour) Colour {
		c.Disable = false
		return c
	}
}

func Disable() Opt {
	return func(c Colour) Colour {
		c.Disable = true
		return c
	}
}

func New(opts ...Opt) Colour {
	_, noColor := os.LookupEnv("NO_COLOR")
	colour := Colour{
		Colours: make(map[string]string),
		Disable: noColor,
	}
	maps.Copy(colour.Colours, defaultColours)
	for _, opt := range opts {
		colour = opt(colour)
	}
	return colour
}

func (c Colour) Colour(input string) string {
	for colour, code := range c.Colours {
		if c.Disable {
			input = strings.ReplaceAll(input, colour, "")
		} else {
			input = strings.ReplaceAll(input, colour, code)
		}
	}

	input = rgbPattern.ReplaceAllStringFunc(input, func(match string) string {
		if c.Disable {
			return ""
		}
		parts := rgbPattern.FindStringSubmatch(match)
		r, g, b := parseRGB(parts[1], parts[2], parts[3])
		return RGB(r, g, b)
	})

	input = bgRgbPattern.ReplaceAllStringFunc(input, func(match string) string {
		if c.Disable {
			return ""
		}
		parts := bgRgbPattern.FindStringSubmatch(match)
		r, g, b := parseRGB(parts[1], parts[2], parts[3])
		return BgRGB(r, g, b)
	})

	return input
}

func (c Colour) Colourf(format string, args ...any) string {
	return fmt.Sprintf(c.Colour(format), args...)
}

func parseRGB(rs, gs, bs string) (uint8, uint8, uint8) {
	r, err := strconv.ParseUint(rs, 10, 8)
	if err != nil {
		panic(fmt.Sprintf("invalid RGB red value %q: %s", rs, err))
	}
	g, err := strconv.ParseUint(gs, 10, 8)
	if err != nil {
		panic(fmt.Sprintf("invalid RGB green value %q: %s", gs, err))
	}
	b, err := strconv.ParseUint(bs, 10, 8)
	if err != nil {
		panic(fmt.Sprintf("invalid RGB blue value %q: %s", bs, err))
	}
	return uint8(r), uint8(g), uint8(b)
}
