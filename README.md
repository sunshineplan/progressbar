# progressbar

A **modern, customizable, concurrent-safe** progress bar library for Go (Golang).  
Supports dynamic terminal width detection, templating, byte-size formatting, speed calculation, ETA, additional info, and works seamlessly with `io.Reader`.

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]

[godev]: https://pkg.go.dev/github.com/sunshineplan/progressbar

## Features

- **Fully customizable template** using Go's `text/template`
- **Automatic terminal width detection** (Windows + Unix)
- **Byte size formatting** (`1.23 MB`, `45.6 KiB/s`, etc.)
- **Accurate speed calculation** (total + interval-based)
- **ETA / Time left** with animated "calculating..." dots
- **Additional dynamic text** support
- **Message channel** for temporary notifications
- **Concurrent-safe** – multiple goroutines can call `Add()` safely
- **Generic** – works with `int` or `int64` totals
- **io.Reader integration** via `FromReader()`
- **Cancelable** with graceful "Cancelled" output

## Installation

```bash
go get github.com/sunshineplan/progressbar
```

## Quick Start

```go
package main

import (
    "os"
    "time"

    "github.com/sunshineplan/progressbar"
)

func main() {
    // Create a progress bar with total = 1000
    pb := progressbar.New(1000)

    // Optional customizations
    pb.SetWidth(50).
       SetRefreshInterval(2 * time.Second).
       SetRenderInterval(100 * time.Millisecond)

    // Start rendering
    if err := pb.Start(); err != nil {
        panic(err)
    }

    // Simulate work
    for i := 0; i < 1000; i++ {
        pb.Add(1)
        pb.Additional("processing user #" + string(rune(i+1)))
        time.Sleep(5 * time.Millisecond)
    }

    // Wait for completion (optional if you want to block)
    pb.Wait()
}
```

## Beautiful Adaptive Output (example)

```
[=========================>                  ]  15.67 MB/s  784.20 MB(78.42%) of 1.00 GB [downloading chunk 89]  Elapsed: 50s  Left: 14s
```

When terminal is narrow, it gracefully degrades to shorter formats, and finally to just the bar + additional text.

## Customization

### Custom Render Function (Recommended)

```go
pb.SetRender(func(w io.Writer, f progressbar.Frame) {
    fmt.Fprintf(w, "[%s] %.1f%% %s/s",
        progressbar.Bar(f.Current, f.Total, f.BarWidth),
        progressbar.Percent(f.Current, f.Total),
        progressbar.Speed(f.Speed, f.Unit),
    )
})
```

### Or Custom Template (Advanced)

```go
pb.SetTemplate(template.Must(template.New("").Funcs(progressbar.DefaultFuncMap()).Parse(`[{{bar .Current .Total .BarWidth}}] {{percent .Current .Total | printf "%.1f%%"}} {{speed .Speed .Unit}}`)))
```

### Byte Downloads (with io.Reader)

```go
pb := progressbar.New(totalBytes)
pb.SetUnit("bytes") // enables ByteSize formatting

pb.FromReader(http.Body, os.Stdout) // starts automatically, returns written bytes
```

### Display Temporary Messages

```go
pb.Message("Downloading update...")
time.Sleep(2 * time.Second)
pb.Message("Verifying checksum...")
```

### Cancel Gracefully

```go
go func() {
    time.Sleep(10 * time.Second)
    pb.Cancel() // prints "\nCancelled\n"
}()
```

## API Reference

```go
// Creation
pb := progressbar.New[T int|int64](total T) *ProgressBar[T]

// Configuration
pb.SetWidth(width int)
pb.SetRefreshInterval(d time.Duration) // speed calculation interval
pb.SetRenderInterval(d time.Duration)  // UI refresh rate
pb.SetRender(fn func(w io.Writer, f Frame)) error
pb.SetTemplate(t *template.Template) error
pb.SetUnit(unit string)               // "bytes" enables ByteSize, otherwise numeric

// Runtime
pb.Start() error
pb.Add(n T)
pb.Additional(text string)
pb.Message(msg string) error
pb.Total() int64
pb.Current() int64
pb.Elapsed() time.Duration
pb.Speed() float64
pb.Wait()
pb.Cancel()

// io.Reader helper
pb.FromReader(r io.Reader, w io.Writer) (written int64, err error)

// template helper functions
progressbar.Bar(current, total int64, width int) string
progressbar.Format(n int64, unit string) string
progressbar.Percent(current, total int64) float64
progressbar.Speed(v float64, unit string) string
progressbar.Left(current, total int64, speed float64) time.Duration
```


## License

MIT © SunshinePlan

---

Enjoy a clean, informative progress bar in your CLI tools!