package main

import (
	"encoding/json"
	//"fmt"
	"log"

	"github.com/mithorium/secure-fasthttp"
	"github.com/valyala/fasthttp"
)

//19 de octubre a las 16:45

type Palabras struct {
	I uint32 `json:"I"`
	T uint32 `json:"T"`
	N string `json:"N"`
}

type Filtro struct {
	Nombre  string    `json:"Nombre"`
	Valores []Valores `json:"Valores"`
}
type Valores struct {
	Nombre string `json:"Nombre"`
}

func requestHandler(ctx *fasthttp.RequestCtx) {

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/auto":

			ctx.Response.Header.Set("Content-Type", "application/json")
			val := string(ctx.QueryArgs().Peek("c"))

			if val == "ab" {
				ab := []Palabras{Palabras{I: 1, T: 1, N: "dsa"}, Palabras{I: 1, T: 1, N: "dma"}, Palabras{I: 1, T: 1, N: "dxf"}}
				json.NewEncoder(ctx).Encode(ab)
			}
			if val == "abc" {
				abc := []Palabras{Palabras{I: 1, T: 1, N: "amda"}, Palabras{I: 1, T: 1, N: "edma"}, Palabras{I: 1, T: 1, N: "idse"}}
				json.NewEncoder(ctx).Encode(abc)
			}
			if val == "abcd" {
				abcd := []Palabras{Palabras{I: 1, T: 1, N: "ocda"}, Palabras{I: 1, T: 1, N: "abra"}, Palabras{I: 1, T: 1, N: "iste"}}
				json.NewEncoder(ctx).Encode(abcd)
			}

		case "/lang":

			ctx.Response.Header.Set("Content-Type", "application/json")
			val := string(ctx.QueryArgs().Peek("c"))
			if val == "zr" {
				zr := []string{"HOLA#MUNDO", "妹妹背著#洋娃娃"}
				json.NewEncoder(ctx).Encode(zr)
			}
			if val == "gr" {
				gr := [][]string{[]string{"HOLA MUNDO", "妹妹背著 洋娃娃"}, []string{"HOLA MUNDO", "妹妹背著 洋娃娃"}}
				json.NewEncoder(ctx).Encode(gr)
			}

		case "/filtro":

			ctx.Response.Header.Set("Content-Type", "application/json")
			val := string(ctx.QueryArgs().Peek("c"))
			if val == "1" {
				zr := []Filtro{Filtro{Nombre: "Marca", Valores: []Valores{Valores{Nombre: "Chevrolet"}, Valores{Nombre: "Audi"}}}, Filtro{Nombre: "Modelo", Valores: []Valores{Valores{Nombre: "Corsa"}, Valores{Nombre: "Astra"}, Valores{Nombre: "A3"}, Valores{Nombre: "A5"}}}}
				json.NewEncoder(ctx).Encode(zr)
			}

		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}

}

func main() {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		//SSLHost:     "localhost:443", // This is optional in production. The default behavior is to just redirect the request to the HTTPS protocol. Example: http://github.com/some_page would be redirected to https://github.com/some_page.
	})

	secureHandler := secureMiddleware.Handler(requestHandler)

	// HTTP
	go func() {
		log.Fatal(fasthttp.ListenAndServe(":80", secureHandler))
	}()

	// HTTPS
	// To generate a development cert and key, run the following from your *nix terminal:
	// go run $GOROOT/src/pkg/crypto/tls/generate_cert.go --host="localhost"
	log.Fatal(fasthttp.ListenAndServeTLS(":443", "/etc/letsencrypt/live/www.redigo.cl/fullchain.pem", "/etc/letsencrypt/live/www.redigo.cl/privkey.pem", secureHandler))
}