package progressbar

import (
	"io"
	"strings"
	"text/template"
	"time"
)

// Frame represents a snapshot of all display-ready fields used to render the
// progress bar at a given moment. It contains the computed visual components
// (e.g., completed blocks, percentages, speed, timing information) that are
// assembled by the renderer into the final output.
type Frame struct {
	Percent        float64
	Done, Undone   string
	Speed          string
	Current, Total string
	Additional     string
	Elapsed        time.Duration
	Left           string
}

var defaultTemplate *template.Template

func init() {
	defaultTemplate =
		template.Must(
			template.Must(
				template.Must(
					template.Must(
						template.New("full").Parse(
							`[{{.Done}}{{.Undone}}]  {{.Speed}}  {{.Current}}({{printf "%.2f%%" .Percent}}) of {{.Total}}{{if .Additional}} [{{.Additional}}]{{end}}  Elapsed: {{.Elapsed }}  {{if eq .Current .Total}}Complete{{else}}Left: {{.Left}}{{end}} `,
						),
					).New("standard").Parse(
						`[{{.Done}}{{.Undone}}] {{.Speed}} {{.Current}}/{{.Total}}({{printf "%.1f%%" .Percent}}){{if .Additional}} [{{.Additional}}]{{end}} ET: {{.Elapsed}} {{if eq .Current .Total}}Done{{else}}LT: {{.Left}}{{end}} `,
					),
				).New("lite").Parse(
					`[{{.Done}}{{.Undone}}] {{.Speed}} {{.Current}}/{{.Total}}{{if .Additional}} [{{.Additional}}]{{end}} {{if eq .Current .Total}}E: {{.Elapsed}}{{else}}L: {{.Left}}{{end}} `,
				),
			).New("mini").Parse(
				`[{{.Done}}{{.Undone}}] {{.Additional}} `,
			),
		)
}

func defaultRenderFunc(w io.Writer, f Frame) {
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
