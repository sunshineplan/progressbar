package progressbar

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/sunshineplan/utils/unit"
)

// Frame represents a snapshot of all display-ready fields used to render the
// progress bar at a given moment. It contains the computed visual components
// (e.g., completed blocks, percentages, speed, timing information) that are
// assembled by the renderer into the final output.
type Frame struct {
	Unit           string
	Percent        float64
	Done, Undone   string
	Speed          float64
	Current, Total int64
	Additional     string
	Elapsed        time.Duration
	Left           time.Duration
}

var defaultTemplate *template.Template

var (
	dot  atomic.Int64
	dots = []string{".  ", ".. ", "..."}
)

var defaultFuncMap = template.FuncMap{
	"speed": func(speed float64, t string) string {
		if speed == 0 {
			return "--/s"
		} else if t == "bytes" {
			return fmt.Sprintf("%s/s", unit.ByteSize(speed))
		} else {
			return fmt.Sprintf("%.2f/s", speed)
		}
	},
	"count": func(n int64, t string) string {
		if t == "bytes" {
			return unit.ByteSize(n).String()
		} else {
			return strconv.FormatInt(n, 10)
		}
	},
	"calc": func() string {
		return "calculating" + dots[dot.Add(1)%3]
	},
}

func init() {
	t := template.New("")
	t.Funcs(defaultFuncMap)
	defaultTemplate =
		template.Must(
			template.Must(
				template.Must(
					template.Must(
						t.New("full").Parse(
							`[{{.Done}}{{.Undone}}]  {{speed .Speed .Unit}}  {{count .Current .Unit}}({{printf "%.2f%%" .Percent}}) of {{count .Total .Unit}}{{if .Additional}} [{{.Additional}}]{{end}}  Elapsed: {{.Elapsed }}  {{if eq .Current .Total}}Complete{{else}}Left: {{if eq .Speed 0.0}}{{calc}}{{else}}{{.Left}}{{end}}{{end}} `,
						),
					).New("standard").Parse(
						`[{{.Done}}{{.Undone}}] {{speed .Speed .Unit}} {{count .Current .Unit}}/{{count .Total .Unit}}({{printf "%.1f%%" .Percent}}){{if .Additional}} [{{.Additional}}]{{end}} ET: {{.Elapsed}} {{if eq .Current .Total}}Done{{else}}LT: {{if eq .Speed 0.0}}{{calc}}{{else}}{{.Left}}{{end}}{{end}} `,
					),
				).New("lite").Parse(
					`[{{.Done}}{{.Undone}}] {{speed .Speed .Unit}} {{count .Current .Unit}}/{{count .Total .Unit}}{{if .Additional}} [{{.Additional}}]{{end}} {{if eq .Current .Total}}E: {{.Elapsed}}{{else}}L: {{if eq .Speed 0.0}}{{calc}}{{else}}{{.Left}}{{end}}{{end}} `,
				),
			).New("mini").Parse(
				`[{{.Done}}{{.Undone}}] {{.Additional}} `,
			),
		)
}

// DefaultRenderFunc is the default function used to render the progress bar.
var DefaultRenderFunc = func(w io.Writer, f Frame) {
	winsize := GetWinsize()
	n := winsize - len(f.Done) - len(f.Undone) - len(f.Additional)
	switch {
	case n >= 60:
		defaultTemplate.ExecuteTemplate(w, "full", f)
	case n >= 40:
		defaultTemplate.ExecuteTemplate(w, "standard", f)
	case n >= 20:
		defaultTemplate.ExecuteTemplate(w, "lite", f)
	case n > 5 && len(f.Additional) > 0:
		defaultTemplate.ExecuteTemplate(w, "mini", f)
	default:
		width := winsize - 5
		w.Write([]byte("["))
		done := int(float64(width) * f.Percent / 100)
		if f.Percent == 100 {
			w.Write([]byte(strings.Repeat("=", done)))
		} else {
			w.Write([]byte(strings.Repeat("=", done-1) + ">"))
			w.Write([]byte(strings.Repeat(" ", width-done)))
		}
		w.Write([]byte("]"))
	}
}
