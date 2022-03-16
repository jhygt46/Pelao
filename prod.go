package main

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

type Response struct {
	Op     uint8  `json:"Op"`
	Msg    string `json:"Msg"`
	Reload int    `json:"Reload"`
	Page   string `json:"Page"`
	Tipo   string `json:"Tipo"`
	Titulo string `json:"Titulo"`
	Texto  string `json:"Texto"`
}
type Giros struct {
	Titulo string `json:"Titulo"`
}
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type MyHandler struct {
	Conf Config `json:"Conf"`
}

type TemplateConf struct {
	Titulo          string  `json:"Titulo"`
	SubTitulo       string  `json:"SubTitulo"`
	SubTitulo2      string  `json:"SubTitulo"`
	FormId          int     `json:"FormId"`
	FormAccion      string  `json:"FormAccion"`
	FormNombre      string  `json:"FormNombre"`
	FormDescripcion string  `json:"FormDescripcion"`
	TituloLista     string  `json:"TituloLista"`
	PageMod         string  `json:"PageMod"`
	DelAccion       string  `json:"DelAccion"`
	DelObj          string  `json:"DelObj"`
	Lista           []Lista `json:"FormDescripcion"`
}
type TemplateInicio struct {
	Titulo string `json:"Titulo"`
}
type Lista struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}
type Data struct {
	Nombre string `json:"Nombre"`
}

var (
	imgPrefix = []byte("/img/")
	//imgHandler = fasthttp.FSHandler("/var/Go/Pelao/img", 1)
	imgHandler = fasthttp.FSHandler("C:/Pelao/img", 1)

	cssPrefix = []byte("/css/")
	//cssHandler = fasthttp.FSHandler("/var/Go/Pelao/css", 1)
	cssHandler = fasthttp.FSHandler("C:/Pelao/css", 1)

	jsPrefix = []byte("/js/")
	//jsHandler = fasthttp.FSHandler("/var/Go/Pelao/js", 1)
	jsHandler = fasthttp.FSHandler("C:/Pelao/js", 1)
)

func main() {

	pass := &MyHandler{Conf: Config{}}
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
	go func() {
		r := router.New()
		r.GET("/", Index)
		r.GET("/css/{name}", Css)
		r.GET("/js/{name}", Js)
		r.GET("/img/{name}", Img)
		r.GET("/pages/{name}", Pages)
		r.POST("/login", Login)
		r.POST("/save", Save)
		r.POST("/delete", Delete)
		r.POST("/salir", Salir)
		fasthttp.ListenAndServe(":80", r.Handler)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Save(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	resp := Response{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "guardar_empresa":
		nombre := string(ctx.FormValue("nombre"))
		if id == 0 {
			resp = InsertEmpresa(db, nombre)
		}
		if id > 0 {
			resp = UpdateEmpresa(db, nombre, id)
		}
	case "guardar_propiedades":
		nombre := string(ctx.FormValue("nombre"))
		if id == 0 {
			resp = InsertEmpresa(db, nombre)
		}
		if id > 0 {
			resp = UpdateEmpresa(db, nombre, id)
		}
	default:

	}

	json.NewEncoder(ctx).Encode(resp)
}
func Delete(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	resp := Response{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "borrar_empresa":

		resp = BorrarEmpresa(db, id)

	default:

	}

	json.NewEncoder(ctx).Encode(resp)
}
func Login(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	resp := Response{Op: 2}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	user := string(ctx.PostArgs().Peek("user"))

	res, err := db.Query("SELECT id_usr, pass, id_emp, admin FROM usuarios WHERE user = ? AND eliminado=0", user)
	defer res.Close()
	ErrorCheck(err)

	if res.Next() {

		var id_usr int
		var pass string
		var id_emp int
		var admin int
		err := res.Scan(&id_usr, &pass, &id_emp, &admin)
		ErrorCheck(err)

		if pass == GetMD5Hash(ctx.PostArgs().Peek("pass")) {

			resp.Op = 1
			resp.Msg = ""
			cookie := randSeq(32)
			cookieset := fmt.Sprintf("%s%v", cookie, id_emp)

			stmt, err := db.Prepare("INSERT INTO sesiones(cookie, id_usr, fecha) VALUES(?,?, NOW())")
			ErrorCheck(err)
			stmt.Exec(cookie, id_usr)

			authcookie := CreateCookie("cu", cookieset, 94608000)
			ctx.Response.Header.SetCookie(authcookie)

		} else {
			resp.Msg = "Usuario Contraseña no existen"
		}

	} else {
		resp.Msg = "Usuario Contraseña no existen"
	}

	json.NewEncoder(ctx).Encode(resp)
}
func Pages(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	name := ctx.UserValue("name")

	switch name {
	case "inicioEmpresa":

		if Permisos(string(ctx.Request.Header.Cookie("cu")), 1) {

			id_emp := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)
			obj := TemplateInicio{}
			aux, found := GetEmpresa(id_emp)
			if found {
				obj.Titulo = aux.Nombre
			}
			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	case "crearEmpresa":

		if SuperAdmin(string(ctx.Request.Header.Cookie("cu"))) {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Empresa", "Subtitulo", "Subtitulo2", "Titulo Lista", "guardar_empresa", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
			lista, found := GetEmpresas()
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetEmpresa(id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	case "crearPropiedad":

		if SuperAdmin(string(ctx.Request.Header.Cookie("cu"))) {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Subtitulo", "Subtitulo2", "Titulo Lista", "guardar_empresa", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
			lista, found := GetEmpresas()
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetEmpresa(id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearMarkers":

		if SuperAdmin(string(ctx.Request.Header.Cookie("cu"))) {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Subtitulo", "Subtitulo2", "Titulo Lista", "guardar_empresa", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
			lista, found := GetEmpresas()
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetEmpresa(id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	default:
		ctx.NotFound()
	}
}
func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	token := string(ctx.Request.Header.Cookie("cu"))
	if len(token) > 32 && GetUser(token) {
		fmt.Fprintf(ctx, showFile("html/inicio.html"))
	} else {
		fmt.Fprintf(ctx, showFile("html/login.html"))
	}
}
func Salir(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DelCookie("cu")
}
func Js(ctx *fasthttp.RequestCtx) {
	jsHandler(ctx)
}
func Css(ctx *fasthttp.RequestCtx) {
	cssHandler(ctx)
}
func Img(ctx *fasthttp.RequestCtx) {
	imgHandler(ctx)
}

// FUNCTION DB //
func GetMySQLDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/pelao")
	return
	//CREATE DATABASE pelao CHARACTER SET utf8 COLLATE utf8_spanish2_ci;
}
func GetUser(token string) bool {

	fmt.Println(token)

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.id_usr FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr", token[0:32])
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		return true
	} else {
		return false
	}
}
func Permisos(token string, n int) bool {

	tkn := token[0:32]
	id_emp, err := strconv.Atoi(token[32:len(token)])
	ErrorCheck(err)

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.id_emp, t1.admin FROM usuarios t1, sesiones t2, usuario_perfil t3, perfil_tarea t4 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.id_usr=t3.id_usr AND t3.id_per=t4.id_per AND t4.id_tar=?", tkn, n)
	defer res.Close()
	ErrorCheck(err)

	var id int
	var admin int

	if res.Next() {

		err := res.Scan(&id, &admin)
		ErrorCheck(err)
		if id == id_emp || admin == 1 {
			return true
		}

	} else {

		res3, err3 := db.Query("SELECT * FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.admin=1", tkn)
		defer res3.Close()
		ErrorCheck(err3)
		if res3.Next() {
			return true
		} else {

			res2, err2 := db.Query("SELECT t1.id_emp, t1.admin FROM usuarios t1, sesiones t2, usuario_tarea t3 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.id_usr=t3.id_usr AND t3.id_tar=?", tkn, n)
			defer res2.Close()
			ErrorCheck(err2)
			if res2.Next() {

				err := res.Scan(&id, &admin)
				ErrorCheck(err)
				if id == id_emp {
					return true
				}

			} else {
				return false
			}

		}

	}
	return false
}
func SuperAdmin(token string) bool {

	tkn := token[0:32]

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT * FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.admin=1", tkn)
	defer res.Close()
	ErrorCheck(err)
	if res.Next() {
		return true
	} else {
		return false
	}
}
func GetEmpresa(id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT nombre FROM empresa WHERE id_emp = ? AND eliminado = ?", id, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var nombre string
		err := res.Scan(&nombre)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = nombre
		return data, true

	} else {
		return data, false
	}
}
func GetEmpresas() ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT id_emp, nombre FROM empresa WHERE eliminado = ?", cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var nombre string

	for res.Next() {

		err := res.Scan(&id, &nombre)
		ErrorCheck(err)
		data = append(data, Lista{Id: id, Nombre: nombre})
		b = true

	}
	return data, b
}

func InsertEmpresa(db *sql.DB, nombre string) Response {

	resp := Response{}
	stmt, err := db.Prepare("INSERT INTO empresa (nombre) VALUES (?)")
	ErrorCheck(err)
	stmt.Exec(nombre)
	if err == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crear"
		resp.Msg = "Empresa ingresada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "La Empresa no pudo ser ingresada"
	}
	return resp
}
func UpdateEmpresa(db *sql.DB, nombre string, id int) Response {

	resp := Response{}
	stmt, err := db.Prepare("UPDATE empresa SET nombre = ? WHERE id_emp = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(nombre, id)
	ErrorCheck(e)
	if e == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crear"
		resp.Msg = "Empresa actualizada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "La Empresa no pudo ser actualizada"
	}
	return resp
}
func BorrarEmpresa(db *sql.DB, id int) Response {

	del := 1
	resp := Response{}
	stmt, err := db.Prepare("UPDATE empresa SET eliminado = ? WHERE id_emp = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(del, id)
	ErrorCheck(e)
	if e == nil {
		resp.Tipo = "success"
		resp.Reload = 1
		resp.Page = "crear"
		resp.Titulo = "Empresa eliminada"
		resp.Texto = "Empresa eliminada correctamente"
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar empresa"
		resp.Texto = "La empresa no pudo ser eliminada"
	}
	return resp
}

// FUNCTION DB //

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 2 * time.Second
	//fmt.Println("DAEMON")
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

// DAEMON //
func TemplatePage(v string) (*template.Template, error) {

	t, err := template.ParseFiles(v)
	if err != nil {
		log.Print(err)
		return t, err
	}
	return t, nil
}
func Read_uint32bytes(data []byte) int {
	var x int
	for _, c := range data {
		x = x*10 + int(c-'0')
	}
	return x
}
func GetMD5Hash(text []byte) string {
	hasher := md5.New()
	hasher.Write(text)
	return hex.EncodeToString(hasher.Sum(nil))
}
func CreateCookie(key string, value string, expire int) *fasthttp.Cookie {
	if strings.Compare(key, "") == 0 {
		key = "GoLog-Token"
	}
	fmt.Println("CreateCookie | Key: ", key, " | Val: ", value)
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(key)
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}
func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func showFile(file string) string {

	dat, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
}
func GetTemplateConf(titulo string, subtitulo string, subtitulo2 string, titulolista string, formaccion string, pagemod string, delaccion string, delobj string) TemplateConf {
	return TemplateConf{Titulo: titulo, SubTitulo: subtitulo, SubTitulo2: subtitulo2, TituloLista: titulolista, FormAccion: formaccion, PageMod: pagemod, DelAccion: delaccion, DelObj: delobj}
}
