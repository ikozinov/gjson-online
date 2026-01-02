package main

import (
	"log"
	"net/http"
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tidwall/gjson"
)

// GJSONPlayground is the main component of the application.
type GJSONPlayground struct {
	app.Compo
	JSONContent string
	Query       string
	Result      string
}

// OnMount initializes the component with default data.
func (p *GJSONPlayground) OnMount(ctx app.Context) {
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
	if p.Query == "" {
		p.Result = ""
		return
	}
	res := gjson.Get(p.JSONContent, p.Query)
	p.Result = res.String()
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
								Text(p.JSONContent).
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

func main() {
	app.Route("/", &GJSONPlayground{})

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
			Default:    "/web/icon.svg", // Use the custom icon
			AppleTouch: "/web/icon.svg", // Use it for Apple Touch too
		},
		RawHeaders: []string{
			`<script>
				window.addEventListener('error', function(event) {
					console.error("Global error caught:", event.error);
					// You can also display a UI alert here if needed
				});
				window.addEventListener('unhandledrejection', function(event) {
					console.error("Unhandled promise rejection:", event.reason);
				});
			</script>`,
		},
	}

	// Check if the "gen" argument is provided to generate the static site
	if len(os.Args) > 1 && os.Args[1] == "gen" {
		if err := app.GenerateStaticWebsite("web", handler); err != nil {
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
