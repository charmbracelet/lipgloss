package lipgloss

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/colorprofile"
)

func TestWriter(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*bytes.Buffer)
		want string
	}{
		{
			name: "Fprint",
			fn: func(buf *bytes.Buffer) {
				_, _ = Fprint(buf, "hello")
			},
			want: "hello",
		},
		{
			name: "Fprintln",
			fn: func(buf *bytes.Buffer) {
				_, _ = Fprintln(buf, "hello")
			},
			want: "hello\n",
		},
		{
			name: "Fprintf",
			fn: func(buf *bytes.Buffer) {
				_, _ = Fprintf(buf, "%s %d", "hi", 1)
			},
			want: "hi 1",
		},
		{
			name: "Sprint",
			fn: func(buf *bytes.Buffer) {
				s := Sprint("world")
				buf.WriteString(s)
			},
			want: "world",
		},
		{
			name: "Sprintln",
			fn: func(buf *bytes.Buffer) {
				s := Sprintln("world")
				buf.WriteString(s)
			},
			want: "world\n",
		},
		{
			name: "Sprintf",
			fn: func(buf *bytes.Buffer) {
				s := Sprintf("val=%v", 42)
				buf.WriteString(s)
			},
			want: "val=42",
		},
		{
			name: "Print",
			fn: func(buf *bytes.Buffer) {
				old := Writer
				Writer = colorprofile.NewWriter(buf, nil)
				defer func() { Writer = old }()
				_, _ = Print("hello")
			},
			want: "hello",
		},
		{
			name: "Println",
			fn: func(buf *bytes.Buffer) {
				old := Writer
				Writer = colorprofile.NewWriter(buf, nil)
				defer func() { Writer = old }()
				_, _ = Println("hello")
			},
			want: "hello\n",
		},
		{
			name: "Printf",
			fn: func(buf *bytes.Buffer) {
				old := Writer
				Writer = colorprofile.NewWriter(buf, nil)
				defer func() { Writer = old }()
				_, _ = Printf("n=%d", 7)
			},
			want: "n=7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.fn(&buf)
			if got := buf.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
