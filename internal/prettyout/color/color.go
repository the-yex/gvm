package color

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/the-yex/gvm/internal/prettyout/color/colorable"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2025/9/9 下午4:44
* @Package:
 */

var (
	NoColor = noColorIsSet() || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

	// Output defines the standard output of the print functions. By default,
	// os.Stdout is used.
	Output = colorable.NewColorableStdout()
)

// noColorIsSet returns true if the environment variable NO_COLOR is set to a non-empty string.
func noColorIsSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

// Color defines a custom color object which is defined by SGR parameters.
type Color struct {
	params  []Attribute
	noColor *bool
}

// Attribute defines a single SGR Code
type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	ResetBold Attribute = iota + 22
	ResetItalic
	ResetUnderline
	ResetBlinking
	_
	ResetReversed
	ResetConcealed
	ResetCrossedOut
)

var mapResetAttributes map[Attribute]Attribute = map[Attribute]Attribute{
	Bold:         ResetBold,
	Faint:        ResetBold,
	Italic:       ResetItalic,
	Underline:    ResetUnderline,
	BlinkSlow:    ResetBlinking,
	BlinkRapid:   ResetBlinking,
	ReverseVideo: ResetReversed,
	Concealed:    ResetConcealed,
	CrossedOut:   ResetCrossedOut,
}

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite

	// used internally for 256 and 24-bit coloring
	foreground
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite

	// used internally for 256 and 24-bit coloring
	background
)

// New returns a newly created color object.
func New(value ...Attribute) *Color {
	c := &Color{
		params: make([]Attribute, 0),
	}

	if noColorIsSet() {
		c.noColor = boolPtr(true)
	}

	c.Add(value...)
	return c
}

// RGB returns a new foreground color in 24-bit RGB.
func RGB(r, g, b int) *Color {
	return New(foreground, 2, Attribute(r), Attribute(g), Attribute(b))
}

// BgRGB returns a new background color in 24-bit RGB.
func BgRGB(r, g, b int) *Color {
	return New(background, 2, Attribute(r), Attribute(g), Attribute(b))
}

// AddRGB is used to chain foreground RGB SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: .Add(34, 0, 12).Add(255, 128, 0).
func (c *Color) AddRGB(r, g, b int) *Color {
	c.params = append(c.params, foreground, 2, Attribute(r), Attribute(g), Attribute(b))
	return c
}

// AddRGB is used to chain background RGB SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: .Add(34, 0, 12).Add(255, 128, 0).
func (c *Color) AddBgRGB(r, g, b int) *Color {
	c.params = append(c.params, background, 2, Attribute(r), Attribute(g), Attribute(b))
	return c
}

// Unset resets all escape attributes and clears the output. Usually should
// be called after Set().
func Unset() {
	if NoColor {
		return
	}

	fmt.Fprintf(Output, "%s[%dm", escape, Reset)
}

// Set sets the SGR sequence.
func (c *Color) Set() *Color {
	if c.isNoColorSet() {
		return c
	}

	fmt.Fprint(Output, c.format())
	return c
}

func (c *Color) unset() {
	if c.isNoColorSet() {
		return
	}

	Unset()
}

// SetWriter is used to set the SGR sequence with the given io.Writer. This is
// a low-level function, and users should use the higher-level functions, such
// as color.Fprint, color.Print, etc.
func (c *Color) SetWriter(w io.Writer) *Color {
	if c.isNoColorSet() {
		return c
	}

	fmt.Fprint(w, c.format())
	return c
}

// UnsetWriter resets all escape attributes and clears the output with the give
// io.Writer. Usually should be called after SetWriter().
func (c *Color) UnsetWriter(w io.Writer) {
	if c.isNoColorSet() {
		return
	}

	fmt.Fprintf(w, "%s[%dm", escape, Reset)
}

// Add is used to chain SGR parameters. Use as many as parameters to combine
// and create custom color objects. Example: Add(color.FgRed, color.Underline).
func (c *Color) Add(value ...Attribute) *Color {
	c.params = append(c.params, value...)
	return c
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	c.SetWriter(w)
	defer c.UnsetWriter(w)

	return fmt.Fprint(w, a...)
}

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a
// string. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// color.
func (c *Color) Print(a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprint(Output, a...)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	c.SetWriter(w)
	defer c.UnsetWriter(w)

	return fmt.Fprintf(w, format, a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
// This is the standard fmt.Printf() method wrapped with the given color.
func (c *Color) Printf(format string, a ...interface{}) (n int, err error) {
	c.Set()
	defer c.unset()

	return fmt.Fprintf(Output, format, a...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// On Windows, users should wrap w with colorable.NewColorable() if w is of
// type *os.File.
func (c *Color) Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(w, c.wrap(sprintln(a...)))
}

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered. This is the standard fmt.Print() method wrapped with the given
// color.
func (c *Color) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Output, c.wrap(sprintln(a...)))
}

// sequence returns a formatted SGR sequence to be plugged into a "\x1b[...m"
// an example output might be: "1;36" -> bold cyan
func (c *Color) sequence() string {
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(v))
	}

	return strings.Join(format, ";")
}

// wrap wraps the s string with the colors attributes. The string is ready to
// be printed.
func (c *Color) wrap(s string) string {
	if c.isNoColorSet() {
		return s
	}

	return c.format() + s + c.unformat()
}

func (c *Color) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

func (c *Color) unformat() string {
	//return fmt.Sprintf("%s[%dm", escape, Reset)
	//for each element in sequence let's use the specific reset escape, or the generic one if not found
	format := make([]string, len(c.params))
	for i, v := range c.params {
		format[i] = strconv.Itoa(int(Reset))
		ra, ok := mapResetAttributes[v]
		if ok {
			format[i] = strconv.Itoa(int(ra))
		}
	}

	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

func (c *Color) isNoColorSet() bool {
	// check first if we have user set action
	if c.noColor != nil {
		return *c.noColor
	}

	// if not return the global option, which is disabled by default
	return NoColor
}

func boolPtr(v bool) *bool {
	return &v
}

// sprintln is a helper function to format a string with fmt.Sprintln and trim the trailing newline.
func sprintln(a ...interface{}) string {
	return strings.TrimSuffix(fmt.Sprintln(a...), "\n")
}

// Sprintf is just like Printf, but returns a string instead of printing it.
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}
