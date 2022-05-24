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
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
	Titulo             string  `json:"Titulo"`
	SubTitulo          string  `json:"SubTitulo"`
	SubTitulo2         string  `json:"SubTitulo2"`
	FormId             int     `json:"FormId"`
	FormIdRec          int     `json:"FormIdRec"`
	FormAccion         string  `json:"FormAccion"`
	FormNombre         string  `json:"FormNombre"`
	Direccion          string  `json:"Direccion"`
	Nombre             string  `json:"Nombre"`
	Lat                float64 `json:"Lat"`
	Lng                float64 `json:"Lng"`
	Numero             string  `json:"Numero"`
	Comuna             string  `json:"Comuna"`
	Ciudad             string  `json:"Ciudad"`
	Region             string  `json:"Region"`
	Pais               string  `json:"Pais"`
	FormDescripcion    string  `json:"FormDescripcion"`
	FormPrecio         float64 `json:"FormPrecio"`
	TituloLista        string  `json:"TituloLista"`
	PageMod            string  `json:"PageMod"`
	DelAccion          string  `json:"DelAccion"`
	DelObj             string  `json:"DelObj"`
	Lista              []Lista `json:"FormDescripcion"`
	Dominio            int     `json:"Dominio"`
	Dominio2           int     `json:"Dominio2"`
	AtencionPublico    int     `json:"AtencionPublico"`
	Copropiedad        int     `json:"Copropiedad"`
	Destino            int     `json:"Destino"`
	Detalle            int     `json:"Detalle"`
	SupTerreno         string  `json:"SupTerreno"`
	SupEdificada       string  `json:"SupEdificada"`
	SupEdificadaSN     string  `json:"SupEdificadaSN"`
	SupEdificadaBN     string  `json:"SupEdificadaBN"`
	CantPisos          string  `json:"CantPisos"`
	PermEdificacion    int     `json:"PermEdificacion"`
	RecepcionFinal     int     `json:"RecepcionFinal"`
	EspecificarPermiso string  `json:"EspecificarPermiso"`
	Npermiso           string  `json:"Npermiso"`
	FechaPermiso       string  `json:"FechaPermiso"`
	PermDocumento      string  `json:"PermDocumento"`
	Recepcion          int     `json:"Recepcion"`

	ElectricoTe1        int `json:"ElectricoTe1"`
	DotacionAp          int `json:"DotacionAp"`
	DotacionAlcance     int `json:"DotacionAlcance"`
	InstalacionAscensor int `json:"InstalacionAscensor"`
	Te1Ascensor         int `json:"Te1Ascensor"`
	CertificadoAscensor int `json:"CertificadoAscensor"`
	Clima               int `json:"Clima"`
	SeguridadIncendio   int `json:"SeguridadIncendio"`

	TasacionValorComercial  string `json:"TasacionValorComercial"`
	AnoTasacion             string `json:"AnoTasacion"`
	ContratoArriendo        int    `json:"ContratoArriendo"`
	ValorArriendo           string `json:"ValorArriendo"`
	VencimientoArriendo     string `json:"VencimientoArriendo"`
	RenovacionAutomatica    int    `json:"RenovacionAutomatica"`
	ContratoSubArriendo     int    `json:"ContratoSubArriendo"`
	ValorSubArriendo        string `json:"ValorSubArriendo"`
	VencimientoSubArriendo  string `json:"VencimientoSubArriendo"`
	RenovacionAutomaticaSub int    `json:"RenovacionAutomaticaSub"`

	DomNomPropietario string `json:"DomNomPropietario"`
	Gp                int    `json:"RenovacionAutomaticaSub"`
	PlazosArchivos    int    `json:"RenovacionAutomaticaSub"`

	FiscalSerie        int    `json:"FiscalSerie"`
	FiscalDestino      int    `json:"FiscalDestino"`
	RolManzana         string `json:"RolManzana"`
	RolPredio          string `json:"RolPredio"`
	FiscalExento       int    `json:"FiscalExento"`
	AvaluoFiscal       string `json:"AvaluoFiscal"`
	ContribucionFiscal string `json:"ContribucionFiscal"`

	ValorTerreno              string `json:"ValorTerreno"`
	ValorEdificacion          string `json:"ValorEdificacion"`
	ValorObrasComplementarias string `json:"ValorObrasComplementarias"`
	ValorTotal                string `json:"ValorTotal"`

	CertInfoPrevias        int    `json:"CertInfoPrevias"`
	TipoInstrumento        int    `json:"TipoInstrumento"`
	DetalleTipoInstrumento string `json:"DetalleTipoInstrumento"`
	NormativoDestino       int    `json:"NormativoDestino"`
	ZonaNormativa          string `json:"ZonaNormativa"`
	UsosPermitidos         string `json:"UsosPermitidos"`
	UsosProhibidos         string `json:"UsosProhibidos"`
	Densidad               string `json:"Densidad"`
	CoefConstructibilidad  string `json:"CoefConstructibilidad"`
	CoefOcupacionSuelo     string `json:"CoefOcupacionSuelo"`

	NextPage int `json:"NextPage"`

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
	Titulo string `json:"Titulo"`
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
}
type Data struct {
	Nombre          string  `json:"Nombre"`
	Direccion       string  `json:"Direccion"`
	Lat             float64 `json:"Lat"`
	Lng             float64 `json:"Lng"`
	Dominio         int     `json:"Dominio"`
	Dominio2        int     `json:"Dominio2"`
	Precio          float64 `json:"Precio"`
	AtencionPublico int     `json:"AtencionPublico"`
	Copropiedad     int     `json:"Copropiedad"`
	Destino         int     `json:"Destino"`
	Detalle         int     `json:"Detalle"`
	Numero          string  `json:"Numero"`

	Comuna string `json:"Comuna"`
	Ciudad string `json:"Ciudad"`
	Region string `json:"Region"`
	Pais   string `json:"Pais"`

	SupTerreno     string `json:"SupTerreno"`
	SupEdificada   string `json:"SupEdificada"`
	SupEdificadaSN string `json:"SupEdificadaSN"`
	SupEdificadaBN string `json:"SupEdificadaBN"`
	CantPisos      string `json:"CantPisos"`

	PermEdificacion    int    `json:"PermEdificacion"`
	EspecificarPermiso string `json:"EspecificarPermiso"`
	Npermiso           string `json:"Npermiso"`
	FechaPermiso       string `json:"FechaPermiso"`
	PermDocumento      string `json:"PermDocumento"`
	RecepcionFinal     int    `json:"RecepcionFinal"`
	Recepcion          int    `json:"Recepcion"`

	ElectricoTe1        int `json:"ElectricoTe1"`
	DotacionAp          int `json:"DotacionAp"`
	DotacionAlcance     int `json:"DotacionAlcance"`
	InstalacionAscensor int `json:"InstalacionAscensor"`
	Te1Ascensor         int `json:"Te1Ascensor"`
	CertificadoAscensor int `json:"CertificadoAscensor"`
	Clima               int `json:"Clima"`
	SeguridadIncendio   int `json:"SeguridadIncendio"`

	TasacionValorComercial  string `json:"TasacionValorComercial"`
	AnoTasacion             string `json:"AnoTasacion"`
	ContratoArriendo        int    `json:"ContratoArriendo"`
	ValorArriendo           string `json:"ValorArriendo"`
	VencimientoArriendo     string `json:"VencimientoArriendo"`
	RenovacionAutomatica    int    `json:"RenovacionAutomatica"`
	ContratoSubArriendo     int    `json:"ContratoSubArriendo"`
	ValorSubArriendo        string `json:"ValorSubArriendo"`
	VencimientoSubArriendo  string `json:"VencimientoSubArriendo"`
	RenovacionAutomaticaSub int    `json:"RenovacionAutomaticaSub"`

	DomNomPropietario string `json:"DomNomPropietario"`
	Gp                int    `json:"Gp"`
	PlazosArchivos    int    `json:"PlazosArchivos"`

	FiscalSerie        int    `json:"FiscalSerie"`
	FiscalDestino      int    `json:"FiscalDestino"`
	RolManzana         string `json:"RolManzana"`
	RolPredio          string `json:"RolPredio"`
	FiscalExento       int    `json:"FiscalExento"`
	AvaluoFiscal       string `json:"AvaluoFiscal"`
	ContribucionFiscal string `json:"ContribucionFiscal"`

	ValorTerreno              string `json:"ValorTerreno"`
	ValorEdificacion          string `json:"ValorEdificacion"`
	ValorObrasComplementarias string `json:"ValorObrasComplementarias"`
	ValorTotal                string `json:"ValorTotal"`

	CertInfoPrevias        int    `json:"CertInfoPrevias"`
	TipoInstrumento        int    `json:"TipoInstrumento"`
	DetalleTipoInstrumento string `json:"DetalleTipoInstrumento"`
	NormativoDestino       int    `json:"NormativoDestino"`
	ZonaNormativa          string `json:"ZonaNormativa"`
	UsosPermitidos         string `json:"UsosPermitidos"`
	UsosProhibidos         string `json:"UsosProhibidos"`
	Densidad               string `json:"Densidad"`
	CoefConstructibilidad  string `json:"CoefConstructibilidad"`
	CoefOcupacionSuelo     string `json:"CoefOcupacionSuelo"`

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
	Paises            []Pais      `json:"Paises"`
	Regiones          []Region    `json:"Regiones"`
	Ciudades          []Ciudad    `json:"Ciudades"`
	Comunas           []Comuna    `json:"Comunas"`
	Propiedades       []Propiedad `json:"Propiedades"`
	Titulo            string      `json:"Titulo"`
	SubTitulo         string      `json:"SubTitulo"`
	SubTitulo2        string      `json:"SubTitulo2"`
	PaisesString      string      `json:"PaisesString"`
	RegionesString    string      `json:"RegionesString"`
	CiudadesString    string      `json:"CiudadesString"`
	ComunasString     string      `json:"ComunasString"`
	PropiedadesString string      `json:"PropiedadesString"`
	PaisesCount       int         `json:"PaisesCount"`
	RegionesCount     int         `json:"RegionesCount"`
	CiudadesCount     int         `json:"CiudadesCount"`
	ComunasCount      int         `json:"ComunasCount"`
	PropiedadesCount  int         `json:"PropiedadesCount"`
}
type Propiedad struct {
	Id_pro    int     `json:"Id_pro"`
	Nombre    string  `json:"Nombre"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Direccion string  `json:"Direccion"`
	Numero    int     `json:"Numero"`
	Id_com    int     `json:"Id_com"`
	Id_ciu    int     `json:"Id_ciu"`
	Id_reg    int     `json:"Id_reg"`
	Id_pai    int     `json:"Id_pai"`
}
type Pais struct {
	Id_pai int    `json:"Id_pai"`
	Nombre string `json:"Nombre"`
}
type Region struct {
	Id_reg int    `json:"Id_reg"`
	Nombre string `json:"Nombre"`
	Id_pai int    `json:"Id_pai"`
}
type Ciudad struct {
	Id_ciu int    `json:"Id_ciu"`
	Nombre string `json:"Nombre"`
	Id_reg int    `json:"Id_reg"`
	Id_pai int    `json:"Id_pai"`
}
type Comuna struct {
	Id_com int    `json:"Id_com"`
	Nombre string `json:"Nombre"`
	Id_ciu int    `json:"Id_ciu"`
	Id_reg int    `json:"Id_reg"`
	Id_pai int    `json:"Id_pai"`
}
type Rec struct {
	Code string `json:"Code"`
}

var (
	imgHandler fasthttp.RequestHandler
	cssHandler fasthttp.RequestHandler
	jsHandler  fasthttp.RequestHandler
	port       string
)

func main() {

	//SendEmail()
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
		r.GET("/recuperar/{name}", Recuperar)
		r.GET("/js/{name}", Js)
		r.GET("/img/{name}", Img)
		r.GET("/pages/{name}", Pages)
		r.POST("/login", Login)
		r.POST("/nueva", Nueva)
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
		dominio2 := string(ctx.FormValue("dominio2"))
		atencion_publico := string(ctx.FormValue("atencion_publico"))
		copropiedad := string(ctx.FormValue("copropiedad"))
		destino := string(ctx.FormValue("destino"))
		detalle_destino := string(ctx.FormValue("detalle_destino"))

		if id == 0 {
			resp = InsertPropiedad(db, token, nombre, lat, lng, comuna, ciudad, region, pais, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino)
		}
		if id > 0 {
			resp = UpdatePropiedad(db, token, id, nombre, lat, lng, comuna, ciudad, region, pais, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino)
		}
	case "guardar_propiedad2":

		sup_terreno := string(ctx.FormValue("sup_terreno"))
		sup_edificada := string(ctx.FormValue("sup_edificada"))
		sup_edificada_sn := string(ctx.FormValue("sup_edificada_sn"))
		sup_edificada_bn := string(ctx.FormValue("sup_edificada_bn"))
		cant_pisos := string(ctx.FormValue("cant_pisos"))

		if id > 0 {
			resp = UpdatePropiedad2(db, token, id, sup_terreno, sup_edificada, sup_edificada_sn, sup_edificada_bn, cant_pisos)
		}
	case "guardar_propiedad2A":

		id_rec := Read_uint32bytes(ctx.FormValue("id_rec"))

		permiso_edificacion := string(ctx.FormValue("permiso_edificacion"))
		especificar_permiso := string(ctx.FormValue("especificar_permiso"))
		num_permiso := string(ctx.FormValue("num_permiso"))
		fecha_permiso := string(ctx.FormValue("fecha_permiso"))
		recepcion := string(ctx.FormValue("recepcion"))
		recepcion_final := string(ctx.FormValue("recepcion_final"))

		if id_rec > 0 {
			resp = UpdatePropiedad2A(db, token, id, id_rec, permiso_edificacion, especificar_permiso, num_permiso, fecha_permiso, recepcion, recepcion_final)
			header, err := ctx.FormFile("documento")
			if err == nil && header != nil {
				UpdateFile(db, token, "perm_edificacion", "permiso_edificacion", "documento", "id_rec", id_rec, header)
			}
		} else {
			var documento string
			header, err := ctx.FormFile("documento")
			if err == nil && header != nil {
				documento = InsertFile(token, "perm_edificacion", header)
			}
			resp = InsertPropiedad2A(db, token, id, permiso_edificacion, especificar_permiso, num_permiso, fecha_permiso, recepcion, recepcion_final, documento)
		}
	case "guardar_propiedad3":

		electrico_te1 := string(ctx.FormValue("electrico_te1"))
		dotacion_ap := string(ctx.FormValue("dotacion_ap"))
		dotacion_alcance := string(ctx.FormValue("dotacion_alcance"))
		instalacion_ascensor := string(ctx.FormValue("instalacion_ascensor"))
		te1_ascensor := string(ctx.FormValue("te1_ascensor"))
		certificado_ascensor := string(ctx.FormValue("certificado_ascensor"))
		clima := string(ctx.FormValue("clima"))
		seguridad_incendio := string(ctx.FormValue("seguridad_incendio"))

		if id > 0 {
			resp = UpdatePropiedad3(db, token, id, electrico_te1, dotacion_ap, dotacion_alcance, instalacion_ascensor, te1_ascensor, certificado_ascensor, clima, seguridad_incendio)
			header1, err := ctx.FormFile("doc_electrico_te1")
			if err == nil && header1 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_electrico_te1", "id_pro", id, header1)
			}
			header2, err := ctx.FormFile("doc_dotacion_ap")
			if err == nil && header2 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_dotacion_ap", "id_pro", id, header2)
			}
			header3, err := ctx.FormFile("doc_dotacion_alcance")
			if err == nil && header3 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_dotacion_alcance", "id_pro", id, header3)
			}
			header4, err := ctx.FormFile("doc_instalacion_ascensor")
			if err == nil && header4 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_instalacion_ascensor", "id_pro", id, header4)
			}
			header5, err := ctx.FormFile("doc_te1_ascensor")
			if err == nil && header5 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_te1_ascensor", "id_pro", id, header5)
			}
			header6, err := ctx.FormFile("doc_certificado_ascensor")
			if err == nil && header6 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_certificado_ascensor", "id_pro", id, header6)
			}
			header7, err := ctx.FormFile("doc_clima")
			if err == nil && header7 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_clima", "id_pro", id, header7)
			}
			header8, err := ctx.FormFile("doc_seguridad_incendio")
			if err == nil && header8 != nil {
				UpdateFile(db, token, "propiedades", "situaciontecnica", "doc_seguridad_incendio", "id_pro", id, header8)
			}
		}
	case "guardar_propiedad4":

		tasacion_valor_comercial := string(ctx.FormValue("tasacion_valor_comercial"))
		ano_tasacion := string(ctx.FormValue("ano_tasacion"))
		contrato_arriendo := string(ctx.FormValue("contrato_arriendo"))
		valor_arriendo := string(ctx.FormValue("valor_arriendo"))
		vencimiento_arriendo := string(ctx.FormValue("vencimiento_arriendo"))
		renovacion_automatica := string(ctx.FormValue("renovacion_automatica"))

		contrato_subarriendo := string(ctx.FormValue("contrato_subarriendo"))
		valor_subarriendo := string(ctx.FormValue("valor_subarriendo"))
		vencimiento_subarriendo := string(ctx.FormValue("vencimiento_subarriendo"))
		renovacion_automaticasub := string(ctx.FormValue("renovacion_automaticasub"))

		if id > 0 {
			resp = UpdatePropiedad4(db, token, id, tasacion_valor_comercial, ano_tasacion, contrato_arriendo, valor_arriendo, vencimiento_arriendo, renovacion_automatica, contrato_subarriendo, valor_subarriendo, vencimiento_subarriendo, renovacion_automaticasub)
		}
	case "guardar_propiedad5":

		dom_nom_propietario := string(ctx.FormValue("dom_nom_propietario"))
		gp := string(ctx.FormValue("gp"))
		planos_archivados := string(ctx.FormValue("planos_archivados"))

		if id > 0 {
			resp = UpdatePropiedad5(db, token, id, dom_nom_propietario, gp, planos_archivados)
		}
	case "guardar_propiedad6":

		fiscal_serie := string(ctx.FormValue("fiscal_serie"))
		fiscal_destino := string(ctx.FormValue("fiscal_destino"))
		rol_manzana := string(ctx.FormValue("rol_manzana"))
		rol_predio := string(ctx.FormValue("rol_predio"))
		fiscal_exento := string(ctx.FormValue("fiscal_exento"))
		avaluo_fiscal := string(ctx.FormValue("avaluo_fiscal"))
		contribucion_fiscal := string(ctx.FormValue("contribucion_fiscal"))

		if id > 0 {
			resp = UpdatePropiedad6(db, token, id, fiscal_serie, fiscal_destino, rol_manzana, rol_predio, fiscal_exento, avaluo_fiscal, contribucion_fiscal)
		}
	case "guardar_propiedad7":

		valor_terreno := string(ctx.FormValue("valor_terreno"))
		valor_edificacion := string(ctx.FormValue("valor_edificacion"))
		valor_obras_complementarias := string(ctx.FormValue("valor_obras_complementarias"))
		valor_total := string(ctx.FormValue("valor_total"))

		if id > 0 {
			resp = UpdatePropiedad7(db, token, id, valor_terreno, valor_edificacion, valor_obras_complementarias, valor_total)
		}
	case "guardar_propiedad8":

		cert_info_previas := string(ctx.FormValue("cert_info_previas"))
		tipo_instrumento := string(ctx.FormValue("tipo_instrumento"))
		detalle_tipo_instrumento := string(ctx.FormValue("detalle_tipo_instrumento"))
		normativo_destino := string(ctx.FormValue("normativo_destino"))
		zona_normativa := string(ctx.FormValue("zona_normativa"))
		usos_permitidos := string(ctx.FormValue("usos_permitidos"))
		usos_prohibidos := string(ctx.FormValue("usos_prohibidos"))
		densidad := string(ctx.FormValue("densidad"))
		coef_constructibilidad := string(ctx.FormValue("coef_constructibilidad"))
		coef_ocupacion_suelo := string(ctx.FormValue("coef_ocupacion_suelo"))

		if id > 0 {
			resp = UpdatePropiedad8(db, token, id, cert_info_previas, tipo_instrumento, detalle_tipo_instrumento, normativo_destino, zona_normativa, usos_permitidos, usos_prohibidos, densidad, coef_constructibilidad, coef_ocupacion_suelo)
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
func FileExist(path string, file string) string {

	pathfile := fmt.Sprintf("%s/%s", path, file)
	if _, err := os.Stat(pathfile); err == nil {
		return fmt.Sprintf("%s_%s", randSeq(8), file)
	} else {
		return file
	}
}
func InsertFile(token string, folder string, header *multipart.FileHeader) string {

	path := fmt.Sprintf("./pdf/%v/%s", GetIdEmp(token), folder)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	file := FileExist(path, header.Filename)
	fasthttp.SaveMultipartFile(header, fmt.Sprintf("%s/%s", path, file))
	return file
}
func UpdateFile(db *sql.DB, token string, folder string, tabla string, campo string, key string, id int, header *multipart.FileHeader) {

	id_emp := GetIdEmp(token)
	path := fmt.Sprintf("./pdf/%v/%s", id_emp, folder)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	documento := FileExist(path, header.Filename)
	fasthttp.SaveMultipartFile(header, fmt.Sprintf("%s/%s", path, documento))

	sql := fmt.Sprintf("UPDATE %v SET %v = ? WHERE %v = ? AND id_emp = ?", tabla, campo, key)
	stmt, err := db.Prepare(sql)
	ErrorCheck(err)
	_, e := stmt.Exec(documento, id, id_emp)
	ErrorCheck(e)
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

			obj := GetTemplateConf("Crear Propiedad", "Datos Generales", "Completar los datos", "Lista de Propiedades", "guardar_propiedad1", fmt.Sprintf("/pages/%s", name), "borrar_empresa", "Empresa")
			lista, found := GetPropiedades(token)
			if found {
				obj.Lista = lista
			}

			if id > 0 {
				aux, found := GetPropiedad(token, id)
				if found {
					obj.Nombre = aux.Nombre
					obj.FormId = id
					obj.Dominio = aux.Dominio
					obj.Dominio2 = aux.Dominio2
					obj.AtencionPublico = aux.AtencionPublico
					obj.Copropiedad = aux.Copropiedad
					obj.Destino = aux.Destino
					obj.Detalle = aux.Detalle
					obj.Direccion = aux.Direccion
					obj.Lat = aux.Lat
					obj.Lng = aux.Lng
					obj.Numero = aux.Numero
					obj.Comuna = aux.Comuna
					obj.Ciudad = aux.Ciudad
					obj.Region = aux.Region
					obj.Pais = aux.Pais
				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad2":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Datos Municipales", "Completar los datos", "", "guardar_propiedad2", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad2(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.SupTerreno = aux.SupTerreno
					obj.SupEdificada = aux.SupEdificada
					obj.SupEdificadaSN = aux.SupEdificadaSN
					obj.SupEdificadaBN = aux.SupEdificadaBN
					obj.CantPisos = aux.CantPisos

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad2PermisoEdificacion":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			id_rec := Read_uint32bytes(ctx.QueryArgs().Peek("id_rec"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Permisos de Edificación", "Completar los datos", "Lista de Permisos de Edificación", "guardar_propiedad2A", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad2(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					lista, found := PermisosEdificacion(id)
					if found {
						obj.Lista = lista
					}

					if id_rec > 0 {
						aux2, found2 := GetPropiedad2A(token, id_rec)
						if found2 {

							obj.FormIdRec = id_rec
							obj.PermEdificacion = aux2.PermEdificacion
							obj.EspecificarPermiso = aux2.EspecificarPermiso
							obj.Npermiso = aux2.Npermiso
							obj.FechaPermiso = aux2.FechaPermiso
							obj.PermDocumento = aux2.PermDocumento
							obj.Recepcion = aux2.Recepcion
							obj.RecepcionFinal = aux2.RecepcionFinal

						}
					}

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad3":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Técnica", "Completar los datos", "", "guardar_propiedad3", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad3(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.ElectricoTe1 = aux.ElectricoTe1
					obj.DotacionAp = aux.DotacionAp
					obj.DotacionAlcance = aux.DotacionAlcance
					obj.InstalacionAscensor = aux.InstalacionAscensor
					obj.Te1Ascensor = aux.Te1Ascensor
					obj.CertificadoAscensor = aux.CertificadoAscensor
					obj.Clima = aux.Clima
					obj.SeguridadIncendio = aux.SeguridadIncendio

					if IsArrendado(id) {
						obj.NextPage = 4
					} else {
						obj.NextPage = 5
					}

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad4":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Comercial", "Completar los datos", "", "guardar_propiedad4", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad4(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.TasacionValorComercial = aux.TasacionValorComercial
					obj.AnoTasacion = aux.AnoTasacion
					obj.ContratoArriendo = aux.ContratoArriendo
					obj.ValorArriendo = aux.ValorArriendo
					obj.VencimientoArriendo = aux.VencimientoArriendo
					obj.RenovacionAutomatica = aux.RenovacionAutomatica
					obj.ContratoSubArriendo = aux.ContratoSubArriendo
					obj.ValorSubArriendo = aux.ValorSubArriendo
					obj.VencimientoSubArriendo = aux.VencimientoSubArriendo
					obj.RenovacionAutomaticaSub = aux.RenovacionAutomaticaSub

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad5":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Legal", "Completar los datos", "", "guardar_propiedad5", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad5(token, id)
				if found {

					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.DomNomPropietario = aux.DomNomPropietario
					obj.Gp = aux.Gp
					obj.PlazosArchivos = aux.PlazosArchivos

					if IsArrendado(id) {
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

		}
	case "crearPropiedad6":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Situación Fiscal", "Completar los datos", "", "guardar_propiedad6", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad6(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.FiscalSerie = aux.FiscalSerie
					obj.FiscalDestino = aux.FiscalDestino
					obj.RolManzana = aux.RolManzana
					obj.RolPredio = aux.RolPredio
					obj.FiscalExento = aux.FiscalExento
					obj.AvaluoFiscal = aux.AvaluoFiscal
					obj.ContribucionFiscal = aux.ContribucionFiscal

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad7":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Avalúo Comercial", "Completar los datos", "", "guardar_propiedad7", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad7(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.ValorTerreno = aux.ValorTerreno
					obj.ValorEdificacion = aux.ValorTerreno
					obj.ValorObrasComplementarias = aux.ValorTerreno
					obj.ValorTotal = aux.ValorTerreno

				}
			} else {
				obj.FormId = 0
			}

			err = t.Execute(ctx, obj)
			ErrorCheck(err)

		}
	case "crearPropiedad8":

		if found, _ := Permisos(token, 1); found {

			id := Read_uint32bytes(ctx.QueryArgs().Peek("id"))
			t, err := TemplatePage(fmt.Sprintf("html/%s.html", name))
			ErrorCheck(err)

			obj := GetTemplateConf("Crear Propiedad", "Normativo", "Completar los datos", "", "guardar_propiedad8", fmt.Sprintf("/pages/%s", name), "", "")

			if id > 0 {
				aux, found := GetPropiedad8(token, id)
				if found {
					obj.FormNombre = aux.Nombre
					obj.FormId = id

					obj.CertInfoPrevias = aux.CertInfoPrevias
					obj.TipoInstrumento = aux.TipoInstrumento
					obj.DetalleTipoInstrumento = aux.DetalleTipoInstrumento
					obj.NormativoDestino = aux.NormativoDestino
					obj.ZonaNormativa = aux.ZonaNormativa
					obj.UsosPermitidos = aux.UsosPermitidos
					obj.UsosProhibidos = aux.UsosProhibidos
					obj.Densidad = aux.Densidad
					obj.CoefConstructibilidad = aux.CoefConstructibilidad
					obj.CoefOcupacionSuelo = aux.CoefOcupacionSuelo

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
func Recuperar(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	name := ctx.UserValue("name")

	str, ok := name.(string)
	if ok {
		fmt.Printf("%T %v", str, str)

		t, err := TemplatePage("html/recuperar.html")
		ErrorCheck(err)
		var x Rec
		x.Code = str
		err = t.Execute(ctx, x)
		ErrorCheck(err)
	}
}
func Index(ctx *fasthttp.RequestCtx) {

	fmt.Println(GetUF())

	//SendEmail()
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
	res, err := db.Query("SELECT nombre, direccion, numero, lat, lng, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, id_com, id_ciu, id_reg, id_pai FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var nombre string
		var direccion string
		var numero string
		var lat float64
		var lng float64
		var dominio int
		var dominio2 int
		var atencion_publico int
		var copropiedad int
		var destino int
		var detalle_destino int
		var id_com int
		var id_ciu int
		var id_reg int
		var id_pai int
		err := res.Scan(&nombre, &direccion, &numero, &lat, &lng, &dominio, &dominio2, &atencion_publico, &copropiedad, &destino, &detalle_destino, &id_com, &id_ciu, &id_reg, &id_pai)
		if err != nil {
			log.Fatal(err)
		}
		data.Nombre = nombre
		data.Direccion = direccion
		data.Numero = numero
		data.Lat = lat
		data.Lng = lng
		data.Dominio = dominio
		data.Dominio2 = dominio2
		data.AtencionPublico = atencion_publico
		data.Copropiedad = copropiedad
		data.Destino = destino
		data.Detalle = detalle_destino
		data.Comuna = GetComunaNombre(db, id_com, id_ciu, id_reg, id_pai)
		data.Ciudad = GetCiudadNombre(db, id_ciu, id_reg, id_pai)
		data.Region = GetRegionNombre(db, id_reg, id_pai)
		data.Pais = GetPaisNombre(db, id_pai)
		return data, true

	} else {
		return data, false
	}
}
func GetCiudadNombre(db *sql.DB, id_ciu int, id_reg int, id_pai int) string {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT nombre FROM ciudades WHERE id_ciu = ? AND id_reg = ? AND id_pai = ?", id_ciu, id_reg, id_pai)
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
		return nombre
	} else {
		return ""
	}
}
func GetRegionNombre(db *sql.DB, id_reg int, id_pai int) string {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT nombre FROM regiones WHERE id_reg = ? AND id_pai = ?", id_reg, id_pai)
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
		return nombre
	} else {
		return ""
	}
}
func GetComunaNombre(db *sql.DB, id_com int, id_ciu int, id_reg int, id_pai int) string {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT nombre FROM comunas WHERE id_com = ? AND id_ciu = ? AND id_reg = ? AND id_pai = ?", id_com, id_ciu, id_reg, id_pai)
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
		return nombre
	} else {
		return ""
	}
}
func GetPaisNombre(db *sql.DB, id_pai int) string {

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT nombre FROM paises WHERE id_pai = ?", id_pai)
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
		return nombre
	} else {
		return ""
	}
}

func GetPropiedad2(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT sup_terreno, sup_edificada, sup_edificada_sn, sup_edificada_bn, cant_pisos FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var sup_terreno string
		var sup_edificada string
		var sup_edificada_sn string
		var sup_edificada_bn string
		var cant_pisos string

		err := res.Scan(&sup_terreno, &sup_edificada, &sup_edificada_sn, &sup_edificada_bn, &cant_pisos)
		if err != nil {
			log.Fatal(err)
		}

		data.SupTerreno = sup_terreno
		data.SupEdificada = sup_edificada
		data.SupEdificadaSN = sup_edificada_sn
		data.SupEdificadaBN = sup_edificada_bn
		data.CantPisos = cant_pisos

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad2A(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT tipo, especificar_tipo, numero, fecha, documento, recepcion, recepcion_total FROM permiso_edificacion WHERE id_rec = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var tipo int
		var especificar_tipo string
		var numero string
		var fecha string
		var documento string
		var recepcion int
		var recepcion_total int

		err := res.Scan(&tipo, &especificar_tipo, &numero, &fecha, &documento, &recepcion, &recepcion_total)
		if err != nil {
			log.Fatal(err)
		}

		data.PermEdificacion = tipo
		data.EspecificarPermiso = especificar_tipo
		data.Npermiso = numero
		data.FechaPermiso = fecha
		data.PermDocumento = fecha
		data.Recepcion = recepcion
		data.RecepcionFinal = recepcion_total

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad3(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT electrico_te1, dotacion_ap, dotacion_alcance, instalacion_ascensor, te1_ascensor, certificado_ascensor, clima, seguridad_incendio FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var electrico_te1 int
		var dotacion_ap int
		var dotacion_alcance int
		var instalacion_ascensor int
		var te1_ascensor int
		var certificado_ascensor int
		var clima int
		var seguridad_incendio int

		err := res.Scan(&electrico_te1, &dotacion_ap, &dotacion_alcance, &instalacion_ascensor, &te1_ascensor, &certificado_ascensor, &clima, &seguridad_incendio)
		if err != nil {
			log.Fatal(err)
		}

		data.ElectricoTe1 = electrico_te1
		data.DotacionAp = dotacion_ap
		data.DotacionAlcance = dotacion_alcance
		data.InstalacionAscensor = instalacion_ascensor
		data.Te1Ascensor = te1_ascensor
		data.CertificadoAscensor = certificado_ascensor
		data.Clima = clima
		data.SeguridadIncendio = seguridad_incendio

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad4(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT tasacion_valor_comercial, ano_tasacion, contrato_arriendo, valor_arriendo, vencimiento_arriendo, renovacion_automatica, contrato_subarriendo, valor_subarriendo, vencimiento_subarriendo, renovacion_automaticasub FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var tasacion_valor_comercial string
		var ano_tasacion string
		var contrato_arriendo int
		var valor_arriendo string
		var vencimiento_arriendo string
		var renovacion_automatica int
		var contrato_subarriendo int
		var valor_subarriendo string
		var vencimiento_subarriendo string
		var renovacion_automaticasub int

		err := res.Scan(&tasacion_valor_comercial, &ano_tasacion, &contrato_arriendo, &valor_arriendo, &vencimiento_arriendo, &renovacion_automatica, &contrato_subarriendo, &valor_subarriendo, &vencimiento_subarriendo, &renovacion_automaticasub)
		if err != nil {
			log.Fatal(err)
		}

		data.TasacionValorComercial = tasacion_valor_comercial
		if ano_tasacion == "" {
			data.AnoTasacion = "2020"
		} else {
			data.AnoTasacion = ano_tasacion
		}
		data.ContratoArriendo = contrato_arriendo
		data.ValorArriendo = valor_arriendo
		data.VencimientoArriendo = vencimiento_arriendo
		data.RenovacionAutomatica = renovacion_automatica
		data.ContratoSubArriendo = contrato_subarriendo
		data.ValorSubArriendo = valor_subarriendo
		data.VencimientoSubArriendo = vencimiento_subarriendo
		data.RenovacionAutomaticaSub = renovacion_automaticasub

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad5(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT domnompropietario, gp, plazosarchivos FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var domnompropietario string
		var gp int
		var plazosarchivos int

		err := res.Scan(&domnompropietario, &gp, &plazosarchivos)
		if err != nil {
			log.Fatal(err)
		}

		data.DomNomPropietario = domnompropietario
		data.Gp = gp
		data.PlazosArchivos = plazosarchivos

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad6(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT fiscal_serie, fiscal_destino, rol_manzana, rol_predio, fiscal_exento, avaluo_fiscal, contribucion_fiscal FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var fiscal_serie int
		var fiscal_destino int
		var rol_manzana string
		var rol_predio string
		var fiscal_exento int
		var avaluo_fiscal string
		var contribucion_fiscal string

		err := res.Scan(&fiscal_serie, &fiscal_destino, &rol_manzana, &rol_predio, &fiscal_exento, &avaluo_fiscal, &contribucion_fiscal)
		if err != nil {
			log.Fatal(err)
		}

		data.FiscalSerie = fiscal_serie
		data.FiscalDestino = fiscal_destino
		data.RolManzana = rol_manzana
		data.RolPredio = rol_predio
		data.FiscalExento = fiscal_exento
		data.AvaluoFiscal = avaluo_fiscal
		data.ContribucionFiscal = contribucion_fiscal

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad7(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT valor_terreno, valor_edificacion, valor_obras_complementarias, valor_total FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var valor_terreno string
		var valor_edificacion string
		var valor_obras_complementarias string
		var valor_total string

		err := res.Scan(&valor_terreno, &valor_edificacion, &valor_obras_complementarias, &valor_total)
		if err != nil {
			log.Fatal(err)
		}

		data.ValorTerreno = valor_terreno
		data.ValorEdificacion = valor_edificacion
		data.ValorObrasComplementarias = valor_obras_complementarias
		data.ValorTotal = valor_total

		return data, true

	} else {
		return data, false
	}
}
func GetPropiedad8(token string, id int) (Data, bool) {

	data := Data{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	cn := 0
	res, err := db.Query("SELECT cert_info_previas, tipo_instrumento, detalle_tipo_instrumento, normativo_destino, zona_normativa, usos_permitidos, usos_prohibidos, densidad, coef_constructibilidad, coef_ocupacion_suelo FROM propiedades WHERE id_pro = ? AND eliminado = ? AND id_emp = ?", id, cn, GetIdEmp(token))
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Next() {

		var cert_info_previas int
		var tipo_instrumento int
		var detalle_tipo_instrumento string
		var normativo_destino int
		var zona_normativa string
		var usos_permitidos string
		var usos_prohibidos string
		var densidad string
		var coef_constructibilidad string
		var coef_ocupacion_suelo string

		err := res.Scan(&cert_info_previas, &tipo_instrumento, &detalle_tipo_instrumento, &normativo_destino, &zona_normativa, &usos_permitidos, &usos_prohibidos, &densidad, &coef_constructibilidad, &coef_ocupacion_suelo)
		if err != nil {
			log.Fatal(err)
		}

		data.CertInfoPrevias = cert_info_previas
		data.TipoInstrumento = tipo_instrumento
		data.DetalleTipoInstrumento = detalle_tipo_instrumento
		data.NormativoDestino = normativo_destino
		data.ZonaNormativa = zona_normativa
		data.UsosPermitidos = usos_permitidos
		data.UsosProhibidos = usos_prohibidos
		data.Densidad = densidad
		data.CoefConstructibilidad = coef_constructibilidad
		data.CoefOcupacionSuelo = coef_ocupacion_suelo

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
		fmt.Println(err)
		return 0, false
	}
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

	return Localidades{Paises: paises, Regiones: regiones, Ciudades: ciudades, Comunas: comunas, Propiedades: propiedades}
}

func InsertPropiedad(db *sql.DB, token string, nombre string, lat string, lng string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, dominio2 string, atencion_publico string, copropiedad string, destino string, detalle_destino string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {

		if nombre != "" {
			id_pai, b1 := GetPais(db, pais)
			id_reg, b2 := GetRegion(db, region, id_pai)
			id_ciu, b3 := GetCiudad(db, ciudad, id_pai, id_reg)
			id_com, b4 := GetComuna(db, comuna, id_pai, id_reg, id_ciu)

			if b1 && b2 && b3 && b4 {
				stmt, err := db.Prepare("INSERT INTO propiedades (nombre, lat, lng, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, id_emp) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
				ErrorCheck(err)
				defer stmt.Close()
				r, err := stmt.Exec(nombre, lat, lng, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, id_emp)
				ErrorCheck(err)
				if err == nil {

					id, err := r.LastInsertId()
					if err == nil {
						resp.Op = 1
						resp.Reload = 1
						resp.Page = fmt.Sprintf("crearPropiedad2?id=%v", id)
						resp.Msg = "Propiedad ingresada correctamente"
					}

				} else {
					resp.Msg = "La Propiedad no pudo ser ingresada"
				}
			} else {
				resp.Msg = "Error al ingresar posicion"
			}
		} else {
			resp.Msg = "Debe ingresar nombre"
		}

	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad(db *sql.DB, token string, id int, nombre string, lat string, lng string, comuna string, ciudad string, region string, pais string, direccion string, numero string, dominio string, dominio2 string, atencion_publico string, copropiedad string, destino string, detalle_destino string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {

		id_pai, b1 := GetPais(db, pais)
		id_reg, b2 := GetRegion(db, region, id_pai)
		id_ciu, b3 := GetCiudad(db, ciudad, id_pai, id_reg)
		id_com, b4 := GetComuna(db, comuna, id_pai, id_reg, id_ciu)

		if b1 && b2 && b3 && b4 {

			stmt, err := db.Prepare("UPDATE propiedades SET nombre = ?, lat = ?, lng = ?, id_ciu = ?, id_com = ?, id_reg = ?, id_pai = ?, direccion = ?, numero = ?, dominio = ?, dominio2 = ?, atencion_publico = ?, copropiedad = ?, destino = ?, detalle_destino = ? WHERE id_pro = ? AND id_emp = ?")
			ErrorCheck(err)
			_, e := stmt.Exec(nombre, lat, lng, id_ciu, id_com, id_reg, id_pai, direccion, numero, dominio, dominio2, atencion_publico, copropiedad, destino, detalle_destino, id, id_emp)
			ErrorCheck(e)
			if e == nil {
				resp.Op = 1
				resp.Reload = 1
				resp.Page = fmt.Sprintf("crearPropiedad2?id=%v", id)
				resp.Msg = "Empresa actualizada correctamente"
			} else {
				resp.Msg = "La empresa no pudo ser actualizada"
			}

		} else {
			resp.Msg = "Error al ingresar posicion"
		}

	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad2(db *sql.DB, token string, id int, sup_terreno string, sup_edificada string, sup_edificada_sn string, sup_edificada_bn string, cant_pisos string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET sup_terreno = ?, sup_edificada = ?, sup_edificada_sn = ?, sup_edificada_bn = ?, cant_pisos = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(sup_terreno, sup_edificada, sup_edificada_sn, sup_edificada_bn, cant_pisos, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad2PermisoEdificacion?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad2A(db *sql.DB, token string, id int, id_rec int, tipo string, especificar_permiso string, num_permiso string, fecha_permiso string, recepcion string, recepcion_final string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE permiso_edificacion SET tipo = ?, especificar_tipo = ?, numero = ?, fecha = ?, recepcion = ?, recepcion_total = ? WHERE id_rec = ? AND id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(tipo, especificar_permiso, num_permiso, fecha_permiso, recepcion, recepcion_final, id_rec, id, id_emp)
		ErrorCheck(e)
		if e == nil {

			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad3?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func InsertPropiedad2A(db *sql.DB, token string, id int, tipo string, especificar_permiso string, num_permiso string, fecha_permiso string, recepcion string, recepcion_final string, documento string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {

		stmt, err := db.Prepare("INSERT INTO permiso_edificacion (tipo, especificar_tipo, numero, fecha, documento, recepcion, recepcion_total, id_pro, id_emp) VALUES (?,?,?,?,?,?,?,?,?)")
		ErrorCheck(err)
		defer stmt.Close()
		_, e := stmt.Exec(tipo, especificar_permiso, num_permiso, fecha_permiso, documento, recepcion, recepcion_final, id, id_emp)
		ErrorCheck(e)
		if e == nil {

			//id, err := r.LastInsertId()
			//if err == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad3?id=%v", id)
			resp.Msg = "Propiedad ingresada correctamente"
			//}

		} else {
			resp.Msg = "La Propiedad no pudo ser ingresada"
		}

	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad3(db *sql.DB, token string, id int, electrico_te1 string, dotacion_ap string, dotacion_alcance string, instalacion_ascensor string, te1_ascensor string, certificado_ascensor string, clima string, seguridad_incendio string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET electrico_te1 = ?, dotacion_ap = ?, dotacion_alcance = ?, instalacion_ascensor = ?, te1_ascensor = ?, certificado_ascensor = ?, clima = ?, seguridad_incendio = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(electrico_te1, dotacion_ap, dotacion_alcance, instalacion_ascensor, te1_ascensor, certificado_ascensor, clima, seguridad_incendio, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			if IsArrendado(id) {
				resp.Page = fmt.Sprintf("crearPropiedad4?id=%v", id)
			} else {
				resp.Page = fmt.Sprintf("crearPropiedad5?id=%v", id)
			}
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
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
func UpdatePropiedad4(db *sql.DB, token string, id int, tasacion_valor_comercial string, ano_tasacion string, contrato_arriendo string, valor_arriendo string, vencimiento_arriendo string, renovacion_automatica string, contrato_subarriendo string, valor_subarriendo string, vencimiento_subarriendo string, renovacion_automaticasub string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET tasacion_valor_comercial = ?, ano_tasacion = ?, contrato_arriendo = ?, valor_arriendo = ?, vencimiento_arriendo = ?, renovacion_automatica = ?, contrato_subarriendo = ?, valor_subarriendo = ?, vencimiento_subarriendo = ?, renovacion_automaticasub = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(tasacion_valor_comercial, ano_tasacion, contrato_arriendo, valor_arriendo, vencimiento_arriendo, renovacion_automatica, contrato_subarriendo, valor_subarriendo, vencimiento_subarriendo, renovacion_automaticasub, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad5?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad5(db *sql.DB, token string, id int, domnompropietario string, gp string, plazosarchivos string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET domnompropietario = ?, gp = ?, plazosarchivos = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(domnompropietario, gp, plazosarchivos, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad6?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad6(db *sql.DB, token string, id int, fiscal_serie string, fiscal_destino string, rol_manzana string, rol_predio string, fiscal_exento string, avaluo_fiscal string, contribucion_fiscal string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET fiscal_serie = ?, fiscal_destino = ?, rol_manzana = ?, rol_predio = ?, fiscal_exento = ?, avaluo_fiscal = ?, contribucion_fiscal = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(fiscal_serie, fiscal_destino, rol_manzana, rol_predio, fiscal_exento, avaluo_fiscal, contribucion_fiscal, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad7?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad7(db *sql.DB, token string, id int, valor_terreno string, valor_edificacion string, valor_obras_complementarias string, valor_total string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET valor_terreno = ?, valor_edificacion = ?, valor_obras_complementarias = ?, valor_total = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(valor_terreno, valor_edificacion, valor_obras_complementarias, valor_total, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = fmt.Sprintf("crearPropiedad8?id=%v", id)
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdatePropiedad8(db *sql.DB, token string, id int, cert_info_previas string, tipo_instrumento string, detalle_tipo_instrumento string, normativo_destino string, zona_normativa string, usos_permitidos string, usos_prohibidos string, densidad string, coef_constructibilidad string, coef_ocupacion_suelo string) Response {

	resp := Response{}
	resp.Op = 2
	if found, id_emp := Permisos(token, 1); found {
		stmt, err := db.Prepare("UPDATE propiedades SET cert_info_previas = ?, tipo_instrumento = ?, detalle_tipo_instrumento = ?, normativo_destino = ?, zona_normativa = ?, usos_permitidos = ?, usos_prohibidos = ?, densidad = ?, coef_constructibilidad = ?, coef_ocupacion_suelo = ? WHERE id_pro = ? AND id_emp = ?")
		ErrorCheck(err)
		_, e := stmt.Exec(cert_info_previas, tipo_instrumento, detalle_tipo_instrumento, normativo_destino, zona_normativa, usos_permitidos, usos_prohibidos, densidad, coef_constructibilidad, coef_ocupacion_suelo, id, id_emp)
		ErrorCheck(e)
		if e == nil {
			resp.Op = 1
			resp.Reload = 1
			resp.Page = "crearPropiedad"
			resp.Msg = "Empresa actualizada correctamente"
		} else {
			resp.Msg = "La empresa no pudo ser actualizada"
		}
	} else {
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
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar propiedad"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}
func InsertUsuario(db *sql.DB, token string, nombre string, p0 string, p1 string, p2 string, p3 string, p4 string, p5 string, p6 string, p7 string, p8 string, p9 string) Response {

	code := randSeq(32)
	resp := Response{}
	stmt, err := db.Prepare("INSERT INTO usuarios (user, admin, pass, code, id_emp, p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, eliminado) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	ErrorCheck(err)
	defer stmt.Close()
	pass := ""
	admin := 0
	_, e := stmt.Exec(nombre, admin, pass, code, GetIdEmp(token), p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, admin)
	ErrorCheck(e)
	if e == nil {
		SendEmail(code)
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
	if SuperAdmin(token) {
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
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar usuario"
		resp.Texto = "No tiene los permisos"
	}
	return resp
}

func InsertEmpresa(db *sql.DB, token string, nombre string, precio string) Response {

	resp := Response{}
	resp.Op = 2
	if SuperAdmin(token) {
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
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func UpdateEmpresa(db *sql.DB, token string, id int, nombre string, precio string) Response {

	resp := Response{}
	resp.Op = 2
	if SuperAdmin(token) {
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
	} else {
		resp.Msg = "No tiene permisos"
	}
	return resp
}
func BorrarEmpresa(db *sql.DB, token string, id int) Response {

	resp := Response{}
	if SuperAdmin(token) {
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
	} else {
		resp.Tipo = "error"
		resp.Titulo = "Error al eliminar empresa"
		resp.Texto = "No tiene los permisos"
	}
	return resp
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
	Code string
}

func getHTMLTemplate(code string) string {
	var templateBuffer bytes.Buffer
	data := EmailData{
		Code: code,
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
func GenerateSESTemplate(code string) (template *ses.SendEmailInput) {

	sender := "diego.gomez.bezmalinovic@gmail.com"
	receiver := "diego.gomez.bezmalinovic@gmail.com"
	html := getHTMLTemplate(code)
	title := "Nuevo Usuarios"
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
func SendEmail(code string) {

	region := "us-east-1"

	emailTemplate := GenerateSESTemplate(code)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
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

func SendEmail2() {

	region := "us-east-1"
	svc := ses.New(session.New(&aws.Config{Region: aws.String(region)}))
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

//CREATE DATABASE pelao CHARACTER SET utf8 COLLATE utf8_spanish2_ci;
