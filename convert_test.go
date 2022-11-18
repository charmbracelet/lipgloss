package lipgloss

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/pmezard/go-difflib/difflib"
)

func TestImport(t *testing.T) {
	emptyStyle := NewStyle()
	td := []struct {
		src    Style
		in     string
		out    string
		expErr string
	}{
		{emptyStyle, ``, ``, ``},
		{emptyStyle, `invalid`, ``, `invalid syntax: "invalid"`},
		{emptyStyle, `set-foo: bar`, ``, `in "set-foo: bar": don't use 'set-xx: foo;'  use 'xx: foo;' instead`},
		{emptyStyle, `unset-foo: bar`, ``, `in "unset-foo: bar": don't use 'unset-xx: foo;' use 'xx: unset;' instead`},
		{emptyStyle, `get-foo: bar`, ``, `in "get-foo: bar": property not supported: "get-foo"`},
		{emptyStyle, `unsupported: foo`, ``, `in "unsupported: foo": property not supported: "unsupported"`},
		{emptyStyle, `render: foo`, ``, `in "render: foo": method "Render" exists but does not return Style`},
		{emptyStyle.PaddingLeft(11), `padding-left:22`, `padding-left: 22;`, ``},
		{emptyStyle, `padding-left:aaa`, ``, `in "padding-left:aaa": no value found`},
		{emptyStyle, `padding-left:9999999999999999999999`, ``, `in "padding-left:9999999999999999999999": strconv.Atoi: parsing "9999999999999999999999": value out of range`},
		{emptyStyle, `bold: true`, `bold: true;`, ``},
		{emptyStyle, `bold: aa`, ``, `in "bold: aa": no value found`},
		{emptyStyle, `bold: true extra`, ``, `in "bold: true extra": excess values at end: ...extra`},
		{emptyStyle.Foreground(Color("11")), `foreground: unset`, ``, ``},
		{emptyStyle, `align-horizontal: left`, ``, ``},
		{emptyStyle, `align: left`, ``, ``},
		{emptyStyle, `align: xx`, ``, `in "align: xx": no value found`},
		{emptyStyle, `align: center`, `align-horizontal: 0.5;`, ``},
		{emptyStyle, `align: right`, `align-horizontal: 1;`, ``},
		{emptyStyle, `align: 1.0`, `align-horizontal: 1;`, ``},
		{emptyStyle, `align-horizontal: right`, `align-horizontal: 1;`, ``},
		{emptyStyle, `align-vertical: top`, ``, ``},
		{emptyStyle, `align-vertical: center`, `align-vertical: 0.5;`, ``},
		{emptyStyle, `align-vertical: bottom`, `align-vertical: 1;`, ``},
		{emptyStyle, `align: bottom right`, `align-horizontal: 1;
align-vertical: 1;`, ``},
		{emptyStyle.Foreground(Color("11")), `foreground: none`, ``, ``},
		{emptyStyle.Foreground(Color("11")), `clear`, ``, ``},
		{emptyStyle.Foreground(Color("11")), `background: 12; clear`, ``, ``},
		{emptyStyle, `foreground: 11`, `foreground: 11;`, ``},
		{emptyStyle, `foreground: #123`, `foreground: #123;`, ``},
		{emptyStyle, `foreground: #123456`, `foreground: #123456;`, ``},
		{emptyStyle, `foreground: #axxa`, ``, `in "foreground: #axxa": color not recognized`},
		{emptyStyle, `foreground: adaptive(1,2)`, `foreground: adaptive(1,2);`, ``},
		{emptyStyle, `foreground: complete(#111, 22, 3)`, `foreground: complete(#111,22,3);`, ``},
		{emptyStyle, `foreground: adaptive(complete(#111, 22, 3), complete(#444,55,6))`, `foreground: adaptive(complete(#111,22,3),complete(#444,55,6));`, ``},
		{emptyStyle, `foreground: adaptive(a,b)`, ``, `in "foreground: adaptive(a,b)": color not recognized: "a"`},
		{emptyStyle, `foreground: adaptive(1,b)`, ``, `in "foreground: adaptive(1,b)": color not recognized: "b"`},
		{emptyStyle, `foreground: complete(1,1,b)`, ``, `in "foreground: complete(1,1,b)": color not recognized: "b"`},
		{emptyStyle, `foreground: adaptive(complete(1,1,b),complete(2,2,b))`, ``, `in "foreground: adaptive(complete(1,1,b),complete(2,2,b))": color not recognized: "b"`},
		{emptyStyle, `margin: 10`, `margin-bottom: 10;
margin-left: 10;
margin-right: 10;
margin-top: 10;`, ``},
		{emptyStyle, `margin: 10 20`, `margin-bottom: 10;
margin-left: 20;
margin-right: 20;
margin-top: 10;`, ``},
		{emptyStyle, `margin: 10 20 30 40`, `margin-bottom: 30;
margin-left: 40;
margin-right: 20;
margin-top: 10;`, ``},
		{emptyStyle, `border-style: border("","","","","","","","")`, ``, ``},
		{emptyStyle, `border-style: xx`, ``, `in "border-style: xx": no valid border value found`},
		{emptyStyle,
			`border-style: border("a","b","c","d","e","f","g","h")`,
			`border-style: border("a","b","c","d","e","f","g","h");`, ``},
		{emptyStyle,
			`border-style: border("\"","\x41","\102","\u004a","\U0000004a","abc","a\"b","\\")`,
			`border-style: border("\"","A","B","J","J","abc","a\"b","\\");`, ``},
		{emptyStyle,
			`border: border("a","b","c","d","e","f","g","h") true false`,
			`border-bottom: true;
border-style: border("a","b","c","d","e","f","g","h");
border-top: true;`, ``},
		{emptyStyle,
			`border: border("a","b","c","d","e","f","g","h")`,
			`border-bottom: true;
border-left: true;
border-right: true;
border-style: border("a","b","c","d","e","f","g","h");
border-top: true;`,
			``},
		{emptyStyle,
			`border: border("a","b","c","d","e","f","g","h") true xx`,
			``,
			`in "border: border(\"a\",\"b\",\"c\",\"d\",\"e\",\"f\",\"g\",\"h\") true xx": no value found`},
		{emptyStyle,
			`border-style: rounded`,
			`border-style: border("─","─","│","│","╭","╮","╯","╰");`, ``},
		{emptyStyle,
			`border-style: normal`,
			`border-style: border("─","─","│","│","┌","┐","┘","└");`, ``},
		{emptyStyle,
			`border-style: thick`,
			`border-style: border("━","━","┃","┃","┏","┓","┛","┗");`, ``},
		{emptyStyle,
			`border-style: hidden`,
			`border-style: border(" "," "," "," "," "," "," "," ");`, ``},
		{emptyStyle,
			`border-style: double`,
			`border-style: border("═","═","║","║","╔","╗","╝","╚");`, ``},
		{emptyStyle,
			`border-style: block`,
			`border-style: border("█","█","█","█","█","█","█","█");`, ``},
		{emptyStyle,
			`border-style: inner-half-block`,
			`border-style: border("▄","▀","▐","▌","▗","▖","▘","▝");`, ``},
		{emptyStyle,
			`border-style: outer-half-block`,
			`border-style: border("▀","▄","▌","▐","▛","▜","▟","▙");`, ``},
	}

	for i, tc := range td {
		t.Run(fmt.Sprintf("%d: %s", i, tc.in), func(t *testing.T) {
			result, err := Import(tc.src, tc.in)
			if err != nil {
				if tc.expErr != "" {
					if err.Error() != tc.expErr {
						t.Fatalf("expected error:\n%q\ngot:\n%q", tc.expErr, err)
					}
					return
				} else {
					t.Fatal(err)
				}
			}
			if tc.expErr != "" {
				t.Fatalf("expected error %q, got no error", tc.expErr)
			}
			t.Logf("%# v", pretty.Formatter(result))
			actual := Export(result, WithSeparator("\n"))
			if actual != tc.out {
				expectedLines := difflib.SplitLines(tc.out)
				actualLines := difflib.SplitLines(actual)
				diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
					Context: 5,
					A:       expectedLines,
					B:       actualLines,
				})
				if err != nil {
					t.Fatal(err)
				}

				t.Fatalf("mismatch:\n%s\ndiff:\n%s", actual, diff)
			}
		})
	}
}

func TestExport(t *testing.T) {
	style := NewStyle().
		Bold(true).
		Align(Center).
		Foreground(Color("#FAFAFA")).
		Background(AdaptiveColor{Light: "#7D56F4", Dark: "#112233"}).
		BorderTopForeground(Color("12")).
		BorderStyle(RoundedBorder()).
		PaddingTop(2).
		PaddingLeft(4).
		Width(22)

	t.Run("shortened", func(t *testing.T) {
		exp := `align-horizontal: 0.5;
background: adaptive(#7D56F4,#112233);
bold: true;
border-style: border("─","─","│","│","╭","╮","╯","╰");
border-top-foreground: 12;
foreground: #FAFAFA;
padding-left: 4;
padding-top: 2;
width: 22;`
		result := Export(style, WithSeparator("\n"))
		if result != exp {
			expectedLines := difflib.SplitLines(exp)
			actualLines := difflib.SplitLines(result)
			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				Context: 5,
				A:       expectedLines,
				B:       actualLines,
			})
			if err != nil {
				t.Fatal(err)
			}

			t.Errorf("mismatch:\n%s\ndiff:\n%s", result, diff)
		}
	})

	t.Run("full", func(t *testing.T) {
		exp := `align-horizontal: 0.5;
align-vertical: 0;
background: adaptive(#7D56F4,#112233);
blink: false;
bold: true;
border-bottom: false;
border-bottom-background: none;
border-bottom-foreground: none;
border-left: false;
border-left-background: none;
border-left-foreground: none;
border-right: false;
border-right-background: none;
border-right-foreground: none;
border-style: border("─","─","│","│","╭","╮","╯","╰");
border-top: false;
border-top-background: none;
border-top-foreground: 12;
color-whitespace: false;
faint: false;
foreground: #FAFAFA;
height: 0;
inline: false;
italic: false;
margin-bottom: 0;
margin-left: 0;
margin-right: 0;
margin-top: 0;
max-height: 0;
max-width: 0;
padding-bottom: 0;
padding-left: 4;
padding-right: 0;
padding-top: 2;
reverse: false;
strikethrough: false;
strikethrough-spaces: false;
underline: false;
underline-spaces: false;
width: 22;`
		result := Export(style, WithExportDefaults(), WithSeparator("\n"))
		if result != exp {
			expectedLines := difflib.SplitLines(exp)
			actualLines := difflib.SplitLines(result)
			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				Context: 5,
				A:       expectedLines,
				B:       actualLines,
			})
			if err != nil {
				t.Fatal(err)
			}

			t.Errorf("mismatch:\n%s\ndiff:\n%s", result, diff)
		}
	})
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"hello-world", "HelloWorld"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			res := camelCase(tc.in)
			if res != tc.out {
				t.Errorf("expected %q, got %q", tc.out, res)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"Hello", "hello"},
		{"HelloWorld", "hello-world"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			res := snakeCase(tc.in)
			if res != tc.out {
				t.Errorf("expected %q, got %q", tc.out, res)
			}
		})
	}
}
