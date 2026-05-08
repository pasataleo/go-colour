# v0.1.0

## FEATURES

- ANSI colour and style constants for terminal output
- `Colour` type with inline placeholder syntax (`{red}`, `{bold}`, etc.)
- `RGB` and `BgRGB` functions for 24-bit foreground and background colours
- `{rgb:R,G,B}` and `{bg-rgb:R,G,B}` placeholder support in format strings
- `WithColour` and `WithColours` options for registering custom placeholders
- `NO_COLOR` environment variable support — colours are automatically disabled when set (see [no-color.org](https://no-color.org))
- `Enable` option to force colours on regardless of `NO_COLOR`
- `Disable` option to force colours off regardless of `NO_COLOR`
- `Colourf` for applying placeholders to format strings before `fmt.Sprintf`
