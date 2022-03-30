package main

import (
	"bytes"
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
	"runtime"
	"strings"
	"syscall"
	"time"
	"io/ioutil"

	"github.com/fasthttp/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	FormPrecio      float64     `json:"FormPrecio"`
	TituloLista     string  `json:"TituloLista"`
	PageMod         string  `json:"PageMod"`
	DelAccion       string  `json:"DelAccion"`
	DelObj          string  `json:"DelObj"`
	Lista           []Lista `json:"FormDescripcion"`
	Dominio      	int  	`json:"Dominio"`
	AtencionPublico int 	`json:"AtencionPublico"`
	Copropiedad 	int 	`json:"Copropiedad"`
	Destino 		int 	`json:"Destino"`
	Detalle 		int 	`json:"Detalle"`
	P0  			bool 	`json:"P0"`
	P1  			bool 	`json:"P1"`
	P2  			bool 	`json:"P2"`
	P3  			bool 	`json:"P3"`
	P4  			bool 	`json:"P4"`
	P5  			bool 	`json:"P5"`
	P6  			bool 	`json:"P6"`
	P7  			bool 	`json:"P7"`
	P8  			bool 	`json:"P8"`
	P9  			bool 	`json:"P9"`
}
type TemplateInicio struct {
	Titulo string `json:"Titulo"`
}
type UfRes struct {
	Version string `json:"version"`
	Autor string `json:"autor"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
	Unidad_medida string `json:"unidad_medida"`
	Serie []UfSerie `json:"serie"`
}
type UfSerie struct {
	Fecha string `json:"fecha"`
	Valor float64 `json:"valor"`
}
type Lista struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}
type Data struct {
	Nombre string `json:"Nombre"`
	Direccion string `json:"Direccion"`
	Lat float64 `json:"Lat"`
	Lng float64 `json:"Lng"`
	Dominio int `json:"Dominio"`
	Precio float64 `json:"Precio"`
	AtencionPublico int `json:"AtencionPublico"`
	Copropiedad int `json:"Copropiedad"`
	Destino int `json:"Destino"`
	Detalle int `json:"Detalle"`
	P0 bool `json:"P0"`
	P1 bool `json:"P1"`
	P2 bool `json:"P2"`
	P3 bool `json:"P3"`
	P4 bool `json:"P4"`
	P5 bool `json:"P5"`
	P6 bool `json:"P6"`
	P7 bool `json:"P7"`
	P8 bool `json:"P8"`
	P9 bool `json:"P9"`
}
type PermisoUser struct {
	Bool  bool `json:"Bool"`
	Admin bool `json:"Admin"`
	Idemp bool `json:"Idemp"`
}
type Localidades struct {
	Paises []Pais `json:"Paises"`
	Regiones []Region `json:"Regiones"`
	Ciudades []Ciudad `json:"Ciudades"`
	Comunas []Comuna `json:"Comunas"`
	Propiedades []Propiedad `json:"Propiedades"`
	Titulo string `json:"Titulo"`
	SubTitulo string `json:"SubTitulo"`
	SubTitulo2 string `json:"SubTitulo2"`
	PaisesString string `json:"PaisesString"`
	RegionesString string `json:"RegionesString"`
	CiudadesString string `json:"CiudadesString"`
	ComunasString string `json:"ComunasString"`
	PropiedadesString string `json:"PropiedadesString"`
	PaisesCount int `json:"PaisesCount"`
	RegionesCount int `json:"RegionesCount"`
	CiudadesCount int `json:"CiudadesCount"`
	ComunasCount int `json:"ComunasCount"`
	PropiedadesCount int `json:"PropiedadesCount"`

}
type Propiedad struct {
	Id_pro int `json:"Id_pro"`
	Nombre string `json:"Nombre"`
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Direccion string `json:"Direccion"`
	Numero int `json:"Numero"`
	Id_com int `json:"Id_com"`
	Id_ciu int `json:"Id_ciu"`
	Id_reg int `json:"Id_reg"`
	Id_pai int `json:"Id_pai"`
}
type Pais struct {
	Id_pai int `json:"Id_pai"`
	Nombre string `json:"Nombre"`
}
type Region struct {
	Id_reg int `json:"Id_reg"`
	Nombre string `json:"Nombre"`
	Id_pai int `json:"Id_pai"`
}
type Ciudad struct {
	Id_ciu int `json:"Id_ciu"`
	Nombre string `json:"Nombre"`
	Id_reg int `json:"Id_reg"`
	Id_pai int `json:"Id_pai"`
}
type Comuna struct {
	Id_com int `json:"Id_com"`
	Nombre string `json:"Nombre"`
	Id_ciu int `json:"Id_ciu"`
	Id_reg int `json:"Id_reg"`
	Id_pai int `json:"Id_pai"`
}

var (
	imgHandler fasthttp.RequestHandler
	cssHandler fasthttp.RequestHandler
	jsHandler  fasthttp.RequestHandler
	port       string
)

func main() {

	SendEmail()
	//fmt.Println(GetUF())
	//SendEmail2()

	if runtime.GOOS == "windows" {
		imgHandler = fasthttp.FSHandler("C:/Go/Pelao/img", 1)
		cssHandler = fasthttp.FSHandler("C:/Go/Pelao/css", 1)
		jsHandler = fasthttp.FSHandler("C:/Go/Pelao/js", 1)
		port = ":81"
	} else {
		imgHandler = fasthttp.FSHandler("/var/Pelao/img", 1)
		cssHandler = fasthttp.FSHandler("/var/Pelao/css", 1)
		jsHandler = fasthttp.FSHandler("/var/Pelao/js", 1)
		port = ":80"
	}

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
		r.GET("/salir", Salir)
		r.GET("/SetEmpresa/{name}", SetEmpresa)
		fasthttp.ListenAndServe(port, r.Handler)
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func Save(ctx *fasthttp.RequestCtx) {

	resp := Response{}

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	token := string(ctx.Request.Header.Cookie("cu"))

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "guardar_empresa":

		nombre := string(ctx.FormValue("nombre"))
		precio := string(ctx.FormValue("precio"))
		if id == 0 {
			resp = InsertEmpresa(db, token, nombre, precio)
		}
		if id > 0 {
			resp = UpdateEmpresa(db, token, id, nombre, precio)
		}

	case "guardar_propiedad1":

		nombre := string(ctx.FormValue("nombre"))
		lat := string(ctx.FormValue("lat"))
		lng := string(ctx.FormValue("lng"))
		comuna := string(ctx.FormValue("comuna"))
		ciudad := string(ctx.FormValue("ciudad"))
		region := string(ctx.FormValue("region"))
		pais := string(ctx.FormValue("pais"))
		direccion := string(ctx.FormValue("direccion"))
		numero := string(ctx.FormValue("numero"))
		dominio := string(ctx.FormValue("dominio"))
		atencion_publico := string(ctx.FormValue("atencion_publico"))
		copropiedad := string(ctx.FormValue("copropiedad"))
		destino := string(ctx.FormValue("destino"))
		detalle_destino := string(ctx.FormValue("detalle_destino"))

		if id == 0 {
			resp = InsertPropiedad(db, token, nombre, lat, lng, comuna, ciudad, region, pais, direccion, numero, dominio, atencion_publico, copropiedad, destino, detalle_destino)
		}
		if id > 0 {
			resp = UpdatePropiedad1(db, token, id, nombre, lat, lng, comuna, ciudad, region, pais, direccion, numero, dominio, atencion_publico, copropiedad, destino, detalle_destino)
		}

	case "guardar_usuarios":

		nombre := string(ctx.FormValue("nombre"))
		p0 := string(ctx.FormValue("p0"))
		p1 := string(ctx.FormValue("p1"))
		p2 := string(ctx.FormValue("p2"))
		p3 := string(ctx.FormValue("p3"))
		p4 := string(ctx.FormValue("p4"))
		p5 := string(ctx.FormValue("p5"))
		p6 := string(ctx.FormValue("p6"))
		p7 := string(ctx.FormValue("p7"))
		p8 := string(ctx.FormValue("p8"))
		p9 := string(ctx.FormValue("p9"))

		if id == 0 {
			resp = InsertUsuario(db, token, nombre, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9)
		}
		if id > 0 {
			resp = UpdateUsuario(db, token, id, nombre, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9)
		}

	default:

	}

	json.NewEncoder(ctx).Encode(resp)
}
func Delete(ctx *fasthttp.RequestCtx) {

	resp := Response{}

	ctx.Response.Header.Set("Content-Type", "application/json")
	id := Read_uint32bytes(ctx.FormValue("id"))
	token := string(ctx.Request.Header.Cookie("cu"))
	
	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "borrar_empresa":
		resp = BorrarEmpresa(db, token, id)
	case "borrar_propiedad":
		resp = BorrarPropiedad(db, token, id)
	case "borrar_usuarios":
		resp = BorrarUsuario(db, token, id)
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
			cookieset := fmt.Sprintf("%s", cookie)

			stmt, err := db.Prepare("INSERT INTO sesiones(cookie, id_usr, fecha) VALUES(?,?, NOW())")
			ErrorCheck(err)
			defer stmt.Close()
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
	token := string(ctx.Request.Header.Cookie("cu"))

	switch name {
	case "inicioEmpresa":

		if found, _ := Permisos(token, 1); found {

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

		if SuperAdmin(token) {

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
				fmt.Println(aux.Precio)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormPrecio = aux.Precio
					obj.FormId = id
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	case "crearUsuarios":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Usuarios", "Subtitulo", "Subtitulo2", "Titulo Usuarios", "guardar_usuarios", fmt.Sprintf("/pages/%s", name), "borrar_usuario", "Usuario")
			lista, found := GetUsuarios(token)
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetUsuario(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.P0 = aux.P0
					obj.P1 = aux.P1
					obj.P2 = aux.P2
					obj.P3 = aux.P3
					obj.P4 = aux.P4
					obj.P5 = aux.P5
					obj.P6 = aux.P6
					obj.P7 = aux.P7
					obj.P8 = aux.P8
					obj.P9 = aux.P9

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	case "crearPropiedad":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Subtitulo", "Subtitulo2", "Titulo Lista", "guardar_propiedad1", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
			lista, found := GetPropiedades(token)
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetPropiedad(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id
					obj.Dominio = 1
					obj.AtencionPublico = 1
					obj.Copropiedad = 1
					obj.Destino = 1
					obj.Detalle = 1
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	case "buscarPropiedades":

		if found, id_emp := Permisos(token, 1); found {

			db, err := GetMySQLDB()
			defer db.Close()
			ErrorCheck(err)

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetLocalidades(db, id_emp)

			obj.PaisesCount = len(obj.Paises)
			obj.RegionesCount = len(obj.Regiones)
			obj.CiudadesCount = len(obj.Ciudades)
			obj.ComunasCount = len(obj.Comunas)
			obj.PropiedadesCount = len(obj.Propiedades)

			paises, err := json.Marshal(obj.Paises)
			ErrorCheck(err)
			obj.PaisesString = string(paises)

			regiones, err := json.Marshal(obj.Regiones)
			ErrorCheck(err)
			obj.RegionesString = string(regiones)

			ciudades, err := json.Marshal(obj.Ciudades)
			ErrorCheck(err)
			obj.CiudadesString = string(ciudades)

			comunas, err := json.Marshal(obj.Comunas)
			ErrorCheck(err)
			obj.ComunasString = string(comunas)

			propiedades, err := json.Marshal(obj.Propiedades)
			ErrorCheck(err)
			obj.PropiedadesString = string(propiedades)

			obj.Titulo = "Titulo"
			obj.SubTitulo = "Subtitulo"
			obj.SubTitulo2 = "Subtitulo2"

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}

	default:
		ctx.NotFound()
	}
}
func Index(ctx *fasthttp.RequestCtx) {

	SendEmail()
	//SendEmail2()
	ctx.SetContentType("text/html; charset=utf-8")
	token := string(ctx.Request.Header.Cookie("cu"))
	gpu := GetPermisoUser(token)

	if gpu.Bool {

		t, err := TemplatePage("html/inicio.html")
		ErrorCheck(err)
		err = t.Execute(ctx, gpu)
		ErrorCheck(err)

	} else {
		fmt.Fprintf(ctx, showFile("html/login.html"))
	}
}
func Salir(ctx *fasthttp.RequestCtx) {

	token := string(ctx.Request.Header.Cookie("cu"))
	//ctx.Response.Header.DelCookie("cu")
	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	delForm, err := db.Prepare("DELETE FROM sesiones WHERE cookie=?")
	ErrorCheck(err)
	delForm.Exec(token)
	defer db.Close()

	ctx.Redirect("/", 200)
}
func SetEmpresa(ctx *fasthttp.RequestCtx) {

	token := string(ctx.Request.Header.Cookie("cu"))

	if SuperAdmin(token) {

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		cn := 1
		id_emp := ctx.UserValue("name")

		stmt, err := db.Prepare("UPDATE usuarios SET id_emp = ? WHERE admin = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(id_emp, cn)
		ErrorCheck(e)
		if e == nil {
			ctx.Redirect("/", 200)
		}

	}
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
func GetPermisoUser(tkn string) PermisoUser {

	Pu := PermisoUser{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.admin, t1.id_emp FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr", tkn)
	defer res.Close()
	ErrorCheck(err)

	var admin int
	var id_emp int

	if res.Next() {

		err := res.Scan(&admin, &id_emp)
		ErrorCheck(err)

		if id_emp > 0 {
			Pu.Idemp = true
		} else {
			Pu.Idemp = false
		}

		Pu.Bool = true
		Pu.Admin = true

	} else {

		Pu.Bool = false
		Pu.Admin = false

	}
	return Pu
}

// FUNCTION DB //
func GetMySQLDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/pelao")
	return
	//CREATE DATABASE pelao CHARACTER SET utf8 COLLATE utf8_spanish2_ci;
}
func GetUser(token string) bool {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.id_usr FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr", token)
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
func Permisos(token string, n int) (bool, int) {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	sql := fmt.Sprintf("SELECT t1.p%v, t1.id_emp, t1.admin FROM usuarios t1, sesiones t2 WHERE t2.cookie = ? AND t2.id_usr=t1.id_usr", n)

	res, err := db.Query(sql, token)
	defer res.Close()
	ErrorCheck(err)

	var p int
	var id_emp int
	var admin int

	if res.Next() {

		err := res.Scan(&p, &id_emp, &admin)
		ErrorCheck(err)
		if p == 1 || admin == 1 {
			return true, id_emp
		}

	}
	return false, 0
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
	res, err := db.Query("SELECT nombre, precio FROM empresa WHERE id_emp = ? AND eliminado = ?", id, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var nombre string
		var precio float64
		err := res.Scan(&nombre, &precio)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = nombre
		data.Precio = precio
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
func GetIdEmp(token string) int {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t2.id_emp FROM sesiones t1, usuarios t2 WHERE t1.cookie = ? AND t2.id_usr=t1.id_usr", token)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.Next() {

		var id_emp int
		err := res.Scan(&id_emp)
		if err != nil {
			log.Fatal(err)
		}
		return id_emp

	} else {
		return 0
	}
}
func GetUsuarios(token string) ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	res, err := db.Query("SELECT id_usr, user FROM usuarios WHERE id_emp = ? AND eliminado = ? AND admin = ?", GetIdEmp(token), cn, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	if res.Next() {

		var id_usr int
		var user string
		err := res.Scan(&id_usr, &user)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Lista{Id: id_usr, Nombre: user})
		b = true

	}
	return data, b
}
func GetUsuario(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT user, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9 FROM usuarios WHERE id_usr = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var user string
		var p0 bool
		var p1 bool
		var p2 bool
		var p3 bool
		var p4 bool
		var p5 bool
		var p6 bool
		var p7 bool
		var p8 bool
		var p9 bool

		err := res.Scan(&user, &p0, &p1, &p2, &p3, &p4, &p5, &p6, &p7, &p8, &p9)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = user
		data.P0 = p0
		data.P1 = p1
		data.P2 = p2
		data.P3 = p3
		data.P4 = p4
		data.P5 = p5
		data.P6 = p6
		data.P7 = p7
		data.P8 = p8
		data.P9 = p9

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedades(token string) ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	res, err := db.Query("SELECT id_pro, nombre FROM propiedades WHERE id_emp = ? AND eliminado = ?", GetIdEmp(token), cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		var id_pro int
		var nombre string
		err := res.Scan(&id_pro, &nombre)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Lista{Id: id_pro, Nombre: nombre})
		b = true

	}
	return data, b
}
func GetPropiedad(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT nombre, direccion, lat, lng, dominio FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var nombre string
		var direccion string
		var lat float64
		var lng float64
		var dominio int
		err := res.Scan(&nombre, &direccion, &lat, &lng, &dominio)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = nombre
		data.Direccion = direccion
		data.Lat = lat
		data.Lng = lng
		data.Dominio = dominio
		return data, true

	} else {
		return data, false
	}
}

func GetUF() int {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)
	
	cn := 1
	res, err := db.Query("SELECT valor, ano, mes, dia FROM uf WHERE id = ?", cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var valor int
	var ano int
	var mes int
	var dia int

	if res.Next() {

		err := res.Scan(&valor, &ano, &mes, &dia)
		if err != nil {
			log.Fatal(err)
		}

		start := time.Date(ano, GetMonth(mes - 1), dia, 0, 0, 0, 0, time.UTC)
		duration := time.Now().Sub(start)
		if int(duration.Hours()/24) > 0 {

			val, ok := GetHttpUF()
			if ok {
				valor = val
				UpdateUF(val)
			}

		}  

	}
	return valor
}
func GetHttpUF() (int, bool) {

	req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    req.SetRequestURI("https://mindicador.cl/api/uf")

	resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(req, resp)
    if err != nil {
        fmt.Printf("Client get failed: %s\n", err)
        return 0, false
    }
	if resp.StatusCode() != fasthttp.StatusOK {
        fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
        return 0, false
    }
	var res UfRes
	body := resp.Body()
	
	if err := json.Unmarshal(body, &res); err == nil {
		return int(res.Serie[0].Valor), true
	}else{
		fmt.Println(err)
		return 0, false
	}

}
func UpdateUF(valor int){

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	year, month, day := time.Now().Date()

	id := 1
	stmt, err := db.Prepare("UPDATE uf SET valor = ?, ano = ?, mes = ?, dia = ? WHERE id = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(valor, year, int(month), day, id)
	ErrorCheck(e)

}

func GetLocalidades(db *sql.DB, id_emp int) Localidades {

	paises := []Pais{}
	regiones := []Region{}
	ciudades := []Ciudad{}
	comunas := []Comuna{}
	propiedades := []Propiedad{}

	res1, err := db.Query("SELECT DISTINCT(t1.id_pai), t1.nombre FROM paises t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_pai=t1.id_pai", id_emp)
	defer res1.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id_pai int
	var nombrepais string

	for res1.Next() {
		err := res1.Scan(&id_pai, &nombrepais)
		if err != nil {
			log.Fatal(err)
		}
		paises = append(paises, Pais{Id_pai: id_pai, Nombre: nombrepais})
	}

	res2, err := db.Query("SELECT DISTINCT(t1.id_reg), t1.nombre, t1.id_pai FROM regiones t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_reg=t1.id_reg", id_emp)
	defer res2.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id_reg int
	var nombreregion string

	for res2.Next() {
		err := res2.Scan(&id_reg, &nombreregion, &id_pai)
		if err != nil {
			log.Fatal(err)
		}
		regiones = append(regiones, Region{Id_reg: id_reg, Nombre: nombreregion, Id_pai: id_pai})
	}

	res3, err := db.Query("SELECT DISTINCT(t1.id_ciu), t1.nombre, t1.id_reg, t1.id_pai FROM ciudades t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_ciu=t1.id_ciu", id_emp)
	defer res3.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id_ciu int
	var nombreciudad string

	for res3.Next() {
		err := res3.Scan(&id_ciu, &nombreciudad, &id_reg, &id_pai)
		if err != nil {
			log.Fatal(err)
		}
		ciudades = append(ciudades, Ciudad{Id_ciu: id_ciu, Nombre: nombreciudad, Id_reg: id_reg, Id_pai: id_pai})
	}

	res4, err := db.Query("SELECT DISTINCT(t1.id_com), t1.nombre, t1.id_ciu, t1.id_reg, t1.id_pai FROM comunas t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_com=t1.id_com", id_emp)
	defer res4.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id_com int
	var nombrecomuna string

	for res4.Next() {
		err := res4.Scan(&id_com, &nombrecomuna, &id_ciu, &id_reg, &id_pai)
		if err != nil {
			log.Fatal(err)
		}
		comunas = append(comunas, Comuna{Id_com: id_com, Nombre: nombrecomuna, Id_ciu: id_ciu, Id_reg: id_reg, Id_pai: id_pai})
	}

	cn := 0
	res0, err := db.Query("SELECT id_pro, nombre, lat, lng, direccion, numero, id_com, id_ciu, id_reg, id_pai FROM propiedades WHERE eliminado = ? AND id_emp = ?", cn, id_emp)
	defer res0.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id_pro int
	var nombrepropiedad string
	var lat float64
	var lng float64
	var direccion string
	var numero int

	for res0.Next() {
		err := res0.Scan(&id_pro, &nombrepropiedad, &lat, &lng, &direccion, &numero, &id_com, &id_ciu, &id_reg, &id_pai)
		if err != nil {
			log.Fatal(err)
		}
		propiedades = append(propiedades, Propiedad{Id_pro: id_pro, Nombre: nombrepropiedad, Lat: lat, Lng: lng, Direccion: direccion, Numero: numero, Id_com: id_com, Id_ciu: id_ciu, Id_reg: id_reg, Id_pai: id_pai})
	}


	return Localidades{ Paises: paises, Regiones: regiones, Ciudades: ciudades, Comunas: comunas, Propiedades: propiedades }

}

func InsertPropiedad(db *sql.DB, token string, nombre string, lat string, lng string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, atencion_publico string, copropiedad string, destino string, detalle_destino string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {

		id_pai, b1 := GetPais(db, pais)
		id_reg, b2 := GetRegion(db, region, id_pai)
		id_ciu, b3 := GetCiudad(db, ciudad, id_pai, id_reg)
		id_com, b4 := GetComuna(db, comuna, id_pai, id_reg, id_ciu)

		if b1 && b2 && b3 && b4 {
			stmt, err := db.Prepare("INSERT INTO propiedades (nombre, lat, lng, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, atencion_publico, copropiedad, destino, detalle_destino, id_emp) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			ErrorCheck(err)
			defer stmt.Close()
			stmt.Exec(nombre, lat, lng, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, atencion_publico, copropiedad, destino, detalle_destino, id_emp)
			if err == nil {
				resp.Op = 1
				resp.Reload = 1
				resp.Page = "crearPropiedad"
				resp.Msg = "Propiedad ingresada correctamente"
			} else {
				resp.Msg = "La Propiedad no pudo ser ingresada"
			}
		}else{
			resp.Msg = "Error al ingresar posicion"
		}
		
	}else{
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad1(db *sql.DB, token string, id int, nombre string, lat string, lng string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, atencion_publico string, copropiedad string, destino string, detalle_destino string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET nombre = ?, lat = ?, lng = ?, ciudad = ?, comuna = ?, region = ?, pais = ?, direccion = ?, numero = ?, dominio = ?, atencion_publico = ?, copropiedad = ?, destino = ?, detalle_destino = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(nombre, lat, lng, ciudad, comuna, region, pais, direccion, numero, dominio, atencion_publico, copropiedad, destino, detalle_destino, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = "crearPropiedad"
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	}else{
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func BorrarPropiedad(db *sql.DB, token string, id int) Response {

	resp := Response{}
	if found, id_emp := Permisos(token, 1); found {
		del := 1
		stmt, err := db.Prepare("UPDATE propiedades SET eliminado = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(del, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Tipo = "success"
			resp.Reload = 1
			resp.Page = "crearPropiedad"
			resp.Titulo = "Propiedad eliminada"
			resp.Texto = "Propiedad eliminada correctamente"
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar propiedad"
			resp.Texto = "La propiedad no pudo ser eliminada"
		}
	}else{
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar propiedad"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}

func InsertUsuario(db *sql.DB, token string, nombre string, p0 string, p1 string, p2 string, p3 string, p4 string, p5 string, p6 string, p7 string, p8 string, p9 string) Response {

	resp := Response{}
	stmt, err := db.Prepare("INSERT INTO usuarios (user, id_emp, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	stmt.Exec(nombre, GetIdEmp(token), p0, p1, p2, p3, p4, p5, p6, p7, p8, p9)
	if err == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crearUsuarios"
		resp.Msg = "Usuario ingresada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "El usuario no pudo ser ingresada"
	}
	return resp
}
func UpdateUsuario(db *sql.DB, token string, id int, nombre string, p0 string, p1 string, p2 string, p3 string, p4 string, p5 string, p6 string, p7 string, p8 string, p9 string) Response {

	resp := Response{}
	stmt, err := db.Prepare("UPDATE usuarios SET user = ?, p0 = ?, p1 = ?, p2 = ?, p3 = ?, p4 = ?, p5 = ?, p6 = ?, p7 = ?, p8 = ?, p9 = ? WHERE id_usr = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(nombre, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, id)
	ErrorCheck(e)
	if e == nil {
		resp.Op = 1
		resp.Reload = 1
		resp.Page = "crearUsuarios"
		resp.Msg = "Usuario actualizada correctamente"
	} else {
		resp.Op = 2
		resp.Msg = "El usuario no pudo ser actualizada"
	}
	return resp
}
func BorrarUsuario(db *sql.DB, token string, id int) Response {

	resp := Response{}
	if SuperAdmin(token){
		del := 1
		stmt, err := db.Prepare("UPDATE usuarios SET eliminado = ? WHERE id_usr = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(del, id)
		ErrorCheck(e)
		if e == nil {
			resp.Tipo = "success"
			resp.Reload = 1
			resp.Page = "crearUsuarios"
			resp.Titulo = "Usuario eliminada"
			resp.Texto = "Usuario eliminada correctamente"
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar usuario"
			resp.Texto = "El usuario no pudo ser eliminada"
		}
	}else{
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar usuario"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}

func InsertEmpresa(db *sql.DB, token string, nombre string, precio string) Response {

	resp := Response{}
	resp.Op = 2
	if SuperAdmin(token){
		stmt, err := db.Prepare("INSERT INTO empresa (nombre, precio) VALUES (?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		stmt.Exec(nombre, precio)
		if err == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = "crearEmpresa"
			resp.Msg = "Empresa ingresada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser ingresada"
		}
	}else{
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdateEmpresa(db *sql.DB, token string, id int, nombre string, precio string) Response {

	resp := Response{}
	resp.Op = 2
	if SuperAdmin(token){
		stmt, err := db.Prepare("UPDATE empresa SET nombre = ?, precio = ? WHERE id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(nombre, precio, id)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = "crearEmpresa"
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	}else{
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func BorrarEmpresa(db *sql.DB, token string, id int) Response {

	resp := Response{}
	if SuperAdmin(token){
		del := 1
		stmt, err := db.Prepare("UPDATE empresa SET eliminado = ? WHERE id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(del, id)
		ErrorCheck(e)
		if e == nil {
			resp.Tipo = "success"
			resp.Reload = 1
			resp.Page = "crearEmpresa"
			resp.Titulo = "Empresa eliminada"
			resp.Texto = "Empresa eliminada correctamente"
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar empresa"
			resp.Texto = "La empresa no pudo ser eliminada"
		}
	}else{
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar empresa"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}

func GetPais(db *sql.DB, nombre string) (int64, bool) {

	res, err := db.Query("SELECT id_pai FROM paises WHERE nombre = ?", nombre)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		return id, true
	} else {
		stmt, err := db.Prepare("INSERT INTO paises (nombre) VALUES (?)")
		ErrorCheck(err)
		defer stmt.Close()
		r, err := stmt.Exec(nombre)
		if err == nil {
			idx, err := r.LastInsertId()
			if err == nil {
				return idx, true
			}else{
				return 0, false
			}
		} else {
			return 0, false
		}
	}
}
func GetRegion(db *sql.DB, nombre string, id_pai int64) (int64, bool) {

	res, err := db.Query("SELECT id_reg FROM regiones WHERE nombre = ? AND id_pai = ?", nombre, id_pai)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		return id, true
	} else {
		stmt, err := db.Prepare("INSERT INTO regiones (nombre, id_pai) VALUES (?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		r, err := stmt.Exec(nombre, id_pai)
		if err == nil {
			idx, err := r.LastInsertId()
			if err == nil {
				return idx, true
			}else{
				return 0, false
			}
		} else {
			return 0, false
		}
	}
}
func GetCiudad(db *sql.DB, nombre string, id_pai int64, id_reg int64) (int64, bool) {

	res, err := db.Query("SELECT id_ciu FROM ciudades WHERE nombre = ? AND id_pai = ? AND id_reg = ?", nombre, id_pai, id_reg)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		return id, true
	} else {
		stmt, err := db.Prepare("INSERT INTO ciudades (nombre, id_pai, id_reg) VALUES (?,?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		r, err := stmt.Exec(nombre, id_pai, id_reg)
		if err == nil {
			idx, err := r.LastInsertId()
			if err == nil {
				return idx, true
			}else{
				return 0, false
			}
		} else {
			return 0, false
		}
	}
}
func GetComuna(db *sql.DB, nombre string, id_pai int64, id_reg int64, id_ciu int64) (int64, bool) {

	res, err := db.Query("SELECT id_com FROM comunas WHERE nombre = ? AND id_pai = ? AND id_reg = ? AND id_ciu = ?", nombre, id_pai, id_reg, id_ciu)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		return id, true
	} else {
		stmt, err := db.Prepare("INSERT INTO comunas (nombre, id_pai, id_reg, id_ciu) VALUES (?,?,?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		r, err := stmt.Exec(nombre, id_pai, id_reg, id_ciu)
		if err == nil {
			idx, err := r.LastInsertId()
			if err == nil {
				return idx, true
			}else{
				return 0, false
			}
		} else {
			return 0, false
		}
	}
}

// FUNCTION DB //

// DAEMON //
func (h *MyHandler) StartDaemon() {
	h.Conf.Tiempo = 200 * time.Second
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
func GetMonth(m int) time.Month {

	var t time.Month

	switch m {
		case 0:
			t = time.January
		case 1:
			t = time.February
		case 2:
			t = time.March
		case 3:
			t = time.April
		case 4:
			t = time.May
		case 5:
			t = time.June
		case 6:
			t = time.July
		case 7:
			t = time.August
		case 8:
			t = time.September
		case 9:
			t = time.October
		case 10:
			t = time.November
		case 11:
			t = time.December
	}
	return t
}

type EmailData struct {
	FirstName string
	LastName  string
}
func getHTMLTemplate() string {
	var templateBuffer bytes.Buffer
	data := EmailData{
	   FirstName: "John",
	   LastName:  "Doe",
	}
	htmlData, err := ioutil.ReadFile("email/recuperar.html")
	htmlTemplate := template.Must(template.New("email.html").Parse(string(htmlData)))
	err = htmlTemplate.ExecuteTemplate(&templateBuffer, "email.html", data)
	if err != nil {
	   log.Fatal(err)
	   return ""
	}
	return templateBuffer.String()
}
func GenerateSESTemplate() (template *ses.SendEmailInput) {

	sender := "diego.gomez.bezmalinovic@gmail.com"
	receiver := "diego.gomez.bezmalinovic@gmail.com"
	html := getHTMLTemplate()
	title := "Sample Email"
	template = &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(receiver),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(sender),
	}
	return
}
func SendEmail() {

	region := "us-east-1"

	emailTemplate := GenerateSESTemplate()
	sess, err := session.NewSession(&aws.Config{
	   Region:      aws.String(region),
	})
	if err != nil {
	   log.Fatal(err)
	}
	service := ses.New(sess)
	_, err = service.SendEmail(emailTemplate)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Fatal(aerr.Error())
		} else {
			log.Fatal(err)
		}
	}
}

func SendEmail2(){

	region := "us-east-1"
	svc := ses.New(session.New(&aws.Config{ Region: aws.String(region) }))
	input := &ses.SendRawEmailInput{
		FromArn: aws.String(""),
		RawMessage: &ses.RawMessage{
			Data: []byte("From: diego.gomez.bezmalinovic@gmail.com\\nTo: diego.gomez.bezmalinovic@gmail.com\\nSubject: Test email (contains an attachment)\\nMIME-Version: 1.0\\nContent-type: Multipart/Mixed; boundary=\"NextPart\"\\n\\n--NextPart\\nContent-Type: text/plain\\n\\nThis is the message body.\\n\\n--NextPart\\nContent-Type: text/plain;\\nContent-Disposition: attachment; filename=\"attachment.txt\"\\n\\nThis is the text in the attachment.\\n\\n--NextPart--"),
		},
		ReturnPathArn: aws.String(""),
		Source:        aws.String(""),
		SourceArn:     aws.String(""),
	}

	result, err := svc.SendRawEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}