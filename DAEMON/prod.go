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
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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
	h.Conf.Tiempo = 15 * time.Second
	if !Request() {
		h.StartProcess2()
	}
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

func Request() bool {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	r, err := client.Get("https://localhost/RCPG47D4F1AZS5")
	if err != nil {
		fmt.Println("Error al realizar la solicitud HTTP:", err)
		return false
	}
	defer r.Body.Close() // Cerrar el cuerpo de la respuesta al finalizar la función

	// Verificar el código de estado de la respuesta
	if r.StatusCode != http.StatusOK {
		fmt.Printf("Respuesta no exitosa. Código de estado: %d\n", r.StatusCode)
		return false
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return false
	}

	if len(bodyBytes) == 1 {
		if bodyBytes[0] == 1 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (h *MyHandler) StartProcess2() {

	cmd := exec.Command("prod")

	// Iniciar el comando
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error al iniciar el comando: %s\n", err)
		return
	}

	// Configurar un nuevo grupo de procesos para el proceso hijo
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Utilizar un canal para esperar a que el comando termine
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Esperar con un temporizador
	select {
	case <-time.After(time.Second * 10): // Puedes ajustar el tiempo límite según tus necesidades
		fmt.Println("El programa se ejecutó durante mucho tiempo. Terminándolo...")
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			fmt.Printf("Error al matar el proceso: %s\n", err)
		}
	case err := <-done:
		if err != nil {
			fmt.Printf("Error al esperar al comando: %s\n", err)
		} else {
			fmt.Println("El comando se ejecutó correctamente.")
		}
	}

}

func (h *MyHandler) StartProcess() {

	if !h.Running {

		cmd := exec.Command("prod", "&")
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
func SolicitarSSL() bool {
	cmd := exec.Command("certbot", "certonly", "--standalone", "-d", "redigo.cl", "-d", "www.redigo.cl", "--noninteractive", "--agree-tos", "--register-unsafely-without-email")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error al ejecutar el comando: %s\n", err)
		return false
	}
	fmt.Println(string(out))
	return true
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
