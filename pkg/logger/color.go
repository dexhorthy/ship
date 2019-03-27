// The MIT License (MIT)
//
// Copyright (c) 2015 Peter Bourgon
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
//
// This file copied in part from
// https://github.com/go-kit/kit/blob/master/log/term/colorlogger.go
package logger

import "fmt"

// Color represents an ANSI color. The zero value is Default.
type Color uint8

// ANSI colors.
const (
	Default = Color(iota)

	Black
	DarkRed
	DarkGreen
	Brown
	DarkBlue
	DarkMagenta
	DarkCyan
	Gray

	DarkGray
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	numColors
)

// For more on ANSI escape codes see
// https://en.wikipedia.org/wiki/ANSI_escape_code. See in particular
// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors.

var (
	resetColorBytes = []byte("\x1b[39;49;22m")
	fgColorBytes    [][]byte
	bgColorBytes    [][]byte
)

func init() {
	// Default
	fgColorBytes = append(fgColorBytes, []byte("\x1b[39m"))
	bgColorBytes = append(bgColorBytes, []byte("\x1b[49m"))

	// dark colors
	for color := Black; color < DarkGray; color++ {
		fgColorBytes = append(fgColorBytes, []byte(fmt.Sprintf("\x1b[%dm", 30+color-Black)))
		bgColorBytes = append(bgColorBytes, []byte(fmt.Sprintf("\x1b[%dm", 40+color-Black)))
	}

	// bright colors
	for color := DarkGray; color < numColors; color++ {
		fgColorBytes = append(fgColorBytes, []byte(fmt.Sprintf("\x1b[%d;1m", 30+color-DarkGray)))
		bgColorBytes = append(bgColorBytes, []byte(fmt.Sprintf("\x1b[%d;1m", 40+color-DarkGray)))
	}
}
