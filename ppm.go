package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

type Pixel struct {
	R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error) {
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier")
		return nil, error
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	magicNumber := scanner.Text()

	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
	}

	if magicNumber == "P3" {

		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			byteCaseRGB := strings.Fields(line)
			for j := 0; j < width; j++ {
				r, _ := strconv.Atoi(byteCaseRGB[j*3])
				g, _ := strconv.Atoi(byteCaseRGB[j*3+1])
				b, _ := strconv.Atoi(byteCaseRGB[j*3+2])
				data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
			}
		}
	}
	fmt.Println("Data:")
	for _, row := range data {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println()
	}
	return &PPM{data, width, height, magicNumber, max}, nil
}

func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[x][y]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

func (ppm *PPM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	fmt.Fprintf(fileSave, "%s\n%d %d\n %d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	for i := range ppm.data {
		for j := range ppm.data[i] {
			fmt.Fprintf(fileSave, "%d %d %d ", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

func (ppm *PPM) Invert() {
	for i := range ppm.data {
		for j := range ppm.data[i] {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i]
	}
}

func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		count := ppm.width - 1
		for j := 0; j < ppm.width/2; j++ {
			valTemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][count]
			ppm.data[i][count] = valTemp
			count--
		}
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {

}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if maxValue <= 255 || maxValue >= 1 {
		newMax := float64(maxValue) / float64(ppm.max)
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				ppm.data[i][j].R = uint8(math.Round(float64(ppm.data[i][j].R) * float64(newMax)))
				ppm.data[i][j].G = uint8(math.Round(float64(ppm.data[i][j].G) * float64(newMax)))
				ppm.data[i][j].B = uint8(math.Round(float64(ppm.data[i][j].B) * float64(newMax)))
			}
		}
	} else {
		fmt.Println("Veuillez mettre une valeure comprise entre 1 et 255")
	}
}

func (ppm *PPM) Rotate90CW() {
	rotateData := make([][]Pixel, ppm.width)
	for i := range rotateData {
		rotateData[i] = make([]Pixel, ppm.height)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			d := ppm.height - j - 1
			rotateData[i][d] = ppm.data[j][i]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotateData
}

//func (ppm *PPM) ToPGM() *PGM {
//
//}

func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		magicNumber: "P1",
		width:       ppm.width,
		height:      ppm.height,
	}

	data := make([][]bool, ppm.height)
	for i := range data {
		data[i] = make([]bool, ppm.width)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j].R == 0 && ppm.data[i][j].G == 0 && ppm.data[i][j].B == 0 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}
	return pbm
}

/* func main() {
	ppm, _ := ReadPPM("test.ppm")
	ppm.Flip()
} */
