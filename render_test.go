package progressbar

import (
	"strings"
	"testing"
	"time"
)

func TestDefaultRenderFunc(t *testing.T) {
	oldWinsize := GetWinsize
	defer func() { GetWinsize = oldWinsize }()

	tests := []struct {
		name     string
		winsize  int
		frame    Frame
		expected string
	}{
		{
			name:    "full template with additional",
			winsize: 130, // n = 130 - 40 - 20 = 70 ≥ 60 → full template
			frame: Frame{
				Percent:    50.0,
				Done:       strings.Repeat("=", 19) + ">",
				Undone:     strings.Repeat(" ", 20),
				Speed:      15.67,
				Current:    500,
				Total:      1000,
				Additional: "downloading chunk 42",
				Elapsed:    95 * time.Second,
				Left:       65 * time.Second,
			},
			expected: "[===================>                    ]  15.67/s  500(50.00%) of 1000 [downloading chunk 42]  Elapsed: 1m35s  Left: 1m5s ",
		},
		{
			name:    "full template without additional",
			winsize: 130,
			frame: Frame{
				Percent:    75.5,
				Done:       strings.Repeat("=", 29) + ">",
				Undone:     strings.Repeat(" ", 10),
				Speed:      2.50,
				Current:    750,
				Total:      1000,
				Additional: "",
				Elapsed:    200 * time.Second,
				Left:       53 * time.Second,
			},
			expected: "[=============================>          ]  2.50/s  750(75.50%) of 1000  Elapsed: 3m20s  Left: 53s ",
		},
		{
			name:    "standard template",
			winsize: 110, // n = 110 - 40 - 20 = 50 < 60, ≥ 40 → standard
			frame: Frame{
				Percent:    50.0,
				Done:       strings.Repeat("=", 19) + ">",
				Undone:     strings.Repeat(" ", 20),
				Speed:      15.67,
				Current:    500,
				Total:      1000,
				Additional: "downloading chunk 42",
				Elapsed:    95 * time.Second,
				Left:       65 * time.Second,
			},
			expected: "[===================>                    ] 15.67/s 500/1000(50.0%) [downloading chunk 42] ET: 1m35s LT: 1m5s ",
		},
		{
			name:    "lite template",
			winsize: 90, // n = 90 - 40 - 20 = 30 < 60, ≥ 40 → lite
			frame: Frame{
				Percent:    50.0,
				Done:       strings.Repeat("=", 19) + ">",
				Undone:     strings.Repeat(" ", 20),
				Speed:      15.67,
				Current:    500,
				Total:      1000,
				Additional: "downloading chunk 42",
				Elapsed:    95 * time.Second,
				Left:       65 * time.Second,
			},
			expected: "[===================>                    ] 15.67/s 500/1000 [downloading chunk 42] L: 1m5s ",
		},
		{
			name:    "mini template",
			winsize: 66, // n = 66 - 40 - 20 = 6 > 5 → mini
			frame: Frame{
				Percent:    50.0,
				Done:       strings.Repeat("=", 19) + ">",
				Undone:     strings.Repeat(" ", 20),
				Additional: "downloading chunk 42",
			},
			expected: "[===================>                    ] downloading chunk 42 ",
		},
		{
			name:    "fallback incomplete",
			winsize: 50, // n = 50 - 40 - 0 = 10 ≤ 20 and Additional empty → fallback
			frame: Frame{
				Percent: 50.0,
				Done:    strings.Repeat("=", 19) + ">",
				Undone:  strings.Repeat(" ", 20),
			},
			// width = 50-5 = 45, done = int(45*0.5) = 22
			// output: [ + "="*21 + ">" + " "*23 + ]
			expected: "[" + strings.Repeat("=", 21) + ">" + strings.Repeat(" ", 23) + "]",
		},
		{
			name:    "fallback complete",
			winsize: 50,
			frame: Frame{
				Percent: 100.0,
				Done:    strings.Repeat("=", 40),
			},
			expected: "[" + strings.Repeat("=", 50-5) + "]",
		},
		{
			name:    "complete in full template",
			winsize: 130,
			frame: Frame{
				Percent:    100.0,
				Done:       strings.Repeat("=", 40),
				Undone:     "",
				Speed:      20.00,
				Current:    1000,
				Total:      1000,
				Additional: "finished",
				Elapsed:    95 * time.Second,
				Left:       0,
			},
			expected: "[========================================]  20.00/s  1000(100.00%) of 1000 [finished]  Elapsed: 1m35s  Complete ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetWinsize = func() int { return tt.winsize }

			var b strings.Builder
			DefaultRenderFunc(&b, tt.frame)

			got := b.String()

			if got != tt.expected {
				t.Errorf("\n  want: %q\n  got : %q", tt.expected, got)
			}
		})
	}
}
