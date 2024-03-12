package main

import (
	"fmt"
	"log"
	"net/http"

	"places/filescanner"
	"places/httpserver"
	"places/search"
	"places/search/graphhopper"
	"places/search/opentripmap"
	"places/search/openweather"
)

func getSiteFiles() ([]string, error) {
	css, err := filescanner.ScanForSuffix("./site", ".css")
	if err != nil {
		return nil, fmt.Errorf("get site files: %w", err)
	}
	js, err := filescanner.ScanForSuffix("./site", ".js")
	if err != nil {
		return nil, fmt.Errorf("get site files: %w", err)
	}
	png, err := filescanner.ScanForSuffix("./site", ".png")
	if err != nil {
		return nil, fmt.Errorf("get site files: %w", err)
	}
	return append(append(css, js...), png...), nil
}

func main() {
	// insert your keys for that services
	// it should work using command line arguments, but now it doesn't
	gKey := "your-key"
	tKey := "your-key"
	wKey := "your-key"

	if gKey == "your-key" || tKey == "your-key" || wKey == "your-key" {
		log.Println("WARN: possibly missing keys for apis (one of value is \"your-key\")")
	}

	g := graphhopper.NewGraphHopper(gKey)
	t := opentripmap.NewOpenTripMap(tKey)
	w := openweather.NewOpenWeather(wKey)
	placer := search.NewPlacer(g, w, t)
	addons, err := getSiteFiles()
	if err != nil {
		log.Fatal(err)
	}

	httpServer, err := httpserver.NewPlacesHttpServer("./site", "./site/main.html", addons, placer)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", httpServer.GetOnRoot)
	http.HandleFunc("/ws", httpServer.WebSocket)

	// uncomment code below and set your domainName to enable HTTPS
	// domainName := "your-dns"
	// m := &autocert.Manager{
	// 	Cache:      autocert.DirCache("secret-dir"),
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist(domainName),
	// }
	// s := &http.Server{
	// 	Addr:      ":6969",
	// 	TLSConfig: m.TLSConfig(),
	// }
	// log.Fatal(s.ListenAndServeTLS("", ""))

	// comment this line if you want to enable HTTPS
	log.Fatal(http.ListenAndServe(":6969", nil))
}
