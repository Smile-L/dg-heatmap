package main

// load csv data in the range
// change axis
// get log
// smoothed
// paint
// paste to original picture

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readRawData(fileName string) [][]string {

	csvfile, err := os.Open(fileName)
	check(err)
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rawCsvData, err := reader.ReadAll()
	check(err)
	return rawCsvData
}

func readParameter(fileName string) ([9]float64, [3]float64) {
	fd, err := os.Open(fileName)
	check(err)
	defer fd.Close()

	var o2 int64
	var f2 float32

	var r [9]float64
	for i := 0; i < 9; i++ {
		o2, err = fd.Seek(int64(46*4+i*4), 0)
		check(err)
		// err = binary.Read(fd, binary.BigEndian, &f2)
		err = binary.Read(fd, binary.LittleEndian, &f2)
		check(err)
		r[i] = float64(f2)
		_ = o2
		// fmt.Println(r[i], o2/4)
	}

	var t [3]float64
	for i := 0; i < 3; i++ {
		o2, err = fd.Seek(int64(55*4+i*4), 0)
		check(err)
		// err = binary.Read(fd, binary.BigEndian, &f2)
		err = binary.Read(fd, binary.LittleEndian, &f2)
		check(err)
		t[i] = float64(f2)
		// fmt.Println(t[i], o2/4)
		_ = o2
	}
	return r, t
}

func getPicture(SensorId_selected string) image.Image {
	resp, err := http.Get("http://54.223.151.143:8082/api/slices/live.thumbnail?sensorid=" + SensorId_selected + "&cachetime=10")
	// resp, err := http.Get("http://i.imgur.com/Peq1U1u.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")
	// fmt.Println(string(lines[0]))

	url_picture := "http://54.223.151.143:8082/api/slices/thumbnail/" + lines[0]
	fmt.Println(url_picture)
	resp, err = http.Get(url_picture)
	// resp, err := http.Get("http://i.imgur.com/Peq1U1u.jpg")
	if err != nil {
		log.Fatal(err)
	}
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func getPoints(r [9]float64, t [3]float64, rawCsvData [][]string, sensorSelected string, startTime int, endTime int) []heatmap.DataPoint {
	points := []heatmap.DataPoint{}
	for _, data := range rawCsvData {

		if data[1] != sensorSelected {
			continue
		}
		time, err := strconv.Atoi(data[3])
		check(err)
		if time < startTime {
			continue
		}
		if time > endTime {
			continue
		}

		xFloor, err := strconv.ParseFloat(data[4], 64)
		check(err)
		yFloor, err := strconv.ParseFloat(data[5], 64)
		check(err)

		xCam := xFloor*r[0] + yFloor*r[3] + t[0]
		yCam := xFloor*r[1] + yFloor*r[4] + t[1]
		zCam := xFloor*r[2] + yFloor*r[5] + t[2]

		coeffX := 640 / 1.12213
		coeffY := 480 / 0.84160
		xImage := 320 - coeffX*(xCam/zCam)
		yImage := 240 + coeffY*(yCam/zCam)

		if yImage > 0 {
			points = append(points, heatmap.P(xImage, yImage))
		}

	}
	return points
}

func pasetePicture(m image.Image, img image.Image, sensorSelected string) {
	b := m.Bounds()
	mm := image.NewNRGBA(b)
	draw.Draw(mm, b, m, image.ZP, draw.Src)
	draw.Draw(mm, img.Bounds(), img, image.ZP, draw.Over)
	fmt.Println("save picuture")
	savePicName := "./output_" + sensorSelected + ".png"
	out, err := os.Create(savePicName)
	check(err)
	err = png.Encode(out, mm)
	check(err)
}

func main() {

	sensorSelected := "a1f20e44503436343300000500c90016"
	calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib4.dat"

	// sensorSelected := "a1f20e44503436343300000500a60028"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib8.dat"

	// sensorSelected := "a1f20e445034363433000005005c0022"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib14.dat"

	// sensorSelected := "a1f20e44503436343300000500880029"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib15.dat"

	// sensorSelected := "a1f20e44503436343300000500b60021"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib16.dat"

	// sensorSelected := "a1f20e445034363433000005005c0022"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib98.dat"

	// sensorSelected := "a1f20e44503436343300000500540029"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib99.dat"

	// sensorSelected := "a1f20e44503436343300000500640029"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib103.dat"

	// sensorSelected := "a1f20e44503436343300000500d80021"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib104.dat"

	// sensorSelected := "a1f20e44503436343300000500630021"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib105.dat"

	// sensorSelected := "a1f20e44503436343300000500db0028"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib106.dat"

	// fileName := "output_20150603_dr_2d.csv"
	fileName := "D:/Desktop/go/output_20150603.csv"
	rawCsvData := readRawData(fileName)
	startTime := 1433244616921
	endTime := 1433382753447
	// 1433382753447 1433244616921

	r, t := readParameter(calibFileName)
	fmt.Println(r, t)

	points := getPoints(r, t, rawCsvData, sensorSelected, startTime, endTime)

	scheme := schemes.Classic
	img := heatmap.Heatmap(image.Rect(0, 0, 640, 480),
		points, 40, 128, scheme)

	m := getPicture(sensorSelected)
	fmt.Println(m.Bounds())

	pasetePicture(m, img, sensorSelected)
}
