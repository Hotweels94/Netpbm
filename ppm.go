package main

import (
	"bufio"
	"fmt"
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

	fmt.Printf("%+v\n", &PPM{data, width, height, magicNumber, max})
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

func main() {
	ppm, _ := ReadPPM("test.ppm")
	ppm.Save("save.ppm")
}
