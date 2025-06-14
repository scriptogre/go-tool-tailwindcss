# Tailwind CSS CLI as a Go Tool

A dead simple way to run [Tailwind CSS](https://tailwindcss.com/) CLI as a Go tool.

No more **Node.JS**. No more **manual downloads of the CLI**.

## Install

```bash
go get -tool github.com/scriptogre/go-tool-tailwindcss@latest
```

## Use

1. Create `input.css` in your project:
   ```css
   @import 'tailwindcss';
   ```

2. Run TailwindCSS:
   ```bash
   go tool tailwindcss -i input.css -o output.css --watch
   ```

That's it!

**Note:** Downloaded TailwindCSS binary is cached in `~/.cache/go-tool-tailwindcss/`
