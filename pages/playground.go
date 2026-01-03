package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/tidwall/gjson"
)

// GJSONPlayground is the main component of the application.
type GJSONPlayground struct {
	app.Compo
	JSONContent string
	Query       string
	Result      string
	JSONError   string
	PathFound   bool
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
	p.updateResult(ctx)
}

func (p *GJSONPlayground) updateResult(ctx app.Context) {
	p.calculateResult()
	// Trigger UI update
	ctx.Update()
}

func (p *GJSONPlayground) calculateResult() {
	if !gjson.Valid(p.JSONContent) {
		p.JSONError = "Invalid JSON"
		p.Result = ""
		p.PathFound = false
	} else {
		p.JSONError = "" // Clear error if valid
		if p.Query == "" {
			p.Result = ""
			p.PathFound = true // Empty query is considered "found" in terms of no error
		} else {
			res := gjson.Get(p.JSONContent, p.Query)
			if res.Exists() {
				p.Result = res.String()
				p.PathFound = true
			} else {
				p.Result = ""
				p.PathFound = false
			}
		}
	}
}

// OnJSONChange handles changes in the JSON input textarea.
func (p *GJSONPlayground) OnJSONChange(ctx app.Context, e app.Event) {
	p.JSONContent = ctx.JSSrc().Get("value").String()
	p.updateResult(ctx)
}

// OnQueryChange handles changes in the GJSON path input.
func (p *GJSONPlayground) OnQueryChange(ctx app.Context, e app.Event) {
	p.Query = ctx.JSSrc().Get("value").String()
	p.updateResult(ctx)
}

// Render describes the UI.
func (p *GJSONPlayground) Render() app.UI {
	textareaClass := "form-control"
	if p.JSONError != "" {
		textareaClass += " is-invalid"
	}

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
							app.Textarea().
								Class(textareaClass).
								ID("json-area").
								Style("height", "70vh").
								Style("font-family", "monospace").
								// Value(p.JSONContent) was undefined
								Body(app.Text(p.JSONContent)).
								OnInput(p.OnJSONChange),
							app.If(p.JSONError != "", func() app.UI {
								return app.Div().Class("invalid-feedback").Text(p.JSONError)
							}),
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
								Body(app.Text(p.Result)),
							app.If(!p.PathFound && p.JSONError == "" && p.Query != "", func() app.UI {
								return app.Div().Class("text-danger").Text("Value not found")
							}),
						),
					),
				),
			),
		),
	)
}
