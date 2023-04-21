package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

func main() {

	BarChart := ChartBar{}
	BarChart.Type = "bar"
	BarChart.Data = DataChart{Labels: []uint32{2020, 2021, 2022, 2023}}
	BarChart.Data.Datasets = []DataSets{DataSets{Label: "Users", Data: []uint32{10, 11, 12, 13}}}

	b, n, _ := ImageChart("chart1", BarChart)
	if b {
		fmt.Println(n)
	}

	RemoveFile(n)

}
func RemoveFile(file string) {

	e := os.Remove(file)
	if e != nil {
		fmt.Println(e)
	}

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
