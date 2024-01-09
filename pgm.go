package main

import (
	"bufio"
	"fmt"
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
		if len(scanner.Text()) > 0 && scanner.Text()[0] == '#' {
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

	fmt.Println("Data:")
	for _, row := range data {
		for _, value := range row {
			fmt.Print(value, " ")
		}
		fmt.Println()
	}

	return &PGM{data, width, height, magicNumber, max}, nil

}
