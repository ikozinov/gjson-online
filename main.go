package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ikozinov/gjson-online/pages"
	"github.com/ikozinov/gjson-online/utils"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const assetsDir = "web"

func main() {
	app.Route("/", func() app.Composer {
		return &pages.GJSONPlayground{}
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
		if err := utils.CopyAssetsToDist(assetsDir); err != nil {
			log.Fatal(err)
		}
		if err := utils.UpdateWasmExec(); err != nil {
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
