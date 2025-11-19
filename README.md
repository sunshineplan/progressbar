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

## Default Render

Example output:

```
[=========================>                   ]  12.34 MB/s  567.89 MB (56.79%) of 1 GB [processing chunk 42]  Elapsed: 45s  Left: 36s
```

## Customization

### Change Template

```go
pb.SetTemplate(`[{{.Done}}{{.Undone}}] {{.Percent}} | {{.Speed}} | {{.Elapsed}} | {{.Left}}`)
```

Available fields:

| Field        | Description                              |
|--------------|------------------------------------------|
| `{{.Done}}`      | Completed bar (`=`)                      |
| `{{.Undone}}`    | Remaining bar (` `)                      |
| `{{.Speed}}`     | Current speed (auto-unit)                |
| `{{.Current}}`   | Current value (formatted)                |
| `{{.Total}}`     | Total value (formatted)                  |
| `{{.Percent}}`   | Percentage                               |
| `{{.Additional}}`| Custom additional text                   |
| `{{.Elapsed}}`   | Time elapsed                             |
| `{{.Left}}`      | Estimated time remaining                 |

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
pb.SetTemplate(tmplt string) error
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
```


## License

MIT © SunshinePlan

---

Enjoy a clean, informative progress bar in your CLI tools!