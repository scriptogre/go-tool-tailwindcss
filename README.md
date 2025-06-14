# Tailwind CSS CLI as a Go Tool

A dead simple way to run Tailwind CSS standalone CLI using `go tool` [released with Go 1.2.4](https://tip.golang.org/doc/go1.24#tools).

No more **Node.JS**. No more **manual downloads**.

## Usage

1. Install it:
   ```bash
   go get -tool github.com/scriptogre/tailwindcss-go-tool@latest
   ```

2. Create `input.css` in your project:
   ```css
   @import 'tailwindcss';
   ```

3. Run it:
   ```bash
   go tool tailwindcss -i input.css -o output.css --watch
   ```

That's it!

**Note:** Downloaded TailwindCSS binary is cached in `~/.cache/tailwindcss-go-tool/`
