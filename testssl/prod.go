package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mithorium/secure-fasthttp"
	"github.com/valyala/fasthttp"
)

type Idioma struct {
	Page int
	Text []string
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n")

	if string(ctx.Method()) == "GET" {
		switch string(ctx.Path()) {
		case "/autoCuad":
			/*
				var p []int32
				if err := json.Unmarshal(ctx.QueryArgs().Peek("c"), &p); err == nil {

					var bn []int32
					var key []byte
					var b strings.Builder
					b.Write([]byte{91})

					bn = p[0:2]
					key = GetKey2(bn, ParamBytes(ctx.QueryArgs().Peek("u")))
					val, _ := h.Db.Get(key)
					if len(val) > 0 {
						WriteResponse(&val, 0, &b, p[2:len(p)])
					} else {
						fmt.Println("NOT FOUND DB-CUAD KEY", key)
					}

					b.Write([]byte{93})
					fmt.Fprint(ctx, b.String())

				}
			*/
		case "/lang":

			//now := time.Now()
			val := string(ctx.QueryArgs().Peek("c"))
			if val == "zr" {
				zr := []Idioma{Idioma{Page: 1, Text: []string{"apple", "peach", "pear"}}, Idioma{Page: 1, Text: []string{"apple", "peach", "pear"}}}
				zr_txt, _ := json.Marshal(zr)
				fmt.Fprint(ctx, string(zr_txt))
			}

			//fmt.Println("time elapse:", time.Since(now))

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
