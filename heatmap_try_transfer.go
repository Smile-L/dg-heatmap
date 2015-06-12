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
	"github.com/lisiminy/go-heatmap"
	"github.com/lisiminy/go-heatmap/schemes"
	"image"
	// "image/draw"
	"image/color"
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

func readParameterReverse(fileName string) ([9]float64, [3]float64) {
	fd, err := os.Open(fileName)
	check(err)
	defer fd.Close()

	var o2 int64
	var f2 float32

	var r [9]float64
	for i := 0; i < 9; i++ {
		o2, err = fd.Seek(int64(34*4+i*4), 0)
		check(err)
		// err = binary.Read(fd, binary.BigEndian, &f2)
		err = binary.Read(fd, binary.LittleEndian, &f2)
		check(err)
		r[i] = float64(f2)
		fmt.Println(r[i], o2/4)
	}

	var t [3]float64
	for i := 0; i < 3; i++ {
		o2, err = fd.Seek(int64(43*4+i*4), 0)
		check(err)
		// err = binary.Read(fd, binary.BigEndian, &f2)
		err = binary.Read(fd, binary.LittleEndian, &f2)
		check(err)
		t[i] = float64(f2)
		fmt.Println(t[i], o2/4)
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
	fmt.Println(string(lines[0]))

	url_picture := "http://54.223.151.143:8082/api/slices/thumbnail/" + lines[0]
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

func main() {

	// fileName := "output_20150603_dr_2d.csv"
	fileName := "output_20150603.csv"
	rawCsvData := readRawData(fileName)
	// startTIme := 1433244616921
	// endTime := 1433382753447
	// 1433382753447 1433244616921
	// sensorSelected := "a1f20e44503436343300000500c90016"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib4.dat"

	// sensorSelected := "a1f20e44503436343300000500a60028"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib8.dat"

	// sensorSelected := "a1f20e445034363433000005005c0022"
	// calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib14.dat"

	sensorSelected := "a1f20e44503436343300000500880029"
	calibFileName := "D:/Desktop/heatmap_deepglint/floor_calib15.dat"

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

	r, t := readParameterReverse(calibFileName)
	fmt.Println(r, t)

	points := []heatmap.DataPoint{}
	var xMax, yMax float64 = 0, 0
	for _, data := range rawCsvData {
		if data[1] == sensorSelected {

			// if data[3] < startTIme {
			// 	continue
			// }
			// if data[3] > endTime {
			// 	continue
			// }

			xFloor, err := strconv.ParseFloat(data[4], 64)
			xFloor = xFloor
			xFloor = (xFloor + 1000) / 10
			check(err)
			yFloor, err := strconv.ParseFloat(data[5], 64)
			yFloor = yFloor / 10
			check(err)
			// xCam := xFloor*r[0] + yFloor*r[3] + t[0]
			// yCam := xFloor*r[1] + yFloor*r[4] + t[1]
			// zCam := xFloor*r[2] + yFloor*r[5] + t[2]

			// coeffX := 640 / 1.12213
			// coeffY := 480 / 0.84160
			// xImage := 320 - coeffX*(xCam/zCam)
			// yImage := 240 + coeffY*(yCam/zCam)
			// if yFloor > 0 {
			// if xFloor > 0 {
			points = append(points, heatmap.P(xFloor, yFloor))

			// }
			if xFloor > xMax {
				xMax = xFloor
			}
			if yFloor > yMax {
				yMax = yFloor
			}
			// }
		}
	}
	fmt.Println(xMax, yMax)
	scheme := schemes.Classic
	img := heatmap.Heatmap(image.Rect(0, 0, int(xMax), int(yMax)),
		points, 20, 128, scheme)
	//////////////////////////
	b := img.Bounds()
	m := image.NewNRGBA(b)

	// draw.Draw(mm, b, m, image.ZP, draw.Src)
	draw.Draw(m, b, img, image.ZP, draw.Over)

	savePicNamem := "./output_" + sensorSelected + "_before_transfer.png"
	outm, err := os.Create(savePicNamem)
	check(err)
	err = png.Encode(outm, m)
	check(err)
	///////////////////////////////////////////////

	imgTransfer := image.NewNRGBA(image.Rect(0, 0, 640, 480))

	coeffX := 640 / 1.12213
	coeffY := 480 / 0.84160
	for x := 0; x < 640; x++ {
		for y := 0; y < 480; y++ {
			xImage := float64(320 - x)
			yImage := float64(y - 240) ////
			// fmt.Println(y)
			zCam := -t[2] / (r[8] + r[2]*xImage/coeffX + r[5]*yImage/coeffY)
			// fmt.Println(zCam)
			xCam := xImage * zCam / coeffX
			yCam := yImage * zCam / coeffY
			xFloor := xCam*r[0] + yCam*r[3] + zCam*r[6] + t[0]
			yFloor := xCam*r[1] + yCam*r[4] + zCam*r[7] + t[1]
			// zCam := -t[2] / (r[8] + r[6]*xImage/coeffX + r[7]*yImage/coeffY)
			// fmt.Println(zCam)
			// xCam := xImage * zCam / coeffX
			// yCam := yImage * zCam / coeffY
			// xFloor := xCam*r[0] + yCam*r[1] + zCam*r[2] + t[0]
			// yFloor := xCam*r[3] + yCam*r[4] + zCam*r[5] + t[1]
			xFloorInt := (int(xFloor) + 1000) / 10 ///
			yFloorInt := int(yFloor) / 10
			// fmt.Println(xFloorInt, yFloorInt) ///
			r, g, b, a := img.At(xFloorInt, yFloorInt).RGBA()
			imgTransfer.Set(x, y, color.RGBA{
				// R: unit8(r),
				// G: unit8(g),
				// B: unit8(b),
				// A: unit8(a),
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}

	// m := getPicture(sensorSelected)
	// fmt.Println(m.Bounds())
	b2 := imgTransfer.Bounds()
	mm := image.NewNRGBA(b2)

	// draw.Draw(mm, b, m, image.ZP, draw.Src)
	draw.Draw(mm, b2, imgTransfer, image.ZP, draw.Over)

	savePicName := "./output_" + sensorSelected + "_transfer.png"
	out, err := os.Create(savePicName)
	check(err)
	err = png.Encode(out, mm)
	check(err)
}
