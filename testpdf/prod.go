package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//"image/color"

	_ "github.com/go-sql-driver/mysql"
	qrcode "github.com/skip2/go-qrcode"

	col "github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type Boletin struct {
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
}
type ListaCotizacion struct {
	Propiedad   string  `json:"Propiedad"`
	NombreAle   string  `json:"NombreAle"`
	Descripcion string  `json:"Descripcion"`
	Precio      float32 `json:"Precio"`
}

type ChartBar struct {
	Type string    `json:"type"`
	Data DataChart `json:"data"`
}
type DataChart struct {
	Labels   []uint32   `json:"labels"`
	Datasets []DataSets `json:"datasets"`
}
type DataSets struct {
	Label string   `json:"label"`
	Data  []uint32 `json:"data"`
}

type Resumen struct {
	Prods               map[int]ResumenProds   `json:"Prods"`
	Alertas             map[int]ResumenAlertas `json:"Alertas"`
	Notificaciones      map[int]ResumenAlertas `json:"Alertas"`
	TotalAlertas        int                    `json:"TotalAlertas"`
	TotalNotificaciones int                    `json:"TotalNotificaciones"`
	Localidades         Localidades            `json:"Localidades"`
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

func main() {

	b1, n1 := CreateCotizacion(1)
	if b1 {
		fmt.Println(n1)
	}

	b2, n2 := CreateBoletin(8)
	if b2 {
		fmt.Println(n2)
	}

}

func CreateBoletin(id int) (bool, string) {

	datosBoletin, err := GetResumenPropiedades(id)
	if err {
		fmt.Println("ERROR")
	}

	BarChart := ChartBar{}
	BarChart.Type = "bar"
	BarChart.Data = DataChart{Labels: []uint32{2020, 2021, 2022, 2023}}
	BarChart.Data.Datasets = []DataSets{DataSets{Label: "Users", Data: []uint32{10, 11, 12, 13}}}

	b1, chart1, _ := ImageChart("chart1", BarChart)
	if b1 {
		fmt.Println(chart1)
	}

	Mes := "Noviembre"

	pdffile := fmt.Sprintf("./tmp/boletin_%v.pdf", id)
	titulo := fmt.Sprintf("Boletin %v", Mes)

	darkGrayColor := col.Color{Red: 55, Green: 55, Blue: 55}
	//grayColor := col.Color{Red: 220, Green: 220, Blue: 220}
	//whiteColor := col.NewWhite()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 10, 10)

	// PDF
	fmt.Println(datosBoletin)

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
					Top:   10,
					Size:  14,
					Style: consts.Bold,
					Align: consts.Center,
				})
			})
			m.ColSpace(3)
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

	err1 := m.OutputFileAndClose(pdffile)
	//RemoveFile(chart1)

	if err1 != nil {
		return false, fmt.Sprintf("No se pudo guarda el Pdf: err:%v", err)
	} else {
		return true, ""
	}
}
func DatosCotizacion(id int) Cotizacion {

	cotizacion := Cotizacion{}

	db, err := GetMySQLDB()
	defer db.Close()
	ErrorCheck(err)

	res, err := db.Query("SELECT t1.precio_uf, t3.id_ale, t3.nombre as nombreale, t3.descripcion, t3.precio, t4.nombre as nombreprop, t1.fecha, t5.nombre as nombreemp FROM cotizaciones t1, cotizacion_detalle t2, alertas t3, propiedades t4, empresa t5 WHERE t1.id_cot = ? AND t1.id_cot=t2.id_cot AND t2.id_ale=t3.id_ale AND t2.id_pro=t4.id_pro AND t1.id_emp=t5.id_emp", id)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	var id_ale int
	var nombreale string
	var descripcion string
	var precio float32
	var nombreprop string
	var precio_uf float32

	var fecha string
	var nombreemp string

	for res.Next() {

		err := res.Scan(&precio_uf, &id_ale, &nombreale, &descripcion, &precio, &nombreprop, &fecha, &nombreemp)
		ErrorCheck(err)

		cotizacion.Uf = precio_uf
		cotizacion.Fecha = fecha
		cotizacion.NombreEmp = nombreemp

		cotizacion.Lista = append(cotizacion.Lista, ListaCotizacion{Propiedad: nombreprop, Descripcion: descripcion, Precio: precio, NombreAle: nombreale})
		cotizacion.TotalUf = cotizacion.TotalUf + precio

	}

	//cotizacion.Fecha = FormatDateString(cotizacion.Fecha)
	cotizacion.Fecha = ""
	cotizacion.Subtotal = cotizacion.TotalUf * cotizacion.Uf
	cotizacion.Iva = cotizacion.Subtotal * 0.19
	cotizacion.Total = cotizacion.Subtotal + cotizacion.Iva

	return cotizacion
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
func FormatDateTime(fecha string) time.Time {

	date0 := strings.Split(fecha, " ")
	date1 := strings.Split(date0[0], "-")
	date2 := strings.Split(date0[1], ":")

	Year, _ := strconv.Atoi(date1[0])
	Month, _ := strconv.Atoi(date1[1])
	Day, _ := strconv.Atoi(date1[2])

	Hour, _ := strconv.Atoi(date2[0])
	Minute, _ := strconv.Atoi(date2[1])
	Second, _ := strconv.Atoi(date2[2])

	return time.Date(Year, GetMonth(Month), Day, Hour, Minute, Second, 0, time.UTC)
}
func GetMonth(m int) time.Month {
	switch m {
	case 1:
		return time.January
	case 2:
		return time.February
	case 3:
		return time.March
	case 4:
		return time.April
	case 5:
		return time.May
	case 6:
		return time.June
	case 7:
		return time.July
	case 8:
		return time.August
	case 9:
		return time.September
	case 10:
		return time.October
	case 11:
		return time.November
	default:
		return time.December
	}
}
func CreateCotizacion(id int) (bool, string) {

	datosCot := DatosCotizacion(1)

	b, qrfile := CreateQr(id)
	if b {
		fmt.Println(qrfile)
	}

	pdffile := fmt.Sprintf("./tmp/cotizacion_%v.pdf", id)
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

	err := m.OutputFileAndClose(pdffile)
	//RemoveFile(qrfile)

	if err != nil {
		return false, fmt.Sprintf("No se pudo guarda el Pdf: err:%v", err)
	} else {
		return true, ""
	}
}
func RemoveFile(file string) {

	e := os.Remove(file)
	if e != nil {
		fmt.Println(e)
	}
}
func CreateQr(key int) (bool, string) {

	url := "https://www.redigo.cl/cotizacion/"
	urlqr := fmt.Sprintf("%v/%v", url, key)

	q, err := qrcode.New(urlqr, qrcode.Medium)
	if err != nil {
		return false, ""
	}

	name := fmt.Sprintf("./tmp/%vqr.png", key)

	q.DisableBorder = true
	q.BackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	q.ForegroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	err = q.WriteFile(128, name)
	if err != nil {
		return false, ""
	}
	return true, name
}
func GetMySQLDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/pelao")
	return
}
func ErrorCheck(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
	}
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
func ImageChart(name string, v interface{}) (bool, string, int64) {

	u, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	imagename := fmt.Sprintf("%v.jpg", name)
	file := fmt.Sprintf("./tmp/%v", imagename)

	img, err1 := os.Create(file)
	defer img.Close()
	if err1 != nil {
		return false, "", 0
	}

	chart := fmt.Sprintf("https://quickchart.io/chart?c=%v", string(u))
	fmt.Println(chart)

	resp, err2 := http.Get(chart)
	defer resp.Body.Close()
	if err2 != nil {
		return false, "", 0
	}

	b, err3 := io.Copy(img, resp.Body)
	if err3 != nil {
		return false, "", 0
	}

	return true, file, b
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

	resp.Localidades = GetLocalidades(db, id_emp)
	fmt.Println(resp.Localidades.Paises)

	res, err := db.Query("SELECT t1.id_pro as id_pro, t1.nombre as nombre, t2.id_ale as id_ale, t3.nombre as nombre_ale, t3.notificacion as tipo_notificacion FROM propiedades t1, propiedad_alerta t2, alertas t3 WHERE t1.id_emp = ? AND t1.eliminado = ? AND t1.id_pro=t2.id_pro AND t2.id_ale=t3.id_ale AND t3.eliminado = ?", id_emp, cn, cn)
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	}

	for res.Next() {

		var id_pro int
		var nombre string
		var id_ale int
		var tipo_notificacion int
		var nombre_ale string
		err := res.Scan(&id_pro, &nombre, &id_ale, &nombre_ale, &tipo_notificacion)
		if err != nil {
			fmt.Println(err)
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
func GetLocalidades(db *sql.DB, id_emp int) Localidades {

	paises := []Pais{}
	regiones := []Region{}
	ciudades := []Ciudad{}
	comunas := []Comuna{}
	propiedades := []Propiedad{}

	res1, err := db.Query("SELECT DISTINCT(t1.id_pai), t1.nombre FROM paises t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_pai=t1.id_pai", id_emp)
	defer res1.Close()
	if err != nil {
		ErrorCheck(err)
	}

	var id_pai int
	var nombrepais string

	for res1.Next() {
		err := res1.Scan(&id_pai, &nombrepais)
		if err != nil {
			ErrorCheck(err)
		}
		paises = append(paises, Pais{Id_pai: id_pai, Nombre: nombrepais})
	}

	res2, err := db.Query("SELECT DISTINCT(t1.id_reg), t1.nombre, t1.id_pai FROM regiones t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_reg=t1.id_reg", id_emp)
	defer res2.Close()
	if err != nil {
		ErrorCheck(err)
	}

	var id_reg int
	var nombreregion string

	for res2.Next() {
		err := res2.Scan(&id_reg, &nombreregion, &id_pai)
		if err != nil {
			ErrorCheck(err)
		}
		regiones = append(regiones, Region{Id_reg: id_reg, Nombre: nombreregion, Id_pai: id_pai})
	}

	res3, err := db.Query("SELECT DISTINCT(t1.id_ciu), t1.nombre, t1.id_reg, t1.id_pai FROM ciudades t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_ciu=t1.id_ciu", id_emp)
	defer res3.Close()
	if err != nil {
		ErrorCheck(err)
	}

	var id_ciu int
	var nombreciudad string

	for res3.Next() {
		err := res3.Scan(&id_ciu, &nombreciudad, &id_reg, &id_pai)
		if err != nil {
			ErrorCheck(err)
		}
		ciudades = append(ciudades, Ciudad{Id_ciu: id_ciu, Nombre: nombreciudad, Id_reg: id_reg, Id_pai: id_pai})
	}

	res4, err := db.Query("SELECT DISTINCT(t1.id_com), t1.nombre, t1.id_ciu, t1.id_reg, t1.id_pai FROM comunas t1, propiedades t2 WHERE t2.id_emp = ? AND t2.id_com=t1.id_com", id_emp)
	defer res4.Close()
	if err != nil {
		ErrorCheck(err)
	}

	var id_com int
	var nombrecomuna string

	for res4.Next() {
		err := res4.Scan(&id_com, &nombrecomuna, &id_ciu, &id_reg, &id_pai)
		if err != nil {
			ErrorCheck(err)
		}
		comunas = append(comunas, Comuna{Id_com: id_com, Nombre: nombrecomuna, Id_ciu: id_ciu, Id_reg: id_reg, Id_pai: id_pai})
	}

	cn := 0
	res0, err := db.Query("SELECT id_pro, nombre, lat, lng, direccion, numero, id_com, id_ciu, id_reg, id_pai FROM propiedades WHERE eliminado = ? AND id_emp = ?", cn, id_emp)
	defer res0.Close()
	if err != nil {
		ErrorCheck(err)
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
			ErrorCheck(err)
		}
		propiedades = append(propiedades, Propiedad{Id_pro: id_pro, Nombre: nombrepropiedad, Lat: lat, Lng: lng, Direccion: direccion, Numero: numero, Id_com: id_com, Id_ciu: id_ciu, Id_reg: id_reg, Id_pai: id_pai})
	}

	return Localidades{Paises: paises, Regiones: regiones, Ciudades: ciudades, Comunas: comunas, Propiedades: propiedades}
}
