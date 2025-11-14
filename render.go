package progressbar

import (
	"io"
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

var fullTemplate = template.Must(template.New("ProgressBar").Parse(`[{{.Done}}{{.Undone}}]  {{.Speed}}  {{.Current}}({{printf "%.2f%%" .Percent}}) of {{.Total}}{{if .Additional}} [{{.Additional}}]{{end}}  Elapsed: {{.Elapsed }}  {{if eq .Current .Total}}Complete{{else}}Left: {{.Left}}{{end}} `))

func defaultRenderFn(w io.Writer, f Frame) {
	fullTemplate.Execute(w, f)
}
