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
	"image"
	"image/color"
	"image/png"
	"net/smtp"
	"reflect"

	//"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/image/draw"

	"github.com/fasthttp/router"
	_ "github.com/go-sql-driver/mysql"

	"github.com/valyala/fasthttp"

	qrcode "github.com/skip2/go-qrcode"

	col "github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	"github.com/xuri/excelize/v2"
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
type Config struct {
	Tiempo time.Duration `json:"Tiempo"`
}
type ListaPropAlert struct {
	Id_pro int `json:"Id_pro"`
	Id_emp int `json:"Id_emp"`
}
type Cotizacion struct {
	Op        uint8             `json:"Op"`
	Lista     []ListaCotizacion `json:"Lista"`
	Total     float32           `json:"Total"`
	TotalUf   float32           `json:"TotalUf"`
	Subtotal  float32           `json:"Subtotal"`
	Iva       float32           `json:"Iva"`
	Uf        float32           `json:"Uf"`
	Fecha     string            `json:"Fecha"`
	NombreEmp string            `json:"NombreEmp"`
	IdCot     int               `json:"IdCot"`
}
type ListaCotizacion struct {
	Propiedad   string  `json:"Propiedad"`
	NombreAle   string  `json:"NombreAle"`
	Descripcion string  `json:"Descripcion"`
	Precio      float32 `json:"Precio"`
	IdAle       int     `json:"IdAle"`
	IdPro       int     `json:"IdPro"`
}
type MyHandler struct {
	Conf            Config           `json:"Conf"`
	ListaPropAlerts []ListaPropAlert `json:"ListaPropAlerts"`
	DeleteFiles     []string         `json:"DeleteFiles"`
	Passwords       Passwords        `json:"Passwords"`
}
type Passwords struct {
	PassDb    string `json:"PassDb"`
	PassEmail string `json:"PassEmail"`
	Gmapkey   string `json:"Gmapkey"`
}
type TemplateConf struct {
	CamposPropiedades CamposPropiedades `json:"Titulo"`
	Alertas           Alertas           `json:"Alerta"`

	Id_emp          int     `json:"Id_emp"`
	Titulo          string  `json:"Titulo"`
	SubTitulo       string  `json:"SubTitulo"`
	SubTitulo2      string  `json:"SubTitulo2"`
	FormId          int     `json:"FormId"`
	FormIdRec       int     `json:"FormIdRec"`
	FormAccion      string  `json:"FormAccion"`
	FormNombre      string  `json:"FormNombre"`
	FormDescripcion string  `json:"FormDescripcion"`
	FormPrecio      float64 `json:"FormPrecio"`
	TituloLista     string  `json:"TituloLista"`
	PageMod         string  `json:"PageMod"`
	DelAccion       string  `json:"DelAccion"`
	DelObj          string  `json:"DelObj"`
	Lista           []Lista `json:"Lista"`
	Lista2          []Lista `json:"Lista2"`
	Lista3          []Lista `json:"Lista3"`

	NextPage int `json:"NextPage"`

	Is_Arrendado int `json:"Is_Arrendado"`

	Descripcion string `json:"Descripcion"`

	Valor int `json:"Valor"`

	ValorCampo string `json:"ValorCampo"`
	FormIdAle  int    `json:"FormIdAle"`
	PrecioUf   int    `json:"PrecioUf"`

	Permisos ListaPermisos `json:"Permisos"`
}
type ListaPermisos struct {
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
type TemplateInicio struct {
	Nombre   string        `json:"Nombre"`
	Precio   float64       `json:"Precio"`
	UF       int           `json:"UF"`
	Resp     Resumen       `json:"Resp"`
	Permisos ListaPermisos `json:"Permisos"`
}
type UfRes struct {
	Version       string    `json:"version"`
	Autor         string    `json:"autor"`
	Codigo        string    `json:"codigo"`
	Nombre        string    `json:"nombre"`
	Unidad_medida string    `json:"unidad_medida"`
	Serie         []UfSerie `json:"serie"`
}
type UfSerie struct {
	Fecha string  `json:"fecha"`
	Valor float64 `json:"valor"`
}
type Lista struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
	Aux    int    `json:"Aux"`
}
type Empresa struct {
	Nombre string  `json:"Nombre"`
	Precio float64 `json:"Precio"`
}
type PermisoUser struct {
	Bool          bool   `json:"Bool"`
	Admin         bool   `json:"Admin"`
	Idemp         bool   `json:"Idemp"`
	Id_usr        int    `json:"Idusr"`
	Id_emp        int    `json:"Id_emp"`
	Pass_impuesta int    `json:"Pass_impuesta"`
	Function      string `json:"Function"`
}
type CamposPropiedades struct {
	Id                   int              `json:"Id"`
	Nombre               string           `json:"Nombre"`
	Direccion            string           `json:"Direccion"`
	Lat                  float64          `json:"Lat"`
	Lng                  float64          `json:"Lng"`
	RangoPos             string           `json:"RangoPos"`
	Numero               string           `json:"Numero"`
	Dominio              int              `json:"Dominio"`
	Dominio2             int              `json:"Dominio2"`
	Atencion_publico     int              `json:"AtencionPublico"`
	Copropiedad          int              `json:"Copropiedad"`
	Destino              int              `json:"Destino"`
	Detalle_destino      int              `json:"Detalle"`
	Detalle_destino_otro string           `json:"DetalleDestinoOtro"`
	Imagenes             map[int][]string `json:"Imagenes"`
	Id_com               int              `json:"Id_com"`
	Id_ciu               int              `json:"Id_ciu"`
	Id_reg               int              `json:"Id_reg"`
	Id_pai               int              `json:"Id_pai"`

	P1 int `json:"P1"`
	P2 int `json:"P2"`
	P3 int `json:"P3"`
	P4 int `json:"P4"`
	P5 int `json:"P5"`
	P6 int `json:"P6"`
	P7 int `json:"P7"`
	P8 int `json:"P8"`

	Is_Arrendado int `json:"Is_Arrendado"`

	PermisosEdificacion      []PermisoEdificacion    `json:"PermisosEdificacion"`
	PermisosEdificacionIndex PermisoEdificacion      `json:"PermisosEdificacionIndex"`
	ListaArchivos            map[int][]ListaArchivos `json:"ListaArchivos"`
	UltimosArchivos          map[int]ListaArchivos   `json:"UltimosArchivos"`

	Electrico_te1        int `json:"Electrico_te1"`
	Dotacion_ap          int `json:"Dotacion_ap"`
	Dotacion_alcance     int `json:"Dotacion_alcance"`
	Instalacion_ascensor int `json:"Instalacion_ascensor"`
	Te1_ascensor         int `json:"Te1_ascensor"`
	Certificado_ascensor int `json:"Certificado_ascensor"`
	Clima                int `json:"Clima"`
	Seguridad_incendio   int `json:"Seguridad_incendio"`

	Tasacion_valor_comercial string `json:"Tasacion_valor_comercial"`
	Ano_tasacion             string `json:"Ano_tasacion"`
	Contrato_arriendo        int    `json:"Contrato_arriendo"`
	Contrato_subarriendo     int    `json:"Contrato_subarriendo"`

	Nom_propietario_conservador string `json:"Nom_propietario_conservador"`
	Posee_gp                    int    `json:"Posee_gp"`
	Posee_ap                    int    `json:"Posee_ap"`

	Fiscal_serie          int    `json:"Fiscal_serie"`
	Fiscal_destino        int    `json:"Fiscal_destino"`
	Rol_manzana           string `json:"Rol_manzana"`
	Rol_predio            string `json:"Rol_predio"`
	Fiscal_exento         int    `json:"Fiscal_exento"`
	Fiscal_avaluo         string `json:"Fiscal_avaluo"`
	Fiscal_contribucion   string `json:"Fiscal_contribucion"`
	Fiscal_sup_terreno    string `json:"Fiscal_sup_terreno"`
	Fiscal_sup_edificada  string `json:"Fiscal_sup_edificada"`
	Fiscal_sup_pavimentos string `json:"Fiscal_sup_pavimentos"`

	Valor_terreno               string `json:"Valor_terreno"`
	Valor_edificacion           string `json:"Valor_edificacion"`
	Valor_obras_complementarias string `json:"Valor_obras_complementarias"`
	Valor_total                 string `json:"Valor_total"`

	Cert_info_previas              int    `json:"Cert_info_previas"`
	Tipo_instrumento               int    `json:"Tipo_instrumento"`
	Especificar_tipo_instrumento   string `json:"Especificar_tipo_instrumento"`
	Indicar_area                   int    `json:"Indicar_area"`
	Zona_normativa_pregulador      string `json:"Zona_normativa_pregulador"`
	Area_riesgo                    int    `json:"Area_riesgo"`
	Area_proteccion                int    `json:"Area_proteccion"`
	Zona_conservacion_historica    int    `json:"Zona_conservacion_historica"`
	Zona_tipica                    int    `json:"Zona_tipica"`
	Monumento_nacional             int    `json:"Monumento_nacional"`
	Zona_uso_suelo                 string `json:"Zona_uso_suelo"`
	Usos_permitidos                string `json:"Usos_permitidos"`
	Usos_prohibidos                string `json:"Usos_prohibidos"`
	Superficie_predial_minima      int    `json:"Superficie_predial_minima"`
	Densidad_maxima_bruta          int    `json:"Densidad_maxima_bruta"`
	Densidad_maxima_neta           int    `json:"Densidad_maxima_neta"`
	Altura_maxima                  int    `json:"Altura_maxima"`
	Sistema_agrupamiento           int    `json:"Sistema_agrupamiento"`
	Coef_constructibilidad         int    `json:"Coef_constructibilidad"`
	Coef_ocupacionsuelo            int    `json:"Coef_ocupacionsuelo"`
	Coef_ocupacionsuelopsuperiores int    `json:"Coef_ocupacionsuelopsuperiores"`
	Rasante                        int    `json:"Rasante"`
	Adosamiento                    int    `json:"Adosamiento"`
	Distanciamiento                string `json:"Distanciamiento"`
	Cierres_perialtura             int    `json:"Cierres_perialtura"`
	Cierres_peritransparencia      int    `json:"Cierres_peritransparencia"`
	Ochavos                        int    `json:"Ochavos"`
	Ochavos_metros                 int    `json:"Ochavos_metros"`
	Estado_UrbaEjecutada           int    `json:"Estado_UrbaEjecutada"`
	Estado_UrbaRecibida            int    `json:"Estado_UrbaRecibida"`
	Estado_UrbaGarantizada         int    `json:"Estado_UrbaGarantizada"`

	Comuna string `json:"Comuna"`
	Ciudad string `json:"Ciudad"`
	Region string `json:"Region"`
	Pais   string `json:"Pais"`
}
type PermisoEdificacion struct {
	Id_rec                               int                                  `json:"Id_rec"`
	Nombre                               string                               `json:"Nombre"`
	Sup_Terreno                          string                               `json:"SupTerreno"`
	Posee_Permiso_Edificacion            int                                  `json:"Posee_Permiso_Edificacion"`
	Tipo_Permiso_Edificacion             int                                  `json:"Tipo_Permiso_Edificacion"`
	Especificar_Tipo_Permiso_Edificacion string                               `json:"Especificar_Tipo_Permiso_Edificacion"`
	Numero_Permiso                       string                               `json:"Numero_Permiso"`
	Fecha_Permiso                        string                               `json:"Fecha_Permiso"`
	Cant_Pisos_Sobre_Nivel               int                                  `json:"Cant_Pisos_Sobre_Nivel"`
	Cant_Pisos_Bajo_Nivel                int                                  `json:"Cant_Pisos_Bajo_Nivel"`
	Sup_Edificada_Sobre_Nivel            int                                  `json:"Sup_Edificada_Sobre_Nivel"`
	Sup_Edificada_Bajo_Nivel             int                                  `json:"Sup_Edificada_Bajo_Nivel"`
	Aco_Art_Esp_Transitorio              int                                  `json:"Aco_Art_Esp_Transitorio"`
	Recepcion_Definitiva                 int                                  `json:"Recepcion_Definitiva"`
	RecepcionTotalAcoge                  int                                  `json:"RecepcionTotalAcoge"`
	RecepcionParcialAcoge                int                                  `json:"RecepcionParcialAcoge"`
	ObraP_Faena                          int                                  `json:"ObraP_Faena"`
	ObraP_Grua                           int                                  `json:"ObraP_Grua"`
	ObraP_Excavacion                     int                                  `json:"ObraP_Excavacion"`
	Op0                                  int                                  `json:"Op0"`
	Op1                                  int                                  `json:"Op1"`
	Op2                                  int                                  `json:"Op2"`
	Op3                                  int                                  `json:"Op3"`
	Op4                                  int                                  `json:"Op4"`
	Op5                                  int                                  `json:"Op5"`
	Op6                                  int                                  `json:"Op6"`
	Op7                                  int                                  `json:"Op7"`
	Op8                                  int                                  `json:"Op8"`
	Op9                                  int                                  `json:"Op9"`
	Op10                                 int                                  `json:"Op10"`
	Op11                                 int                                  `json:"Op11"`
	Op12                                 int                                  `json:"Op12"`
	Archivos                             map[int][]ArchivosPermisoEdificacion `json:"Archivos"`
	UltimosArchivos                      map[int]ArchivosPermisoEdificacion   `json:"UltimosArchivos"`
}
type ArchivosPermisoEdificacion struct {
	Id_arc        int    `json:"Id_rec"`
	Nombre        string `json:"Nombre"`
	Nombre2       string `json:"Nombre2"`
	Tipo          int    `json:"Tipo"`
	Indicar_acoge int    `json:"Ano"`
	Fecha         string `json:"Fecha"`
	Fecha_insert  string `json:"Fecha_insert"`
}
type ListaArchivos struct {
	Id_arc          int    `json:"Id_rec"`
	Nombre          string `json:"Nombre"`
	Fojas           string `json:"Fojas"`
	Numero          string `json:"Numero"`
	Ano             string `json:"Ano"`
	Tipo            int    `json:"Tipo"`
	Fecha           string `json:"Fecha"`
	Fecha_insert    string `json:"Fecha_insert"`
	Valor_arriendo  string `json:"Valor_arriendo"`
	Renovacion_auto string `json:"Renovacion_auto"`
	Tipo_de_Plano   string `json:"Tipo_de_Plano"`
}
type buscarPropiedades struct {
	Titulo            string        `json:"Titulo"`
	SubTitulo         string        `json:"SubTitulo"`
	SubTitulo2        string        `json:"SubTitulo2"`
	SubTitulo3        string        `json:"SubTitulo3"`
	SubTitulo4        string        `json:"SubTitulo4"`
	PropiedadesString string        `json:"PropiedadesString"`
	Permisos          ListaPermisos `json:"Permisos"`
}
type Rec struct {
	Rec  bool   `json:"Rec"`
	Id   string `json:"Id"`
	Code string `json:"Code"`
}
type Resumen struct {
	Prods               map[int]ResumenProds   `json:"Prods"`
	Alertas             map[int]ResumenAlertas `json:"Alertas"`
	Notificaciones      map[int]ResumenAlertas `json:"Alertas"`
	TotalAlertas        int                    `json:"TotalAlertas"`
	TotalNotificaciones int                    `json:"TotalNotificaciones"`
	Localidades         []CamposPropiedades    `json:"Localidades"`
}
type ResumenProds struct {
	Nombre string          `json:"Nombre"`
	Lista  []ResumenAlerta `json:"Lista"`
}
type ResumenAlerta struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}
type ResumenAlertas struct {
	Nombre string        `json:"Nombre"`
	Lista  []ResumenProd `json:"Lista"`
}
type ResumenProd struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}
type Alerta struct {
	Id_ale       int      `json:"Id_ale"`
	Alerta       int      `json:"Alerta"`
	Notificacion int      `json:"Notificacion"`
	Campos       []string `json:"Campos"`
	Valores      []string `json:"Valores"`
}
type EmailData struct {
	Code string
}
type Alertas struct {
	Id_ale       int     `json:"Id_alr"`
	Nombre       string  `json:"Nombre"`
	Descripcion  string  `json:"Descripcion"`
	Alerta       int     `json:"Alerta"`
	Notificacion int     `json:"Notificacion"`
	Precio       float32 `json:"Precio"`
	Reglas       []Regla `json:"Reglas"`
	ReglasIndex  Regla   `json:"ReglasIndex"`
}
type Regla struct {
	Id_alr int    `json:"Id_alr"`
	Nombre string `json:"Nombre"`
	Pagina int    `json:"Pagina"`
	Tipo   int    `json:"Tipo"`
	Campo  string `json:"Campo"`
	Valor  string `json:"Valor"`
}
type ExcelImage struct {
	Nombre string `json:"Nombre"`
	Tipo   string `json:"Tipo"`
}

var (
	imgHandler fasthttp.RequestHandler
	cssHandler fasthttp.RequestHandler
	jsHandler  fasthttp.RequestHandler
	port       string
)
var pass = &MyHandler{Conf: Config{}}

func main() {

	if runtime.GOOS == "windows" {
		imgHandler = fasthttp.FSHandler("C:/Go/Pelao_No_Git/img", 1)
		cssHandler = fasthttp.FSHandler("C:/Go/Pelao_No_Git/css", 1)
		jsHandler = fasthttp.FSHandler("C:/Go/Pelao_No_Git/js", 1)
		port = ":81"
	} else {
		imgHandler = fasthttp.FSHandler("/var/Pelao/img", 1)
		cssHandler = fasthttp.FSHandler("/var/Pelao/css", 1)
		jsHandler = fasthttp.FSHandler("/var/Pelao/js", 1)
		port = ":80"
	}

	passwords, err := os.ReadFile("../password_redigo.json")
	if err == nil {
		if err := json.Unmarshal(passwords, &pass.Passwords); err == nil {
			fmt.Println(pass.Passwords)
		}
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
	go func() {
		r := router.New()
		r.GET("/", Index)
		r.GET("/css/{name}", Css)
		r.GET("/recuperar", Recuperar)
		r.GET("/recuperar/{name}/{id}", Recuperar2)
		r.GET("/js/{name}", Js)
		r.GET("/img/{name}", Img)
		r.GET("/images/{id}/{name}", Images)
		r.GET("/pages/{name}", Pages)
		r.POST("/login", Login)
		r.POST("/cart", Cart)
		r.POST("/nueva", Nueva)
		r.POST("/save", Save)
		r.POST("/delete", Delete)
		r.POST("/acciones", Acciones)
		r.GET("/salir", Salir)
		r.GET("/cotizacion/{name}", Cotizacionfunc)
		r.GET("/SetEmpresa/{name}", SetEmpresa)
		r.GET("/video/{name}", Video)

		// ANTES
		fasthttp.ListenAndServe(port, r.Handler)

		// DESPUES
		/*
			go func() {
				fasthttp.ListenAndServe(":80", redirectHTTP)
			}()
			server := &fasthttp.Server{Handler: r.Handler}
			server.ListenAndServeTLS(":443", "/etc/letsencrypt/live/www.redigo.cl/fullchain.pem", "/etc/letsencrypt/live/www.redigo.cl/privkey.pem")
		*/
	}()
	if err := run(con, pass, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func redirectHTTP(ctx *fasthttp.RequestCtx) {
	redirectURL := fmt.Sprintf("https://%v%v", string(ctx.Host()), string(ctx.URI().RequestURI()))
	ctx.Redirect(redirectURL, fasthttp.StatusMovedPermanently)
}
func Video(ctx *fasthttp.RequestCtx) {

	// Lee el archivo de video
	video, err := ioutil.ReadFile(fmt.Sprintf("/var/Pelao/video/%v", ctx.UserValue("name")))
	if err != nil {
		ctx.Error("Error al leer el archivo de video", fasthttp.StatusInternalServerError)
		return
	}

	// Establece las cabeceras para el video
	ctx.Response.Header.Set("Content-Type", "video/mp4")
	ctx.Response.Header.Set("Content-Length", fmt.Sprint(len(video)))

	// Escribe el video en el cuerpo de la respuesta
	ctx.Write(video)
}
func Save(ctx *fasthttp.RequestCtx) {

	resp := Response{}
	resp.Op = 2
	resp.Msg = "Error Inesperado"

	id := Read_uint32bytes(ctx.FormValue("id"))
	token := string(ctx.Request.Header.Cookie("cu"))

	db, err := GetMySQLDB()
	defer db.Close()
	if err != nil {
		resp.Msg = "Error Base de Datos"
	}

	switch string(ctx.FormValue("accion")) {
	case "guardar_empresa":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, _ := SuperAdmin(token); found {

			nombre := string(ctx.FormValue("nombre"))
			precio := string(ctx.FormValue("precio"))
			if id > 0 {
				resp.Op, resp.Msg = UpdateEmpresa(db, id, nombre, precio)
			}
			if id == 0 {
				resp.Op, resp.Msg = InsertEmpresa(db, nombre, precio)
			}
			if resp.Op == 1 {
				resp.Page = "crearEmpresa"
				resp.Reload = 1
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_alerta":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, _ := SuperAdmin(token); found {
			nombre := string(ctx.FormValue("nombre"))
			descripcion := string(ctx.FormValue("descripcion"))
			alerta := string(ctx.FormValue("tipo_alerta"))
			notificacion := string(ctx.FormValue("notificacion"))
			precio := string(ctx.FormValue("precio"))
			if id > 0 {
				resp.Op, resp.Msg = UpdateAlerta(db, id, nombre, descripcion, alerta, notificacion, precio)
			}
			if id == 0 {
				resp.Op, resp.Msg = InsertAlerta(db, nombre, descripcion, alerta, notificacion, precio)
			}
			if resp.Op == 1 {
				resp.Page = "crearAlerta"
				resp.Reload = 1
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_regla":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, _ := SuperAdmin(token); found {

			id_ale := Read_uint32bytes(ctx.FormValue("id_ale"))
			nombre := string(ctx.FormValue("nombre"))
			tipo := Read_uint32bytes(ctx.FormValue("tipo"))
			pagina := Read_uint32bytes(ctx.FormValue("pagina"))

			if tipo > 0 && tipo < 3 {
				if pagina > 0 && pagina < 9 {
					campo := string(ctx.FormValue(fmt.Sprintf("campo%v%v", pagina, tipo)))
					valor := string(ctx.FormValue(fmt.Sprintf("valor%v", tipo)))
					if id > 0 {
						resp.Op, resp.Msg = UpdateRegla(db, id, nombre, tipo, pagina, campo, valor, id_ale)
					}
					if id == 0 {
						resp.Op, resp.Msg = InsertRegla(db, nombre, tipo, pagina, campo, valor, id_ale)
					}
					if resp.Op == 1 {
						resp.Page = fmt.Sprintf("crearRegla?id_ale=%v", id_ale)
						resp.Reload = 1
					}
				} else {
					resp.Msg = "Error al ingresar Regla"
				}
			} else {
				resp.Msg = "Error al ingresar Regla"
			}

		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad1":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				nombre := string(ctx.FormValue("nombre"))
				lat := string(ctx.FormValue("lat"))
				lng := string(ctx.FormValue("lng"))
				rangopos := string(ctx.FormValue("rangopos"))
				comuna := string(ctx.FormValue("comuna"))
				ciudad := string(ctx.FormValue("ciudad"))
				region := string(ctx.FormValue("region"))
				pais := string(ctx.FormValue("pais"))
				direccion := string(ctx.FormValue("direccion"))
				numero := string(ctx.FormValue("numero"))
				dominio := string(ctx.FormValue("dominio"))
				dominio2 := string(ctx.FormValue("dominio2"))
				atencion_publico := string(ctx.FormValue("atencion_publico"))
				copropiedad := string(ctx.FormValue("copropiedad"))
				destino := string(ctx.FormValue("destino"))
				detalle_destino := string(ctx.FormValue("detalle_destino"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad(db, id_emp, id, nombre, lat, lng, rangopos, comuna, ciudad, region, pais, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino)
				}
				if id == 0 {
					resp.Op, resp.Msg, id = InsertPropiedad(db, id_emp, nombre, lat, lng, rangopos, comuna, ciudad, region, pais, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino)
				}
				pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
				files, err := ctx.MultipartForm()
				if err == nil {
					errup1, list1 := UploadFile(fmt.Sprintf("./files/images/%v", id_emp), files.File["foto_principal"], true, []string{"JPG"}, "")
					if errup1 {
						for _, x := range list1 {
							SavePhotoDb(db, x, 1, id, id_emp)
						}
					}
					errup2, list2 := UploadFile(fmt.Sprintf("./files/images/%v", id_emp), files.File["foto_interior"], true, []string{"JPG"}, "")
					if errup2 {
						for _, x := range list2 {
							SavePhotoDb(db, x, 2, id, id_emp)
						}
					}
					errup3, list3 := UploadFile(fmt.Sprintf("./files/images/%v", id_emp), files.File["foto_exterior"], true, []string{"JPG"}, "")
					if errup3 {
						for _, x := range list3 {
							SavePhotoDb(db, x, 3, id, id_emp)
						}
					}
				} else {
					resp.Msg = "Error Interno"
				}
				if resp.Op == 1 {
					resp.Page = fmt.Sprintf("crearPropiedad2PermisoEdificacion?id=%v", id)
					resp.Reload = 1
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad2A":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				id_rec := Read_uint32bytes(ctx.FormValue("id_rec"))

				sup_terreno := string(ctx.FormValue("sup_terreno"))
				posee_permiso_edificacion := string(ctx.FormValue("posee_permiso_edificacion"))
				tipo_permiso_edificacion := string(ctx.FormValue("tipo_permiso_edificacion"))
				especificar := string(ctx.FormValue("especificar"))
				numero_permiso_edificacion := string(ctx.FormValue("numero_permiso_edificacion"))
				fecha_permiso := string(ctx.FormValue("fecha_permiso"))
				cant_pisos_sobre_nivel := string(ctx.FormValue("cant_pisos_sobre_nivel"))
				cant_pisos_bajo_nivel := string(ctx.FormValue("cant_pisos_bajo_nivel"))
				superficie_edificada_sobre_nivel := string(ctx.FormValue("superficie_edificada_sobre_nivel"))
				superficie_edificada_bajo_nivel := string(ctx.FormValue("superficie_edificada_bajo_nivel"))
				aco_art_esp_transitorio := string(ctx.FormValue("aco_art_esp_transitorio"))
				recepcion_definitiva := string(ctx.FormValue("recepcion_definitiva"))

				obra_p0 := string(ctx.FormValue("obra_p0"))
				obra_p1 := string(ctx.FormValue("obra_p1"))
				obra_p2 := string(ctx.FormValue("obra_p2"))

				op0 := string(ctx.FormValue("op0"))
				op1 := string(ctx.FormValue("op1"))
				op2 := string(ctx.FormValue("op2"))
				op3 := string(ctx.FormValue("op3"))
				op4 := string(ctx.FormValue("op4"))
				op5 := string(ctx.FormValue("op5"))
				op6 := string(ctx.FormValue("op6"))
				op7 := string(ctx.FormValue("op7"))
				op8 := string(ctx.FormValue("op8"))
				op9 := string(ctx.FormValue("op9"))
				op10 := string(ctx.FormValue("op10"))
				op11 := string(ctx.FormValue("op11"))
				op12 := string(ctx.FormValue("op12"))

				if id_rec > 0 {
					resp.Op, resp.Msg = UpdatePropiedad2A(db, id_emp, id, id_rec, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar, numero_permiso_edificacion, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obra_p0, obra_p1, obra_p2, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12)
				}
				if id_rec == 0 {
					resp.Op, resp.Msg, id_rec = InsertPropiedad2A(db, id_emp, id, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar, numero_permiso_edificacion, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obra_p0, obra_p1, obra_p2, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12)
				}
				pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
				files, err := ctx.MultipartForm()
				if err == nil {
					errup1, list1 := UploadFile(fmt.Sprintf("./files/pdf2/%v/1", id_emp), files.File["permiso"], false, []string{"PDF"}, "")
					if errup1 {
						for _, x := range list1 {
							SaveFileDb2(db, x, "", 1, "", "", id_rec, id, id_emp)
						}
					}
					errup2, list2 := UploadFile(fmt.Sprintf("./files/pdf2/%v/2", id_emp), files.File["plano"], false, []string{"PDF"}, "")
					if errup2 {
						for _, x := range list2 {
							SaveFileDb2(db, x, "", 2, "", "", id_rec, id, id_emp)
						}
					}

					errup3, list3 := UploadFile(fmt.Sprintf("./files/pdf2/%v/3", id_emp), files.File["doc_recepcion_parcial"], false, []string{"PDF"}, "")
					errup4, list4 := UploadFile(fmt.Sprintf("./files/pdf2/%v/4", id_emp), files.File["doc2_recepcion_parcial"], false, []string{"PDF"}, "")

					if errup3 || errup4 {
						fecha_recepcion_parcial := string(ctx.FormValue("fecha_recepcion_parcial"))
						recepcion_parcial_acoge := string(ctx.FormValue("recepcion_parcial_acoge"))
						SaveFileDb2(db, list3[0], list4[0], 3, recepcion_parcial_acoge, fecha_recepcion_parcial, id_rec, id, id_emp)
					}

					errup5, list5 := UploadFile(fmt.Sprintf("./files/pdf2/%v/5", id_emp), files.File["doc_recepcion_total"], false, []string{"PDF"}, "")
					errup6, list6 := UploadFile(fmt.Sprintf("./files/pdf2/%v/6", id_emp), files.File["doc2_recepcion_total"], false, []string{"PDF"}, "")

					if errup5 || errup6 {
						fecha_recepcion_total := string(ctx.FormValue("fecha_recepcion_total"))
						recepcion_total_acoge := string(ctx.FormValue("recepcion_total_acoge"))
						SaveFileDb2(db, list5[0], list6[0], 4, recepcion_total_acoge, fecha_recepcion_total, id_rec, id, id_emp)
					}

					errup7, list7 := UploadFile(fmt.Sprintf("./files/pdf2/%v/7", id_emp), files.File["doc_demolicion"], false, []string{"PDF"}, "")
					errup8, list8 := UploadFile(fmt.Sprintf("./files/pdf2/%v/8", id_emp), files.File["doc2_demolicion"], false, []string{"PDF"}, "")

					if errup7 || errup8 {
						demolicion_fecha := string(ctx.FormValue("demolicion_fecha"))
						demolicion_estado := string(ctx.FormValue("demolicion_estado"))
						SaveFileDb2(db, list7[0], list8[0], 5, demolicion_estado, demolicion_fecha, id_rec, id, id_emp)
					}

					errup9, list9 := UploadFile(fmt.Sprintf("./files/pdf2/%v/9", id_emp), files.File["doc_demolicion"], false, []string{"PDF"}, "")
					errup10, list10 := UploadFile(fmt.Sprintf("./files/pdf2/%v/10", id_emp), files.File["doc2_demolicion"], false, []string{"PDF"}, "")

					if errup9 || errup10 {
						fusion_fecha := string(ctx.FormValue("fusion_fecha"))
						fusion_estado := string(ctx.FormValue("fusion_estado"))
						SaveFileDb2(db, list9[0], list10[0], 6, fusion_estado, fusion_fecha, id_rec, id, id_emp)
					}

					errup11, list11 := UploadFile(fmt.Sprintf("./files/pdf2/%v/11", id_emp), files.File["doc_obrap_faena"], false, []string{"PDF"}, "")
					errup12, list12 := UploadFile(fmt.Sprintf("./files/pdf2/%v/12", id_emp), files.File["doc2_obrap_faena"], false, []string{"PDF"}, "")

					if errup11 || errup12 {
						fusion_fecha := string(ctx.FormValue("obrap_faena_fecha"))
						SaveFileDb2(db, list11[0], list12[0], 7, "", fusion_fecha, id_rec, id, id_emp)
					}

					errup13, list13 := UploadFile(fmt.Sprintf("./files/pdf2/%v/13", id_emp), files.File["doc_obrap_grua"], false, []string{"PDF"}, "")
					errup14, list14 := UploadFile(fmt.Sprintf("./files/pdf2/%v/14", id_emp), files.File["doc2_obrap_grua"], false, []string{"PDF"}, "")

					if errup13 || errup14 {
						grua_fecha := string(ctx.FormValue("obrap_grua_fecha"))
						SaveFileDb2(db, list13[0], list14[0], 8, "", grua_fecha, id_rec, id, id_emp)
					}

					errup15, list15 := UploadFile(fmt.Sprintf("./files/pdf2/%v/15", id_emp), files.File["doc_obrap_excanacion"], false, []string{"PDF"}, "")
					errup16, list16 := UploadFile(fmt.Sprintf("./files/pdf2/%v/16", id_emp), files.File["doc2_obrap_excanacion"], false, []string{"PDF"}, "")

					if errup15 || errup16 {
						excanacion_fecha := string(ctx.FormValue("obrap_excanacion_fecha"))
						SaveFileDb2(db, list15[0], list16[0], 9, "", excanacion_fecha, id_rec, id, id_emp)
					}

					errup17, list17 := UploadFile(fmt.Sprintf("./files/pdf2/%v/17", id_emp), files.File["opf1"], false, []string{"PDF"}, "")
					if errup17 {
						SaveFileDb2(db, list17[0], "", 10, "", "", id_rec, id, id_emp)
					}
					errup18, list18 := UploadFile(fmt.Sprintf("./files/pdf2/%v/18", id_emp), files.File["opf2"], false, []string{"PDF"}, "")
					if errup18 {
						SaveFileDb2(db, list18[0], "", 11, "", "", id_rec, id, id_emp)
					}
					errup19, list19 := UploadFile(fmt.Sprintf("./files/pdf2/%v/19", id_emp), files.File["opf3"], false, []string{"PDF"}, "")
					if errup19 {
						SaveFileDb2(db, list19[0], "", 12, "", "", id_rec, id, id_emp)
					}
					errup20, list20 := UploadFile(fmt.Sprintf("./files/pdf2/%v/20", id_emp), files.File["opf4"], false, []string{"PDF"}, "")
					if errup20 {
						SaveFileDb2(db, list20[0], "", 13, "", "", id_rec, id, id_emp)
					}
					errup21, list21 := UploadFile(fmt.Sprintf("./files/pdf2/%v/21", id_emp), files.File["opf5"], false, []string{"PDF"}, "")
					if errup21 {
						SaveFileDb2(db, list21[0], "", 14, "", "", id_rec, id, id_emp)
					}
					errup22, list22 := UploadFile(fmt.Sprintf("./files/pdf2/%v/22", id_emp), files.File["opf6"], false, []string{"PDF"}, "")
					if errup22 {
						SaveFileDb2(db, list22[0], "", 15, "", "", id_rec, id, id_emp)
					}
					errup23, list23 := UploadFile(fmt.Sprintf("./files/pdf2/%v/23", id_emp), files.File["opf7"], false, []string{"PDF"}, "")
					if errup23 {
						SaveFileDb2(db, list23[0], "", 16, "", "", id_rec, id, id_emp)
					}
					errup24, list24 := UploadFile(fmt.Sprintf("./files/pdf2/%v/24", id_emp), files.File["opf8"], false, []string{"PDF"}, "")
					if errup24 {
						SaveFileDb2(db, list24[0], "", 17, "", "", id_rec, id, id_emp)
					}
					errup25, list25 := UploadFile(fmt.Sprintf("./files/pdf2/%v/25", id_emp), files.File["opf9"], false, []string{"PDF"}, "")
					if errup25 {
						SaveFileDb2(db, list25[0], "", 18, "", "", id_rec, id, id_emp)
					}
					errup26, list26 := UploadFile(fmt.Sprintf("./files/pdf2/%v/26", id_emp), files.File["opf10"], false, []string{"PDF"}, "")
					if errup26 {
						SaveFileDb2(db, list26[0], "", 19, "", "", id_rec, id, id_emp)
					}
					errup27, list27 := UploadFile(fmt.Sprintf("./files/pdf2/%v/27", id_emp), files.File["opf11"], false, []string{"PDF"}, "")
					if errup27 {
						SaveFileDb2(db, list27[0], "", 20, "", "", id_rec, id, id_emp)
					}
					errup28, list28 := UploadFile(fmt.Sprintf("./files/pdf2/%v/28", id_emp), files.File["opf12"], false, []string{"PDF"}, "")
					if errup28 {
						SaveFileDb2(db, list28[0], "", 21, "", "", id_rec, id, id_emp)
					}

				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		if resp.Op == 1 {
			siguiente_accion := string(ctx.FormValue("siguiente_accion"))
			if siguiente_accion == "0" {
				resp.Page = fmt.Sprintf("crearPropiedad3?id=%v", id)
			} else {
				resp.Page = fmt.Sprintf("crearPropiedad2PermisoEdificacion?id=%v", id)
			}
			resp.Reload = 1
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad3":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				electrico_te1 := string(ctx.FormValue("electrico_te1"))
				dotacion_ap := string(ctx.FormValue("dotacion_ap"))
				dotacion_alcance := string(ctx.FormValue("dotacion_alcance"))
				instalacion_ascensor := string(ctx.FormValue("instalacion_ascensor"))
				te1_ascensor := string(ctx.FormValue("te1_ascensor"))
				certificado_ascensor := string(ctx.FormValue("certificado_ascensor"))
				clima := string(ctx.FormValue("clima"))
				seguridad_incendio := string(ctx.FormValue("seguridad_incendio"))

				if id > 0 {

					var isarrendado uint8 = 0
					resp.Op, resp.Msg, isarrendado = UpdatePropiedad3(db, id_emp, id, electrico_te1, dotacion_ap, dotacion_alcance, instalacion_ascensor, te1_ascensor, certificado_ascensor, clima, seguridad_incendio)
					resp.Reload = 1
					if isarrendado == 1 {
						resp.Page = fmt.Sprintf("crearPropiedad4?id=%v", id)
					} else {
						resp.Page = fmt.Sprintf("crearPropiedad5?id=%v", id)
					}
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
					files, err := ctx.MultipartForm()
					if err == nil {
						errup1, list1 := UploadFile(fmt.Sprintf("./files/pdf/%v/1", id_emp), files.File["doc_electrico_te1"], false, []string{"PDF"}, fmt.Sprintf("t1_%v.pdf", string(ctx.FormValue("fecha_electrico_te1"))))
						fecha_electrico_te1 := string(ctx.FormValue("fecha_electrico_te1"))
						if errup1 {
							for _, x := range list1 {
								SaveFileDb(db, x, "", "", fecha_electrico_te1, "", 1, "", "", "", id, id_emp)
							}
						} else {
							id_doc_01 := Read_uint32bytes(ctx.FormValue("id_doc_01"))
							if id_doc_01 > 0 {
								UpdateFileDb(db, "", "", fecha_electrico_te1, "", "", "", "", id_doc_01, id, id_emp)
							}
						}
						errup2, list2 := UploadFile(fmt.Sprintf("./files/pdf/%v/2", id_emp), files.File["doc_dotacion_ap"], false, []string{"PDF"}, "")
						fecha_dotacion_ap := string(ctx.FormValue("fecha_dotacion_ap"))
						if errup2 {
							for _, x := range list2 {
								SaveFileDb(db, x, "", "", fecha_dotacion_ap, "", 2, "", "", "", id, id_emp)
							}
						} else {
							id_doc_02 := Read_uint32bytes(ctx.FormValue("id_doc_02"))
							if id_doc_02 > 0 {
								UpdateFileDb(db, "", "", fecha_dotacion_ap, "", "", "", "", id_doc_02, id, id_emp)
							}
						}
						errup3, list3 := UploadFile(fmt.Sprintf("./files/pdf/%v/3", id_emp), files.File["doc_dotacion_alcance"], false, []string{"PDF"}, "")
						fecha_dotacion_alcance := string(ctx.FormValue("fecha_dotacion_alcance"))
						if errup3 {
							for _, x := range list3 {
								SaveFileDb(db, x, "", "", fecha_dotacion_alcance, "", 3, "", "", "", id, id_emp)
							}
						} else {
							id_doc_03 := Read_uint32bytes(ctx.FormValue("id_doc_03"))
							if id_doc_03 > 0 {
								UpdateFileDb(db, "", "", fecha_dotacion_alcance, "", "", "", "", id_doc_03, id, id_emp)
							}
						}
						errup4, list4 := UploadFile(fmt.Sprintf("./files/pdf/%v/4", id_emp), files.File["doc_instalacion_ascensor"], false, []string{"PDF"}, "")
						fecha_instalacion_ascensor := string(ctx.FormValue("fecha_instalacion_ascensor"))
						if errup4 {
							for _, x := range list4 {
								SaveFileDb(db, x, "", "", fecha_instalacion_ascensor, "", 4, "", "", "", id, id_emp)
							}
						} else {
							id_doc_04 := Read_uint32bytes(ctx.FormValue("id_doc_04"))
							if id_doc_04 > 0 {
								UpdateFileDb(db, "", "", fecha_instalacion_ascensor, "", "", "", "", id_doc_04, id, id_emp)
							}
						}
						errup5, list5 := UploadFile(fmt.Sprintf("./files/pdf/%v/5", id_emp), files.File["doc_te1_ascensor"], false, []string{"PDF"}, "")
						fecha_te1_ascensor := string(ctx.FormValue("fecha_te1_ascensor"))
						if errup5 {
							for _, x := range list5 {
								SaveFileDb(db, x, "", "", fecha_te1_ascensor, "", 5, "", "", "", id, id_emp)
							}
						} else {
							id_doc_05 := Read_uint32bytes(ctx.FormValue("id_doc_05"))
							if id_doc_05 > 0 {
								UpdateFileDb(db, "", "", fecha_te1_ascensor, "", "", "", "", id_doc_05, id, id_emp)
							}
						}
						errup6, list6 := UploadFile(fmt.Sprintf("./files/pdf/%v/6", id_emp), files.File["doc_certificado_ascensor"], false, []string{"PDF"}, "")
						fecha_certificado_ascensor := string(ctx.FormValue("fecha_certificado_ascensor"))
						if errup6 {
							for _, x := range list6 {
								SaveFileDb(db, x, "", "", fecha_certificado_ascensor, "", 6, "", "", "", id, id_emp)
							}
						} else {
							id_doc_06 := Read_uint32bytes(ctx.FormValue("id_doc_06"))
							if id_doc_06 > 0 {
								UpdateFileDb(db, "", "", fecha_certificado_ascensor, "", "", "", "", id_doc_06, id, id_emp)
							}
						}
						errup7, list7 := UploadFile(fmt.Sprintf("./files/pdf/%v/7", id_emp), files.File["doc_clima"], false, []string{"PDF"}, "")
						fecha_clima := string(ctx.FormValue("fecha_clima"))
						if errup7 {
							for _, x := range list7 {
								SaveFileDb(db, x, "", "", fecha_clima, "", 7, "", "", "", id, id_emp)
							}
						} else {
							id_doc_07 := Read_uint32bytes(ctx.FormValue("id_doc_07"))
							if id_doc_07 > 0 {
								UpdateFileDb(db, "", "", fecha_clima, "", "", "", "", id_doc_07, id, id_emp)
							}
						}
						errup8, list8 := UploadFile(fmt.Sprintf("./files/pdf/%v/8", id_emp), files.File["doc_seguridad_incendio"], false, []string{"PDF"}, "")
						fecha_seguridad_incendio := string(ctx.FormValue("fecha_seguridad_incendio"))
						if errup8 {
							for _, x := range list8 {
								SaveFileDb(db, x, "", "", fecha_seguridad_incendio, "", 8, "", "", "", id, id_emp)
							}
						} else {
							id_doc_08 := Read_uint32bytes(ctx.FormValue("id_doc_08"))
							if id_doc_08 > 0 {
								UpdateFileDb(db, "", "", fecha_seguridad_incendio, "", "", "", "", id_doc_08, id, id_emp)
							}
						}
					}
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad4":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				tasacion_valor_comercial := string(ctx.FormValue("tasacion_valor_comercial"))
				ano_tasacion := string(ctx.FormValue("ano_tasacion"))
				contrato_arriendo := string(ctx.FormValue("contrato_arriendo"))
				contrato_subarriendo := string(ctx.FormValue("contrato_subarriendo"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad4(db, id_emp, id, tasacion_valor_comercial, ano_tasacion, contrato_arriendo, contrato_subarriendo)
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
					files, err := ctx.MultipartForm()
					if err == nil {
						errup1, list1 := UploadFile(fmt.Sprintf("./files/pdf/%v/9", id_emp), files.File["doc_arriendo"], false, []string{"PDF"}, "")
						valor_arriendo := string(ctx.FormValue("valor_arriendo"))
						vencimiento_arriendo := string(ctx.FormValue("vencimiento_arriendo"))
						renovacion_automatica := string(ctx.FormValue("renovacion_automatica"))
						if errup1 {
							for _, x := range list1 {
								SaveFileDb(db, x, "", "", "", vencimiento_arriendo, 9, valor_arriendo, renovacion_automatica, "", id, id_emp)
							}
						} else {
							id_doc_09 := Read_uint32bytes(ctx.FormValue("id_doc_09"))
							if id_doc_09 > 0 {
								UpdateFileDb(db, "", "", vencimiento_arriendo, "", valor_arriendo, renovacion_automatica, "", id_doc_09, id, id_emp)
							}
						}
						errup2, list2 := UploadFile(fmt.Sprintf("./files/pdf/%v/10", id_emp), files.File["valor_subarriendo"], false, []string{"PDF"}, "")
						valor_subarriendo := string(ctx.FormValue("valor_subarriendo"))
						vencimiento_subarriendo := string(ctx.FormValue("vencimiento_subarriendo"))
						renovacion_automaticasub := string(ctx.FormValue("renovacion_automaticasub"))
						if errup2 {
							for _, x := range list2 {
								SaveFileDb(db, x, "", "", "", vencimiento_subarriendo, 10, valor_subarriendo, renovacion_automaticasub, "", id, id_emp)
							}
						} else {
							id_doc_10 := Read_uint32bytes(ctx.FormValue("id_doc_10"))
							if id_doc_10 > 0 {
								UpdateFileDb(db, "", "", vencimiento_subarriendo, "", valor_subarriendo, renovacion_automaticasub, "", id_doc_10, id, id_emp)
							}
						}
					}
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		if resp.Op == 1 {
			resp.Page = fmt.Sprintf("crearPropiedad5?id=%v", id)
			resp.Reload = 1
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad5":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				nom_propietario_conservador := string(ctx.FormValue("nom_propietario_conservador"))
				posee_gp := string(ctx.FormValue("posee_gp"))
				posee_ap := string(ctx.FormValue("posee_ap"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad5(db, id_emp, id, nom_propietario_conservador, posee_gp, posee_ap)
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
					files, err := ctx.MultipartForm()
					if err == nil {
						errup1, list1 := UploadFile(fmt.Sprintf("./files/pdf/%v/11", id_emp), files.File["doc_dom_nom_propietario"], false, []string{"PDF"}, "")
						foja_dom_vigencia := string(ctx.FormValue("foja_dom_vigencia"))
						numero_dom_vigencia := string(ctx.FormValue("numero_dom_vigencia"))
						ano_dom_vigencia := string(ctx.FormValue("ano_dom_vigencia"))
						if errup1 {
							for _, x := range list1 {
								SaveFileDb(db, x, foja_dom_vigencia, numero_dom_vigencia, "", ano_dom_vigencia, 11, "", "", "", id, id_emp)
							}
						} else {
							id_doc_11 := Read_uint32bytes(ctx.FormValue("id_doc_11"))
							if id_doc_11 > 0 {
								UpdateFileDb(db, foja_dom_vigencia, numero_dom_vigencia, "", ano_dom_vigencia, "", "", "", id_doc_11, id, id_emp)
							}
						}

						errup2, list2 := UploadFile(fmt.Sprintf("./files/pdf/%v/12", id_emp), files.File["doc_gp"], false, []string{"PDF"}, "")
						foja_gp := string(ctx.FormValue("foja_gp"))
						numero_gp := string(ctx.FormValue("numero_gp"))
						ano_gp := string(ctx.FormValue("ano_gp"))
						if errup2 {
							for _, x := range list2 {
								SaveFileDb(db, x, foja_gp, numero_gp, "", ano_gp, 12, "", "", "", id, id_emp)
							}
						} else {
							id_doc_12 := Read_uint32bytes(ctx.FormValue("id_doc_12"))
							if id_doc_12 > 0 {
								UpdateFileDb(db, foja_gp, numero_gp, "", ano_gp, "", "", "", id_doc_12, id, id_emp)
							}
						}

						errup3, list3 := UploadFile(fmt.Sprintf("./files/pdf/%v/13", id_emp), files.File["doc_planos_archivos"], false, []string{"PDF"}, "")
						tipo_de_plano := string(ctx.FormValue("tipo_de_plano"))
						numero_planos_archivos := string(ctx.FormValue("numero_planos_archivos"))
						if errup3 {
							for _, x := range list3 {
								SaveFileDb(db, x, "", numero_planos_archivos, "", "", 13, "", "", tipo_de_plano, id, id_emp)
							}
						} else {
							id_doc_13 := Read_uint32bytes(ctx.FormValue("id_doc_13"))
							if id_doc_13 > 0 {
								UpdateFileDb(db, "", numero_planos_archivos, "", "", "", "", tipo_de_plano, id_doc_13, id, id_emp)
							}
						}

						errup4, list4 := UploadFile(fmt.Sprintf("./files/pdf/%v/14", id_emp), files.File["reglamento_copropiedad"], false, []string{"PDF"}, "")
						foja_reglamento_copropiedad := string(ctx.FormValue("foja_reglamento_copropiedad"))
						numero_reglamento_copropiedad := string(ctx.FormValue("numero_reglamento_copropiedad"))
						ano_reglamento_copropiedad := string(ctx.FormValue("ano_reglamento_copropiedad"))
						if errup4 {
							for _, x := range list4 {
								SaveFileDb(db, x, foja_reglamento_copropiedad, numero_reglamento_copropiedad, "", ano_reglamento_copropiedad, 14, "", "", "", id, id_emp)
							}
						} else {
							id_doc_14 := Read_uint32bytes(ctx.FormValue("id_doc_14"))
							if id_doc_14 > 0 {
								UpdateFileDb(db, foja_reglamento_copropiedad, numero_reglamento_copropiedad, "", ano_reglamento_copropiedad, "", "", "", id_doc_14, id, id_emp)
							}
						}
					}
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		if resp.Op == 1 {
			resp.Page = fmt.Sprintf("crearPropiedad6?id=%v", id)
			resp.Reload = 1
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad6":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				fiscal_serie := string(ctx.FormValue("fiscal_serie"))
				fiscal_destino := string(ctx.FormValue("fiscal_destino"))
				rol_manzana := string(ctx.FormValue("rol_manzana"))
				rol_predio := string(ctx.FormValue("rol_predio"))
				fiscal_exento := string(ctx.FormValue("fiscal_exento"))
				avaluo_fiscal := string(ctx.FormValue("avaluo_fiscal"))
				contribucion_fiscal := string(ctx.FormValue("contribucion_fiscal"))
				superficie_terreno := string(ctx.FormValue("superficie_terreno"))
				superficie_edificada := string(ctx.FormValue("superficie_edificada"))
				superficie_pavimentos := string(ctx.FormValue("superficie_pavimentos"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad6(db, id_emp, id, fiscal_serie, fiscal_destino, rol_manzana, rol_predio, fiscal_exento, avaluo_fiscal, contribucion_fiscal, superficie_terreno, superficie_edificada, superficie_pavimentos)
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		if resp.Op == 1 {
			resp.Page = fmt.Sprintf("crearPropiedad7?id=%v", id)
			resp.Reload = 1
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad7":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				valor_terreno := string(ctx.FormValue("valor_terreno"))
				valor_edificacion := string(ctx.FormValue("valor_edificacion"))
				valor_obras_complementarias := string(ctx.FormValue("valor_obras_complementarias"))
				valor_total := string(ctx.FormValue("valor_total"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad7(db, id_emp, id, valor_terreno, valor_edificacion, valor_obras_complementarias, valor_total)
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
					files, err := ctx.MultipartForm()
					fecha_informe_tasacion := string(ctx.FormValue("fecha_informe_tasacion"))
					if err == nil {
						errup1, list1 := UploadFile(fmt.Sprintf("./files/images/%v/15", id_emp), files.File["doc_informe_tasacion"], false, []string{"PDF"}, "")
						if errup1 {
							for _, x := range list1 {
								SaveFileDb(db, x, "", "", fecha_informe_tasacion, "", 15, "", "", "", id, id_emp)
							}
						}
					} else {
						id_doc_15 := Read_uint32bytes(ctx.FormValue("id_doc_15"))
						if id_doc_15 > 0 {
							UpdateFileDb(db, "", "", fecha_informe_tasacion, "", "", "", "", id_doc_15, id, id_emp)
						}
					}
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		if resp.Op == 1 {
			resp.Page = fmt.Sprintf("crearPropiedad8?id=%v", id)
			resp.Reload = 1
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_propiedad8":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {

				cert_info_previas := string(ctx.FormValue("cert_info_previas"))
				tipo_instrumento := string(ctx.FormValue("tipo_instrumento"))
				especificar_tipo_instrumento := string(ctx.FormValue("especificar_tipo_instrumento"))

				indicar_area := string(ctx.FormValue("indicar_area"))
				zona_normativa_plan_regulador := string(ctx.FormValue("zona_normativa_plan_regulador"))
				area_riesgo := string(ctx.FormValue("area_riesgo"))
				area_proteccion := string(ctx.FormValue("area_proteccion"))
				zona_conservacion_historica := string(ctx.FormValue("zona_conservacion_historica"))
				zona_tipica := string(ctx.FormValue("zona_tipica"))
				monumento_nacional := string(ctx.FormValue("monumento_nacional"))
				zona_uso_suelo := string(ctx.FormValue("zona_uso_suelo"))
				usos_permitidos := string(ctx.FormValue("usos_permitidos"))
				usos_prohibidos := string(ctx.FormValue("usos_prohibidos"))
				superficie_predial_minima := string(ctx.FormValue("superficie_predial_minima"))
				densidad_maxima_bruta := string(ctx.FormValue("densidad_maxima_bruta"))
				densidad_maxima_neta := string(ctx.FormValue("densidad_maxima_neta"))
				altura_maxima := string(ctx.FormValue("altura_maxima"))
				sistema_agrupamiento := string(ctx.FormValue("sistema_agrupamiento"))
				coef_constructibilidad := string(ctx.FormValue("coef_constructibilidad"))
				coef_ocupacion_suelo := string(ctx.FormValue("coef_ocupacion_suelo"))
				coef_ocupacion_suelo_psuperiores := string(ctx.FormValue("coef_ocupacion_suelo_psuperiores"))
				rasante := string(ctx.FormValue("rasante"))
				adosamiento := string(ctx.FormValue("adosamiento"))
				distanciamiento := string(ctx.FormValue("distanciamiento"))
				cierres_perimetrales_altura := string(ctx.FormValue("cierres_perimetrales_altura"))
				cierres_perimetrales_transparencia := string(ctx.FormValue("cierres_perimetrales_transparencia"))
				ochavos := string(ctx.FormValue("ochavos"))
				ochavos_metros := string(ctx.FormValue("ochavos_metros"))
				estado_urbanizacion_ejecutada := string(ctx.FormValue("estado_urbanizacion_ejecutada"))
				estado_urbanizacion_recibida := string(ctx.FormValue("estado_urbanizacion_recibida"))
				estado_urbanizacion_garantizada := string(ctx.FormValue("estado_urbanizacion_garantizada"))

				if id > 0 {
					resp.Op, resp.Msg = UpdatePropiedad8(db, id_emp, id, cert_info_previas, tipo_instrumento, especificar_tipo_instrumento, indicar_area, zona_normativa_plan_regulador, area_riesgo, area_proteccion, zona_conservacion_historica, zona_tipica, monumento_nacional, zona_uso_suelo, usos_permitidos, usos_prohibidos, superficie_predial_minima, densidad_maxima_bruta, densidad_maxima_neta, altura_maxima, sistema_agrupamiento, coef_constructibilidad, coef_ocupacion_suelo, coef_ocupacion_suelo_psuperiores, rasante, adosamiento, distanciamiento, cierres_perimetrales_altura, cierres_perimetrales_transparencia, ochavos, ochavos_metros, estado_urbanizacion_ejecutada, estado_urbanizacion_recibida, estado_urbanizacion_garantizada)
					pass.ListaPropAlerts = append(pass.ListaPropAlerts, ListaPropAlert{Id_pro: id, Id_emp: id_emp})
					files, err := ctx.MultipartForm()
					if err == nil {
						errup1, list1 := UploadFile(fmt.Sprintf("./files/images/%v/16", id_emp), files.File["doc_cert_info_previas"], false, []string{"PDF"}, "")
						if errup1 {
							for _, x := range list1 {
								SaveFileDb(db, x, "", "", "", "", 16, "", "", "", id, id_emp)
							}
						}
					}
				}
				if resp.Op == 1 {
					resp.Page = "crearPropiedad"
					resp.Reload = 1
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_usuarios":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				nombre := string(ctx.FormValue("nombre"))
				pass := string(ctx.FormValue("pass"))
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

				if id > 0 {
					resp.Op, resp.Msg = UpdateUsuario(db, id_emp, id, nombre, pass, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9)
				}
				if id == 0 {
					resp.Op, resp.Msg = InsertUsuario(db, id_emp, nombre, pass, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9)
				}
				if resp.Op == 1 {
					resp.Page = "crearUsuarios"
					resp.Reload = 1
				}
			} else {
				resp.Msg = "No tiene permisos"
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_detalle_cotizacion":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, _ := SuperAdmin(token); found {

			descripcion := string(ctx.FormValue("descripcion"))
			precio := string(ctx.FormValue("precio"))
			insmod := string(ctx.FormValue("insmod"))
			id_pro := Read_uint32bytes(ctx.FormValue("id_pro"))
			id_ale := Read_uint32bytes(ctx.FormValue("id_ale"))
			if insmod == "0" {
				resp.Op, resp.Msg = InsertDetalleCotizacion(db, id, descripcion, precio, id_pro, id_ale)
			} else {
				resp.Op, resp.Msg = UpdateDetalleCotizacion(db, id, descripcion, precio, id_pro, id_ale)
			}
			if resp.Op == 1 {
				resp.Page = fmt.Sprintf("confCotizacion?id_cot=%v", id)
				resp.Reload = 1
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "guardar_admin_cotizacion":

		ctx.Response.Header.Set("Content-Type", "application/json")
		if found, _ := SuperAdmin(token); found {

			id_emp := Read_uint32bytes(ctx.FormValue("id_emp"))
			uf := string(ctx.FormValue("precio"))

			if id > 0 {
				resp.Op, resp.Msg = UpdateCotizacion(db, id, uf, id_emp)
			}
			if id == 0 {
				resp.Op, resp.Msg = InsertCotizacion(db, uf, id_emp)
			}
			if resp.Op == 1 {
				resp.Page = "AdminCotizacion"
				resp.Reload = 1
			}
		} else {
			resp.Msg = "No tiene permisos"
		}
		json.NewEncoder(ctx).Encode(resp)
	case "descargar_custom_pdf":

		ctx.SetContentType("application/octet-stream")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				aux, found := GetPropiedad(id_emp, id, true)
				if found {

					p1a := Read_uint32bytes(ctx.FormValue("p1a"))
					p1b := Read_uint32bytes(ctx.FormValue("p1b"))
					p2 := Read_uint32bytes(ctx.FormValue("p2"))
					p3 := Read_uint32bytes(ctx.FormValue("p3"))
					p4 := Read_uint32bytes(ctx.FormValue("p4"))
					p5 := Read_uint32bytes(ctx.FormValue("p5"))
					p6 := Read_uint32bytes(ctx.FormValue("p6"))
					p7 := Read_uint32bytes(ctx.FormValue("p7"))
					p8 := Read_uint32bytes(ctx.FormValue("p8"))

					titulo := fmt.Sprintf("Resumen Propiedad %v", aux.Nombre)

					urlqr := fmt.Sprintf("https://www.redigo.cl/?p=detalle_propiedad&i=%v", id)
					b, qrfile := CreateQr(urlqr, id, "detalle")

					if b {

						addPage := false

						m := pdf.NewMaroto(consts.Portrait, consts.A4)
						m.SetPageMargins(10, 10, 10)

						blackColor := col.Color{Red: 0, Green: 0, Blue: 0}
						darkGrayColor := col.Color{Red: 55, Green: 55, Blue: 55}
						darkGrayColor2 := col.Color{Red: 220, Green: 220, Blue: 220}
						whiteColor := col.Color{Red: 255, Green: 255, Blue: 255}

						m.RegisterHeader(func() {
							m.Row(20, func() {
								m.Col(3, func() {
									_ = m.FileImage("./logo.png", props.Rect{
										Center:  true,
										Percent: 100,
									})
								})
								m.Col(6, func() {
									m.Text(titulo, props.Text{
										Top:   6,
										Size:  18,
										Style: consts.Bold,
										Align: consts.Center,
									})
								})
								m.Col(3, func() {
									_ = m.FileImage(qrfile, props.Rect{
										Center:  true,
										Percent: 100,
									})
								})
							})
						})
						m.RegisterFooter(func() {
							m.Row(4, func() {
								m.Col(12, func() {
									m.Text("www.redigo.cl", props.Text{
										Top:   16,
										Style: consts.BoldItalic,
										Size:  10,
										Align: consts.Center,
										Color: darkGrayColor,
									})
								})
							})
						})

						// PAGINA 1
						if p1a == 1 {
							addPage = true
							m.Row(8, func() {})
							m.Row(7, func() {
								m.Col(12, func() {
									m.Text("Datos Generales", props.Text{
										Top:   1.5,
										Size:  18,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.Row(8, func() {})

							// PASO 1 DIRECION-NUMERO-COMUNA
							m.SetBackgroundColor(darkGrayColor)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text("Direccion", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Numero", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Comuna", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(darkGrayColor2)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text(aux.Direccion, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(aux.Numero, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(aux.Comuna, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(whiteColor)
							m.Row(1, func() {})

							// PASO 1 CIUDAD-REGION-PAIS
							m.SetBackgroundColor(darkGrayColor)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text("Ciudad", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Region", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Pais", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(darkGrayColor2)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text(aux.Ciudad, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(aux.Region, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(aux.Pais, props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(whiteColor)
							m.Row(1, func() {})

							// PASO 1 DOMINIO-DOMINIO2-ATENCIONPUBLICO
							m.SetBackgroundColor(darkGrayColor)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text("Dominio", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Dominio2", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Atención a publico", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(darkGrayColor2)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text(PdfStr(0, aux.Dominio, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(PdfStr(1, aux.Dominio2, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(PdfStr(2, aux.Atencion_publico, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(whiteColor)
							m.Row(1, func() {})

							// PASO 1 COPROPIEDAD-USOODESTINO-DETALLEUSOODESTINO
							m.SetBackgroundColor(darkGrayColor)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text("Copropiedad", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Uso o destino", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text("Detalle uso o destino", props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewWhite(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(darkGrayColor2)
							m.Row(7, func() {
								m.Col(4, func() {
									m.Text(PdfStr(3, aux.Copropiedad, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(PdfStr(4, aux.Destino, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
								m.Col(4, func() {
									m.Text(PdfStr(5, aux.Detalle_destino, ""), props.Text{
										Top:   1.5,
										Size:  9,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.SetBackgroundColor(whiteColor)
							m.Row(1, func() {})

							_, nommap := CreateMapImage("-33.4397852", "-70.6169508", 1200, 300)

							m.Row(90, func() {
								m.Col(12, func() {
									_ = m.FileImage(nommap, props.Rect{
										Center:  true,
										Percent: 100,
									})
								})
							})
							m.SetBackgroundColor(whiteColor)
							m.Row(1, func() {})
						}
						// PAGINA 1 FOTOS
						if p1b == 1 {
							addPage = true
							limages := ExcelImages(aux.Imagenes)
							for _, arr := range limages {
								m.Row(35, func() {
									for _, obj := range arr {
										m.Col(2, func() {
											_ = m.FileImage(fmt.Sprintf("files/images/%v/%v", id_emp, obj.Nombre), props.Rect{
												Center:  true,
												Percent: 100,
											})
										})
									}
								})
								m.Row(10, func() {
									for _, obj := range arr {
										m.Col(2, func() {
											m.Text(obj.Tipo, props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Center,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									}
								})
							}
						}

						if addPage {
							m.AddPage()
							addPage = false
						}

						if p2 == 1 {
							addPage = true
							m.Row(7, func() {
								m.Col(12, func() {
									m.Text("Situación Municipal", props.Text{
										Top:   1.5,
										Size:  18,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.Row(8, func() {})
							for _, Perms := range aux.PermisosEdificacion {

								m.SetBackgroundColor(blackColor)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text(Perms.Nombre, props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
								})

								m.SetBackgroundColor(whiteColor)
								m.Row(1, func() {})

								m.SetBackgroundColor(darkGrayColor)
								m.Row(7, func() {
									m.Col(4, func() {
										m.Text("Superficie Terreno", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
									m.Col(4, func() {
										m.Text("Numero de Permiso", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
									m.Col(4, func() {
										m.Text("Fecha", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
								})

								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(4, func() {
										m.Text(Perms.Sup_Terreno, props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
									m.Col(4, func() {
										m.Text(Perms.Numero_Permiso, props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
									m.Col(4, func() {
										m.Text(Perms.Fecha_Permiso, props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})

								m.SetBackgroundColor(whiteColor)
								m.Row(1, func() {})

								if Perms.Tipo_Permiso_Edificacion == 7 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Archivos Demolición", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 1 {
											if len(x) > 0 {
												m.Row(7, func() {
													m.Col(3, func() {
														m.Text("Estado", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Fecha", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Permiso", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Plano Aprobado", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
												for _, z := range x {
													m.Row(7, func() {
														m.Col(3, func() {
															m.Text(PdfStr(6, z.Indicar_acoge, ""), props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Fecha, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Nombre2, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Tipo_Permiso_Edificacion == 9 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Archivos Fusión y Subdivisión", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 1 {
											if len(x) > 0 {
												m.Row(7, func() {
													m.Col(3, func() {
														m.Text("Estado", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Fecha", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Permiso", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(3, func() {
														m.Text("Plano Aprobado", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
												for _, z := range x {
													m.Row(7, func() {
														m.Col(3, func() {
															m.Text(PdfStr(6, z.Indicar_acoge, ""), props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Fecha, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(3, func() {
															m.Text(z.Nombre2, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Tipo_Permiso_Edificacion == 8 {
									if Perms.ObraP_Faena == 1 {
										m.SetBackgroundColor(darkGrayColor2)
										m.Row(7, func() {
											m.Col(12, func() {
												m.Text("Obra Preliminar Intalación de Faenas", props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
										})
										for i, x := range Perms.Archivos {
											if i == 1 {
												if len(x) > 0 {
													m.Row(7, func() {
														m.Col(4, func() {
															m.Text("Fecha", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Permiso", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Plano Aprobado", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
													for _, z := range x {
														m.Row(7, func() {
															m.Col(4, func() {
																m.Text(z.Fecha, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre2, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
														})
													}
												}
											}
										}
									}
									if Perms.ObraP_Grua == 1 {
										m.SetBackgroundColor(darkGrayColor2)
										m.Row(7, func() {
											m.Col(12, func() {
												m.Text("Obra Preliminar Instalación de Grúa", props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
										})
										for i, x := range Perms.Archivos {
											if i == 1 {
												if len(x) > 0 {
													m.Row(7, func() {
														m.Col(4, func() {
															m.Text("Fecha", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Permiso", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Plano Aprobado", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
													for _, z := range x {
														m.Row(7, func() {
															m.Col(4, func() {
																m.Text(z.Fecha, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre2, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
														})
													}
												}
											}
										}
									}
									if Perms.ObraP_Excavacion == 1 {
										m.SetBackgroundColor(darkGrayColor2)
										m.Row(7, func() {
											m.Col(12, func() {
												m.Text("Obra Preliminar Ejecución excavaciones, entibaciones y socalzado", props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
										})
										for i, x := range Perms.Archivos {
											if i == 1 {
												if len(x) > 0 {
													m.Row(7, func() {
														m.Col(4, func() {
															m.Text("Fecha", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Permiso", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text("Plano Aprobado", props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
													for _, z := range x {
														m.Row(7, func() {
															m.Col(4, func() {
																m.Text(z.Fecha, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
															m.Col(4, func() {
																m.Text(z.Nombre2, props.Text{
																	Top:   1.5,
																	Size:  9,
																	Style: consts.Bold,
																	Align: consts.Left,
																	Color: col.NewBlack(),
																	Left:  2.0,
																})
															})
														})
													}
												}
											}
										}
									}
								}
								if Perms.Op0 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Proyecto de Telecomunicaciones", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 1 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op1 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Estudio de carga Combustible. Art. 4.3.4 OGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 2 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op2 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Estudio de Seguridad. Art. 4.2.13, 4.2.14, 4.2.15, 4.3.1, 4.3.2, 4.3.6 OGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 3 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op3 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Mecánica de Suelo. Art. 1.2.14 OGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 4 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op4 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Estudio de Evacuación. Art. 1.2.14 OGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 5 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op5 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Plano y memoria de Accesibilidad cuando corresponda (Según Art. 5.1.6 N° 14 OGUC)", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 6 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op6 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Estudio Impacto sobre Sistema de Transporte Urbano (EISTU) Art. 2.4.3, 4.5.4, 4.8.3, 4.13.4 OGUC (Exigible conforme a plazos del Articulo primero transitorio de la Ley N° 20.958)", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 7 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op7 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Autorización Previa del Consejo de Monumentos Nacionales (Zona Típica) Ley 17.288 y sus modificaciones artículo 30 N°1", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 8 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op8 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Autorización SEREMI MINVU inciso segundo Art. 60 LGUC según corresponda", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 9 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op9 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Estudio especifico de riego que incluye medidas y obras de mitigación Art. 2.1.17 OGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 10 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op10 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Informe SEREMI, Art. 60 LGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 11 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op11 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Construcciones en el área rural (Autorización MINAGRI (en caso de loteos) o Informes favorables SAG y SEREMI-MINVU en caso de construcciones) Art. 55 LGUC", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 12 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Op12 == 1 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Otro", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 13 {
											if len(x) > 0 {
												for _, z := range x {
													m.Row(7, func() {
														m.Col(12, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}

								m.SetBackgroundColor(darkGrayColor)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text("Cantidad de Pisos Sobre Nivel", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text("Cantidad de Pisos Bajo Nivel", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
								})
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Cant_Pisos_Sobre_Nivel), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Cant_Pisos_Bajo_Nivel), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})

								m.SetBackgroundColor(darkGrayColor)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text("Superficie Edificada Sobre Nivel", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text("Superficie Edificada Bajo Nivel", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
								})
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Sup_Edificada_Sobre_Nivel), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Sup_Edificada_Bajo_Nivel), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})

								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("Permiso", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range Perms.Archivos {
									if i == 13 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(12, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}

								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("Planos", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range Perms.Archivos {
									if i == 13 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(12, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}

								m.SetBackgroundColor(darkGrayColor)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text("Acogido a Artículos Especiales Transitorios", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text("Recepción Definitiva", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewWhite(),
											Left:  2.0,
										})
									})
								})
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Aco_Art_Esp_Transitorio), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
									m.Col(6, func() {
										m.Text(fmt.Sprintf("%v", Perms.Recepcion_Definitiva), props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})

								if Perms.Recepcion_Definitiva == 2 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Recepcion Total", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 1 {
											if len(x) > 0 {
												m.Row(7, func() {
													m.Col(4, func() {
														m.Text("Documento", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(4, func() {
														m.Text("Fecha", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(4, func() {
														m.Text("Indicar si se acoge a Art. 5.2.8. OGUC", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
												for _, z := range x {
													m.Row(7, func() {
														m.Col(4, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text(z.Fecha, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text(fmt.Sprintf("%v", z.Indicar_acoge), props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}
								if Perms.Recepcion_Definitiva == 3 {
									m.SetBackgroundColor(darkGrayColor2)
									m.Row(7, func() {
										m.Col(12, func() {
											m.Text("Recepcion Parcial", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for i, x := range Perms.Archivos {
										if i == 1 {
											if len(x) > 0 {
												m.Row(7, func() {
													m.Col(4, func() {
														m.Text("Documento", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(4, func() {
														m.Text("Fecha", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(4, func() {
														m.Text("Indicar si se acoge a Art. 5.2.8. OGUC", props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
												for _, z := range x {
													m.Row(7, func() {
														m.Col(4, func() {
															m.Text(z.Nombre, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text(z.Fecha, props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
														m.Col(4, func() {
															m.Text(fmt.Sprintf("%v", z.Indicar_acoge), props.Text{
																Top:   1.5,
																Size:  9,
																Style: consts.Bold,
																Align: consts.Left,
																Color: col.NewBlack(),
																Left:  2.0,
															})
														})
													})
												}
											}
										}
									}
								}

								m.SetBackgroundColor(whiteColor)
								m.Row(6, func() {})
							}
						}

						if addPage {
							m.AddPage()
							addPage = false
						}

						if p3 == 1 {

							m.Row(7, func() {
								m.Col(12, func() {
									m.Text("Situación Técnica", props.Text{
										Top:   1.5,
										Size:  18,
										Style: consts.Bold,
										Align: consts.Left,
										Color: col.NewBlack(),
										Left:  2.0,
									})
								})
							})
							m.Row(8, func() {})
							if aux.Electrico_te1 == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("ELECTRICO TE1", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 1 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Dotacion_ap == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("DOTACION AP", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 2 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Dotacion_alcance == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("DOTACION ALCANCE", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 3 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Instalacion_ascensor == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("INSTALACION ASCENSOR", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 4 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Te1_ascensor == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("TE1 ASCENSOR", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 5 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Certificado_ascensor == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("CERTIFICADO ASCENSOR", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 6 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Clima == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("CLIMA", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 7 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
							if aux.Seguridad_incendio == 1 {
								m.SetBackgroundColor(darkGrayColor2)
								m.Row(7, func() {
									m.Col(12, func() {
										m.Text("SEGURIDAD INCENDIO", props.Text{
											Top:   1.5,
											Size:  9,
											Style: consts.Bold,
											Align: consts.Left,
											Color: col.NewBlack(),
											Left:  2.0,
										})
									})
								})
								for i, x := range aux.ListaArchivos {
									if i == 8 {
										if len(x) > 0 {
											for _, z := range x {
												m.Row(7, func() {
													m.Col(6, func() {
														m.Text(z.Nombre, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
													m.Col(6, func() {
														m.Text(z.Fecha, props.Text{
															Top:   1.5,
															Size:  9,
															Style: consts.Bold,
															Align: consts.Left,
															Color: col.NewBlack(),
															Left:  2.0,
														})
													})
												})
											}
										}
									}
								}
								m.Row(1, func() {})
							}
						}

						m.AddPage()

						m.Row(7, func() {
							m.Col(12, func() {
								m.Text("Situación Comercial", props.Text{
									Top:   1.5,
									Size:  18,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})
						m.Row(8, func() {})

						m.SetBackgroundColor(darkGrayColor)
						m.Row(7, func() {
							m.Col(6, func() {
								m.Text("Tasacion Valor Comercial", props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewWhite(),
									Left:  2.0,
								})
							})
							m.Col(6, func() {
								m.Text("Año Tasacion", props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewWhite(),
									Left:  2.0,
								})
							})
						})
						m.SetBackgroundColor(darkGrayColor2)
						m.Row(7, func() {
							m.Col(6, func() {
								m.Text(fmt.Sprintf("%v", aux.Tasacion_valor_comercial), props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
							m.Col(6, func() {
								m.Text(fmt.Sprintf("%v", aux.Ano_tasacion), props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})

						m.SetBackgroundColor(darkGrayColor2)
						m.Row(7, func() {
							m.Col(12, func() {
								m.Text("Posee Contrato de Arriendo", props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})
						for i, x := range aux.ListaArchivos {
							if i == 1 {
								if len(x) > 0 {
									m.Row(7, func() {
										m.Col(3, func() {
											m.Text("Documento", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Valor Arriendo en Pesos", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Vencimiento", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Renovacion Automatica", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for _, z := range x {
										m.Row(7, func() {
											m.Col(3, func() {
												m.Text(z.Nombre, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Valor_arriendo, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Fecha, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Renovacion_auto, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
										})
									}
								}
							}
						}

						m.SetBackgroundColor(darkGrayColor2)
						m.Row(7, func() {
							m.Col(12, func() {
								m.Text("Posee Contrato de Subarriendo", props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})
						for i, x := range aux.ListaArchivos {
							if i == 1 {
								if len(x) > 0 {
									m.Row(7, func() {
										m.Col(3, func() {
											m.Text("Documento", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Valor Arriendo en Pesos", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Vencimiento", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
										m.Col(3, func() {
											m.Text("Renovacion Automatica", props.Text{
												Top:   1.5,
												Size:  9,
												Style: consts.Bold,
												Align: consts.Left,
												Color: col.NewBlack(),
												Left:  2.0,
											})
										})
									})
									for _, z := range x {
										m.Row(7, func() {
											m.Col(3, func() {
												m.Text(z.Nombre, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Valor_arriendo, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Fecha, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
											m.Col(3, func() {
												m.Text(z.Renovacion_auto, props.Text{
													Top:   1.5,
													Size:  9,
													Style: consts.Bold,
													Align: consts.Left,
													Color: col.NewBlack(),
													Left:  2.0,
												})
											})
										})
									}
								}
							}
						}

						m.AddPage()

						m.Row(7, func() {
							m.Col(12, func() {
								m.Text("Situación Legal", props.Text{
									Top:   1.5,
									Size:  18,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})
						m.Row(8, func() {})

						m.Row(7, func() {
							m.Col(12, func() {
								m.Text("Nombre propietario inscrito en Conservador", props.Text{
									Top:   1.5,
									Size:  18,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})
						m.Row(7, func() {
							m.Col(12, func() {
								m.Text(aux.Nom_propietario_conservador, props.Text{
									Top:   1.5,
									Size:  18,
									Style: consts.Bold,
									Align: consts.Left,
									Color: col.NewBlack(),
									Left:  2.0,
								})
							})
						})

						if p2 == 1 {
						}
						if p3 == 1 {
						}
						if p4 == 1 {
						}
						if p5 == 1 {
						}
						if p6 == 1 {
						}
						if p7 == 1 {
						}
						if p8 == 1 {
						}

						pdf, err := m.Output()

						if err != nil {
							ErrorCheck(err)
							return
						} else {
							ctx.SetBody(pdf.Bytes())
						}
					}
				}
			} else {
				ctx.SetBody([]byte{})
			}
		} else {
			ctx.SetBody([]byte{})
		}
	case "descargar_lista_propiedades":

		ctx.SetContentType("application/octet-stream")
		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				lista := ctx.FormValue("lista")
				var s []int
				if err := json.Unmarshal([]byte(lista), &s); err == nil {
					output, err := CreateExcel(id_emp, s)
					if err == nil {
						ctx.SetBody(output.Bytes())
					} else {
						ctx.SetBody([]byte{})
					}
				}
			} else {
				ctx.SetBody([]byte{})
			}
		} else {
			ctx.SetBody([]byte{})
		}
	default:
		ctx.Response.Header.Set("Content-Type", "application/json")
		json.NewEncoder(ctx).Encode(resp)
	}
}
func ExcelImages(images map[int][]string) [][]ExcelImage {
	str := make([][]ExcelImage, 0)
	count := 0
	for tipoImg, listImgs := range images {
		for _, image := range listImgs {
			fmt.Println(count / 6)
			if count%6 == 0 {
				str = append(str, make([]ExcelImage, 0))
			}
			str[count/6] = append(str[count/6], ExcelImage{image, GetTipo(tipoImg)})
			count++
		}
	}
	return str
}
func GetTipo(t int) string {
	if t == 1 {
		return "Principal"
	} else if t == 2 {
		return "Interior"
	} else {
		return "Exterior"
	}
}
func PdfStr(i int, j int, str string) string {

	strs := make([][]string, 0)
	strs = append(strs, []string{"Sin Seleccionar", "Propio", "Arrendado"})
	strs = append(strs, []string{"Sin Seleccionar", "Uso Propio", "Arrendado a Terceros"})
	strs = append(strs, []string{"Sin Seleccionar", "Si", "No"})
	strs = append(strs, []string{"Sin Seleccionar", "Si", "No"})
	strs = append(strs, []string{"Sin Seleccionar", "Retail", "Servicios", "Industrial", "Similar al industrial", "Salud", "Educacional", "Transporte", "Otros"})
	strs = append(strs, []string{"Sin Seleccionar", "Local comercial", "Restaurant", "Oficina", "Bodega", "Hospital", "Clínica", "Colegio", "Universidad", "Jardín infantil", "Terminal buses", "Estación servicio", "Industria", "Estacionamiento", "Otros"})
	strs = append(strs, []string{"Sin Seleccionar", "Parcial", "Total"})

	if len(strs) > i {
		if len(strs[i]) > j {
			return strs[i][j]
		}
	}
	return str
}
func CreateExcel(id_emp int, lista []int) (*bytes.Buffer, error) {

	/*
		index, err := f.NewSheet("Hola")
		if err != nil {
			fmt.Println(err)
			return
		}

		f.SetSheetName("Sheet1", "Mundo")
		f.SetActiveSheet(index)
	*/

	f := excelize.NewFile()

	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 20, Color: "6d64e8"}, Fill: excelize.Fill{Type: "pattern", Color: []string{"FF0000"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "center"}})
	if err != nil {
		panic(err)
	}

	f.SetCellValue("Sheet1", "B1", "Datos Generales")
	f.MergeCell("Sheet1", "B1", "R1")
	f.SetCellStyle("Sheet1", "B1", "B1", style)

	f.SetCellValue("Sheet1", "S1", "Situación Municipal")
	f.MergeCell("Sheet1", "S1", "S1")
	f.SetCellStyle("Sheet1", "S1", "S1", style)

	f.SetCellValue("Sheet1", "T1", "Situación Técnica")
	f.MergeCell("Sheet1", "T1", "AA1")
	f.SetCellStyle("Sheet1", "T1", "T1", style)

	f.SetCellValue("Sheet1", "AB1", "Situación Comercial")
	f.MergeCell("Sheet1", "AB1", "AE1")
	f.SetCellStyle("Sheet1", "AB1", "AB1", style)

	f.SetCellValue("Sheet1", "AF1", "Situación Legal")
	f.MergeCell("Sheet1", "AF1", "AH1")
	f.SetCellStyle("Sheet1", "AF1", "AF1", style)

	f.SetCellValue("Sheet1", "AI1", "Situación Avaluo Fiscal")
	f.MergeCell("Sheet1", "AI1", "AR1")
	f.SetCellStyle("Sheet1", "AI1", "AI1", style)

	f.SetCellValue("Sheet1", "AS1", "Avaluo Comercial")
	f.MergeCell("Sheet1", "AS1", "AV1")
	f.SetCellStyle("Sheet1", "AS1", "AS1", style)

	f.SetCellValue("Sheet1", "AW1", "Normativo")
	f.MergeCell("Sheet1", "AW1", "BF1")
	f.SetCellStyle("Sheet1", "AW1", "AW1", style)

	listaLetras := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ", "CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ"}
	j := 1
	m := 2

	//DATOS GENERALES 1//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Link") //b
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Nombre") //c
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Direccion") //d
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Numero") //e
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Lat") //f
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Lng") //g
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Pais") //h
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Region") //i
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Ciudad") //j
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Comuna") //k
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Dominio") //l
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Dominio2") //m
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Atencion a público") //n
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Copropiedad") //o
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Uso o destino") //p
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Detalle uso o destino") //q
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Especificar Destino") //r
	j++

	//SITUACION MUNICIPAL 2//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Cantidad") //s
	j++

	//SITUACION TECNICA 3//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Electrico TE1") //t
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Dotacion AP") //u
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Dotacion Alcance") //v
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Instalacion Ascensor") //w
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "TE1 Ascensor") //x
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Certificado Ascensor") //y
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Clima (HVAC)") //z
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Seguridad Contraincendio (CPI)") //aa
	j++

	//SITUACION COMERCIAL 4//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Tasacion Valor Comercial") //ab
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Año Tasacion") //ac
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Posee Contrato de Arriendo") //ad
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Posee Contrato de Subarriendo") //ae
	j++

	//SITUACION LEGAL 5//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Nombre propietario inscrito en Conservador") //af
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "El inmueble posee Gravámenes y/o prohibiciones") //ag
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Posee Planos Archivados") //ah
	j++

	//SITUACION AVALUO FISCAL 6//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Serie") //ai
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Destino") //aj
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Rol Manzana") //ak
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Rol Predio") //al
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Exento") //am
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Avaluo Fiscal") //an
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Contribucion") //ao
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Superficie Terreno") //ap
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Superficie Edificada") //aq
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Superficie Pavimentos") //ar
	j++

	//AVALUO COMERCIAL 7//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Valor Terreno") //as
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Valor Edificacion") //at
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Valor Obras Complementarias") //au
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Valor Total") //av
	j++

	//Normativo 8//
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Certificado Informaciones Previas") //aw
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Tipo de Instrumento") //ax
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Especificar tipo instrumento") //ay
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Destino") //az
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Zona Normativa") //ba
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Usos Permitidos") //bb
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Usos Prohibidos") //bc
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Densidad") //bd
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Coef Constructibilidad") //be
	j++
	f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], m), "Coef Ocupacion Suelo") //bf
	j++

	dominio := []string{"Sin Seleccionar", "Propio", "Arrendado"}
	dominio2 := []string{"Sin Seleccionar", "Uso Propio", "Arrendado a Terceros"}
	atencion_publico := []string{"Sin Seleccionar", "Si", "No"}
	copropiedad := []string{"Sin Seleccionar", "Si", "No"}
	destino := []string{"Sin Seleccionar", "Retail", "Servicios", "Industrial", "Similar al industrial", "Salud", "Educacional", "Transporte", "Otros"}
	detalle_destino := []string{"Sin Seleccionar", "Local comercial", "Restaurant", "Oficina", "Bodega", "Hospital", "Clínica", "Colegio", "Universidad", "Jardín infantil", "Terminal buses", "Estación servicio", "Industria", "Estacionamiento", "Otros"}

	generico_situacion_tecnica := []string{"Sin Seleccionar", "Si", "No", "No Aplica"}

	for i, x := range lista {
		pro, found := GetPropiedad(id_emp, x, false)
		if found {

			//DATOS GENERALES 1//
			display, tooltip := fmt.Sprintf("http://www.redigo.cl/%v", pro.Id), pro.Nombre
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), "Ver")
			f.SetCellHyperLink("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), fmt.Sprintf("http://www.redigo.cl/%v", pro.Id), "External", excelize.HyperlinkOpts{Display: &display, Tooltip: &tooltip}) //b
			j++

			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Nombre) //c
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Direccion) //d
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Numero) //e
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Lat) //f
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Lng) //g
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Pais) //h
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Region) //i
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Ciudad) //j
			j++
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Comuna) //k
			j++
			if pro.Dominio < len(dominio) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), dominio[pro.Dominio]) //l
			}
			j++
			if pro.Dominio2 < len(dominio2) && pro.Dominio == 1 {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), dominio2[pro.Dominio2]) //m
			}
			j++
			if pro.Atencion_publico < len(atencion_publico) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), atencion_publico[pro.Atencion_publico]) //n
			}
			j++
			if pro.Copropiedad < len(copropiedad) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), copropiedad[pro.Copropiedad]) //o
			}
			j++
			if pro.Destino < len(destino) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), destino[pro.Destino]) //p
			}
			j++
			if pro.Detalle_destino < len(detalle_destino) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), detalle_destino[pro.Detalle_destino]) //q
			}
			j++
			if pro.Detalle_destino_otro != "" && pro.Detalle_destino == len(detalle_destino) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Detalle_destino_otro) //r
			}
			j++

			//SITUACION MUNICIPAL 2//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), fmt.Sprintf("%v", len(pro.PermisosEdificacion))) //s
			f.SetCellHyperLink("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), "Sheet2!B2", "Location")
			j++

			//SITUACION TECNICA 3//
			if pro.Electrico_te1 < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Electrico_te1]) //t
			}
			if pro.Dotacion_ap < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Dotacion_ap]) //u
			}
			if pro.Dotacion_alcance < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Dotacion_alcance]) //v
			}
			if pro.Instalacion_ascensor < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Instalacion_ascensor]) //w
			}
			if pro.Te1_ascensor < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Te1_ascensor]) //x
			}
			if pro.Certificado_ascensor < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Certificado_ascensor]) //y
			}
			if pro.Clima < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Clima]) //z
			}
			if pro.Seguridad_incendio < len(generico_situacion_tecnica) {
				f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), generico_situacion_tecnica[pro.Seguridad_incendio]) //aa
			}

			//SITUACION COMERCIAL 4//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Tasacion_valor_comercial) //ab
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Ano_tasacion)             //ac
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Contrato_arriendo)        //ad
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Contrato_subarriendo)     //ae

			//SITUACION LEGAL 5//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Nom_propietario_conservador) //af
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Posee_gp)                    //ag
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Posee_ap)                    //ah

			//SITUACION AVALUO FISCAL 6//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_serie)          //ai
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_destino)        //aj
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Rol_manzana)           //ak
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Rol_predio)            //al
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_exento)         //am
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_avaluo)         //an
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_contribucion)   //ao
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_sup_terreno)    //ap
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_sup_edificada)  //aq
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Fiscal_sup_pavimentos) //ar

			//AVALUO COMERCIAL 7//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Valor_terreno)               //as
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Valor_edificacion)           //at
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Valor_obras_complementarias) //au
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Valor_total)                 //av

			//NORMATIVO 8//
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Cert_info_previas)            //aw
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Tipo_instrumento)             //ax
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Especificar_tipo_instrumento) //ay             //ba
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Usos_permitidos)              //bb
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Usos_prohibidos)              //bc
			f.SetCellValue("Sheet1", fmt.Sprintf("%v%v", listaLetras[j], i+m), pro.Coef_constructibilidad)       //be

		}
	}

	return f.WriteToBuffer()
}
func Delete(ctx *fasthttp.RequestCtx) {

	resp := Response{}
	ctx.Response.Header.Set("Content-Type", "application/json")
	token := string(ctx.Request.Header.Cookie("cu"))

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	switch string(ctx.FormValue("accion")) {
	case "borrar_empresa":

		if found, _ := SuperAdmin(token); found {
			id := Read_uint32bytes(ctx.FormValue("id"))
			resp.Tipo, resp.Titulo, resp.Texto = BorrarEmpresa(db, id)
			if resp.Tipo == "success" {
				resp.Reload = 1
				resp.Page = "crearEmpresa"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Empresa", "No tiene permiso para esta acción"
		}
	case "borrar_propiedad":

		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				id := Read_uint32bytes(ctx.FormValue("id"))
				resp.Tipo, resp.Titulo, resp.Texto = BorrarPropiedad(db, id_emp, id)
				if resp.Tipo == "success" {
					resp.Reload = 1
					resp.Page = "crearPropiedad"
				}
			} else {
				resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Propiedad", "No tiene permiso para esta acción"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Propiedad", "No tiene permiso para esta acción"
		}
	case "borrar_permiso":

		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				id := string(ctx.FormValue("id"))
				var idx string
				resp.Tipo, resp.Titulo, resp.Texto, idx = BorrarPermiso(db, id_emp, id)
				if resp.Tipo == "success" {
					resp.Reload = 1
					resp.Page = fmt.Sprintf("crearPropiedad2PermisoEdificacion?id=%v", idx)
				}
			} else {
				resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Propiedad", "No tiene permiso para esta acción"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Propiedad", "No tiene permiso para esta acción"
		}
	case "borrar_usuarios":

		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				id := Read_uint32bytes(ctx.FormValue("id"))
				resp.Tipo, resp.Titulo, resp.Texto = BorrarUsuario(db, id_emp, id)
				if resp.Tipo == "success" {
					resp.Reload = 1
					resp.Page = "crearUsuarios"
				}
			} else {
				resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Usuario", "No tiene permiso para esta acción"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Usuario", "No tiene permiso para esta acción"
		}
	case "borrar_cotizacion":

		if found, id_emp, listaPermisos := Permisos(token); found {
			if listaPermisos.P0 {
				id := Read_uint32bytes(ctx.FormValue("id"))
				resp.Tipo, resp.Titulo, resp.Texto = BorrarCotizacion(db, id_emp, id)
				if resp.Tipo == "success" {
					resp.Reload = 1
					resp.Page = "misCotizaciones"
				}
			} else {
				resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Cotizacion", "No tiene permiso para esta acción"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Cotizacion", "No tiene permiso para esta acción"
		}
	case "borrar_cotizacion_admin":

		if found, _ := SuperAdmin(token); found {

			id := Read_uint32bytes(ctx.FormValue("id"))
			resp.Tipo, resp.Titulo, resp.Texto = BorrarCotizacionAdmin(db, token, id)
			if resp.Tipo == "success" {
				resp.Reload = 1
				resp.Page = "AdminCotizacion"
			}
		} else {
			resp.Tipo, resp.Titulo, resp.Texto = "error", "Error al eliminar Cotizacion", "No tiene permiso para esta acción"
		}
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

	res, err := db.Query("SELECT id_usr, pass FROM usuarios WHERE user = ? AND eliminado=0", user)
	defer res.Close()
	ErrorCheck(err)

	if res.Next() {

		var id_usr int
		var pass string
		err := res.Scan(&id_usr, &pass)
		ErrorCheck(err)

		if pass == GetMD5Hash(ctx.PostArgs().Peek("pass")) {

			resp.Op = 1
			resp.Msg = ""
			cookie := randSeq(32)

			stmt, err := db.Prepare("INSERT INTO sesiones(cookie, id_usr, fecha) VALUES(?,?, NOW())")
			ErrorCheck(err)
			defer stmt.Close()
			r, err := stmt.Exec(string(cookie), id_usr)
			ErrorCheck(err)
			id_usr, err := r.LastInsertId()
			cookieset := fmt.Sprintf("%v-%v", string(cookie), id_usr)
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
func Cart(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	token := string(ctx.Request.Header.Cookie("cu"))
	Pu := GetPermisoUser(token)

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	accion := string(ctx.PostArgs().Peek("accion"))
	id_cot := ParamInt(ctx.PostArgs().Peek("id_cot"))

	if accion == "add" {

		if id_cot == 0 {

			uf := GetUF()
			stmt, err := db.Prepare("INSERT INTO cotizaciones(id_usr, fecha, precio_uf, id_emp) VALUES(?,NOW(),?,?)")
			ErrorCheck(err)
			defer stmt.Close()
			r, err := stmt.Exec(Pu.Id_usr, uf, Pu.Id_emp)
			ErrorCheck(err)
			idx, err := r.LastInsertId()
			id_cot = int(idx)

		}

		id_pro := ParamInt(ctx.PostArgs().Peek("id_pro"))
		id_ale := ParamInt(ctx.PostArgs().Peek("id_ale"))

		stmt, err := db.Prepare("INSERT INTO cotizacion_detalle(id_cot, id_pro, id_ale) VALUES(?,?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		stmt.Exec(id_cot, id_pro, id_ale)
	}
	if accion == "rm" {

		id_pro := ParamInt(ctx.PostArgs().Peek("id_pro"))
		id_ale := ParamInt(ctx.PostArgs().Peek("id_ale"))

		delForm, err := db.Prepare("DELETE FROM cotizacion_detalle WHERE id_cot=? AND id_pro=? AND id_ale=?")
		ErrorCheck(err)
		delForm.Exec(id_cot, id_pro, id_ale)
		defer db.Close()
	}
	if id_cot > 0 {
		json.NewEncoder(ctx).Encode(DatosCotizacion(id_cot))
	}
}
func Nueva(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	resp := Response{Op: 2}

	pass1 := string(ctx.PostArgs().Peek("pass_01"))
	pass2 := string(ctx.PostArgs().Peek("pass_02"))

	if pass1 == pass2 {

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		code := string(ctx.PostArgs().Peek("code"))
		cn := 0
		res, err := db.Query("SELECT id_usr FROM usuarios WHERE code = ? AND eliminado = ?", code, cn)
		defer res.Close()
		ErrorCheck(err)

		if res.Next() {

			pass := GetMD5Hash(ctx.PostArgs().Peek("pass_01"))

			var id_usr int
			err := res.Scan(&id_usr)
			ErrorCheck(err)
			st := ""
			stmt, err := db.Prepare("UPDATE usuarios SET pass = ?, code = ? WHERE id_usr = ?")
			ErrorCheck(err)
			_, e := stmt.Exec(pass, st, id_usr)
			ErrorCheck(e)
			if e == nil {
				resp.Op = 1
				resp.Msg = ""
			}

		} else {
			resp.Msg = "Se produjo un error"
		}
	} else {
		resp.Msg = "Se produjo un error"
	}

	json.NewEncoder(ctx).Encode(resp)
}
func Cotizacionfunc(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/pdf")
	//token := string(ctx.Request.Header.Cookie("cu"))

	name := ctx.UserValue("name")
	str, ok := name.(string)
	if ok {

		aux := strings.Split(str, "_")
		if len(aux) == 2 {
			aux2 := strings.Split(aux[1], ".")
			if len(aux2) == 2 {

				id, ers := strconv.Atoi(aux2[0])
				ErrorCheck(ers)
				datosCot := DatosCotizacion(id)

				urlqr := fmt.Sprintf("https://www.redigo.cl/cotizacion/%v", id)
				b, qrfile := CreateQr(urlqr, id, "cotizacion")
				if b {

					titulo := fmt.Sprintf("Cotizacion #%v", id)

					darkGrayColor := col.Color{Red: 55, Green: 55, Blue: 55}
					grayColor := col.Color{Red: 220, Green: 220, Blue: 220}
					//grayColor2 := col.Color{Red: 200, Green: 200, Blue: 200}
					whiteColor := col.NewWhite()

					m := pdf.NewMaroto(consts.Portrait, consts.A4)
					m.SetPageMargins(10, 10, 10)

					m.RegisterHeader(func() {
						m.Row(20, func() {
							m.Col(3, func() {
								_ = m.FileImage("./logo.png", props.Rect{
									Center:  true,
									Percent: 100,
								})
							})
							m.Col(6, func() {
								m.Text(titulo, props.Text{
									Top:   6,
									Size:  18,
									Style: consts.Bold,
									Align: consts.Center,
								})
							})
							m.Col(3, func() {
								_ = m.FileImage(qrfile, props.Rect{
									Center:  true,
									Percent: 100,
								})
							})
						})
					})
					m.RegisterFooter(func() {
						m.Row(4, func() {
							m.Col(12, func() {
								m.Text("Esta cotización tiene un vigencia de 30 días", props.Text{
									Top:   16,
									Style: consts.BoldItalic,
									Size:  10,
									Align: consts.Center,
									Color: darkGrayColor,
								})
							})
						})
					})

					m.Row(8, func() {})

					m.Row(5, func() {
						m.ColSpace(1)
						m.Col(2, func() {
							m.Text("Empresa: ", props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: darkGrayColor,
							})
						})
						m.Col(9, func() {
							m.Text(datosCot.NombreEmp, props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: darkGrayColor,
							})
						})
					})
					m.Row(5, func() {
						m.ColSpace(1)
						m.Col(2, func() {
							m.Text("Fecha: ", props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: darkGrayColor,
							})
						})
						m.Col(9, func() {
							m.Text(datosCot.Fecha, props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: darkGrayColor,
							})
						})
					})

					m.Row(8, func() {})

					m.SetBackgroundColor(darkGrayColor)
					m.Row(7, func() {
						m.ColSpace(1)
						m.Col(9, func() {
							m.Text("Lista Servicios", props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: col.NewWhite(),
							})
						})
						m.Col(2, func() {
							m.Text("Precio UF", props.Text{
								Top:   1.5,
								Size:  9,
								Style: consts.Bold,
								Align: consts.Left,
								Color: col.NewWhite(),
							})
						})
					})
					m.SetBackgroundColor(col.Color{Red: 245, Green: 245, Blue: 245})
					m.Row(1, func() {})

					for i := 0; i < len(datosCot.Lista); i++ {

						m.SetBackgroundColor(grayColor)
						m.Row(6, func() {
							m.ColSpace(1)
							m.Col(9, func() {
								m.Text(datosCot.Lista[i].NombreAle, props.Text{
									Top:   2.0,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
								})
							})
							m.Col(2, func() {
								m.Text(fmt.Sprintf("%v", datosCot.Lista[i].Precio), props.Text{
									Top:   1.5,
									Size:  9,
									Style: consts.Bold,
									Align: consts.Left,
								})
							})
						})
						m.Row(4, func() {
							m.ColSpace(1)
							m.Col(9, func() {
								m.Text(datosCot.Lista[i].Propiedad, props.Text{
									Top:   0.2,
									Size:  7,
									Style: consts.Bold,
									Align: consts.Left,
								})
							})
							m.ColSpace(2)
						})
						m.Row(16, func() {
							m.ColSpace(1)
							m.Col(8, func() {
								m.Text(datosCot.Lista[i].Descripcion, props.Text{
									Top:   1.2,
									Size:  8,
									Style: consts.Bold,
									Align: consts.Left,
								})
							})
							m.ColSpace(3)
						})

						m.SetBackgroundColor(col.Color{Red: 235, Green: 235, Blue: 235})
						m.Row(1, func() {})
					}

					m.SetBackgroundColor(whiteColor)

					m.Row(8, func() {})
					m.Row(5, func() {
						m.ColSpace(7)
						m.Col(2, func() {
							m.Text("Total UF:", props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Right,
							})
						})
						m.Col(3, func() {
							m.Text(fmt.Sprintf("%v", datosCot.TotalUf), props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Center,
							})
						})
					})
					m.Row(5, func() {
						m.ColSpace(7)
						m.Col(2, func() {
							m.Text("Valor UF:", props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Right,
							})
						})
						m.Col(3, func() {
							m.Text(SeparadordeMiles(int(datosCot.Uf)), props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Center,
							})
						})
					})
					m.Row(5, func() {
						m.ColSpace(7)
						m.Col(2, func() {
							m.Text("Subtotal:", props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Right,
							})
						})
						m.Col(3, func() {
							m.Text(SeparadordeMiles(int(datosCot.Subtotal)), props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Center,
							})
						})
					})
					m.Row(5, func() {
						m.ColSpace(7)
						m.Col(2, func() {
							m.Text("Iva 19%:", props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Right,
							})
						})
						m.Col(3, func() {
							m.Text(SeparadordeMiles(int(datosCot.Iva)), props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Center,
							})
						})
					})
					m.Row(5, func() {
						m.ColSpace(7)
						m.Col(2, func() {
							m.Text("Total:", props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Right,
							})
						})
						m.Col(3, func() {
							m.Text(SeparadordeMiles(int(datosCot.Total)), props.Text{
								Top:   0,
								Style: consts.Bold,
								Size:  9,
								Align: consts.Center,
							})
						})
					})

					pdf, err := m.Output()

					if err != nil {
						ErrorCheck(err)
						return
					} else {
						ctx.SetBody(pdf.Bytes())
					}

				}
			}
		}
	}
}
func Pages(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	name := ctx.UserValue("name")
	token := string(ctx.Request.Header.Cookie("cu"))

	switch name {
	case "inicioEmpresa":

		if found, id_emp, listaPermisos := Permisos(token); found {

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)
			obj := TemplateInicio{}
			obj.Permisos = listaPermisos
			aux, found := GetEmpresa(id_emp)
			if found {

				obj.Nombre = aux.Nombre
				obj.Precio = aux.Precio
				obj.UF = GetUF()

				res, op := GetResumenPropiedades(id_emp)
				if op {
					obj.Resp = res
				}
			}
			err = t.Execute(ctx, obj)
			ErrorCheck(err)
		}
	case "crearEmpresa":

		if found, _ := SuperAdmin(token); found {

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
					obj.FormPrecio = aux.Precio
					obj.FormId = id
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearAlerta":

		if found, _ := SuperAdmin(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Alerta", "Nueva Alerta", "Configurar", "Titulo Lista", "guardar_alerta", fmt.Sprintf("/pages/%s", name), "borrar_alerta", "Alerta")
			lista, found := GetAlertas()
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetAlerta(id, 0)
				if found {
					obj.FormId = id
					obj.Alertas = aux
				}
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearRegla":

		if found, _ := SuperAdmin(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			id_ale := Read_uint32bytes(ctx.QueryArgs().Peek("id_ale"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Regla Alerta", "Configurar Regla", "Configurar", "Lista de Reglas", "guardar_regla", fmt.Sprintf("/pages/%s", name), "borrar_regla", "Regla")

			if id_ale > 0 {
				aux, found := GetAlerta(id_ale, id)
				if found {
					obj.FormId = id
					obj.FormIdAle = id_ale
					obj.Alertas = aux
				}
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearUsuarios":

		if found, id_emp, listaPermisos := Permisos(token); found {

			if listaPermisos.P0 {
				id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
				t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
				ErrorCheck(err)

				obj := GetTemplateConf("Crear Usuarios", "Subtitulo", "Subtitulo2", "Titulo Usuarios", "guardar_usuarios", fmt.Sprintf("/pages/%s", name), "borrar_usuario", "Usuario")
				obj.Permisos = listaPermisos
				lista, found := GetUsuarios(id_emp)
				if found {
					obj.Lista = lista
				}

				if id > 0 {
					aux, nombre, found := GetUsuario(id_emp, id)
					if found {
						obj.FormNombre = nombre
						obj.FormId = id
						obj.Permisos = aux

					}
				}
				fmt.Println(obj.Permisos)

				err = t.Execute(ctx, obj)
				ErrorCheck(err)
			} else {
				t, _ := TemplatePage("html/ErrorPermisos.html")
				_ = t.Execute(ctx, nil)
			}

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "AdminCotizacion":

		if found, _ := SuperAdmin(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Mis Cotizaciones", "Subtitulo", "Subtitulo2", "Titulo Usuarios", "guardar_admin_cotizacion", fmt.Sprintf("/pages/%s", name), "borrar_cotizacion_admin", "Cotizacion")
			obj.Lista = GetAllListaCotizaciones()
			emps, found := GetEmpresas()
			if found {
				obj.Lista2 = emps
			}

			if id > 0 {
				obj.FormIdRec, obj.PrecioUf = GetCotizacion(id)
				obj.FormId = id
			} else {
				obj.FormId = 0
				obj.PrecioUf = GetUF()
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "confCotizacion":

		if found, _ := SuperAdmin(token); found {

			id_cot := Read_uint32bytes(ctx.QueryArgs().Peek("id_cot"))

			id_ale := Read_uint32bytes(ctx.QueryArgs().Peek("id_ale"))
			id_pro := Read_uint32bytes(ctx.QueryArgs().Peek("id_pro"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			lista, id_emp, nombreemp := GetAlertasCotizaciones(id_cot)

			obj := GetTemplateConf("Cotizacion", fmt.Sprintf("Cotizacion #%v", id_cot), nombreemp, "Titulo Usuarios", "guardar_detalle_cotizacion", fmt.Sprintf("/pages/%s", name), "borrar_detalle_cotizacion", "Cotizacion")

			obj.Lista = lista
			obj.FormId = id_cot

			fmt.Println(id_emp)
			listaP, found := GetPropiedades(id_emp)
			if found {
				obj.Lista2 = listaP
			}

			listaA, found2 := GetAlertas()
			if found2 {
				obj.Lista3 = listaA
			}

			obj.FormIdAle = id_ale
			obj.FormIdRec = id_pro

			if id_ale > 0 && id_pro > 0 {

				obj.Descripcion, obj.FormPrecio = GetDetalleCotizacion(id_cot, id_ale, id_pro)
				obj.Valor = 1

			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "envCotizacion":

		if found, _ := SuperAdmin(token); found {

			id_cot := Read_uint32bytes(ctx.QueryArgs().Peek("id_cot"))

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Mis Cotizaciones", "Subtitulo", "Subtitulo2", "Titulo Usuarios", "", fmt.Sprintf("/pages/%s", name), "borrar_cotizacion", "Cotizacion")

			obj.Lista, obj.FormId = GetUserFromCot(id_cot)

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "misCotizaciones":

		if found, id_emp, listaPermisos := Permisos(token); found {

			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Mis Cotizaciones", "Subtitulo", "Subtitulo2", "Titulo Usuarios", "", fmt.Sprintf("/pages/%s", name), "borrar_cotizacion", "Cotizacion")
			obj.Permisos = listaPermisos
			obj.Lista = GetListaCotizaciones(id_emp)

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Datos Generales", "Completar los datos", "Lista de Propiedades", "guardar_propiedad1", fmt.Sprintf("/pages/%s", name), "borrar_propiedad", "Propiedad")
			obj.Permisos = listaPermisos
			lista, found := GetPropiedades(id_emp)
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					obj.Titulo = "Editar Propiedad"
				}
			} else {
				obj.FormId = 0
			}
			err = t.Execute(ctx, obj)
			ErrorCheck(err)
		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad2PermisoEdificacion":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			id_rec := Read_uint32bytes(ctx.QueryArgs().Peek("id_rec"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Municipal", "Completar los datos", "Lista de Permisos de Edificación", "guardar_propiedad2A", fmt.Sprintf("/pages/%s", name), "borrar_permiso", "Permiso Edificación")
			obj.Permisos = listaPermisos

			if id > 0 {
				obj.FormId = id
				obj.FormIdRec = id_rec
				CamposPropiedades, found := GetPropiedad2A(id, id_emp, id_rec)
				if found {
					obj.CamposPropiedades = CamposPropiedades
					fmt.Println(CamposPropiedades)
				}
			}
			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad3":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Técnica", "Completar los datos", "", "guardar_propiedad3", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					if aux.Is_Arrendado == 1 {
						obj.NextPage = 4
					} else {
						obj.NextPage = 5
					}
					if aux.P3 == 1 {
						obj.Titulo = "Editar Propiedad"
					}
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)
		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad4":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Comercial", "Completar los datos", "", "guardar_propiedad4", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					if aux.P4 == 1 {
						obj.Titulo = "Editar Propiedad"
					}
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad5":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Legal", "Completar los datos", "", "guardar_propiedad5", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {

					if aux.P5 == 1 {
						obj.Titulo = "Editar Propiedad"
					}

					obj.CamposPropiedades = aux
					obj.FormId = id

					if aux.Is_Arrendado == 1 {
						obj.NextPage = 4
					} else {
						obj.NextPage = 3
					}

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad6":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Avalúo Fiscal", "Completar los datos", "", "guardar_propiedad6", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					if aux.P6 == 1 {
						obj.Titulo = "Editar Propiedad"
					}
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad7":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Avalúo Comercial", "Completar los datos", "", "guardar_propiedad7", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					if aux.P7 == 1 {
						obj.Titulo = "Editar Propiedad"
					}
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "crearPropiedad8":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Normativo", "Completar los datos", "", "guardar_propiedad8", fmt.Sprintf("/pages/%s", name), "", "")
			obj.Permisos = listaPermisos

			if id > 0 {
				aux, found := GetPropiedad(id_emp, id, false)
				if found {
					obj.CamposPropiedades = aux
					obj.FormId = id
					if aux.P8 == 1 {
						obj.Titulo = "Editar Propiedad"
					}
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "buscarPropiedades":

		if found, id_emp, listaPermisos := Permisos(token); found {

			listaPropiedades, errp := GetPropiedadesCompleto(id_emp)
			if errp {
				t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
				ErrorCheck(err)

				obj := buscarPropiedades{}
				obj.Permisos = listaPermisos

				propiedades, err := json.Marshal(listaPropiedades)
				ErrorCheck(err)
				obj.PropiedadesString = string(propiedades)

				obj.Titulo = "Titulo"
				obj.SubTitulo = "Subtitulo"
				obj.SubTitulo2 = "Subtitulo2"
				obj.SubTitulo3 = "Subtitulo3"
				obj.SubTitulo4 = "Subtitulo4"

				err = t.Execute(ctx, obj)
				ErrorCheck(err)
			} else {
				t, _ := TemplatePage("html/ErrorPermisos.html")
				_ = t.Execute(ctx, nil)
			}

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "detallePropiedad":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("", "Datos Generales", "Completar los datos", "Lista de Propiedades", "guardar_propiedad1", fmt.Sprintf("/pages/%s", name), "borrar_propiedad", "Propiedad")
			obj.Permisos = listaPermisos

			aux, found := GetPropiedad(id_emp, id, true)
			if found {
				obj.Id_emp = id_emp
				obj.Titulo = aux.Nombre
				obj.FormId = id
				obj.CamposPropiedades = aux
				err = t.Execute(ctx, obj)
				ErrorCheck(err)
			} else {
				t, _ := TemplatePage("html/ErrorPermisos.html")
				_ = t.Execute(ctx, nil)
			}

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	case "descargarPropiedad":

		if found, id_emp, listaPermisos := Permisos(token); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("", "Datos Generales", "Completar los datos", "Lista de Propiedades", "descargar_custom_pdf", fmt.Sprintf("/pages/%s", name), "borrar_propiedad", "Propiedad")
			obj.Permisos = listaPermisos

			aux, found := GetPropiedad(id_emp, id, false)
			if found {
				obj.Titulo = aux.Nombre
				obj.FormId = id
				obj.Is_Arrendado = aux.Is_Arrendado
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		} else {
			t, _ := TemplatePage("html/ErrorPermisos.html")
			_ = t.Execute(ctx, nil)
		}
	default:
		ctx.NotFound()
	}
}
func Acciones(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set("Content-Type", "application/json")
	resp := Response{Op: 2}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	accion := string(ctx.PostArgs().Peek("accion"))

	switch accion {
	case "recuperar_password":

		user := string(ctx.PostArgs().Peek("user"))

		res, err := db.Query("SELECT id_usr FROM usuarios WHERE user = ? AND eliminado=0", user)
		defer res.Close()
		ErrorCheck(err)

		if res.Next() {

			var id_usr int
			err := res.Scan(&id_usr)
			ErrorCheck(err)
			cm, cu := GetMail(db, id_usr)

			if cm < 500 && cu < 2 {

				tipo := 1
				rand := string(randSeq(32))
				stmt, err := db.Prepare("INSERT INTO correo_enviados (tipo, code, fecha, id_usr) VALUES (?,?,NOW(),?)")
				ErrorCheck(err)
				defer stmt.Close()
				r, err := stmt.Exec(tipo, rand, id_usr)
				ErrorCheck(err)
				if err == nil {
					idcor, err := r.LastInsertId()
					ErrorCheck(err)
					body := fmt.Sprintf("<a href='https://www.redigo.cl/recuperar/%v/%v'>Recuperar</a>", rand, idcor)
					fmt.Println(body)
					/*
						if SendEmail(user, "Recuperar Contraseña", body) {
							resp.Op = 1
						}
					*/
					resp.Op = 1
					resp.Msg = "Correo ha sido enviada"
				} else {
					resp.Msg = "Se produjo un error"
				}
			} else {
				resp.Msg = "Se ha excedido la cantidad de correos"
			}
		} else {
			resp.Msg = "Usuario no existe"
		}
	case "nueva_password":

		id := string(ctx.PostArgs().Peek("id"))
		code := string(ctx.PostArgs().Peek("code"))

		pass1 := string(ctx.PostArgs().Peek("pass_01"))
		pass2 := ctx.PostArgs().Peek("pass_02")

		if pass1 == string(pass2) {
			if len(code) == 32 {

				res, err := db.Query("SELECT id_usr FROM correo_enviados WHERE id_cor = ? AND code = ? AND tipo=1", id, code)
				defer res.Close()
				ErrorCheck(err)

				if res.Next() {

					var id_usr int
					err := res.Scan(&id_usr)
					ErrorCheck(err)

					rand := ""
					pass := GetMD5Hash(pass2)

					fmt.Println(pass2)
					fmt.Println(string(pass2))
					fmt.Println(pass)

					stmt1, err1 := db.Prepare("UPDATE correo_enviados SET code = ? WHERE id_cor = ?")
					ErrorCheck(err1)
					_, e1 := stmt1.Exec(rand, id)
					ErrorCheck(e1)

					stmt2, err2 := db.Prepare("UPDATE usuarios SET pass = ? WHERE id_usr = ?")
					ErrorCheck(err2)
					_, e2 := stmt2.Exec(pass, id_usr)
					ErrorCheck(e2)

					resp.Op = 1
					resp.Msg = "Contraseña establecida"
				} else {
					resp.Msg = "Se produjo un error con el codigo"
				}
			} else {
				resp.Msg = "Error con el codigo"
			}
		} else {
			resp.Msg = "Contraseñas no coinciden"
		}
	default:
	}

	json.NewEncoder(ctx).Encode(resp)
}
func Recuperar(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	t, err := TemplatePage("html/recuperar.html")
	ErrorCheck(err)
	var x Rec
	x.Rec = false
	err = t.Execute(ctx, x)
	ErrorCheck(err)
}
func Recuperar2(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	name := ctx.UserValue("name")
	id := ctx.UserValue("id")

	strname, ok1 := name.(string)
	strid, ok2 := id.(string)
	if ok1 && ok2 {
		t, err := TemplatePage("html/recuperar.html")
		ErrorCheck(err)
		var x Rec
		x.Rec = true
		x.Code = strname
		x.Id = strid
		err = t.Execute(ctx, x)
		ErrorCheck(err)
	}
}
func Index(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	token := string(ctx.Request.Header.Cookie("cu"))
	gpu := GetPermisoUser(token)

	if gpu.Bool {

		param := string(ctx.FormValue("p"))
		id := string(ctx.FormValue("i"))
		if param == "detalle_propiedad" && id != "" {
			gpu.Function = fmt.Sprintf("detallePropiedad?id=%v", id)
		}

		t, err := TemplatePage("html/inicio.html")
		ErrorCheck(err)
		err = t.Execute(ctx, gpu)
		ErrorCheck(err)

	} else {
		t, _ := TemplatePage("html/login.html")
		_ = t.Execute(ctx, nil)
	}
}
func Salir(ctx *fasthttp.RequestCtx) {

	tkn := string(ctx.Request.Header.Cookie("cu"))

	if len(tkn) > 32 {
		token := tkn[0:32]
		id_ses := Read_uint32bytes([]byte(tkn[33:]))

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		delForm, err := db.Prepare("DELETE FROM sesiones WHERE id_ses=? AND cookie=?")
		ErrorCheck(err)
		delForm.Exec(id_ses, token)
		defer db.Close()
	}

	ctx.Redirect("/", 200)
}
func SetEmpresa(ctx *fasthttp.RequestCtx) {

	token := string(ctx.Request.Header.Cookie("cu"))

	if found, id_usr := SuperAdmin(token); found {

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		cn := 1
		id_emp := ctx.UserValue("name")

		stmt, err := db.Prepare("UPDATE usuarios SET id_emp = ? WHERE admin = ? AND id_usr = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(id_emp, cn, id_usr)
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
func Images(ctx *fasthttp.RequestCtx) {

	name := ctx.UserValue("name")
	id := ctx.UserValue("id")
	//token := string(ctx.Request.Header.Cookie("cu"))
	ctx.SendFile(fmt.Sprintf("./files/images/%v/%v", id, name))
}
func GetPermisoUser(tkn string) PermisoUser {

	Pu := PermisoUser{}

	if len(tkn) > 32 {

		token := tkn[0:32]
		id_ses := Read_uint32bytes([]byte(tkn[33:]))

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		res, err := db.Query("SELECT t1.pass_impuesta, t1.id_usr, t1.admin, t1.id_emp FROM usuarios t1, sesiones t2 WHERE t2.id_ses = ? AND t2.cookie = ? AND t2.id_usr=t1.id_usr", id_ses, token)
		defer res.Close()
		ErrorCheck(err)

		var admin int
		var id_emp int
		var id_usr int
		var pass_impuesta int

		if res.Next() {

			err := res.Scan(&pass_impuesta, &id_usr, &admin, &id_emp)
			ErrorCheck(err)

			if admin == 1 {
				Pu.Admin = true
				Pu.Function = "crearEmpresa"
			} else {
				Pu.Admin = false
				if pass_impuesta == 0 {
					Pu.Function = "inicioEmpresa"
				} else {
					Pu.Function = "cambiarPassword"
				}
			}

			if id_emp > 0 {
				Pu.Idemp = true
			} else {
				Pu.Idemp = false
			}

			Pu.Id_emp = id_emp
			Pu.Id_usr = id_usr

			Pu.Bool = true
		} else {
			Pu.Bool = false
			Pu.Admin = false
		}
	}

	return Pu
}

// FUNCTION DB //
func GetMySQLDB() (db *sql.DB, err error) {
	//CREATE DATABASE redigo CHARACTER SET utf8 COLLATE utf8_spanish2_ci;
	db, err = sql.Open("mysql", fmt.Sprintf("root:%v@tcp(127.0.0.1:3306)/redigo", pass.Passwords.PassDb))
	return
}
func GetUserFromCot(id int) ([]Lista, int) {

	lista := []Lista{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t2.id_usr, t2.nombre, t2.user, t1.id_usr as id_user FROM cotizaciones t1, usuarios t2 WHERE t1.id_cot = ? AND t1.id_emp=t2.id_emp", id)
	defer res.Close()
	ErrorCheck(err)

	var id_usr int
	var nombre string
	var user string
	var id_user int

	for res.Next() {

		err := res.Scan(&id_usr, &nombre, &user, &id_user)
		ErrorCheck(err)
		lista = append(lista, Lista{Id: id_usr, Nombre: fmt.Sprintf("%v (%v)", nombre, user)})

	}

	return lista, id_user
}
func Permisos(tkn string) (bool, int, ListaPermisos) {

	listaPermisos := ListaPermisos{}

	if len(tkn) > 32 {
		token := tkn[0:32]
		id_ses := Read_uint32bytes([]byte(tkn[33:]))

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		res, err := db.Query("SELECT t1.p0, t1.p1, t1.p2, t1.p3, t1.p4, t1.p5, t1.p6, t1.p7, t1.p8, t1.p9, t1.id_emp, t1.admin FROM usuarios t1, sesiones t2 WHERE t2.id_ses = ? AND t2.cookie = ? AND t2.id_usr=t1.id_usr", id_ses, token)
		defer res.Close()
		ErrorCheck(err)

		var p0 int
		var p1 int
		var p2 int
		var p3 int
		var p4 int
		var p5 int
		var p6 int
		var p7 int
		var p8 int
		var p9 int
		var id_emp int
		var admin int

		if res.Next() {

			err := res.Scan(&p0, &p1, &p2, &p3, &p4, &p5, &p6, &p7, &p8, &p9, &id_emp, &admin)
			ErrorCheck(err)
			if admin == 1 {
				listaPermisos.P0 = true
				listaPermisos.P1 = true
				listaPermisos.P2 = true
				listaPermisos.P3 = true
				listaPermisos.P4 = true
				listaPermisos.P5 = true
				listaPermisos.P6 = true
				listaPermisos.P7 = true
				listaPermisos.P8 = true
				listaPermisos.P9 = true
				return true, id_emp, listaPermisos
			} else {
				if p0 == 1 {
					listaPermisos.P0 = true
				}
				if p1 == 1 {
					listaPermisos.P1 = true
				}
				if p2 == 1 {
					listaPermisos.P2 = true
				}
				if p3 == 1 {
					listaPermisos.P3 = true
				}
				if p4 == 1 {
					listaPermisos.P4 = true
				}
				if p5 == 1 {
					listaPermisos.P5 = true
				}
				if p6 == 1 {
					listaPermisos.P6 = true
				}
				if p7 == 1 {
					listaPermisos.P7 = true
				}
				if p8 == 1 {
					listaPermisos.P8 = true
				}
				if p9 == 1 {
					listaPermisos.P9 = true
				}
				return true, id_emp, listaPermisos
			}
		}
	}

	return false, 0, listaPermisos
}
func SuperAdmin(tkn string) (bool, int) {

	if len(tkn) > 32 {
		token := tkn[0:32]
		id_ses := Read_uint32bytes([]byte(tkn[33:]))

		db, err := GetMySQLDB()
		defer db.Close()
		ErrorCheck(err)

		res, err := db.Query("SELECT t1.id_usr FROM usuarios t1, sesiones t2 WHERE t2.id_ses = ? AND t2.cookie = ? AND t2.id_usr=t1.id_usr AND t1.admin=1", id_ses, token)
		defer res.Close()
		ErrorCheck(err)
		var id_usr int

		if res.Next() {
			err := res.Scan(&id_usr)
			ErrorCheck(err)
			return true, id_usr
		} else {
			return false, 0
		}
	} else {
		return false, 0
	}
}
func GetListaCotizaciones(id int) []Lista {

	lista := []Lista{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT id_cot, fecha FROM cotizaciones WHERE id_emp = ? AND eliminado = ?", id, cn)
	defer res.Close()
	ErrorCheck(err)

	for res.Next() {

		var id_cot int
		var fecha string
		err := res.Scan(&id_cot, &fecha)
		ErrorCheck(err)
		lista = append(lista, Lista{Id: id_cot, Nombre: fecha})

	}
	return lista
}
func GetAllListaCotizaciones() []Lista {

	lista := []Lista{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT id_cot, fecha FROM cotizaciones WHERE eliminado = ?", cn)
	defer res.Close()
	ErrorCheck(err)

	for res.Next() {

		var id_cot int
		var fecha string
		err := res.Scan(&id_cot, &fecha)
		ErrorCheck(err)
		lista = append(lista, Lista{Id: id_cot, Nombre: fecha})

	}
	return lista
}
func GetEmpresa(id int) (Empresa, bool) {

	data := Empresa{}

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
func GetAlerta(id_ale int, id_alr int) (Alertas, bool) {

	alerta := Alertas{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT nombre, descripcion, alerta, notificacion, precio FROM alertas WHERE id_ale = ? AND eliminado = ?", id_ale, cn)
	defer res.Close()
	if err != nil {
		ErrorCheck(err)
		return alerta, false
	}
	for res.Next() {

		err := res.Scan(&alerta.Nombre, &alerta.Descripcion, &alerta.Alerta, &alerta.Notificacion, &alerta.Precio)
		if err != nil {
			ErrorCheck(err)
			return alerta, false
		}
		reglas, reglaindex, found := GetReglas(db, id_ale, id_alr)
		if found {
			alerta.ReglasIndex = reglaindex
			alerta.Reglas = reglas
		}
	}
	return alerta, true
}
func GetReglas(db *sql.DB, id_ale int, id_alr int) ([]Regla, Regla, bool) {

	regla := Regla{}
	reglaindex := Regla{}
	reglas := make([]Regla, 0)

	cn := 0
	res, err := db.Query("SELECT id_alr, nombre, pagina, campo, valor, tipo FROM alerta_regla WHERE eliminado = ? AND id_ale = ?", cn, id_ale)
	defer res.Close()
	if err != nil {
		return reglas, reglaindex, false
	}

	for res.Next() {
		err := res.Scan(&regla.Id_alr, &regla.Nombre, &regla.Pagina, &regla.Campo, &regla.Valor, &regla.Tipo)
		if err != nil {
			return reglas, reglaindex, false
		}
		if regla.Id_alr == id_alr {
			reglaindex = regla
		}
		reglas = append(reglas, regla)
	}
	return reglas, reglaindex, true
}
func GetMail(db *sql.DB, id_usr int) (int, int) {

	date := time.Now()
	date2 := date.Add(-24 * time.Hour)
	fecha := date2.Format("2006-01-02 15:04:05")

	res, err := db.Query("SELECT id_usr, tipo FROM correo_enviados WHERE fecha > ?", fecha)
	defer res.Close()
	ErrorCheck(err)

	var count int = 0
	var count2 int = 0

	for res.Next() {

		var dbid_usr int = 0
		var dbtipo int = 0
		err := res.Scan(&dbid_usr, &dbtipo)
		ErrorCheck(err)
		count++

		if dbtipo == 1 && dbid_usr == id_usr {
			count2++
		}
	}

	return count, count2
}
func GetAlertas() ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT id_ale, nombre, alerta FROM alertas WHERE eliminado = ?", cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var nombre string
	var alerta int

	for res.Next() {

		err := res.Scan(&id, &nombre, &alerta)
		ErrorCheck(err)
		data = append(data, Lista{Id: id, Nombre: nombre, Aux: alerta})
		b = true

	}
	return data, b
}
func GetUsuarios(id_emp int) ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	res, err := db.Query("SELECT id_usr, user FROM usuarios WHERE id_emp = ? AND eliminado = ? AND admin = ?", id_emp, cn, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

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
func GetUsuario(id_emp int, id_usr int) (ListaPermisos, string, bool) {

	data := ListaPermisos{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	var user string

	cn := 0
	res, err := db.Query("SELECT user, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9 FROM usuarios WHERE id_usr = ? AND eliminado = ? AND id_emp = ?", id_usr, cn, id_emp)
	defer res.Close()
	if err != nil {
		return data, user, false
	}
	if res.Next() {

		err := res.Scan(&user, &data.P0, &data.P1, &data.P2, &data.P3, &data.P4, &data.P5, &data.P6, &data.P7, &data.P8, &data.P9)
		if err != nil {
			return data, user, false
		}
		return data, user, true

	} else {
		return data, user, false
	}
}
func PermisosEdificacion(id int) ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	res, err := db.Query("SELECT id_rec, tipo, especificar_tipo, fecha FROM permiso_edificacion WHERE id_pro = ? AND eliminado = ?", id, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var nombre string

	for res.Next() {

		var id_pro int
		var tipo int
		var especificar string
		var fecha string
		err := res.Scan(&id_pro, &tipo, &especificar, &fecha)
		if err != nil {
			log.Fatal(err)
		}

		if tipo == 1 {
			nombre = fmt.Sprintf("Obra Nueva %v", fecha)
		}
		if tipo == 2 {
			nombre = fmt.Sprintf("Obra menor alteración (Menor a 100m2) %v", fecha)
		}
		if tipo == 3 {
			nombre = fmt.Sprintf("Obra menor alteración (Mayor a 100m2) %v", fecha)
		}
		if tipo == 4 {
			nombre = fmt.Sprintf("Modificación de Proyecto %v", fecha)
		}
		if tipo == 5 {
			nombre = fmt.Sprintf("Alteración %v", fecha)
		}
		if tipo == 6 {
			nombre = fmt.Sprintf("Recontrucción %v", fecha)
		}
		if tipo == 7 {
			nombre = fmt.Sprintf("%v %v", especificar, fecha)
		}

		data = append(data, Lista{Id: id_pro, Nombre: nombre})
		b = true

	}
	return data, b
}
func GetResumenPropiedades(id_emp int) (Resumen, bool) {

	resp := Resumen{}
	resp.Prods = make(map[int]ResumenProds)
	resp.Alertas = make(map[int]ResumenAlertas)
	resp.Notificaciones = make(map[int]ResumenAlertas)
	resp.TotalAlertas = 0
	resp.TotalNotificaciones = 0

	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	propiedades, errp := GetPropiedadesCompleto(id_emp)
	if errp {
		resp.Localidades = propiedades
	}

	res, err := db.Query("SELECT t1.id_pro as id_pro, t1.nombre as nombre, t2.id_ale as id_ale, t3.nombre as nombre_ale, t3.notificacion as tipo_notificacion FROM propiedades t1, propiedad_alerta t2, alertas t3 WHERE t1.id_emp = ? AND t1.eliminado = ? AND t1.id_pro=t2.id_pro AND t2.id_ale=t3.id_ale AND t3.eliminado = ?", id_emp, cn, cn)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		var id_pro int
		var nombre string
		var id_ale int
		var tipo_notificacion int
		var nombre_ale string
		err := res.Scan(&id_pro, &nombre, &id_ale, &nombre_ale, &tipo_notificacion)
		if err != nil {
			log.Fatal(err)
		}

		if prod, ok := resp.Prods[id_pro]; ok {
			prod.Lista = append(prod.Lista, ResumenAlerta{Id: id_ale, Nombre: nombre_ale})
			resp.Prods[id_pro] = prod
		} else {
			aux := ResumenProds{Nombre: nombre, Lista: make([]ResumenAlerta, 0)}
			aux.Lista = append(aux.Lista, ResumenAlerta{Id: id_ale, Nombre: nombre_ale})
			resp.Prods[id_pro] = aux
		}

		if tipo_notificacion == 1 {

			resp.TotalNotificaciones++
			if ale, ok := resp.Notificaciones[id_ale]; ok {
				ale.Lista = append(ale.Lista, ResumenProd{Id: id_ale, Nombre: nombre_ale})
				resp.Notificaciones[id_ale] = ale
			} else {
				aux := ResumenAlertas{Nombre: nombre_ale, Lista: make([]ResumenProd, 0)}
				aux.Lista = append(aux.Lista, ResumenProd{Id: id_pro, Nombre: nombre})
				resp.Notificaciones[id_ale] = aux
			}
		}
		if tipo_notificacion == 2 {

			resp.TotalAlertas++
			if ale, ok := resp.Alertas[id_ale]; ok {
				ale.Lista = append(ale.Lista, ResumenProd{Id: id_ale, Nombre: nombre_ale})
				resp.Alertas[id_ale] = ale
			} else {
				aux := ResumenAlertas{Nombre: nombre_ale, Lista: make([]ResumenProd, 0)}
				aux.Lista = append(aux.Lista, ResumenProd{Id: id_pro, Nombre: nombre})
				resp.Alertas[id_ale] = aux
			}
		}

		b = true
	}
	return resp, b
}
func GetPropiedades(id_emp int) ([]Lista, bool) {

	data := []Lista{}
	b := false

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0

	res, err := db.Query("SELECT id_pro, nombre FROM propiedades WHERE id_emp = ? AND eliminado = ?", id_emp, cn)
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
func GetPropiedad(id_emp int, id int, perms bool) (CamposPropiedades, bool) {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	data := CamposPropiedades{}
	cn := 0

	res, err := db.Query("SELECT t1.nombre, t1.direccion, t1.numero, t1.lat, t1.lng, t1.rangopos, t1.dominio, t1.dominio2, t1.atencion_publico, t1.copropiedad, t1.destino, t1.detalle_destino, t1.detalle_destino_otro, t1.id_com, t1.id_ciu, t1.id_reg, t1.id_pai, t1.electrico_te1, t1.dotacion_ap, t1.dotacion_alcance, t1.instalacion_ascensor, t1.te1_ascensor, t1.certificado_ascensor, t1.clima, t1.seguridad_incendio, t1.tasacion_valor_comercial, t1.ano_tasacion, t1.contrato_arriendo, t1.contrato_subarriendo, t1.nompropietarioconservador, t1.posee_gp, t1.posee_ap, t1.fiscal_serie, t1.fiscal_destino, t1.rol_manzana, t1.rol_predio, t1.fiscal_exento, t1.fiscal_avaluo, t1.fiscal_contribucion, t1.fiscal_sup_terreno, t1.fiscal_sup_edificada, t1.fiscal_sup_pavimentos, t1.valor_terreno, t1.valor_edificacion, t1.valor_obras_complementarias, t1.valor_total, t1.cert_info_previas, t1.tipo_instrumento, t1.especificar_tipo_instrumento, t1.usos_permitidos, t1.usos_prohibidos, t1.coef_constructibilidad, t2.nombre, t3.nombre, t4.nombre, t5.nombre, t1.p1, t1.p2, t1.p3, t1.p4, t1.p5, t1.p6, t1.p7, t1.p8 FROM propiedades t1, paises t2, regiones t3, ciudades t4, comunas t5 WHERE t1.id_pro = ? AND t1.eliminado = ? AND t1.id_emp = ? AND t1.id_pai=t2.id_pai AND t1.id_reg=t3.id_reg AND t1.id_ciu=t4.id_ciu AND t1.id_com=t5.id_com", id, cn, id_emp)
	defer res.Close()
	if err != nil {
		return data, false
	}

	if res.Next() {
		err := res.Scan(&data.Nombre, &data.Direccion, &data.Numero, &data.Lat, &data.Lng, &data.RangoPos, &data.Dominio, &data.Dominio2, &data.Atencion_publico, &data.Copropiedad, &data.Destino, &data.Detalle_destino, &data.Detalle_destino_otro, &data.Id_com, &data.Id_ciu, &data.Id_reg, &data.Id_pai, &data.Electrico_te1, &data.Dotacion_ap, &data.Dotacion_alcance, &data.Instalacion_ascensor, &data.Te1_ascensor, &data.Certificado_ascensor, &data.Clima, &data.Seguridad_incendio, &data.Tasacion_valor_comercial, &data.Ano_tasacion, &data.Contrato_arriendo, &data.Contrato_subarriendo, &data.Nom_propietario_conservador, &data.Posee_gp, &data.Posee_ap, &data.Fiscal_serie, &data.Fiscal_destino, &data.Rol_manzana, &data.Rol_predio, &data.Fiscal_exento, &data.Fiscal_avaluo, &data.Fiscal_contribucion, &data.Fiscal_sup_terreno, &data.Fiscal_sup_edificada, &data.Fiscal_sup_pavimentos, &data.Valor_terreno, &data.Valor_edificacion, &data.Valor_obras_complementarias, &data.Valor_total, &data.Cert_info_previas, &data.Tipo_instrumento, &data.Especificar_tipo_instrumento, &data.Usos_permitidos, &data.Usos_prohibidos, &data.Coef_constructibilidad, &data.Pais, &data.Region, &data.Ciudad, &data.Comuna, &data.P1, &data.P2, &data.P3, &data.P4, &data.P5, &data.P6, &data.P7, &data.P8)
		if err != nil {
			return data, false
		}
		if data.Dominio == 2 {
			data.Is_Arrendado = 1
		}
		if data.Dominio == 1 && data.Dominio2 == 2 {
			data.Is_Arrendado = 1
		}

		if perms {
			permisos, found := GetPropiedad2A(id, id_emp, 0)
			if found {
				data.PermisosEdificacion = permisos.PermisosEdificacion
				data.PermisosEdificacionIndex = permisos.PermisosEdificacionIndex
			}
		}

		lista_archivos, ultimos_archivos, found := GetArchivosPropiedad(id)
		if found {
			data.ListaArchivos = lista_archivos
			data.UltimosArchivos = ultimos_archivos
		}

		imagenes, foundimg := GetImagenesPropiedad(id)
		if foundimg {
			data.Imagenes = imagenes
		}

		return data, true
	} else {
		return data, false
	}
}
func GetPropiedadesCompleto(id_emp int) ([]CamposPropiedades, bool) {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	resp := make([]CamposPropiedades, 0)
	data := CamposPropiedades{}
	cn := 0

	res, err := db.Query("SELECT t1.id_pro, t1.nombre, t1.direccion, t1.numero, t1.lat, t1.lng, t1.dominio, t1.dominio2, t1.atencion_publico, t1.copropiedad, t1.destino, t1.detalle_destino, t1.detalle_destino_otro, t1.id_com, t1.id_ciu, t1.id_reg, t1.id_pai, t1.electrico_te1, t1.dotacion_ap, t1.dotacion_alcance, t1.instalacion_ascensor, t1.te1_ascensor, t1.certificado_ascensor, t1.clima, t1.seguridad_incendio, t1.tasacion_valor_comercial, t1.ano_tasacion, t1.contrato_arriendo, t1.contrato_subarriendo, t1.nompropietarioconservador, t1.posee_gp, t1.posee_ap, t1.fiscal_serie, t1.fiscal_destino, t1.rol_manzana, t1.rol_predio, t1.fiscal_exento, t1.fiscal_avaluo, t1.fiscal_contribucion, t1.fiscal_sup_terreno, t1.fiscal_sup_edificada, t1.fiscal_sup_pavimentos, t1.valor_terreno, t1.valor_edificacion, t1.valor_obras_complementarias, t1.valor_total, t1.cert_info_previas, t1.tipo_instrumento, t1.especificar_tipo_instrumento, t1.usos_permitidos, t1.usos_prohibidos, t1.coef_constructibilidad, t2.nombre, t3.nombre, t4.nombre, t5.nombre FROM propiedades t1, paises t2, regiones t3, ciudades t4, comunas t5 WHERE t1.eliminado = ? AND t1.id_emp = ? AND t1.id_pai=t2.id_pai AND t1.id_reg=t3.id_reg AND t1.id_ciu=t4.id_ciu AND t1.id_com=t5.id_com", cn, id_emp)
	defer res.Close()
	if err != nil {
		return resp, false
	}

	for res.Next() {

		err := res.Scan(&data.Id, &data.Nombre, &data.Direccion, &data.Numero, &data.Lat, &data.Lng, &data.Dominio, &data.Dominio2, &data.Atencion_publico, &data.Copropiedad, &data.Destino, &data.Detalle_destino, &data.Detalle_destino_otro, &data.Id_com, &data.Id_ciu, &data.Id_reg, &data.Id_pai, &data.Electrico_te1, &data.Dotacion_ap, &data.Dotacion_alcance, &data.Instalacion_ascensor, &data.Te1_ascensor, &data.Certificado_ascensor, &data.Clima, &data.Seguridad_incendio, &data.Tasacion_valor_comercial, &data.Ano_tasacion, &data.Contrato_arriendo, &data.Contrato_subarriendo, &data.Nom_propietario_conservador, &data.Posee_gp, &data.Posee_ap, &data.Fiscal_serie, &data.Fiscal_destino, &data.Rol_manzana, &data.Rol_predio, &data.Fiscal_exento, &data.Fiscal_avaluo, &data.Fiscal_contribucion, &data.Fiscal_sup_terreno, &data.Fiscal_sup_edificada, &data.Fiscal_sup_pavimentos, &data.Valor_terreno, &data.Valor_edificacion, &data.Valor_obras_complementarias, &data.Valor_total, &data.Cert_info_previas, &data.Tipo_instrumento, &data.Especificar_tipo_instrumento, &data.Usos_permitidos, &data.Usos_prohibidos, &data.Coef_constructibilidad, &data.Pais, &data.Region, &data.Ciudad, &data.Comuna)
		if err != nil {
			return resp, false
		}
		if data.Dominio == 2 {
			data.Is_Arrendado = 1
		}
		if data.Dominio == 1 && data.Dominio2 == 2 {
			data.Is_Arrendado = 1
		}
		lista_archivos, ultimos_archivos, found := GetArchivosPropiedad(data.Id)
		if found {
			data.ListaArchivos = lista_archivos
			data.UltimosArchivos = ultimos_archivos
		}
		resp = append(resp, data)
	}
	return resp, true
}
func GetPropiedad2A(id_pro int, id_emp int, id_rec int) (CamposPropiedades, bool) {

	data := CamposPropiedades{}
	perm := PermisoEdificacion{}

	db, err := GetMySQLDB()
	defer db.Close()
	if err != nil {
		return data, false
	}

	cn := 0
	res, err := db.Query("SELECT id_rec, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar_tipo_permiso_edificacion, numero_permiso, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obrap_faena, obrap_grua, obrap_excavacion, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12 FROM permiso_edificacion WHERE eliminado = ? AND id_emp = ? AND id_pro = ?", cn, id_emp, id_pro)
	defer res.Close()
	if err != nil {
		return data, false
	}
	for res.Next() {
		err := res.Scan(&perm.Id_rec, &perm.Sup_Terreno, &perm.Posee_Permiso_Edificacion, &perm.Tipo_Permiso_Edificacion, &perm.Especificar_Tipo_Permiso_Edificacion, &perm.Numero_Permiso, &perm.Fecha_Permiso, &perm.Cant_Pisos_Sobre_Nivel, &perm.Cant_Pisos_Bajo_Nivel, &perm.Sup_Edificada_Sobre_Nivel, &perm.Sup_Edificada_Sobre_Nivel, &perm.Aco_Art_Esp_Transitorio, &perm.Recepcion_Definitiva, &perm.ObraP_Faena, &perm.ObraP_Grua, &perm.ObraP_Excavacion, &perm.Op0, &perm.Op1, &perm.Op2, &perm.Op3, &perm.Op4, &perm.Op5, &perm.Op6, &perm.Op7, &perm.Op8, &perm.Op9, &perm.Op10, &perm.Op11, &perm.Op12)
		if err != nil {
			return data, false
		}

		perm.Nombre = NombrePermiso(perm.Tipo_Permiso_Edificacion, perm.Especificar_Tipo_Permiso_Edificacion)
		auxlista, ultimos_archivos, auxfound := ArchivosPermisosEdificacion(perm.Id_rec)
		if auxfound {
			perm.Archivos = auxlista
			perm.UltimosArchivos = ultimos_archivos
		}

		if perm.Id_rec == id_rec {
			data.PermisosEdificacionIndex = perm
		}

		data.PermisosEdificacion = append(data.PermisosEdificacion, perm)
	}
	return data, true
}
func NombrePermiso(tipo int, especificar string) string {

	if tipo == 1 {
		return "Obra Nueva"
	}
	if tipo == 2 {
		return "Obra menor alteración (Menor a 100m2)"
	}
	if tipo == 3 {
		return "Obra menor alteración (Mayor a 100m2)"
	}
	if tipo == 4 {
		return "Modificación de Proyecto"
	}
	if tipo == 5 {
		return "Alteración"
	}
	if tipo == 6 {
		return "Recontrucción"
	}
	if tipo == 7 {
		return especificar
	}
	return "Sin Nombre"
}
func ArchivosPermisosEdificacion(id int) (map[int][]ArchivosPermisoEdificacion, map[int]ArchivosPermisoEdificacion, bool) {

	lista := make(map[int][]ArchivosPermisoEdificacion, 0)
	for i := 1; i < 28; i++ {
		lista[i] = []ArchivosPermisoEdificacion{}
	}

	aux := ArchivosPermisoEdificacion{}
	ultimos_archivos := make(map[int]ArchivosPermisoEdificacion, 0)

	db, err := GetMySQLDB()
	defer db.Close()
	if err != nil {
		return lista, ultimos_archivos, false
	}

	cn := 0
	res, err := db.Query("SELECT id_arc, nombre, nombre2, tipo, indicar_acoge, fecha FROM permiso_edificacion_archivos WHERE eliminado = ? AND id_rec = ?", cn, id)
	defer res.Close()
	if err != nil {
		return lista, ultimos_archivos, false
	}
	for res.Next() {

		err := res.Scan(&aux.Id_arc, &aux.Nombre, &aux.Nombre2, &aux.Tipo, &aux.Indicar_acoge, &aux.Fecha)
		if err != nil {
			return lista, ultimos_archivos, false
		}
		lista[aux.Tipo] = append(lista[aux.Tipo], aux)
		ultimos_archivos[aux.Tipo] = aux
	}
	return lista, ultimos_archivos, true
}
func GetArchivosPropiedad(id int) (map[int][]ListaArchivos, map[int]ListaArchivos, bool) {

	lista := make(map[int][]ListaArchivos, 0)
	for i := 1; i < 17; i++ {
		lista[i] = []ListaArchivos{}
	}

	aux := ListaArchivos{}
	ultimos_archivos := make(map[int]ListaArchivos, 0)

	db, err := GetMySQLDB()
	defer db.Close()
	if err != nil {
		return lista, ultimos_archivos, false
	}

	cn := 0
	res, err := db.Query("SELECT id_arc, nombre, fojas, numero, ano, tipo, fecha, fecha_insert, valor_arriendo, renovacion_auto, tipo_de_plano FROM propiedad_archivos WHERE eliminado = ? AND id_pro = ?", cn, id)
	defer res.Close()
	if err != nil {
		return lista, ultimos_archivos, false
	}
	for res.Next() {

		err := res.Scan(&aux.Id_arc, &aux.Nombre, &aux.Fojas, &aux.Numero, &aux.Ano, &aux.Tipo, &aux.Fecha, &aux.Fecha_insert, &aux.Valor_arriendo, &aux.Renovacion_auto, &aux.Tipo_de_Plano)
		if err != nil {
			return lista, ultimos_archivos, false
		}
		lista[aux.Tipo] = append(lista[aux.Tipo], aux)
		ultimos_archivos[aux.Tipo] = aux
	}
	return lista, ultimos_archivos, true
}
func GetImagenesPropiedad(id int) (map[int][]string, bool) {

	lista := make(map[int][]string, 0)

	db, err := GetMySQLDB()
	defer db.Close()
	if err != nil {
		return lista, false
	}

	cn := 0
	res, err := db.Query("SELECT nombre, tipo FROM propiedades_imagenes WHERE eliminado = ? AND id_pro = ?", cn, id)
	defer res.Close()
	if err != nil {
		return lista, false
	}
	for res.Next() {

		var nombre string
		var tipo int

		err := res.Scan(&nombre, &tipo)
		if err != nil {
			return lista, false
		}
		lista[tipo] = append(lista[tipo], nombre)
	}
	return lista, true
}
func InsertPropiedad(db *sql.DB, id_emp int, nombre string, lat string, lng string, rangopos string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, dominio2 string, atencion_publico string, copropiedad string, destino string, detalle_destino string) (uint8, string, int) {

	if nombre != "" {
		id_pai, b1 := GetPais(db, pais)
		id_reg, b2 := GetRegion(db, region, id_pai)
		id_ciu, b3 := GetCiudad(db, ciudad, id_pai, id_reg)
		id_com, b4 := GetComuna(db, comuna, id_pai, id_reg, id_ciu)
		if b1 && b2 && b3 && b4 {
			p := 1
			stmt, err := db.Prepare("INSERT INTO propiedades (nombre, lat, lng, rangopos, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, p1, id_emp) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			ErrorCheck(err)
			defer stmt.Close()
			r, err := stmt.Exec(nombre, lat, lng, rangopos, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, p, id_emp)
			ErrorCheck(err)
			if err == nil {
				id, err := r.LastInsertId()
				if err == nil {
					return 1, "Propiedad ingresada correctamente", int(id)
				} else {
					return 2, "La Propiedad no pudo ser ingresada", 0
				}
			} else {
				return 2, "La Propiedad no pudo ser ingresada", 0
			}
		} else {
			return 2, "Error al ingresar posicion", 0
		}
	} else {
		return 2, "Debe ingresar nombre", 0
	}
}
func UpdatePropiedad(db *sql.DB, id_emp int, id int, nombre string, lat string, lng string, rangopos string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, dominio2 string, atencion_publico string, copropiedad string, destino string, detalle_destino string) (uint8, string) {

	id_pai, b1 := GetPais(db, pais)
	id_reg, b2 := GetRegion(db, region, id_pai)
	id_ciu, b3 := GetCiudad(db, ciudad, id_pai, id_reg)
	id_com, b4 := GetComuna(db, comuna, id_pai, id_reg, id_ciu)
	if b1 && b2 && b3 && b4 {
		stmt, err := db.Prepare("UPDATE propiedades SET nombre = ?, lat = ?, lng = ?, rangopos = ?, id_ciu = ?, id_com = ?, id_reg = ?, id_pai = ?, direccion = ?, numero = ?, dominio = ?, dominio2 = ?, atencion_publico = ?, copropiedad = ?, destino = ?, detalle_destino = ? WHERE id_pro = ? AND id_emp = ?")
		if err != nil {
			return 2, "La empresa no pudo ser actualizada"
		}
		_, e := stmt.Exec(nombre, lat, lng, rangopos, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, id, id_emp)
		if e == nil {
			return 1, "Empresa actualizada correctamente"
		} else {
			return 2, "La empresa no pudo ser actualizada"
		}
	} else {
		return 2, "Error al ingresar posicion"
	}
}
func UpdatePropiedad2A(db *sql.DB, id_emp int, id_pro int, id_rec int, sup_terreno string, posee_permiso_edificacion string, tipo_permiso_edificacion string, especificar string, numero_permiso_edificacion string, fecha_permiso string, cant_pisos_sobre_nivel string, cant_pisos_bajo_nivel string, superficie_edificada_sobre_nivel string, superficie_edificada_bajo_nivel string, aco_art_esp_transitorio string, recepcion_definitiva string, obra_p0 string, obra_p1 string, obra_p2 string, op0 string, op1 string, op2 string, op3 string, op4 string, op5 string, op6 string, op7 string, op8 string, op9 string, op10 string, op11 string, op12 string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE permiso_edificacion SET p2 = ?, sup_terreno = ?, posee_permiso_edificacion = ?, tipo_permiso_edificacion = ?, especificar_tipo_permiso_edificacion = ?, numero_permiso = ?, fecha_permiso = ?, cant_pisos_sobre_nivel = ?, cant_pisos_bajo_nivel = ?, superficie_edificada_sobre_nivel = ?, superficie_edificada_bajo_nivel = ?, aco_art_esp_transitorio = ?, recepcion_definitiva = ?, obrap_faena = ?, obrap_grua = ?, obrap_excavacion = ?, op0 = ?, op1 = ?, op2 = ?, op3 = ?, op4 = ?, op5 = ?, op6 = ?, op7 = ?, op8 = ?, op9 = ?, op10 = ?, op11 = ?, op12 = ? WHERE id_rec = ? AND id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar, numero_permiso_edificacion, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obra_p0, obra_p1, obra_p2, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12, id_rec, id_pro, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func InsertPropiedad2A(db *sql.DB, id_emp int, id_pro int, sup_terreno string, posee_permiso_edificacion string, tipo_permiso_edificacion string, especificar string, numero_permiso_edificacion string, fecha_permiso string, cant_pisos_sobre_nivel string, cant_pisos_bajo_nivel string, superficie_edificada_sobre_nivel string, superficie_edificada_bajo_nivel string, aco_art_esp_transitorio string, recepcion_definitiva string, obra_p0 string, obra_p1 string, obra_p2 string, op0 string, op1 string, op2 string, op3 string, op4 string, op5 string, op6 string, op7 string, op8 string, op9 string, op10 string, op11 string, op12 string) (uint8, string, int) {

	p := 1
	stmt, err := db.Prepare("INSERT INTO permiso_edificacion (p2, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar_tipo_permiso_edificacion, numero_permiso, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obrap_faena, obrap_grua, obrap_excavacion, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12, id_pro, id_emp) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	r, err := stmt.Exec(p, sup_terreno, posee_permiso_edificacion, tipo_permiso_edificacion, especificar, numero_permiso_edificacion, fecha_permiso, cant_pisos_sobre_nivel, cant_pisos_bajo_nivel, superficie_edificada_sobre_nivel, superficie_edificada_bajo_nivel, aco_art_esp_transitorio, recepcion_definitiva, obra_p0, obra_p1, obra_p2, op0, op1, op2, op3, op4, op5, op6, op7, op8, op9, op10, op11, op12, id_pro, id_emp)
	ErrorCheck(err)
	if err == nil {
		id, err := r.LastInsertId()
		if err == nil {
			return 1, "Propiedad ingresada correctamente", int(id)
		} else {
			return 2, "El Permiso no pudo ser ingresado", 0
		}
	} else {
		return 2, "El Permiso no pudo ser ingresado", 0
	}
}
func UpdatePropiedad3(db *sql.DB, id_emp int, id int, electrico_te1 string, dotacion_ap string, dotacion_alcance string, instalacion_ascensor string, te1_ascensor string, certificado_ascensor string, clima string, seguridad_incendio string) (uint8, string, uint8) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p3 = ?, electrico_te1 = ?, dotacion_ap = ?, dotacion_alcance = ?, instalacion_ascensor = ?, te1_ascensor = ?, certificado_ascensor = ?, clima = ?, seguridad_incendio = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, electrico_te1, dotacion_ap, dotacion_alcance, instalacion_ascensor, te1_ascensor, certificado_ascensor, clima, seguridad_incendio, id, id_emp)
	if err == nil {
		var arrendado uint8 = 0
		if IsArrendado(id) {
			arrendado = 1
		}
		return 1, "Empresa actualizada correctamente", arrendado
	} else {
		return 2, "La empresa no pudo ser actualizada", 0
	}
}
func IsArrendado(id int) bool {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT dominio, dominio2 FROM propiedades WHERE id_pro = ?", id)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}

	var dominio int
	var dominio2 int

	if res.Next() {

		err := res.Scan(&dominio, &dominio2)
		if err != nil {
			log.Fatal(err)
		}
		if dominio == 2 {
			return true
		}
		if dominio == 1 && dominio2 == 2 {
			return true
		}
		return false

	} else {
		return false
	}
}
func UpdatePropiedad4(db *sql.DB, id_emp int, id int, tasacion_valor_comercial string, ano_tasacion string, contrato_arriendo string, contrato_subarriendo string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p4 = ?, tasacion_valor_comercial = ?, ano_tasacion = ?, contrato_arriendo = ?, contrato_subarriendo = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, tasacion_valor_comercial, ano_tasacion, contrato_arriendo, contrato_subarriendo, id, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func UpdatePropiedad5(db *sql.DB, id_emp int, id int, domnompropietario string, posee_gp string, posee_ap string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p5 = ?, nompropietarioconservador = ?, posee_gp = ?, posee_ap = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, domnompropietario, posee_gp, posee_ap, id, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func UpdatePropiedad6(db *sql.DB, id_emp int, id int, fiscal_serie string, fiscal_destino string, rol_manzana string, rol_predio string, fiscal_exento string, avaluo_fiscal string, contribucion_fiscal string, superficie_terreno string, superficie_edificada string, superficie_pavimentos string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p6 = ?, fiscal_serie = ?, fiscal_destino = ?, rol_manzana = ?, rol_predio = ?, fiscal_exento = ?, fiscal_avaluo = ?, fiscal_contribucion = ?, fiscal_sup_terreno = ?, fiscal_sup_edificada = ?, fiscal_sup_pavimentos = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, fiscal_serie, fiscal_destino, rol_manzana, rol_predio, fiscal_exento, avaluo_fiscal, contribucion_fiscal, superficie_terreno, superficie_edificada, superficie_pavimentos, id, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func UpdatePropiedad7(db *sql.DB, id_emp int, id int, valor_terreno string, valor_edificacion string, valor_obras_complementarias string, valor_total string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p7 = ?, valor_terreno = ?, valor_edificacion = ?, valor_obras_complementarias = ?, valor_total = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, valor_terreno, valor_edificacion, valor_obras_complementarias, valor_total, id, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func UpdatePropiedad8(db *sql.DB, id_emp int, id int, cert_info_previas string, tipo_instrumento string, especificar_tipo_instrumento string, indicar_area string, zona_normativa_plan_regulador string, area_riesgo string, area_proteccion string, zona_conservacion_historica string, zona_tipica string, monumento_nacional string, zona_uso_suelo string, usos_permitidos string, usos_prohibidos string, superficie_predial_minima string, densidad_maxima_bruta string, densidad_maxima_neta string, altura_maxima string, sistema_agrupamiento string, coef_constructibilidad string, coef_ocupacion_suelo string, coef_ocupacion_suelo_psuperiores string, rasante string, adosamiento string, distanciamiento string, cierres_perimetrales_altura string, cierres_perimetrales_transparencia string, ochavos string, ochavos_metros string, estado_urbanizacion_ejecutada string, estado_urbanizacion_recibida string, estado_urbanizacion_garantizada string) (uint8, string) {

	p := 1
	stmt, err := db.Prepare("UPDATE propiedades SET p8 = ?, cert_info_previas = ?, tipo_instrumento = ?, especificar_tipo_instrumento = ?, indicar_area = ?, zona_normativa_plan_regulador = ?, area_riesgo = ?, area_proteccion = ?, zona_conservacion_historica = ?, zona_tipica = ?, monumento_nacional = ?, zona_uso_suelo = ?, usos_permitidos = ?, usos_prohibidos = ?, superficie_predial_minima = ?, densidad_maxima_bruta = ?, densidad_maxima_neta = ?, altura_maxima = ?, sistema_agrupamiento = ?, coef_constructibilidad = ?, coef_ocupacion_suelo = ?, coef_ocupacion_suelo_psuperiores = ?, rasante = ?, adosamiento = ?, distanciamiento = ?, cierres_perimetrales_altura = ?, cierres_perimetrales_transparencia = ?, ochavos = ?, ochavos_metros = ?, estado_urbanizacion_ejecutada = ?, estado_urbanizacion_recibida = ?, estado_urbanizacion_garantizada = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(p, cert_info_previas, tipo_instrumento, especificar_tipo_instrumento, indicar_area, zona_normativa_plan_regulador, area_riesgo, area_proteccion, zona_conservacion_historica, zona_tipica, monumento_nacional, zona_uso_suelo, usos_permitidos, usos_prohibidos, superficie_predial_minima, densidad_maxima_bruta, densidad_maxima_neta, altura_maxima, sistema_agrupamiento, coef_constructibilidad, coef_ocupacion_suelo, coef_ocupacion_suelo_psuperiores, rasante, adosamiento, distanciamiento, cierres_perimetrales_altura, cierres_perimetrales_transparencia, ochavos, ochavos_metros, estado_urbanizacion_ejecutada, estado_urbanizacion_garantizada, id, id_emp)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func BorrarPropiedad(db *sql.DB, id_emp int, id int) (string, string, string) {

	del := 1
	stmt, err := db.Prepare("UPDATE propiedades SET eliminado = ? WHERE id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(del, id, id_emp)
	if err == nil {
		return "success", "Propiedad eliminada", "Propiedad eliminada correctamente"
	} else {
		return "error", "Error al eliminar propiedad", "La propiedad no pudo ser eliminada"
	}
}
func BorrarCotizacion(db *sql.DB, id_emp int, id int) (string, string, string) {

	del := 1
	stmt, err := db.Prepare("UPDATE cotizaciones SET eliminado = ? WHERE id_cot = ? AND id_emp = ?")
	_, err = stmt.Exec(del, id, id_emp)
	if err == nil {
		return "success", "Cotizacion eliminada", "Cotizacion eliminada correctamente"
	} else {
		return "error", "Error al eliminar cotizacion", "La cotizacion no pudo ser eliminada"
	}
}
func BorrarCotizacionAdmin(db *sql.DB, token string, id int) (string, string, string) {

	del := 1
	stmt, err := db.Prepare("UPDATE cotizaciones SET eliminado = ? WHERE id_cot = ?")
	_, err = stmt.Exec(del, id)
	if err == nil {
		return "success", "Cotizacion eliminada", "Cotizacion eliminada correctamente"
	} else {
		return "error", "Error al eliminar cotizacion", "La cotizacion no pudo ser eliminada"
	}
}
func BorrarPermiso(db *sql.DB, id_emp int, id string) (string, string, string, string) {

	s := strings.Split(id, "/")
	if len(s) == 2 {
		del := 1
		stmt, err := db.Prepare("UPDATE permiso_edificacion SET eliminado = ? WHERE id_rec = ? AND id_pro = ? AND id_emp = ?")
		_, err = stmt.Exec(del, s[1], s[0], id_emp)
		if err == nil {
			return "success", "Permiso Edificacion eliminado", "Permiso Edificacion eliminado correctamente", s[0]
		} else {
			return "error", "Error al eliminar Permiso Edificacion", "El Permiso Edificacion no pudo ser eliminada", ""
		}
	} else {
		return "error", "Error al eliminar Permiso Edificacion", "Error inesperado", ""
	}
}
func InsertUsuario(db *sql.DB, id_emp int, nombre string, pass string, p0 string, p1 string, p2 string, p3 string, p4 string, p5 string, p6 string, p7 string, p8 string, p9 string) (uint8, string) {

	cn := 0
	res, err := db.Query("SELECT id_usr FROM usuarios WHERE user = ? AND eliminado = ?", nombre, cn)
	defer res.Close()
	if err != nil {
		return 2, "Error Base de Datos"
	}
	if res.Next() {
		var id_usr int
		err := res.Scan(&id_usr)
		if err != nil {
			return 2, "Error Base de Datos"
		}
		return 2, "El correo ya existe en el sistema"
	} else {

		pass_impuesta := 0
		if pass != "" {
			pass_impuesta = 1
		}

		code := string(randSeq(32))
		stmt, err := db.Prepare("INSERT INTO usuarios (user, admin, pass, pass_impuesta, pass_fecha, code, id_emp, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, eliminado) VALUES (?,?,?,?,Now(),?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		defer stmt.Close()
		admin := 0
		_, err = stmt.Exec(nombre, admin, pass, pass_impuesta, code, id_emp, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, admin)
		if err == nil {
			return 1, "Usuario ingresada correctamente"
		} else {
			return 2, "El usuario no pudo ser ingresada"
		}
	}
}
func UpdateUsuario(db *sql.DB, id_emp int, id int, nombre string, pass string, p0 string, p1 string, p2 string, p3 string, p4 string, p5 string, p6 string, p7 string, p8 string, p9 string) (uint8, string) {

	cn := 0
	res, err := db.Query("SELECT id_usr FROM usuarios WHERE user = ? AND eliminado = ?", nombre, cn)
	defer res.Close()
	if err != nil {
		return 2, "Error Base de Datos"
	}
	if res.Next() {
		var id_usr int
		err := res.Scan(&id_usr)
		if err != nil {
			return 2, "Error Base de Datos"
		}
		return 2, "El correo ya existe en el sistema"
	} else {

		stmt, err := db.Prepare("UPDATE usuarios SET user = ?, pass = ?, p0 = ?, p1 = ?, p2 = ?, p3 = ?, p4 = ?, p5 = ?, p6 = ?, p7 = ?, p8 = ?, p9 = ? WHERE id_usr = ? AND id_emp = ?")
		_, err = stmt.Exec(nombre, pass, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, id, id_emp)
		if err == nil {
			return 1, "Usuario actualizada correctamente"
		} else {
			return 2, "El usuario no pudo ser actualizada"
		}
	}
}
func BorrarUsuario(db *sql.DB, id_emp int, id int) (string, string, string) {

	del := 1
	admin := 0
	stmt, err := db.Prepare("UPDATE usuarios SET eliminado = ? WHERE id_usr = ? AND id_emp = ? AND admin = ?")
	_, err = stmt.Exec(del, id, id_emp, admin)
	if err == nil {
		return "success", "Usuario eliminada", "Usuario eliminada correctamente"
	} else {
		return "error", "Error al eliminar usuario", "El usuario no pudo ser eliminada"
	}
}
func InsertEmpresa(db *sql.DB, nombre string, precio string) (uint8, string) {

	stmt, err := db.Prepare("INSERT INTO empresa (nombre, precio) VALUES (?,?)")
	defer stmt.Close()
	stmt.Exec(nombre, precio)
	if err == nil {
		return 1, "Empresa ingresada correctamente"
	} else {
		return 2, "La empresa no pudo ser ingresada"
	}
}
func UpdateEmpresa(db *sql.DB, id int, nombre string, precio string) (uint8, string) {

	stmt, err := db.Prepare("UPDATE empresa SET nombre = ?, precio = ? WHERE id_emp = ?")
	_, err = stmt.Exec(nombre, precio, id)
	if err == nil {
		return 1, "Empresa actualizada correctamente"
	} else {
		return 2, "La empresa no pudo ser actualizada"
	}
}
func BorrarEmpresa(db *sql.DB, id int) (string, string, string) {

	del := 1
	stmt, err := db.Prepare("UPDATE empresa SET eliminado = ? WHERE id_emp = ?")
	_, err = stmt.Exec(del, id)
	if err == nil {
		return "success", "Empresa eliminada", "Empresa eliminada correctamente"
	} else {
		return "error", "Error al eliminar empresa", "La empresa no pudo ser eliminada"
	}
}
func InsertRegla(db *sql.DB, nombre string, tipo int, pagina int, campo string, valor string, id_ale int) (uint8, string) {

	stmt, err := db.Prepare("INSERT INTO alerta_regla (nombre, tipo, pagina, campo, valor, id_ale) VALUES (?,?,?,?,?,?)")
	defer stmt.Close()
	_, err = stmt.Exec(nombre, tipo, pagina, campo, valor, id_ale)
	if err == nil {
		return 1, "Regla ingresada correctamente"
	} else {
		return 2, "La regla no pudo ser ingresada"
	}
}
func UpdateRegla(db *sql.DB, id int, nombre string, tipo int, pagina int, campo string, valor string, id_ale int) (uint8, string) {

	stmt, err := db.Prepare("UPDATE alerta_regla SET nombre = ?, tipo = ?, pagina = ?, campo = ?, valor = ? WHERE id_alr = ?")
	_, err = stmt.Exec(nombre, tipo, pagina, campo, valor, id)
	if err == nil {
		return 1, "Regla actualizada correctamente"
	} else {
		return 2, "La regla no pudo ser actualizada"
	}
}
func BorrarRegla(db *sql.DB, token string, id string) Response {

	resp := Response{}
	s := strings.Split(id, "/")
	if len(s) == 2 {
		if found, _ := SuperAdmin(token); found {
			del := 1
			stmt, err := db.Prepare("UPDATE alerta_regla SET eliminado = ? WHERE id_alr = ? AND id_ale = ?")
			ErrorCheck(err)
			_, e := stmt.Exec(del, s[1], s[0])
			ErrorCheck(e)
			if e == nil {

				resp.Tipo = "success"
				resp.Reload = 1
				resp.Page = fmt.Sprintf("crearRegla?id=%v", s[0])
				resp.Titulo = "Regla eliminada"
				resp.Texto = "Regla eliminada correctamente"
			} else {
				resp.Tipo = "error"
				resp.Titulo = "Error al eliminar regla"
				resp.Texto = "La regla no pudo ser eliminada"
			}
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar regla"
			resp.Texto = "No tiene los permisos"
		}
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar regla"
		resp.Texto = "No tiene los permisos"
	}

	return resp
}
func InsertAlerta(db *sql.DB, nombre string, descripcion string, alerta string, notificacion string, precio string) (uint8, string) {

	stmt, err := db.Prepare("INSERT INTO alertas (nombre, descripcion, alerta, notificacion, precio) VALUES (?,?,?,?,?)")
	defer stmt.Close()
	stmt.Exec(nombre, descripcion, alerta, notificacion, precio)
	if err == nil {
		return 1, "Alerta ingresada correctamente"
	} else {
		return 2, "La alerta no pudo ser ingresada"
	}
}
func UpdateAlerta(db *sql.DB, id int, nombre string, descripcion string, alerta string, notificacion string, precio string) (uint8, string) {

	stmt, err := db.Prepare("UPDATE alertas SET nombre = ?, descripcion = ?, alerta = ?, notificacion = ?, precio = ? WHERE id_ale = ?")
	_, err = stmt.Exec(nombre, descripcion, alerta, notificacion, precio, id)
	if err == nil {
		return 1, "Alerta actualizada correctamente"
	} else {
		return 2, "La alerta no pudo ser actualizada"
	}
}
func BorrarAlerta(db *sql.DB, token string, id int) Response {

	resp := Response{}
	if found, _ := SuperAdmin(token); found {
		del := 1
		stmt, err := db.Prepare("UPDATE alertas SET eliminado = ? WHERE id_ale = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(del, id)
		ErrorCheck(e)
		if e == nil {
			resp.Tipo = "success"
			resp.Reload = 1
			resp.Page = "crearAlerta"
			resp.Titulo = "Alerta eliminada"
			resp.Texto = "Alerta eliminada correctamente"
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar alerta"
			resp.Texto = "La alerta no pudo ser eliminada"
		}
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar alerta"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}
func InsertCotizacion(db *sql.DB, uf string, id_emp int) (uint8, string) {

	stmt, err := db.Prepare("INSERT INTO cotizaciones (precio_uf, fecha, id_usr, id_emp) VALUES (?,NOW(),1,?)")
	defer stmt.Close()
	_, err = stmt.Exec(uf, id_emp)
	if err == nil {
		return 1, "Cotizacion ingresada correctamente"
	} else {
		return 2, "La cotizacion no pudo ser ingresada"
	}
}
func UpdateCotizacion(db *sql.DB, id int, uf string, id_emp int) (uint8, string) {

	stmt, err := db.Prepare("UPDATE cotizaciones SET precio_uf = ?, id_emp = ? WHERE id_cot = ?")
	_, err = stmt.Exec(uf, id_emp, id)
	if err == nil {
		return 1, "Cotizacion actualizada correctamente"
	} else {
		return 2, "La cotizacion no pudo ser actualizada"
	}
}
func InsertDetalleCotizacion(db *sql.DB, id int, descripcion string, precio string, id_pro int, id_ale int) (uint8, string) {

	stmt, err := db.Prepare("INSERT INTO cotizacion_detalle (id_cot, descripcion, precio, id_pro, id_ale) VALUES (?,?,?,?,?)")
	defer stmt.Close()
	_, err = stmt.Exec(id, descripcion, precio, id_pro, id_ale)
	if err == nil {
		RevisarCotizacion(db, id)
		return 1, "Item cotización ingresado correctamente"
	} else {
		return 2, "El item cotización no pudo ser ingresada"
	}
}
func UpdateDetalleCotizacion(db *sql.DB, id int, descripcion string, precio string, id_pro int, id_ale int) (uint8, string) {

	stmt, err := db.Prepare("UPDATE cotizacion_detalle SET descripcion = ?, precio = ? WHERE id_cot = ? AND id_pro = ? AND id_ale = ?")
	_, err = stmt.Exec(descripcion, precio, id, id_pro, id_ale)
	if err == nil {
		RevisarCotizacion(db, id)
		return 1, "Item cotización actualizada correctamente"
	} else {
		return 2, "El item cotización no pudo ser actualizada"
	}
}
func BorrarDetalleCotizacion(db *sql.DB, token string, id int, id_pro int, id_ale int) Response {

	resp := Response{}
	if found, _ := SuperAdmin(token); found {

		delForm, err := db.Prepare("DELETE FROM cotizacion_detalle WHERE id_cot = ? AND id_pro = ? AND  id_ale = ?")
		ErrorCheck(err)
		_, e := delForm.Exec(id_pro, id_ale)
		defer db.Close()

		ErrorCheck(e)
		if e == nil {
			resp.Tipo = "success"
			resp.Reload = 1
			resp.Page = fmt.Sprintf("confCotizacion?id_cot=%v", id)
			resp.Titulo = "Item cotizacion eliminada"
			resp.Texto = "Cotizacion eliminada correctamente"
		} else {
			resp.Tipo = "error"
			resp.Titulo = "Error al eliminar item cotizacion"
			resp.Texto = "El item cotizacion no pudo ser eliminada"
		}
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar item cotizacion"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}
func RevisarCotizacion(db *sql.DB, id int) {

	res, err := db.Query("SELECT precio FROM cotizacion_detalle WHERE id_cot = ?", id)
	defer res.Close()
	if err != nil {
		ErrorCheck(err)
	}

	var total float32
	var change bool = true

	for res.Next() {
		var precio float32
		err := res.Scan(&precio)
		if err != nil {
			ErrorCheck(err)
		}
		if precio == 0 {
			change = false
		}
		total = total + precio
	}
	if !change {
		total = 0
	}

	stmt, err := db.Prepare("UPDATE cotizaciones SET uf = ? WHERE id_cot = ?")
	ErrorCheck(err)
	_, e := stmt.Exec(total, id)
	ErrorCheck(e)
}
func GetPais(db *sql.DB, nombre string) (int64, bool) {

	if nombre != "" {
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
				} else {
					return 0, false
				}
			} else {
				return 0, false
			}
		}
	} else {
		return 0, false
	}
}
func GetRegion(db *sql.DB, nombre string, id_pai int64) (int64, bool) {

	if nombre != "" && id_pai > 0 {
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
				} else {
					return 0, false
				}
			} else {
				return 0, false
			}
		}
	} else {
		return 0, false
	}
}
func GetCiudad(db *sql.DB, nombre string, id_pai int64, id_reg int64) (int64, bool) {

	if nombre != "" && id_pai > 0 && id_reg > 0 {
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
				} else {
					return 0, false
				}
			} else {
				return 0, false
			}
		}
	} else {
		return 0, false
	}
}
func GetComuna(db *sql.DB, nombre string, id_pai int64, id_reg int64, id_ciu int64) (int64, bool) {

	if nombre != "" && id_pai > 0 && id_reg > 0 && id_ciu > 0 {
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
				} else {
					return 0, false
				}
			} else {
				return 0, false
			}
		}
	} else {
		return 0, false
	}
}

// DAEMON //
func (h *MyHandler) StartDaemon() {

	h.Conf.Tiempo = 2 * time.Second
	for _, x := range pass.ListaPropAlerts {
		DaemonAlertas(x.Id_pro, x.Id_emp)
	}
	pass.ListaPropAlerts = nil
	if len(h.DeleteFiles) > 0 {
		for i, x := range h.DeleteFiles {
			err := os.Remove(x)
			if err == nil {
				h.DeleteFiles = RemoveFiles(h.DeleteFiles, i)
				break
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
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(key)
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}
func GetRandom(min int, max int, id int) int {
	rand.Seed(time.Now().UnixNano() * int64(id))
	return rand.Intn(max-min) + min
}
func randSeq(n int) []byte {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[GetRandom(0, len(letters), i)]
	}
	return b
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
func ParamInt(data []byte) int {
	var x int
	for _, c := range data {
		x = x*10 + int(c-'0')
	}
	return x
}
func DatosCotizacion(id int) Cotizacion {

	cotizacion := Cotizacion{}
	cotizacion.IdCot = id

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.precio_uf, t3.id_ale, t3.nombre as nombreale, t3.descripcion, t2.precio, t4.nombre as nombreprop, t1.fecha, t5.nombre as nombreemp, t2.id_pro FROM cotizaciones t1, cotizacion_detalle t2, alertas t3, propiedades t4, empresa t5 WHERE t1.id_cot = ? AND t1.id_cot=t2.id_cot AND t2.id_ale=t3.id_ale AND t2.id_pro=t4.id_pro AND t1.id_emp=t5.id_emp", id)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	var id_ale int
	var id_pro int
	var nombreale string
	var descripcion string
	var precio float32
	var nombreprop string
	var precio_uf float32

	var fecha string
	var nombreemp string

	for res.Next() {

		err := res.Scan(&precio_uf, &id_ale, &nombreale, &descripcion, &precio, &nombreprop, &fecha, &nombreemp, &id_pro)
		ErrorCheck(err)

		cotizacion.Op = 1
		cotizacion.Uf = precio_uf
		cotizacion.Fecha = fecha
		cotizacion.NombreEmp = nombreemp

		cotizacion.Lista = append(cotizacion.Lista, ListaCotizacion{Propiedad: nombreprop, Descripcion: descripcion, Precio: precio, NombreAle: nombreale, IdAle: id_ale, IdPro: id_pro})
		cotizacion.TotalUf = cotizacion.TotalUf + precio

	}

	if cotizacion.Fecha != "" {
		cotizacion.Fecha = FormatDateString(cotizacion.Fecha)
	}

	cotizacion.Subtotal = cotizacion.TotalUf * cotizacion.Uf
	cotizacion.Iva = cotizacion.Subtotal * 0.19
	cotizacion.Total = cotizacion.Subtotal + cotizacion.Iva

	return cotizacion
}
func GetAlertasCotizaciones(id int) ([]Lista, int, string) {

	lista := []Lista{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res2, err := db.Query("SELECT t1.id_emp, t2.nombre FROM cotizaciones t1, empresa t2 WHERE t1.id_cot = ? AND t1.id_emp=t2.id_emp", id)
	defer res2.Close()
	if err != nil {
		fmt.Println(err)
	}

	var id_emp int
	var nombreemp string

	if res2.Next() {
		err := res2.Scan(&id_emp, &nombreemp)
		ErrorCheck(err)
	}

	res, err := db.Query("SELECT t3.id_ale, t3.nombre as nombreale, t4.nombre as nombreprop, t1.fecha, t2.id_pro, t5.id_emp, t5.nombre as nombreemp FROM cotizaciones t1, cotizacion_detalle t2, alertas t3, propiedades t4, empresa t5 WHERE t1.id_cot = ? AND t1.id_cot=t2.id_cot AND t2.id_ale=t3.id_ale AND t2.id_pro=t4.id_pro AND t1.id_emp=t5.id_emp", id)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	var id_ale int
	var id_pro int
	var nombreale string
	var nombreprop string
	var fecha string

	for res.Next() {

		err := res.Scan(&id_ale, &nombreale, &nombreprop, &fecha, &id_pro, &id_emp, &nombreemp)
		ErrorCheck(err)

		lista = append(lista, Lista{Id: id_ale, Nombre: fmt.Sprintf("%v %v", nombreprop, nombreale), Aux: id_pro})
	}

	return lista, id_emp, nombreemp
}
func GetDetalleCotizacion(id_cot int, id_ale int, id_pro int) (string, float64) {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT descripcion, precio FROM cotizacion_detalle WHERE id_cot=? AND id_ale=? AND id_pro=?", id_cot, id_ale, id_pro)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	var descripcion string
	var precio float64

	for res.Next() {

		err := res.Scan(&descripcion, &precio)
		ErrorCheck(err)
	}

	return descripcion, precio
}
func GetCotizacion(id_cot int) (int, int) {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT id_emp, precio_uf FROM cotizaciones WHERE id_cot=?", id_cot)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	var id_emp int
	var precio_uf int

	for res.Next() {

		err := res.Scan(&id_emp, &precio_uf)
		ErrorCheck(err)
	}

	return id_emp, precio_uf
}
func FormatDateString(fecha string) string {

	date0 := strings.Split(fecha, " ")
	date1 := strings.Split(date0[0], "-")
	Day, _ := strconv.Atoi(date1[2])
	return fmt.Sprintf("%v de %v de %v", Day, GetMonthString(date1[1]), date1[0])
}
func GetMonthString(m string) string {
	switch m {
	case "01":
		return "Enero"
	case "02":
		return "Febrero"
	case "03":
		return "Marzo"
	case "04":
		return "Abril"
	case "05":
		return "Mayo"
	case "06":
		return "Junio"
	case "07":
		return "Julio"
	case "08":
		return "Agosto"
	case "09":
		return "Septiembre"
	case "10":
		return "Octubre"
	case "11":
		return "Noviembre"
	default:
		return "Diciembre"
	}
}
func CreateQr(urlqr string, key int, prefix string) (bool, string) {

	q, err := qrcode.New(urlqr, qrcode.Medium)
	if err != nil {
		return false, ""
	}

	name := fmt.Sprintf("./tmp/%v%vqr.png", prefix, key)

	q.DisableBorder = true
	q.BackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	q.ForegroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	err = q.WriteFile(128, name)
	if err != nil {
		return false, ""
	}
	return true, name
}
func SeparadordeMiles(num int) string {
	numstr := fmt.Sprintf("%v", num)
	res := make([]byte, 0)

	for i, _ := range numstr {
		if i%3 == 0 && i > 0 {
			res = append(res, 46)
		}
		res = append(res, numstr[len(numstr)-i-1])
	}
	res = append(res, 32)
	res = append(res, 36)
	return string(Reverse(res))
}
func Reverse(numbers []uint8) []uint8 {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
func UploadFile(path string, files []*multipart.FileHeader, recortar bool, extends []string, filename string) (bool, []string) {

	list := make([]string, len(files))
	if len(files) == 0 {
		return false, list
	}

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return false, list
		}
		defer file.Close()
		if filename == "" {
			filename = fileHeader.Filename
		}

		filename, errn := FileExist(path, filename)
		if errn {
			return false, list
		}

		if Extension(filename, extends) {
			out, err := os.Create(fmt.Sprintf("%v/_%v", path, filename))
			if err != nil {
				return false, list
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				return false, list
			}
			if recortar {
				if RecortarImagenF(path, filename) {
					pass.DeleteFiles = append(pass.DeleteFiles, fmt.Sprintf("%v/_%v", path, filename))
				}
			}
			list[i] = filename
		} else {
			return false, list
		}
	}
	return true, list
}
func FileExist(path string, file string) (string, bool) {

	if _, serr := os.Stat(path); serr != nil {
		err := os.MkdirAll(path, os.ModeDir)
		if err != nil {
			return "", true
		}
	}

	is, name, ext := NombreExtension(file)
	name1 := name
	count := 1
	if is {
		for {
			pathfile := fmt.Sprintf("%s/%s.%s", path, name, ext)
			if _, err := os.Stat(pathfile); err == nil {
				if count < 10 {
					name = fmt.Sprintf("%s_00%v", name1, count)
				} else if count < 100 {
					name = fmt.Sprintf("%s_0%v", name1, count)
				} else {
					name = fmt.Sprintf("%s_%v", name1, count)
				}
				count++
			} else {
				return fmt.Sprintf("%s.%s", name, ext), false
			}
		}
	} else {
		return "", true
	}
}
func NombreExtension(filename string) (bool, string, string) {

	upper := strings.ToUpper(filename)
	bupper := []byte(upper)
	point := -1
	for i, x := range bupper {
		if x == 46 {
			point = i
		}
	}
	if point > 0 {
		return true, strings.ToLower(string(bupper[:point])), strings.ToLower(string(bupper[point+1:]))
	}
	return false, "", ""
}
func Extension(filename string, extends []string) bool {

	upper := strings.ToUpper(filename)
	bupper := []byte(upper)
	point := -1
	for i, x := range bupper {
		if x == 46 {
			point = i
		}
	}
	if point > 0 {
		for _, x := range extends {
			if x == string(bupper[point+1:]) {
				return true
			}
		}
	}
	return false
}
func SavePhotoDb(db *sql.DB, nombre string, tipo int, id_pro int, id_emp int) bool {

	stmt, err := db.Prepare("INSERT INTO propiedades_imagenes (nombre, tipo, fecha, id_pro, id_emp) VALUES (?,?,Now(),?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(nombre, tipo, id_pro, id_emp)
	ErrorCheck(err)
	if err == nil {
		return true
	} else {
		return false
	}
}
func SaveFileDb(db *sql.DB, nombre string, fojas string, numero string, fecha string, ano string, tipo int, valor_arriendo string, renovacion_auto string, tipo_de_plano string, id_pro int, id_emp int) bool {

	stmt, err := db.Prepare("INSERT INTO propiedad_archivos (nombre, fojas, numero, ano, tipo, fecha, valor_arriendo, renovacion_auto, tipo_de_plano, fecha_insert, id_pro, id_emp) VALUES (?,?,?,?,?,?,?,?,?,Now(),?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(nombre, fojas, numero, ano, tipo, fecha, valor_arriendo, renovacion_auto, tipo_de_plano, id_pro, id_emp)
	ErrorCheck(err)
	if err == nil {
		return true
	} else {
		return false
	}
}
func UpdateFileDb(db *sql.DB, fojas string, numero string, fecha string, ano string, valor_arriendo string, renovacion_auto string, tipo_de_plano string, id_arc int, id_pro int, id_emp int) bool {

	stmt, err := db.Prepare("UPDATE propiedad_archivos SET fojas = ?, numero = ?, ano = ?, fecha = ?, valor_arriendo = ?, renovacion_auto = ?, tipo_de_plano = ? WHERE id_arc = ? AND id_pro = ? AND id_emp = ?")
	_, err = stmt.Exec(fojas, numero, ano, fecha, valor_arriendo, renovacion_auto, tipo_de_plano, id_arc, id_pro, id_emp)
	if err == nil {
		return true
	} else {
		return false
	}
}
func SaveFileDb2(db *sql.DB, nombre1 string, nombre2 string, tipo int, indicar_acoge string, fecha string, id_rec int, id_pro int, id_emp int) bool {

	stmt, err := db.Prepare("INSERT INTO permiso_edificacion_archivos (nombre, nombre2, tipo, indicar_acoge, fecha, fecha_insert, id_rec, id_pro, id_emp) VALUES (?,?,?,?,?,Now(),?,?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(nombre1, nombre2, tipo, indicar_acoge, fecha, id_rec, id_pro, id_emp)
	ErrorCheck(err)
	if err == nil {
		return true
	} else {
		return false
	}
}
func RecortarImagenF(path string, filename string) bool {

	archivoImagen, err := os.Open(fmt.Sprintf("%v/_%v", path, filename))
	if err != nil {
		return false
	}
	defer archivoImagen.Close()

	imgOriginal, _, err := image.Decode(archivoImagen)
	if err != nil {
		return false
	}

	bounds := imgOriginal.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	x, y := 200, 300
	ancho, alto := 400, int(400*height/width)
	imgRecortada := recortarImagen(imgOriginal, x, y, ancho, alto)

	archivoNuevo, err := os.Create(fmt.Sprintf("%v/%v", path, filename))
	if err != nil {
		return false
	}
	defer archivoNuevo.Close()

	err = jpeg.Encode(archivoNuevo, imgRecortada, &jpeg.Options{Quality: 75})
	if err != nil {
		return false
	}
	return true
}
func recortarImagen(img image.Image, x, y, ancho, alto int) image.Image {

	resizedImg := image.NewRGBA(image.Rect(0, 0, ancho, alto))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return resizedImg
}

// ALERTAS //

func DaemonAlertas(id int, id_emp int) {

	propiedad, found := GetPropiedad(id_emp, id, true)
	if !found {
		fmt.Println(propiedad)
	}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	list_ale, found := GetPropiedadAlerta(db, id_emp, id)
	if !found {
		fmt.Println(list_ale)
	}

	alertas, found := GetObjetoAlertas(db)
	if !found {
		fmt.Println(alertas)
	}

	var alertabool bool
	for _, x := range alertas {
		alertabool = true
		for _, y := range x.Reglas {
			if y.Tipo == 1 {
				value, err := ObtenerValorCampo(propiedad, y.Campo)
				if err == nil {
					if value == y.Valor && alertabool {
						alertabool = false
					}
				} else {
					fmt.Println("Error al obtener el valor")
				}
			}
		}
		if !alertabool {
			if !SearchAlert(list_ale, x.Id_ale) {
				InsertAlert(db, id, x.Id_ale)
				fmt.Printf("INSERTAR ALERTA ID_PRO(%v) - ID_ALE(%v)\n", id, x.Id_ale)
			}
		} else {
			if SearchAlert(list_ale, x.Id_ale) {
				DeleteAlert(db, id, x.Id_ale)
				fmt.Printf("DELETE ALERTA ID_PRO(%v) - ID_ALE(%v)\n", id, x.Id_ale)
			}
		}
	}
}
func SearchAlert(list []int, val int) bool {
	for _, x := range list {
		if x == val {
			return true
		}
	}
	return false
}
func ObtenerValorCampo(obj interface{}, campo string) (string, error) {

	refObj := reflect.ValueOf(obj)    // Obtiene el reflejo del objeto
	if refObj.Kind() == reflect.Ptr { // Si el objeto es un puntero, obtenemos el valor apuntado
		refObj = refObj.Elem()
	}
	refCampo := refObj.FieldByName(campo) // Obtiene el reflejo del campo por nombre
	if !refCampo.IsValid() {
		return "", fmt.Errorf("Campo %s no encontrado", campo)
	}
	if reflect.TypeOf(refCampo.Interface()).Kind() == reflect.Int {
		return strconv.Itoa(int(refCampo.Int())), nil
	} else {
		return refCampo.String(), nil // Retorna el valor del campo
	}
}
func GetObjetoAlertas(db *sql.DB) ([]Alertas, bool) {

	alertas := make([]Alertas, 0)
	alerta := Alertas{}

	cn := 0
	res, err := db.Query("SELECT id_ale, alerta, notificacion FROM alertas WHERE eliminado = ?", cn)
	defer res.Close()
	if err != nil {
		return alertas, false
	}
	for res.Next() {
		err := res.Scan(&alerta.Id_ale, &alerta.Alerta, &alerta.Notificacion)
		if err != nil {
			return alertas, false
		}
		aux, found := GetObjetoReglas(db, alerta.Id_ale)
		if found {
			alerta.Reglas = aux
		} else {
			return alertas, false
		}
		alertas = append(alertas, alerta)
	}
	return alertas, true
}
func GetObjetoReglas(db *sql.DB, id int) ([]Regla, bool) {

	reglas := make([]Regla, 0)
	regla := Regla{}

	cn := 0
	res, err := db.Query("SELECT tipo, pagina, campo, valor FROM alerta_regla WHERE id_ale = ? AND eliminado = ?", id, cn)
	defer res.Close()
	if err != nil {
		return reglas, false
	}
	for res.Next() {
		err := res.Scan(&regla.Tipo, &regla.Pagina, &regla.Campo, &regla.Valor)
		if err != nil {
			return reglas, false
		}
		reglas = append(reglas, regla)
	}
	return reglas, true
}
func GetPropiedadAlerta(db *sql.DB, id_emp int, id_pro int) ([]int, bool) {

	resp := []int{}

	cn := 0
	res, err := db.Query("SELECT t2.id_ale FROM propiedades t1, propiedad_alerta t2 WHERE t1.id_pro = ? AND t1.id_emp = ? AND t1.eliminado = ? AND t1.id_pro=t2.id_pro", id_pro, id_emp, cn)
	defer res.Close()
	if err != nil {
		return resp, false
	}
	for res.Next() {
		var id_ale int
		err := res.Scan(&id_ale)
		if err != nil {
			return resp, false
		}
		resp = append(resp, id_ale)
	}
	return resp, true
}
func InsertAlert(db *sql.DB, id_pro int, id_ale int) {

	stmt, err := db.Prepare("INSERT INTO propiedad_alerta (id_pro, id_ale) VALUES (?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	stmt.Exec(id_pro, id_ale)
}
func DeleteAlert(db *sql.DB, id_pro int, id_ale int) {

	delForm, err := db.Prepare("DELETE FROM propiedad_alerta WHERE id_pro=? AND id_ale=?")
	ErrorCheck(err)
	delForm.Exec(id_pro, id_ale)
	defer db.Close()
}
func RemoveFiles(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// ALERTAS //

// UF //
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
		start := time.Date(ano, GetMonth(mes-1), dia, 0, 0, 0, 0, time.UTC)
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
	} else {
		return 0, false
	}
}
func CreateMapImage(lat string, lng string, width int, height int) (bool, string) {

	nom := fmt.Sprintf("mapa_%v_%v.png", lat, lng)
	_, exist := FileExist("files/maps", nom)
	if !exist {
		url := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%v,%v&zoom=13&scale=2&size=%vx%v&maptype=roadmap&key=%v", lat, lng, width, height, pass.Passwords.Gmapkey)

		fmt.Println(url)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(url)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		err := fasthttp.Do(req, resp)
		if err != nil {
			fmt.Printf("Client get failed: %s\n", err)
			return false, ""
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
			return false, ""
		}

		body := resp.Body()

		// Decodifica la respuesta en una imagen
		img, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			fmt.Println("Error al decodificar la imagen:", err)
			return false, ""
		}

		// Guarda la imagen en un archivo PNG
		file, err := os.Create(fmt.Sprintf("files/maps/%v", nom))
		if err != nil {
			fmt.Println("Error al crear el archivo:", err)
			return false, ""
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			fmt.Println("Error al guardar la imagen:", err)
			return false, ""
		}
	}
	return true, fmt.Sprintf("files/maps/%v", nom)
}
func UpdateUF(valor int) {

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

// UF //

func SendEmail(to string, subject string, body string) bool {

	//from := "redigocl@gmail.com"
	from := "valleencantado.cl@gmail.com"
	sub := fmt.Sprintf("From:%v\nTo:%v\nSubject:%v\n", from, to, subject)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass.Passwords.PassEmail, "smtp.gmail.com"), from, []string{to}, []byte(sub+mime+body))
	if err != nil {
		return false
	}
	return true
}
