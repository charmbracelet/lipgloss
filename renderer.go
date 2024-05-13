package lipgloss

import (
	"os"
	"sync"

	"github.com/charmbracelet/x/exp/term"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	renderer     *Renderer
	rendererOnce sync.Once
)

// Renderer is a terminal style renderer that keep track of the color profile
// and background color detection.
type Renderer struct {
	colorProfile      Profile
	mtx               sync.RWMutex
	hasDarkBackground bool
}

// DefaultRenderer returns the default renderer.
func DefaultRenderer() *Renderer {
	rendererOnce.Do(func() {
		hasDarkBackground := true // Assume dark background by default
		if term.IsTerminal(os.Stdout.Fd()) {
			if bg := term.BackgroundColor(os.Stdin, os.Stdout); bg != nil {
				c, ok := colorful.MakeColor(bg)
				if ok {
					_, _, l := c.Hsl()
					hasDarkBackground = l < 0.5
				}
			}

			// Enable support for ANSI on the legacy Windows cmd.exe console. This is a
			// no-op on non-Windows systems and on Windows runs only once.
			// When using a custom renderer, this should be called manually.
			EnableLegacyWindowsANSI(os.Stdout)
		}

		cp := DetectColorProfile(os.Stdout, os.Environ())
		renderer = NewRenderer(cp, hasDarkBackground)
	})
	return renderer
}

// SetDefaultRenderer sets the default global renderer.
func SetDefaultRenderer(r *Renderer) {
	rendererOnce.Do(func() { renderer = r })
}

// NewRenderer creates a new Lip Gloss Renderer.
//
// It takes a color profile and a boolean indicating whether the terminal has a
// dark background.
// These values are then used to determine how to render colors and styles.
func NewRenderer(cp Profile, hasDarkBackground bool) *Renderer {
	return &Renderer{
		colorProfile:      cp,
		hasDarkBackground: hasDarkBackground,
	}
}

// ColorProfile returns the current color profile.
func (r *Renderer) ColorProfile() Profile {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.colorProfile
}

// ColorProfile returns the ren color profile.
func ColorProfile() Profile {
	return DefaultRenderer().ColorProfile()
}

// SetColorProfile sets the color profile on the renderer.
//
// Available color profiles are:
//
//	NoTTY     // no colors or styles
//	Ascii     // no colors, other styles are allowed (bold, italic, underline)
//	ANSI      // 16 colors, 4-bit
//	ANSI256   // 256 colors, 8-bit
//	TrueColor // 16,777,216 colors, 24-bit
//
// This function is thread-safe.
func (r *Renderer) SetColorProfile(p Profile) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.colorProfile = p
}

// SetColorProfile sets the color profile on the renderer.
//
// Available color profiles are:
//
//	NoTTY     // no colors or styles
//	Ascii     // no colors, other styles are allowed (bold, italic, underline)
//	ANSI      // 16 colors, 4-bit
//	ANSI256   // 256 colors, 8-bit
//	TrueColor // 16,777,216 colors, 24-bit
//
// This function is thread-safe.
func SetColorProfile(p Profile) {
	DefaultRenderer().SetColorProfile(p)
}

// HasDarkBackground returns whether or not the renderer has a dark background.
func HasDarkBackground() bool {
	return DefaultRenderer().HasDarkBackground()
}

// HasDarkBackground returns whether or not the renderer has a dark background.
func (r *Renderer) HasDarkBackground() bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.hasDarkBackground
}

// SetHasDarkBackground sets the background color detection value for the
// default renderer.
//
// This function is thread-safe.
func SetHasDarkBackground(b bool) {
	DefaultRenderer().SetHasDarkBackground(b)
}

// SetHasDarkBackground sets the background color detection value for the
// default renderer.
//
// This function is thread-safe.
func (r *Renderer) SetHasDarkBackground(b bool) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.hasDarkBackground = b
}
