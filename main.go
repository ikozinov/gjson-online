package main

import (
	"log"
	"net/http"
	"os"

	"io"
	"path/filepath"

	"runtime"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/tidwall/gjson"
)

// GJSONPlayground is the main component of the application.
type GJSONPlayground struct {
	app.Compo
	JSONContent string
	Query       string
	Result      string
}

func (p *GJSONPlayground) OnNav(ctx app.Context) {
	app.Log("GJSONPlayground navigation")
}

// OnMount initializes the component with default data.
func (p *GJSONPlayground) OnMount(ctx app.Context) {
	app.Log("GJSONPlayground mounted")
	// Example JSON taken from GJSON README or similar common examples
	p.JSONContent = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age": 37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`
	p.Query = "name.last"
	p.updateResult()
}

func (p *GJSONPlayground) updateResult() {
	app.Log("Updating result for query:", p.Query)
	if p.Query == "" {
		p.Result = ""
	} else {
		res := gjson.Get(p.JSONContent, p.Query)
		p.Result = res.String()
	}
}

// OnJSONChange handles changes in the JSON input textarea.
func (p *GJSONPlayground) OnJSONChange(ctx app.Context, e app.Event) {
	p.JSONContent = ctx.JSSrc().Get("value").String()
	p.updateResult()
}

// OnQueryChange handles changes in the GJSON path input.
func (p *GJSONPlayground) OnQueryChange(ctx app.Context, e app.Event) {
	p.Query = ctx.JSSrc().Get("value").String()
	p.updateResult()
}

// Render describes the UI.
func (p *GJSONPlayground) Render() app.UI {
	return app.Div().Class("page-wrapper with-navbar").Body(
		// Navbar
		app.Nav().Class("navbar").Body(
			app.A().Href("#").Class("navbar-brand").Text("GJSON Playground"),
			app.Ul().Class("navbar-nav d-none d-md-flex").Body(
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link").Href("https://github.com/tidwall/gjson").Target("_blank").Text("GJSON Library"),
				),
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link").Href("https://github.com/tidwall/gjson/blob/master/SYNTAX.md").Target("_blank").Text("Syntax"),
				),
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link").Href("https://github.com/ikozinov/gjson-online").Target("_blank").Text("GitHub Repo"),
				),
			),
		),
		// Content
		app.Div().Class("content-wrapper").Body(
			app.Div().Class("container-fluid").Body(
				app.Div().Class("row").Body(
					app.Div().Class("col-12").Body(
						app.Div().Class("p-20").Body( // padding
							app.Label().For("gjson-input").Text("GJSON Path"),
							app.Input().Type("text").Class("form-control").ID("gjson-input").
								Placeholder("Enter GJSON path...").
								Value(p.Query).
								OnInput(p.OnQueryChange),
						),
					),
				),
				app.Div().Class("row").Body(
					app.Div().Class("col-md-6").Body(
						app.Div().Class("p-20").Body(
							app.Label().For("json-area").Text("JSON Input"),
							app.Textarea().Class("form-control").ID("json-area").
								Style("height", "70vh").
								Style("font-family", "monospace").
								// Value(p.JSONContent) was undefined
								Body(app.Text(p.JSONContent)).
								OnInput(p.OnJSONChange),
						),
					),
					app.Div().Class("col-md-6").Body(
						app.Div().Class("p-20").Body(
							app.Label().For("result-area").Text("Result"),
							app.Textarea().Class("form-control").ID("result-area").
								Style("height", "70vh").
								Style("font-family", "monospace").
								Style("background-color", "#f0f0f0"). // Slight visual distinction
								ReadOnly(true).
								Text(p.Result),
						),
					),
				),
			),
		),
	)
}

const assetsDir = "web"

func main() {
	app.Route("/", func() app.Composer {
		return &GJSONPlayground{}
	})

	app.RunWhenOnBrowser()

	// Define the handler with the required configuration
	handler := &app.Handler{
		Name:        "GJSON Playground",
		Description: "Online playground for testing GJSON expressions",
		Styles: []string{
			"https://cdn.jsdelivr.net/npm/halfmoon@1.1.1/css/halfmoon.min.css",
		},
		Title: "GJSON Online",
		Icon: app.Icon{
			Default: "/" + assetsDir + "/icon.svg", // Use the custom icon
			SVG:     "/" + assetsDir + "/icon.svg", // Use the custom icon
		},
	}

	// Check if the "dist" argument is provided to generate the static site
	if len(os.Args) > 1 && os.Args[1] == "dist" {
		if err := app.GenerateStaticWebsite("dist", handler); err != nil {
			log.Fatal(err)
		}
		if err := copyAssetsToDist(assetsDir); err != nil {
			log.Fatal(err)
		}
		if err := updateWasmExec(); err != nil {
			log.Fatal(err)
		}
		return
	}

	http.Handle("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// copyAssetsToDist copies all files and directories from the specified assets directory
// to the dist/web directory.
func copyAssetsToDist(assetsDir string) error {
	dstDir := filepath.Join("dist", assetsDir)

	return filepath.Walk(assetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(assetsDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Create directory in destination
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Open source file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create destination file with same permissions
		dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// Copy content
		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

// updateWasmExec copies the wasm_exec.js file from GOROOT to dist/wasm_exec.js
// ensuring the JS glue code matches the Go version used to build the WASM binary.
func updateWasmExec() error {
	goroot := runtime.GOROOT()

	// Check possible locations for wasm_exec.js
	// Go 1.24+ moved it to lib/wasm, older versions had it in misc/wasm
	locations := []string{
		filepath.Join(goroot, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(goroot, "misc", "wasm", "wasm_exec.js"),
	}

	var src string
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			src = loc
			break
		}
	}

	if src == "" {
		return os.ErrNotExist
	}

	dst := filepath.Join("dist", "wasm_exec.js")

	log.Printf("Updating wasm_exec.js from %s to %s", src, dst)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
