package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
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
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf      Config    `json:"Conf"`
	Passwords Passwords `json:"Passwords"`
	Path      string    `json:"Path"`
	File      string    `json:"File"`
}
type Passwords struct {
	PassDb      string    `json:"PassDb"`
	PassEmail   string    `json:"PassEmail"`
	Gmapkey     string    `json:"Gmapkey"`
	FechaCert   time.Time `json:"FechaCert"`
	UltimoEnvio time.Time `json:"UltimoEnvio"`
}

func main() {

	pass := &MyHandler{}

	if runtime.GOOS == "windows" {
		pass.File = "C:/Go/password_redigo.json"
		pass.Path = "C:/Go/Pelao_No_Git"
	} else {
		pass.File = "/var/password_redigo.json"
		pass.Path = "/var/Pelao"
	}

	passwords, err := os.ReadFile(pass.File)
	if err == nil {
		if err := json.Unmarshal(passwords, &pass.Passwords); err != nil {
			pass.InsertError(23, fmt.Errorf("Error ... Unmarshal datos de configuracion"))
			fmt.Println("Error ... Unmarshal datos de configuracion ", err)
		} else {
			pass.Passwords.FechaCert = time.Now()
			pass.SaveFile()
		}
	} else {
		pass.InsertError(24, fmt.Errorf("Error ... al leer archivo de configuracion"))
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
	h.Conf.Tiempo = 15 * time.Second
	if !h.Request() {
		h.StartProcess()
		h.EnviarError()
	}
	since := time.Since(h.Passwords.FechaCert)
	days := since.Seconds() / 86400
	if days > 175 && time.Now().Hour() == 3 {
		if err := h.SolicitarSSL(); err == nil {
			h.Passwords.FechaCert = time.Now()
			h.SaveFile()
		} else {
			h.InsertError(22, err)
			fmt.Println("ERROR AL RENOVAR EL CERTIFICADO: ", err)
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

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	r, err := client.Get("https://localhost/RCPG47D4F1AZS5")
	if err != nil {
		h.InsertError(17, err)
		fmt.Println("Error al realizar la solicitud HTTP:", err)
		return false
	}
	defer r.Body.Close() // Cerrar el cuerpo de la respuesta al finalizar la función

	// Verificar el código de estado de la respuesta
	if r.StatusCode != http.StatusOK {
		h.InsertError(18, err)
		fmt.Printf("Respuesta no exitosa. Código de estado: %d\n", r.StatusCode)
		return false
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.InsertError(19, err)
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return false
	}

	if len(bodyBytes) == 1 {
		if bodyBytes[0] == 1 {
			return true
		} else {
			h.InsertError(20, fmt.Errorf("Error Respuesta #1"))
			fmt.Println("Error en respuesta #1")
			return false
		}
	} else {
		h.InsertError(21, fmt.Errorf("Error Respuesta #2"))
		fmt.Println("Error en respuesta #2")
		return false
	}
}
func (h *MyHandler) StartProcess() {

	cmd := exec.Command("sh", fmt.Sprintf("%v/startprocess.sh", h.Path))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

	err := cmd.Start()
	if err != nil {
		h.InsertError(15, err)
		fmt.Println("ERROR AL INICIAR EL SUBPROCESO: ", err)
		return
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			h.InsertError(16, err)
			fmt.Println("ERROR PROCESS: ", err)
			return
		}
	}()
}
func (h *MyHandler) RestarProcess() {
	err := h.KillProcess()
	if err == nil {
		h.InsertError(14, err)
		h.StartProcess()
	} else {
		fmt.Println(err)
	}
}
func (h *MyHandler) EnviarError() error {

	if h.Passwords.UltimoEnvio.IsZero() || time.Since(h.Passwords.UltimoEnvio).Seconds() > 86400 {

		contenido, err := ioutil.ReadFile(fmt.Sprintf("%v/error.log", h.Path))
		if err != nil {
			h.InsertError(10, err)
			return fmt.Errorf("ERROR AL LEER EL ARCHIVO: %s\n", err)
		}

		archivo, err := os.OpenFile(fmt.Sprintf("%v/error.log", h.Path), os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			h.InsertError(11, err)
			return fmt.Errorf("ERROR AL ABRIR EL ARCHIVO PARA TRUNCAR: %s\n", err)
		}
		defer archivo.Close()

		err = archivo.Truncate(0)
		if err != nil {
			h.InsertError(12, err)
			return fmt.Errorf("ERROR AL VACIAR ARCHIVO: %s\n", err)
		}

		if err = h.SendEmail("diego.gomez.bezmalinovic@gmail.com", "Fatal Error", string(contenido)); err == nil {
			h.Passwords.UltimoEnvio = time.Now()
			h.SaveFile()
			fmt.Println("CORREO ENVIADO CON EXITO")
			return nil
		} else {
			h.InsertError(13, err)
			return fmt.Errorf("ERROR AL ENVIAR CORREO: %s\n", err)
		}

	} else {
		return fmt.Errorf("AUN NO PASAN 24hrs PARA ENVIAR OTRO CORREO \n")
	}
}
func (h *MyHandler) SendEmail(to string, subject string, body string) error {

	from := "redigocl@gmail.com"
	sub := fmt.Sprintf("From:%v\nTo:%v\nSubject:%v\n", from, to, subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, h.Passwords.PassEmail, "smtp.gmail.com"), from, []string{to}, []byte(sub+mime+body))
	if err != nil {
		h.InsertError(9, err)
		return err
	}
	return nil
}
func (h *MyHandler) SolicitarSSL() error {

	cmd := exec.Command("sh", fmt.Sprintf("%v/renewcert.sh", h.Path))
	_, err := cmd.Output()
	if err != nil {
		h.InsertError(6, err)
		return fmt.Errorf("ERROR AL EJECUTAR RENEWCERT: %s\n", err)
	}
	h.Passwords.FechaCert = time.Now()
	err = h.SaveFile()
	if err == nil {
		return nil
	} else {
		h.InsertError(7, err)
		return err
	}
}
func (h *MyHandler) SaveFile() error {

	archivo, err := os.Create(h.File)
	if err != nil {
		h.InsertError(4, err)
		return fmt.Errorf("ERROR AL CREAR EL ARCHIVO: %s\n", err)
	}
	defer archivo.Close()

	encoder := json.NewEncoder(archivo)
	err = encoder.Encode(h.Passwords)
	if err != nil {
		h.InsertError(5, err)
		return fmt.Errorf("ERROR AL CODIFICAR LA ESTRUCTURA: %s\n", err)
	}
	return nil
}
func (h *MyHandler) KillProcess() error {
	Pid := h.GetProcess()
	if Pid > 0 {
		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", Pid))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			h.InsertError(2, err)
			return fmt.Errorf("ERROR RUNNING KILL PROCESS: %s\n", err)
		}
		return nil
	} else {
		h.InsertError(3, fmt.Errorf("ERROR PID 0"))
		return fmt.Errorf("ERROR PID 0\n")
	}
}
func (h *MyHandler) GetProcess() int {

	netstatCmd := exec.Command("netstat", "-tulpn")
	var netstatOutput bytes.Buffer
	netstatCmd.Stdout = &netstatOutput
	if err := netstatCmd.Run(); err != nil {
		fmt.Println("Error ejecutando netstat:", err)
		h.InsertError(1, err)
		return 0
	}
	var processID string
	lines := strings.Split(netstatOutput.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, ":80") {
			fields := strings.Fields(line)
			processID = fields[len(fields)-1]
			break
		}
	}
	if processID != "" {
		return GetPNum(processID)
	} else {
		return 0
	}
}
func GetPNum(p string) int {
	var x int
	for _, c := range p {
		if c > 47 && c < 58 {
			x = x*10 + int(c-'0')
		}
		if c == 47 {
			return x
		}
	}
	return x
}
func (h *MyHandler) GetMySQLDB() (db *sql.DB, err error) {
	//CREATE DATABASE redigo CHARACTER SET utf8 COLLATE utf8_spanish2_ci;
	db, err = sql.Open("mysql", fmt.Sprintf("root:%v@tcp(127.0.0.1:3306)/redigo", h.Passwords.PassDb))
	return
}
func (h *MyHandler) InsertError(tipo int, errno error) {

	db, err := h.GetMySQLDB()
	ErrorCheck(err)
	stmt, err := db.Prepare("INSERT INTO errores (tipo, nombre, fecha) VALUES (?,?,Now())")
	ErrorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(tipo, errno)
	ErrorCheck(err)
}
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}
