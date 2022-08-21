// Writes an SVG clockface of the current time to Stdout
package main

import (
	"os"
	"time"

	"github.com/jrvldam/learn-go-with-tests/15-mathematics/svg"
)

func main() {
	t := time.Now()
	svg.Write(os.Stdout, t)
}
