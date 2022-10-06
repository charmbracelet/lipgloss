package lipgloss

import (
	"io"

	"github.com/muesli/termenv"
)

var renderer = NewRenderer()

// Renderer is a lipgloss terminal renderer.
type Renderer struct {
	output            *termenv.Output
	hasDarkBackground *bool
}

// RendererOption is a function that can be used to configure a Renderer.
type RendererOption func(r *Renderer)

// DefaultRenderer returns the default renderer.
func DefaultRenderer() *Renderer {
	return renderer
}

// NewRenderer creates a new Renderer.
func NewRenderer(options ...RendererOption) *Renderer {
	r := &Renderer{
		output: termenv.DefaultOutput(),
	}
	for _, option := range options {
		option(r)
	}
	return r
}

// WithOutput sets the io.Writer to be used for rendering.
func WithOutput(w io.Writer) RendererOption {
	return WithTermenvOutput(termenv.NewOutput(w))
}

// WithTermenvOutput sets the termenv Output to use for rendering.
func WithTermenvOutput(output *termenv.Output) RendererOption {
	return func(r *Renderer) {
		r.output = output
	}
}

// WithDarkBackground can force the renderer to use a light/dark background.
func WithDarkBackground(dark bool) RendererOption {
	return func(r *Renderer) {
		r.SetHasDarkBackground(dark)
	}
}

// WithColorProfile sets the color profile on the renderer. This function is
// primarily intended for testing. For details, see the note on
// [Renderer.SetColorProfile].
func WithColorProfile(p termenv.Profile) RendererOption {
	return func(r *Renderer) {
		r.SetColorProfile(p)
	}
}

// ColorProfile returns the detected termenv color profile.
func (r *Renderer) ColorProfile() termenv.Profile {
	return r.output.Profile
}

// ColorProfile returns the detected termenv color profile.
func ColorProfile() termenv.Profile {
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
// termenv.Ascii (no color, 1-bit)
// termenv.ANSI (16 colors, 4-bit)
// termenv.ANSI256 (256 colors, 8-bit)
// termenv.TrueColor (16,777,216 colors, 24-bit)
//
// This function is thread-safe.
func (r *Renderer) SetColorProfile(p termenv.Profile) {
	r.output.Profile = p
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
// termenv.Ascii (no color, 1-bit)
// termenv.ANSI (16 colors, 4-bit)
// termenv.ANSI256 (256 colors, 8-bit)
// termenv.TrueColor (16,777,216 colors, 24-bit)
//
// This function is thread-safe.
func SetColorProfile(p termenv.Profile) {
	renderer.SetColorProfile(p)
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func HasDarkBackground() bool {
	return renderer.HasDarkBackground()
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func (r *Renderer) HasDarkBackground() bool {
	if r.hasDarkBackground != nil {
		return *r.hasDarkBackground
	}
	return r.output.HasDarkBackground()
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
	r.hasDarkBackground = &b
}
