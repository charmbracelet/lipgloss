package lipgloss

import (
	"os"
	"sync"

	"github.com/charmbracelet/x/exp/term"
	ansi "github.com/charmbracelet/x/exp/term/ansi/style"
	"github.com/muesli/termenv"
)

// We're manually creating the struct here to avoid initializing the output and
// query the terminal multiple times.
var renderer = NewRenderer()

// Renderer is a lipgloss terminal renderer.
type Renderer struct {
	environ           Environ
	hasDarkBackground *bool
	isatty            *bool
	colorProfile      Profile

	mtx sync.RWMutex
}

// RendererOption is a function that can be used to configure a [Renderer].
type RendererOption func(r *Renderer)

// DefaultRenderer returns the default renderer.
func DefaultRenderer() *Renderer {
	return renderer
}

// SetDefaultRenderer sets the default global renderer.
func SetDefaultRenderer(r *Renderer) {
	renderer = r
}

// WithRendererEnvironment sets the terminal environment on the renderer.
func WithRendererEnvironment(env Environ) RendererOption {
	return func(r *Renderer) {
		r.environ = env
	}
}

// WithRendererColorProfile sets the color profile on the renderer.
func WithRendererColorProfile(p Profile) RendererOption {
	return func(r *Renderer) {
		r.colorProfile = p
	}
}

// WithRendererDarkBackground sets the background color detection value on the renderer.
func WithRendererDarkBackground(b bool) RendererOption {
	return func(r *Renderer) {
		r.hasDarkBackground = &b
	}
}

// WithRendererTerminal sets whether the renderer should treat the output as a
// terminal.
func WithRendererTerminal(b bool) RendererOption {
	return func(r *Renderer) {
		r.isatty = &b
	}
}

// NewRenderer creates a new Renderer.
//
// w will be used to determine the terminal's color capabilities.
func NewRenderer(opts ...RendererOption) *Renderer {
	r := &Renderer{
		colorProfile: -1,
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.isatty == nil {
		isatty := term.IsTerminal(os.Stdout.Fd())
		r.isatty = &isatty
	}

	if r.hasDarkBackground == nil {
		t := termenv.HasDarkBackground()
		r.hasDarkBackground = &t
	}

	if r.environ == nil {
		r.environ = OsEnviron{}
	}

	if r.colorProfile < 0 {
		r.colorProfile = envColorProfile(*r.isatty, r.environ)
	}

	return r
}

// ColorProfile returns the detected termenv color profile.
func (r *Renderer) ColorProfile() Profile {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.colorProfile
}

// ColorProfile returns the detected termenv color profile.
func ColorProfile() Profile {
	return renderer.ColorProfile()
}

// SetColorProfile sets the color profile on the renderer. This function exists
// mostly for testing purposes so that you can assure you're testing against
// a specific profile.
//
// Outside of testing you likely won't want to use this function as the color
// profile will detect and cache the terminal's color capabilities and choose
// the best available profile.
//
// Available color profiles are:
//
//	Ascii     // no color, 1-bit
//	ANSI      //16 colors, 4-bit
//	ANSI256   // 256 colors, 8-bit
//	TrueColor // 16,777,216 colors, 24-bit
//
// This function is thread-safe.
func (r *Renderer) SetColorProfile(p Profile) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.colorProfile = p
}

// SetColorProfile sets the color profile on the default renderer. This
// function exists mostly for testing purposes so that you can assure you're
// testing against a specific profile.
//
// Outside of testing you likely won't want to use this function as the color
// profile will detect and cache the terminal's color capabilities and choose
// the best available profile.
//
// Available color profiles are:
//
//	Ascii     // no color, 1-bit
//	ANSI      //16 colors, 4-bit
//	ANSI256   // 256 colors, 8-bit
//	TrueColor // 16,777,216 colors, 24-bit
//
// This function is thread-safe.
func SetColorProfile(p Profile) {
	renderer.SetColorProfile(p)
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func HasDarkBackground() bool {
	return renderer.HasDarkBackground()
}

// HasDarkBackground returns whether or not the renderer will render to a dark
// background. A dark background can either be auto-detected, or set explicitly
// on the renderer.
func (r *Renderer) HasDarkBackground() bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if r.hasDarkBackground != nil {
		return *r.hasDarkBackground
	}

	return true
}

// SetHasDarkBackground sets the background color detection value for the
// default renderer. This function exists mostly for testing purposes so that
// you can assure you're testing against a specific background color setting.
//
// Outside of testing you likely won't want to use this function as the
// backgrounds value will be automatically detected and cached against the
// terminal's current background color setting.
//
// This function is thread-safe.
func SetHasDarkBackground(b bool) {
	renderer.SetHasDarkBackground(b)
}

// SetHasDarkBackground sets the background color detection value on the
// renderer. This function exists mostly for testing purposes so that you can
// assure you're testing against a specific background color setting.
//
// Outside of testing you likely won't want to use this function as the
// backgrounds value will be automatically detected and cached against the
// terminal's current background color setting.
//
// This function is thread-safe.
func (r *Renderer) SetHasDarkBackground(b bool) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.hasDarkBackground = &b
}

type style struct {
	attrs []ansi.Attribute
	p     Profile
}

func (s style) Styled(str string) string {
	if s.p <= Ascii || len(s.attrs) == 0 {
		return str
	}

	return ansi.String(str, s.attrs...)
}

func (s style) Foreground(c ansi.Color) style {
	if c != nil {
		s.attrs = append(s.attrs, ansi.ForegroundColor(c))
	}
	return s
}

func (s style) Background(c ansi.Color) style {
	if c != nil {
		s.attrs = append(s.attrs, ansi.BackgroundColor(c))
	}
	return s
}

func (s style) CrossOut() style {
	s.attrs = append(s.attrs, ansi.Strikethrough)
	return s
}

func (s style) Underline() style {
	s.attrs = append(s.attrs, ansi.Underline)
	return s
}

func (s style) Bold() style {
	s.attrs = append(s.attrs, ansi.Bold)
	return s
}

func (s style) Italic() style {
	s.attrs = append(s.attrs, ansi.Italic)
	return s
}

func (s style) Blink() style {
	s.attrs = append(s.attrs, ansi.SlowBlink)
	return s
}

func (s style) Faint() style {
	s.attrs = append(s.attrs, ansi.Faint)
	return s
}

func (s style) Reverse() style {
	s.attrs = append(s.attrs, ansi.Reverse)
	return s
}
