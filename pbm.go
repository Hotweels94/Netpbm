package main

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
			for j := 0; j < width; j++ {
				data[i][j] = scanner.Scan()
			}
		}
	}

	if magicNumber == "P4" {

	}

	fmt.Printf("%+v\n", PBM{data, width, height, magicNumber})
	return &PBM{data, width, height, magicNumber}, nil

}

func main() {
	ReadPBM("test.pbm")
}
