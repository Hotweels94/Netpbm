package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
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

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	if magicNumber == "P2" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			byteCase := strings.Fields(line)

			if len(byteCase) < width {
				break
			}

			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j])
				data[i][j] = uint8(value)
			}
		}
	}

	if magicNumber == "P5" {

	}

	return &PGM{data, width, height, magicNumber, max}, nil

}

func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

func (pgm *PGM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	fmt.Fprintf(fileSave, "%s\n%d %d\n %d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for i := range pgm.data {
		for j := range pgm.data[i] {
			fmt.Fprintf(fileSave, "%d ", pgm.data[i][j])
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i]
	}
}

func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		count := pgm.width - 1
		for j := 0; j < pgm.width/2; j++ {
			valTemp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][count]
			pgm.data[i][count] = valTemp
			count--
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	if maxValue >= 1 && maxValue <= 255 {
		pgm.max = int(maxValue)

		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {
				pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 255))
			}
		}
	} else {
		fmt.Println("Veuillez mettre une valeur comprise entre 1 et 255")
	}
}

func (pgm *PGM) Rotate90CW() {

	rotateData := make([][]uint8, pgm.width)
	for i := range rotateData {
		rotateData[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			d := pgm.height - j - 1
			rotateData[i][d] = pgm.data[j][i]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotateData
}

func (pgm *PGM) ToPBM() *PBM {
	pbm := new(PBM)
	pbm.magicNumber = "P1"
	pbm.width = pgm.width
	pbm.height = pgm.height

	data := make([][]bool, pgm.height)
	for i := range data {
		data[i] = make([]bool, pgm.width)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] >= 1 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}
	pbm.data = data
	return pbm
}

/* func main() {
	pgm, _ := ReadPGM("test.pgm")
	pgm.ToPBM()
} */
