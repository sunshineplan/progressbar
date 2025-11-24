package progressbar

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/sunshineplan/utils/pool"
	"github.com/sunshineplan/utils/unit"
)

// Frame represents a snapshot of all display-ready fields used to render the
// progress bar at a given moment. It contains the computed visual components
// (e.g., completed blocks, percentages, speed, timing information) that are
// assembled by the renderer into the final output.
type Frame struct {
	Unit           string
	BarWidth       int
	Current, Total int64
	Speed          float64
	Additional     string
	Elapsed        time.Duration
}

var framePool = pool.New[Frame]()

func (pb *ProgressBar[T]) frame() *Frame {
	f := framePool.Get()
	f.Unit = pb.unit.Load()
	f.BarWidth = pb.barWidth.Load()
	f.Current = min(pb.Current(), pb.total)
	f.Total = pb.total
	f.Additional = pb.additional.Load()
	f.Elapsed = pb.Elapsed().Truncate(time.Second)
	if f.Current == pb.total {
		f.Speed = float64(pb.total) / (float64(pb.Elapsed()) / float64(time.Second))
	} else {
		f.Speed = pb.Speed()
	}
	return f
}

var defaultTemplate *template.Template

var (
	spinner      = [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinnerIndex atomic.Int64
)

// Spin returns the next frame of the spinner animation.
func Spin() string {
	return spinner[spinnerIndex.Add(1)%10]
}

// Bar generates a textual progress bar representation based on the current progress.
func Bar(current, total int64, barWidth int) string {
	if current == total {
		return strings.Repeat("=", barWidth)
	} else {
		if done := int(float64(barWidth) * float64(current) / float64(total)); done != 0 {
			return strings.Repeat("=", done-1) + ">" + strings.Repeat(" ", barWidth-done)
		}
		return strings.Repeat(" ", barWidth)
	}
}

// Format formats the given number n according to the specified type t.
func Format(n int64, t string) string {
	if t == "bytes" {
		return unit.ByteSize(n).String()
	} else {
		return strconv.FormatInt(n, 10)
	}
}

// Left estimates the remaining time to complete the progress bar based on the current speed.
func Left(current, total int64, speed float64) time.Duration {
	return (time.Duration(float64(total-current)/speed) * time.Second).Truncate(time.Second)
}

// Percent calculates the completion percentage of the progress bar.
func Percent(current, total int64) float64 {
	return float64(current) * 100 / float64(total)
}

// Speed formats the speed of progress per second, adapting the output based on the specified unit type.
func Speed(speed float64, t string) string {
	if speed == 0 {
		return "--/s"
	} else if t == "bytes" {
		return fmt.Sprintf("%s/s", unit.ByteSize(speed))
	} else {
		return fmt.Sprintf("%.2f/s", speed)
	}
}

// DefaultFuncMap returns the default template function map for progress bar rendering.
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"bar": Bar,
		"calc": func() string {
			return "calculating" + Spin()
		},
		"format":  Format,
		"left":    Left,
		"percent": Percent,
		"speed":   Speed,
	}
}

func init() {
	t := template.New("")
	t.Funcs(DefaultFuncMap())
	defaultTemplate =
		template.Must(
			template.Must(
				template.Must(
					template.Must(
						t.New("full").Parse(
							`[{{bar .Current .Total .BarWidth}}]  {{speed .Speed .Unit}}  {{format .Current .Unit}}({{percent .Current .Total | printf "%.2f%%"}}) of {{format .Total .Unit}}{{if .Additional}} [{{.Additional}}]{{end}}  Elapsed: {{.Elapsed }}  {{if eq .Current .Total}}Complete{{else}}Left: {{if eq .Speed 0.0}}{{calc}}{{else}}{{left .Current .Total .Speed}}{{end}}{{end}} `,
						),
					).New("standard").Parse(
						`[{{bar .Current .Total .BarWidth}}] {{speed .Speed .Unit}} {{format .Current .Unit}}/{{format .Total .Unit}}({{percent .Current .Total | printf "%.1f%%"}}){{if .Additional}} [{{.Additional}}]{{end}} ET: {{.Elapsed}} {{if eq .Current .Total}}Done{{else}}LT: {{if eq .Speed 0.0}}{{calc}}{{else}}{{left .Current .Total .Speed}}{{end}}{{end}} `,
					),
				).New("lite").Parse(
					`[{{bar .Current .Total .BarWidth}}] {{speed .Speed .Unit}} {{format .Current .Unit}}/{{format .Total .Unit}}{{if .Additional}} [{{.Additional}}]{{end}} {{if eq .Current .Total}}E: {{.Elapsed}}{{else}}L: {{if eq .Speed 0.0}}{{calc}}{{else}}{{left .Current .Total .Speed}}{{end}}{{end}} `,
				),
			).New("mini").Parse(
				`[{{bar .Current .Total .BarWidth}}] {{.Additional}} `,
			),
		)
}

// DefaultRenderFunc is the default function used to render the progress bar.
var DefaultRenderFunc = func(w io.Writer, f Frame) {
	winsize := GetWinsize()
	n := winsize - f.BarWidth - len(f.Additional)
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
		if f.Current == f.Total {
			w.Write([]byte(strings.Repeat("=", width)))
		} else {
			done := int(float64(width) * float64(f.Current) / float64(f.Total))
			if done != 0 {
				w.Write([]byte(strings.Repeat("=", done-1) + ">"))
			}
			w.Write([]byte(strings.Repeat(" ", width-done)))
		}
		w.Write([]byte("]"))
	}
}
