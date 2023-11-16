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
	Path      string    `json:"Path"`
	Pid       int       `json:"Pid"`
	Running   bool      `json:"Start"`
}
type Passwords struct {
	PassDb    string `json:"PassDb"`
	PassEmail string `json:"PassEmail"`
	Gmapkey   string `json:"Gmapkey"`
	FechaCert string `json:"FechaCert"`
}

func main() {

	pass := &MyHandler{Running: false}
	var file string

	if runtime.GOOS == "windows" {
		file = "C:/Go/password_redigo.json"
		pass.Path = "C:/Go/Pelao_No_Git"
	} else {
		file = "/var/password_redigo.json"
		pass.Path = "/var/Pelao"
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

	//Request()
	//pass.StartProcess()

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
	h.Conf.Tiempo = 5 * time.Second
	fmt.Println("DAEMON")
	Request()
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
		//StartProcess()
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

}
func (h *MyHandler) StartProcess() {

	if !h.Running {

		cmd := exec.Command("prod")
		errorFile, err := os.Create(fmt.Sprintf("%s/error.log", h.Path))
		if err != nil {
			fmt.Println("Error al abrir el archivo de error:", err)
			return
		}
		defer errorFile.Close()
		cmd.Stderr = errorFile

		err = cmd.Start()
		if err != nil {
			fmt.Println("Error al iniciar el comando:", err)
			return
		}
		h.Pid = cmd.Process.Pid
		h.Running = true
		fmt.Println("Comando ejecutándose en segundo plano. PID:", cmd.Process.Pid)

		// Espera a que el comando termine
		err = cmd.Wait()
		if err != nil {
			fmt.Println("Error al esperar a que el comando termine:", err)
		}

	} else {
		fmt.Println("Error el programa ya esta corriendo")
	}
}
func (h *MyHandler) KillProcess() {

	if h.Running {

		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", h.Pid))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error al enviar la señal SIGKILL: %s\n", err)
			return
		}

		h.Running = false
		fmt.Printf("Señal SIGKILL enviada al proceso con PID %d\n", h.Pid)

	} else {
		fmt.Println("Error el programa no esta corriendo")
	}
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
