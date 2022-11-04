package lipgloss

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Import reads style specifications from the input string
// and sets the corresponding properties in the dst style.
func Import(dst Style, input string) (Style, error) {
	// Syntax: semicolon-separated list of prop: values... pairs.
	assignments := strings.Split(input, ";")
	for _, a := range assignments {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}

		if a == "clear" {
			// Special keyword: reset style.
			dst = NewStyle()
			continue
		}

		pair := strings.SplitN(a, ":", 2)
		if len(pair) != 2 {
			return dst, fmt.Errorf("invalid syntax: %q", a)
		}
		propName, args := pair[0], pair[1]
		propName = strings.TrimSpace(propName)
		args = strings.TrimSpace(args)
		p, err := getProp(propName)
		if err != nil {
			return dst, fmt.Errorf("in %q: %v", a, err)
		}

		dst, err = p.assign(dst, args)
		if err != nil {
			return dst, fmt.Errorf("in %q: %v", a, err)
		}
	}
	return dst, nil
}

type options struct {
	includeDefaults bool
	sep             string
}

// ExportOption represents an option for the Export function.
type ExportOption func(*options)

// WithSeparator sets the separator between directives.
func WithSeparator(sep string) ExportOption {
	return func(e *options) {
		e.sep = sep
	}
}

// WithExportDefaults includes the fields that are set to default values.
func WithExportDefaults() ExportOption {
	return func(e *options) {
		e.includeDefaults = true
	}
}

// Export emits style specifications that represent
// the given style.
// If includeDefaults is set, all the fields set to
// default values are also included in the output.
func Export(s Style, opts ...ExportOption) string {
	opt := options{
		sep: " ",
	}
	for _, o := range opts {
		o(&opt)
	}

	var buf strings.Builder

	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !strings.HasPrefix(m.Name, "Get") {
			continue
		}
		if ignoredMethods[m.Name] {
			continue
		}
		if m.Type.NumIn() != 1 {
			// Method with parameters; not truly a Getter. Ignore.
			continue
		}

		res := m.Func.Call([]reflect.Value{v})

		if !opt.includeDefaults && len(res) == 1 && isDefault(res[0]) {
			// Default value. Don't report anything for this getter.
			continue
		}

		if buf.Len() > 0 {
			buf.WriteString(opt.sep)
		}
		buf.WriteString(snakeCase(strings.TrimPrefix(m.Name, "Get")))
		buf.WriteString(": ")
		for j, v := range res {
			if j > 0 {
				buf.WriteByte(' ')
			}
			printValue(&buf, v)
		}
		buf.WriteByte(';')
	}
	return buf.String()
}

func printValue(buf *strings.Builder, v reflect.Value) {
	switch v.Type().Name() {
	case "TerminalColor":
		tc := v.Interface().(TerminalColor)
		switch c := tc.(type) {
		case NoColor:
			buf.WriteString("none")
		case Color:
			buf.WriteString(string(c))
		case AdaptiveColor:
			fmt.Fprintf(buf, "adaptive(%s,%s)", c.Light, c.Dark)
		case CompleteColor:
			fmt.Fprintf(buf, "complete(%s,%s,%s)", c.TrueColor, c.ANSI256, c.ANSI)
		case CompleteAdaptiveColor:
			fmt.Fprintf(buf, "adaptive(complete(%s,%s,%s),complete(%s,%s,%s))",
				c.Light.TrueColor, c.Light.ANSI256, c.Light.ANSI,
				c.Dark.TrueColor, c.Dark.ANSI256, c.Dark.ANSI,
			)
		default:
			r, g, b, _ := tc.RGBA()
			fmt.Fprintf(buf, "#%02x%02x%02x", r, g, b)
		}
	case "Border":
		b := v.Interface().(Border)
		fmt.Fprintf(buf, "border(%q,%q,%q,%q,%q,%q,%q,%q)",
			b.Top, b.Bottom, b.Left, b.Right,
			b.TopLeft, b.TopRight, b.BottomRight, b.BottomLeft,
		)
	default:
		fmt.Fprintf(buf, "%v", v.Interface())
	}
}

func isDefault(v reflect.Value) bool {
	if v.IsZero() {
		return true
	}
	switch v.Type().Name() {
	case "TerminalColor":
		color := v.Interface().(TerminalColor)
		_, isNoColor := color.(NoColor)
		return isNoColor
	default:
		// Unknown type. Always include in output.
		return false
	}
}

var ignoredMethods = map[string]bool{
	"GetAlign":                true,
	"GetBorder":               true,
	"GetMargin":               true,
	"GetPadding":              true,
	"GetFrameSize":            true,
	"GetVerticalFrameSize":    true,
	"GetHorizontalFrameSize":  true,
	"GetVerticalMargins":      true,
	"GetHorizontalMargins":    true,
	"GetVerticalPadding":      true,
	"GetHorizontalPadding":    true,
	"GetVerticalBorderSize":   true,
	"GetHorizontalBorderSize": true,
	"GetBorderBottomSize":     true,
	"GetBorderTopSize":        true,
	"GetBorderBottomWidth":    true, // facepalm
	"GetBorderTopWidth":       true, // facepalm
	"GetBorderLeftSize":       true,
	"GetBorderRightSize":      true,
}

func getProp(name string) (prop, error) {
	p, ok := propRegistry[name]
	if !ok {
		var err error
		p, err = discoverProp(name)
		if err != nil {
			return prop{}, err
		}
	}
	return p, nil
}

func discoverProp(name string) (prop, error) {
	if strings.HasPrefix(name, "set-") {
		return prop{}, fmt.Errorf("don't use 'set-xx: foo;'  use 'xx: foo;' instead")
	}
	if strings.HasPrefix(name, "unset-") {
		return prop{}, fmt.Errorf("don't use 'unset-xx: foo;' use 'xx: unset;' instead")
	}
	if strings.HasPrefix(name, "get-") {
		return prop{}, fmt.Errorf("property not supported: %q", name)
	}
	propName := name
	name = camelCase(name)
	s := reflect.ValueOf(NewStyle())
	t := s.Type()

	m, hasMethod := t.MethodByName(name)
	if !hasMethod {
		return prop{}, fmt.Errorf("property not supported: %q", propName)
	}
	if m.Type.NumOut() != 1 || m.Type.Out(0) != styleType {
		return prop{}, fmt.Errorf("method %q exists but does not return Style", name)
	}

	var args []argtype
	for i := 1; i < m.Type.NumIn(); i++ {
		argT := m.Type.In(i)

		if m.Type.IsVariadic() && i == m.Type.NumIn()-1 {
			argT = argT.Elem()
		}

		switch {
		case argT.Kind() == reflect.Int:
			args = append(args, inttype{})
		case argT.Kind() == reflect.Bool:
			args = append(args, booltype{})
		case argT.Name() == "Border":
			args = append(args, bordertype{})
		case argT.Name() == "Position":
			args = append(args, postype{})
		case argT.Name() == "TerminalColor":
			args = append(args, colortype{})
		default:
			return prop{}, fmt.Errorf("Style has method %s, but method uses unsupported argument type %s", name, argT)
		}
	}
	p := prop{
		setFn:      m.Func,
		isVariadic: m.Type.IsVariadic(),
		args:       args,
	}

	if um, hasUnsetMethod := t.MethodByName("Unset" + name); hasUnsetMethod &&
		m.Type.NumOut() == 1 && m.Type.Out(0) == styleType {
		p.unsetFn = um.Func
	}

	return p, nil
}

var styleType = reflect.TypeOf(NewStyle())

type argtype interface {
	parse([]byte, int) (int, reflect.Value, error)
}

type inttype struct{}

func (inttype) parse(input []byte, first int) (pos int, val reflect.Value, err error) {
	pos = first
	r := reInt.FindSubmatch(input[pos:])
	if r == nil {
		return pos, val, fmt.Errorf("no value found")
	}
	pos += len(r[0])
	i, err := strconv.Atoi(string(r[1]))
	if err != nil {
		return pos, val, err
	}
	return pos, reflect.ValueOf(i), nil
}

var reInt = regexp.MustCompile(`^\s*([0-9]+)(?:\s+|$)`)

type booltype struct{}

func (booltype) parse(input []byte, first int) (pos int, val reflect.Value, err error) {
	pos = first
	r := reBool.FindSubmatch(input[pos:])
	if r == nil {
		return pos, val, fmt.Errorf("no value found")
	}
	pos += len(r[0])
	b, err := strconv.ParseBool(string(r[1]))
	if err != nil {
		return pos, val, err
	}
	return pos, reflect.ValueOf(b), nil
}

var reBool = regexp.MustCompile(`^\s*(1|[tT]|TRUE|[tT]rue|0|[fF]|FALSE|[fF]alse)(?:\s+|$)`)

type postype struct{}

func (postype) parse(input []byte, first int) (pos int, val reflect.Value, err error) {
	pos = first
	r := rePos.FindSubmatch(input[pos:])
	if r == nil {
		return pos, val, fmt.Errorf("no value found")
	}
	pos += len(r[0])
	word := string(r[1])
	switch word {
	case "top":
		val = reflect.ValueOf(Top)
	case "bottom":
		val = reflect.ValueOf(Bottom)
	case "center":
		val = reflect.ValueOf(Center)
	case "left":
		val = reflect.ValueOf(Left)
	case "right":
		val = reflect.ValueOf(Right)
	default:
		p, err := strconv.ParseFloat(word, 64)
		if err != nil {
			return pos, val, err
		}
		position := Position(p)
		val = reflect.ValueOf(position)
	}
	return pos, val, nil
}

var rePos = regexp.MustCompile(`^\s*(top|bottom|center|left|right|1|1\.0|0\.5|\.5|0|0\.0|\.0)(?:\s+|$)`)

type colortype struct{}

func getColors(rematch [][]byte, cvals []string) error {
	for i := 0; i < len(cvals); i++ {
		val := strings.TrimSpace(string(rematch[i+1]))
		if !reColor.MatchString(val) {
			return fmt.Errorf("color not recognized: %q", val)
		}
		cvals[i] = val
	}
	return nil
}

func (colortype) parse(input []byte, first int) (pos int, val reflect.Value, err error) {
	pos = first
	// possible syntaxes:
	// - adaptive(X, Y)
	// - complete(X, Y, Z)
	// - adaptive(complete(A,B,C), complete(D,E,F))
	// - one word, either "none", just a number or a RGB value
	if r := reAdaptive.FindSubmatch(input[pos:]); r != nil {
		pos += len(r[0])
		var cvals [2]string
		if err := getColors(r, cvals[:]); err != nil {
			return pos, val, err
		}
		c := AdaptiveColor{Light: cvals[0], Dark: cvals[1]}
		val = reflect.ValueOf(c)
		return pos, val, nil
	}
	if r := reComplete.FindSubmatch(input[pos:]); r != nil {
		pos += len(r[0])
		var cvals [3]string
		if err := getColors(r, cvals[:]); err != nil {
			return pos, val, err
		}
		c := CompleteColor{TrueColor: cvals[0], ANSI256: cvals[1], ANSI: cvals[2]}
		val = reflect.ValueOf(c)
		return pos, val, nil
	}
	if r := reCompleteAdaptive.FindSubmatch(input[pos:]); r != nil {
		pos += len(r[0])
		var cvals [6]string
		if err := getColors(r, cvals[:]); err != nil {
			return pos, val, err
		}
		c := CompleteAdaptiveColor{
			Light: CompleteColor{TrueColor: cvals[0], ANSI256: cvals[1], ANSI: cvals[2]},
			Dark:  CompleteColor{TrueColor: cvals[3], ANSI256: cvals[4], ANSI: cvals[5]},
		}
		val = reflect.ValueOf(c)
		return pos, val, nil
	}

	r := reColorOrNone.FindSubmatch(input[pos:])
	if r == nil {
		return pos, val, fmt.Errorf("color not recognized")
	}
	pos += len(r[0])
	word := string(r[1])
	switch word {
	case "none":
		val = reflect.ValueOf(NoColor{})
	default:
		if !reColor.MatchString(word) {
			return pos, val, fmt.Errorf("color not recognized: %q", word)
		}
		val = reflect.ValueOf(Color(word))
	}
	return pos, val, nil
}

var reColor = regexp.MustCompile(`^\s*(\d+|#[0-9a-fA-F]{3}|#[0-9a-fA-F]{6})(?:\s+|$)`)
var reColorOrNone = regexp.MustCompile(`^\s*(none|\d+|#[0-9a-fA-F]{3}|#[0-9a-fA-F]{6})(?:\s+|$)`)
var reAdaptive = regexp.MustCompile(`^\s*(?:adaptive\s*\(([^,]*),([^,]*)\))(?:\s+|$)`)

var reComplete = regexp.MustCompile(`^\s*(?:complete\s*\(([^,]*),([^,]*),([^,]*)\))(?:\s+|$)`)

var reCompleteAdaptive = regexp.MustCompile(`^\s*(?:` +
	`adaptive\s*\(\s*` +
	`complete\s*\(([^,]*),([^,]*),([^,]*)\)` +
	`\s*,\s*` +
	`complete\s*\(([^,]*),([^,]*),([^,]*)\)` +
	`\))(?:\s+|$)`)

type bordertype struct{}

func (bordertype) parse(input []byte, first int) (pos int, val reflect.Value, err error) {
	pos = first
	if r := reSpecialBorder.FindSubmatch(input[pos:]); r != nil {
		pos += len(r[0])
		word := string(r[1])
		var b Border
		switch word {
		case "rounded":
			b = RoundedBorder()
		case "normal":
			b = NormalBorder()
		case "thick":
			b = ThickBorder()
		case "hidden":
			b = HiddenBorder()
		case "double":
			b = DoubleBorder()
		case "block":
			b = BlockBorder()
		case "inner-half-block":
			b = InnerHalfBlockBorder()
		case "outer-half-block":
			b = OuterHalfBlockBorder()
		default:
			return pos, val, fmt.Errorf("unrecognized border name: %q", word)
		}
		return pos, reflect.ValueOf(b), nil
	}
	r := reBorder.FindSubmatch(input[pos:])
	if r == nil {
		return pos, val, fmt.Errorf("no valid border value found")
	}
	pos += len(r[0])
	var b Border
	for i, field := range []*string{
		&b.Top, &b.Bottom, &b.Left, &b.Right,
		&b.TopLeft, &b.TopRight, &b.BottomRight, &b.BottomLeft,
	} {
		word := string(r[i+1])
		word, err := strconv.Unquote(word)
		if err != nil {
			return pos, val, err
		}
		*field = word
	}
	val = reflect.ValueOf(b)
	return pos, val, nil
}

// Example valid border strings:
// "h", "|", etc
// "\"" - the character '"' itself
// "\\" - the character '\' itself
// "\012" - a octal-encoded ascii value
// "\xFF" - a hex-encoded ascii value
// "\u1234" - a hex-encoded rune
// "\U12345678" - a hex-encoded rune
var reBorderStr = `"(?:\\[\\"]|\\[0-7]{3}|\\x[0-9a-fA-F]{2}|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}|[^\\"])*"`

var reBorder = regexp.MustCompile(`^\s*(?:border\s*\(\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*,\s*(` +
	reBorderStr + `)\s*\))(?:\s+|$)`)

var reSpecialBorder = regexp.MustCompile(`^\s*(rounded|normal|thick|hidden|double|block|inner-half-block|outer-half-block)(?:\s+|$)`)

// camelCase converts hello-world to HelloWorld.
func camelCase(s string) string {
	var buf strings.Builder
	cap := true
	for _, r := range s {
		if r == '-' {
			cap = true
			continue
		}
		if cap {
			r = unicode.ToUpper(r)
			cap = false
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

// snakeCase converts HelloWorld to hello-world.
func snakeCase(s string) string {
	var buf strings.Builder
	for _, r := range s {
		if unicode.IsUpper(r) {
			if buf.Len() > 0 {
				buf.WriteByte('-')
			}
			r = unicode.ToLower(r)
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

var propRegistry = map[string]prop{}

type prop struct {
	setFn      reflect.Value
	unsetFn    reflect.Value
	isVariadic bool
	args       []argtype
}

func (p prop) assign(dst Style, args string) (Style, error) {
	if args == "unset" {
		// Special keyword.
		var noValue reflect.Value
		if p.unsetFn == noValue {
			return dst, fmt.Errorf("no unset method defined")
		}
		out := p.unsetFn.Call([]reflect.Value{reflect.ValueOf(dst)})
		return out[0].Interface().(Style), nil
	}

	// Read the arguments from the input string.
	vals := make([]reflect.Value, 0, 1+len(p.args))
	vals = append(vals, reflect.ValueOf(dst))
	pos := 0
	input := []byte(args)
	for i, arg := range p.args {
		if pos >= len(input) {
			if p.isVariadic && i == len(p.args)-1 {
				// It's ok for a variadic arg list to have zero argument.
				break
			}
			return dst, fmt.Errorf("missing value")
		}
		var err error
		var val reflect.Value
		pos, val, err = arg.parse(input, pos)
		if err != nil {
			return dst, err
		}
		vals = append(vals, val)
	}
	if p.isVariadic {
		for pos < len(input) {
			var val reflect.Value
			var err error
			pos, val, err = p.args[len(p.args)-1].parse(input, pos)
			if err != nil {
				return dst, err
			}
			vals = append(vals, val)
		}
	}
	if pos < len(input) {
		return dst, fmt.Errorf("excess values at end: ...%s", string(input[pos:]))
	}

	// Finally call the setter.
	out := p.setFn.Call(vals)
	return out[0].Interface().(Style), nil
}
