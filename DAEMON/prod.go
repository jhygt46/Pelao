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
	"net/smtp"
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
	Conf         Config    `json:"Conf"`
	Passwords    Passwords `json:"Passwords"`
	Path         string    `json:"Path"`
	File         string    `json:"File"`
	Debug        int       `json:"Debug"`
	UltimoEnvio  time.Time `json:"UltimoEnvio"`
	EnviarCorreo bool      `json:"EnviarCorreo"`
}
type Passwords struct {
	PassDb    string    `json:"PassDb"`
	PassEmail string    `json:"PassEmail"`
	Gmapkey   string    `json:"Gmapkey"`
	FechaCert time.Time `json:"FechaCert"`
	Pid       int       `json:"Pid"`
}

func main() {

	pass := &MyHandler{Debug: 0, EnviarCorreo: false}

	if runtime.GOOS == "windows" {
		pass.File = "C:/Go/password_redigo.json"
		pass.Path = "C:/Go/Pelao_No_Git"
	} else {
		pass.File = "/var/password_redigo.json"
		pass.Path = "/var/Pelao"
	}

	passwords, err := os.ReadFile(pass.File)
	if err == nil {
		fmt.Println("Ok ... Archivo de Configuracion leido correctamente")
		if err := json.Unmarshal(passwords, &pass.Passwords); err == nil {
			fmt.Println("Ok ... Unmarshal datos de configuracion")
			if pass.Passwords.Pid == 0 {
				pass.StartProcess()
			} else {
				if pass.Request() {
					fmt.Println("Ok ... Servicio Arriba, Reiniciando ...")
					pass.RestarProcess()
				} else {
					fmt.Println("Error ... Servicio Caido, Iniciando ...")
					pass.StartProcess()
				}
			}
		} else {
			fmt.Println("Error ... Unmarshal datos de configuracion ", err)
		}
	} else {
		fmt.Println("Error ... al leer archivo de configuracion", err)
	}

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
	if h.Debug == 2 {
		fmt.Println("FUNC StartDaemon")
	}
	h.Conf.Tiempo = 15 * time.Second

	if h.EnviarCorreo {
		if h.UltimoEnvio.IsZero() || time.Since(h.UltimoEnvio).Seconds() > 86400 {
			if h.EnviarError() {
				fmt.Println("CORREO ENVIADO")
				h.EnviarCorreo = false
			} else {
				fmt.Println("ERROR AL ENVIAR EL ARCHIVO AL CORREO")
			}
		}
	}

	if !h.Request() {
		h.StartProcess()

	} else {
		since := time.Since(h.Passwords.FechaCert)
		days := since.Seconds() / 86400
		if days > 175 && time.Now().Hour() == 3 {
			if h.SolicitarSSL() {
				h.Passwords.FechaCert = time.Now()
				h.SaveFile()
			}
		}
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
func (h *MyHandler) Request() bool {

	if h.Debug == 2 {
		fmt.Println("FUNC Request")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	r, err := client.Get("https://localhost/RCPG47D4F1AZS5")
	if err != nil {
		if h.Debug == 1 {
			fmt.Println("Error al realizar la solicitud HTTP:", err)
		}
		return false
	}
	defer r.Body.Close() // Cerrar el cuerpo de la respuesta al finalizar la función

	// Verificar el código de estado de la respuesta
	if r.StatusCode != http.StatusOK {
		if h.Debug == 1 {
			fmt.Printf("Respuesta no exitosa. Código de estado: %d\n", r.StatusCode)
		}
		return false
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if h.Debug == 1 {
			fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		}
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
func (h *MyHandler) StartProcess() {

	if h.Debug == 2 || h.Debug == 3 {
		fmt.Println("FUNC StartProcess")
	}

	cmd := exec.Command("sh", fmt.Sprintf("%v/startprocess.sh", h.Path))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

	err := cmd.Start()
	if err != nil {
		if h.Debug == 1 || h.Debug == 3 {
			fmt.Println("Error al iniciar el subproceso:", err)
		}
		return
	}

	fmt.Printf("START PROCESS - PID ANTIGUO (%v) - ", h.Passwords.Pid)

	h.Passwords.Pid = cmd.Process.Pid
	h.SaveFile()

	fmt.Printf("PID NUEVO (%v) \n", cmd.Process.Pid)

	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Println("Error Process: ", err)
			if h.Debug == 1 || h.Debug == 3 {
				fmt.Println("Error al esperar a que el subproceso termine:", err)
			}
			return
		}
		if h.Debug == 1 || h.Debug == 3 {
			fmt.Println("Subproceso completado.")
		}
	}()
}
func (h *MyHandler) RestarProcess() {
	h.KillProcess()
	h.StartProcess()
}
func (h *MyHandler) KillProcess() bool {

	if h.Debug == 2 || h.Debug == 3 {
		fmt.Println("FUNC KillProcess")
	}

	if h.Passwords.Pid > 0 {
		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", h.Passwords.Pid))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			if h.Debug == 1 || h.Debug == 3 {
				fmt.Printf("Error al enviar la señal SIGKILL: %s\n", err)
			}
			return false
		}
		h.Passwords.Pid = 0
		h.SaveFile()
		return true
	}

	return true
}
func (h *MyHandler) SolicitarSSL() bool {
	fmt.Println("SolicitarSSL RE-NEWCERT")
	cmd := exec.Command("sh", fmt.Sprintf("%v/renewcert.sh", h.Path))
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Error al ejecutar el script:", err)
		return false
	}
	h.Passwords.FechaCert = time.Now()
	h.SaveFile()
	return true
}

func (h *MyHandler) SaveFile() bool {

	archivo, err := os.Create(h.File)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return false
	}
	defer archivo.Close()

	encoder := json.NewEncoder(archivo)
	err = encoder.Encode(h.Passwords)
	if err != nil {
		fmt.Println("Error al codificar la estructura:", err)
		return false
	}
	return true
}
func (h *MyHandler) SendEmail(to string, subject string, body string) bool {

	from := "redigocl@gmail.com"
	sub := fmt.Sprintf("From:%v\nTo:%v\nSubject:%v\n", from, to, subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, h.Passwords.PassEmail, "smtp.gmail.com"), from, []string{to}, []byte(sub+mime+body))
	if err != nil {
		return false
	}
	fmt.Println("CORREO ENVIADO A ", to)
	return true
}
func (h *MyHandler) EnviarError() bool {

	contenido, err := ioutil.ReadFile(fmt.Sprintf("%v/error.log", h.Path))
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return false
	}

	archivo, err := os.OpenFile(fmt.Sprintf("%v/error.log", h.Path), os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error al abrir el archivo para truncar:", err)
		return false
	}
	defer archivo.Close()

	err = archivo.Truncate(0)
	if err != nil {
		fmt.Println("Error al truncar el archivo:", err)
		return false
	}

	return h.SendEmail("diego.gomez.bezmalinovic@gmail.com", "Fatal Error", string(contenido))
}
