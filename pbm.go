package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
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
		if len(scanner.Text()) > 0 && scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	data := make(([][]bool), height)
	for i := range data {
		data[i] = make(([]bool), width)
	}

	if magicNumber == "P1" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			byteCase := strings.Fields(line)
			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j])
				if value == 1 {
					data[i][j] = true
				} else {
					data[i][j] = false
				}
			}
		}
	}

	if magicNumber == "P4" {
		reader := bufio.NewReader(file)
		for i := 0; i < height; i++ {
			for j := 0; j < width; j += 8 {
				byteValue, _ := reader.ReadByte()
				for bit := 7; bit >= 0; bit-- {
					pixel := (byteValue >> bit) & 1
					if j+7-bit < width {
						data[i][j+7-bit] = pixel == 1
					}
				}
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	for _, i := range pbm.data {
		for _, j := range i {
			if j {
				fmt.Fprint(fileSave, "1 ")
			} else {
				fmt.Fprint(fileSave, "0 ")
			}
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

func (pbm *PBM) Invert() {
	for i := range pbm.data {
		for j := range pbm.data[i] {
			if pbm.data[i][j] == true {
				pbm.data[i][j] = false
			} else {
				pbm.data[i][j] = true
			}
		}
	}
}

func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		count := pbm.width - 1
		for j := 0; j < pbm.width/2; j++ {
			valTemp := pbm.data[i][j]
			pbm.data[i][j] = pbm.data[i][count]
			pbm.data[i][count] = valTemp
			count--
		}
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

/* func main() {
	ReadPBM("test.pbm")
} */
