# go-colour

A Go library for adding ANSI colour codes to strings using inline placeholders.

## Installation

```sh
go get github.com/pasataleo/go-colour
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/pasataleo/go-colour/pkg/colour"
)

func main() {
	c := colour.New()
	fmt.Println(c.Colour("{red}error:{reset} something went wrong"))
}
```

### Format strings

`Colourf` applies colour placeholders to the format string before passing arguments through `fmt.Sprintf`. Arguments are not processed for colour placeholders.

```go
c := colour.New()
fmt.Println(c.Colourf("{red}error: %s{reset}", "disk full"))
```

### RGB colours

Use `{rgb:R,G,B}` and `{bg-rgb:R,G,B}` for 24-bit foreground and background colours:

```go
c := colour.New()
fmt.Println(c.Colour("{rgb:255,100,0}orange text{reset}"))
fmt.Println(c.Colour("{bg-rgb:0,0,128}navy background{reset}"))
```

### Custom placeholders

Register custom placeholders with `WithColour` or `WithColours`:

```go
c := colour.New(
	colour.WithColour("{highlight}", colour.Yellow+colour.Bold),
	colour.WithColours(map[string]string{
		"{error}":   colour.Red,
		"{success}": colour.Green,
	}),
)
fmt.Println(c.Colour("{error}failed{reset} / {success}passed{reset}"))
```

### TTY detection

Colours are automatically disabled when stdout is not a terminal. Use `Enable` or `Disable` to override:

```go
c := colour.New(colour.Enable())  // force colours on
c := colour.New(colour.Disable()) // force colours off
```

## Built-in placeholders

### Styles

`{reset}`, `{bold}`, `{italic}`, `{underline}`, `{strikethrough}`

### Foreground colours

`{black}`, `{red}`, `{green}`, `{yellow}`, `{blue}`, `{magenta}`, `{cyan}`, `{white}`, `{default}`

### Bright foreground colours

`{grey}`, `{bright-red}`, `{bright-green}`, `{bright-yellow}`, `{bright-blue}`, `{bright-magenta}`, `{bright-cyan}`, `{bright-white}`

### Background colours

`{bg-black}`, `{bg-red}`, `{bg-green}`, `{bg-yellow}`, `{bg-blue}`, `{bg-magenta}`, `{bg-cyan}`, `{bg-white}`, `{bg-default}`

### Bright background colours

`{bg-grey}`, `{bg-bright-red}`, `{bg-bright-green}`, `{bg-bright-yellow}`, `{bg-bright-blue}`, `{bg-bright-magenta}`, `{bg-bright-cyan}`, `{bg-bright-white}`

