package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf      Config    `json:"Conf"`
	Passwords Passwords `json:"Passwords"`
}
type Passwords struct {
	PassDb    string `json:"PassDb"`
	PassEmail string `json:"PassEmail"`
	Gmapkey   string `json:"Gmapkey"`
	FechaCert string `json:"FechaCert"`
}

func main() {

	pass := &MyHandler{}
	var file string

	if runtime.GOOS == "windows" {
		file = "C:/Go/password_redigo.json"
	} else {
		file = "/var/password_redigo.json"
	}

	passwords, err := os.ReadFile(file)
	if err == nil {
		fmt.Println("Ok ... Archivo de Configuracion leido correctamente")
		if err := json.Unmarshal(passwords, &pass.Passwords); err == nil {
			fmt.Println("Ok ... Unmarshal datos de configuracion")
		} else {
			fmt.Println("Error ... Unmarshal datos de configuracion")
		}
	} else {
		fmt.Println("Error ... al leer archivo de configuracion")
	}

	Restart()

	con := context.Background()
	con, cancel := context.WithCancel(con)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					pass.Conf.init()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-con.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 15 * time.Second
	//Request()
	fmt.Println("DAEMON")
}
func (c *Config) init() {
	var tick = flag.Duration("tick", 1*time.Second, "Ticking interval")
	c.Tiempo = *tick
}
func run(con context.Context, c *MyHandler, stdout io.Writer) error {
	c.Conf.init()
	log.SetOutput(os.Stdout)
	for {
		select {
		case <-con.Done():
			return nil
		case <-time.Tick(c.Conf.Tiempo):
			c.StartDaemon()
		}
	}
}

func Request() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	data := url.Values{"code": {"Fs1R5MKsL94AmT2zXd"}}
	body := strings.NewReader(data.Encode())
	r, err := client.Post("https://localhost/RCPG47D4F1AZS5", "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Println(err)
		Restart()
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

}
func Restart() {
	// Comando a ejecutar

	cmd := exec.Command("./prod 2> error.log &")
	cmd.Dir = "/var/Pelao"

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Captura la salida estándar del comando
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convierte la salida en una cadena y muestra el resultado
	outputStr := string(stdout)
	fmt.Println("Salida del comando:")
	fmt.Println(outputStr)
}

/*
func SendEmail(to string, subject string, body string) bool {

	from := "redigocl@gmail.com"
	sub := fmt.Sprintf("From:%v\nTo:%v\nSubject:%v\n", from, to, subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass.Passwords.PassEmail, "smtp.gmail.com"), from, []string{to}, []byte(sub+mime+body))
	if err != nil {
		return false
	}
	return true
}
*/
