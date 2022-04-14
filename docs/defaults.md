# Lipgloss Defaults


## Color Selection

In some cases, your lipgloss UI might not display colors. 
The reason this happens is because Lip Gloss automatically degrades colors to the best available option in the given terminal. 
For example, if you're running tests, they exist in a sub-process and are not attached to a TTY and thus Lip Gloss strips color output entirely.

However! You can force a color profile in your tests with SetColorProfile.

```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/muesli/termenv"
)

lipgloss.SetColorProfile(termenv.TrueColor)
```
