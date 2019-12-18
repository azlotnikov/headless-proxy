package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"proxy/api_server"
	"proxy/proxy_server"
)

var (
	app    = kingpin.New("headless-proxy-server", "Bla bla")
	port   = app.Flag("port", "Api server port").Short('p').Default("8080").Int()
	apiKey = app.Flag("apiKey", "Api key").Short('a').Default("").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	ps, err := proxy_server.NewLocalProxyServer()
	if err != nil {
		log.Fatal(err)
	}
	as := api_server.NewApiServer(ps, *apiKey)
	log.Fatal(as.Serve(*port))
}
