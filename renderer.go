package lipgloss

import (
	"os"
	"sync"

	"github.com/charmbracelet/x/exp/term"
	"github.com/lucasb-eyer/go-colorful"
)

// We're manually creating the struct here to avoid initializing the output and
// query the terminal multiple times.
var (
	renderer     *Renderer
	rendererOnce sync.Once
)

// Renderer is a lipgloss terminal renderer.
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

		cp := EnvColorProfile(os.Stdout, os.Environ())
		renderer = NewRenderer(cp, hasDarkBackground)
	})
	return renderer
}

// SetDefaultRenderer sets the default global renderer.
func SetDefaultRenderer(r *Renderer) {
	rendererOnce.Do(func() { renderer = r })
}

// NewRenderer creates a new Renderer.
//
// The stdout argument is used to detect if the renderer is writing to a
// terminal. If it is nil, the renderer will assume it's not writing to a
// terminal.
// The environ argument is used to detect the color profile based on the
// environment variables. If it's nil, os.Environ() will be used.
// Set hasDarkBackground to true if the terminal has a dark background.
func NewRenderer(cp Profile, hasDarkBackground bool) *Renderer {
	return &Renderer{
		colorProfile:      cp,
		hasDarkBackground: hasDarkBackground,
	}
}

// ColorProfile returns the detected color profile.
func (r *Renderer) ColorProfile() Profile {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.colorProfile
}

// ColorProfile returns the detected color profile.
func ColorProfile() Profile {
	return DefaultRenderer().ColorProfile()
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
	DefaultRenderer().SetColorProfile(p)
}

// HasDarkBackground returns whether or not the terminal has a dark background.
func HasDarkBackground() bool {
	return DefaultRenderer().HasDarkBackground()
}

// HasDarkBackground returns whether or not the renderer will render to a dark
// background. A dark background can either be auto-detected, or set explicitly
// on the renderer.
func (r *Renderer) HasDarkBackground() bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.hasDarkBackground
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
	DefaultRenderer().SetHasDarkBackground(b)
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

	r.hasDarkBackground = b
}
