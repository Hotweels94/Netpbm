package Netpbm

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Creation of the PBM struct
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// Function for reading the file
func ReadPBM(filename string) (*PBM, error) {

	// Open the file
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Error opening the file")
		return nil, error
	}
	defer file.Close() // Close the file once the function is done

	// Create the scanner
	scanner := bufio.NewScanner(file)

	// Get the magicNumber
	scanner.Scan()
	magicNumber := scanner.Text()

	// Check if there is a comment, if yes, skip it
	for scanner.Scan() {
		if len(scanner.Text()) > 0 && scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// Create the scope variable that will scan the width line entirely, then extract width and height independently
	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	// Create the data matrix that will contain each element of our file
	data := make(([][]bool), height)
	for i := range data {
		data[i] = make(([]bool), width)
	}

	// If the magicNumber is equal to P1, iterate through data element by element to retrieve all values from the array
	if magicNumber == "P1" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()           // Scan each line
			byteCase := strings.Fields(line) // Retrieve each element separated by a space
			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j]) // Finally, retrieve the value
				if value == 1 {
					data[i][j] = true // If value = 1, data will be true
				} else {
					data[i][j] = false // Otherwise, it will be false
				}
			}
		}
	}

	// If the MagicNumber is P4
	if magicNumber == "P4" {

		// Find out how many bytes are on a single line
		nbr_byte := 0
		n := 0
		for n <= width {
			nbr_byte += 1
			n += 8

		}
		// Then create a second data matrix in hexadecimal version
		data_hexa := make([][]int, height)
		for g := range data_hexa {
			data_hexa[g] = make([]int, nbr_byte)
		}

		var bin string

		// Create a third matrix to store each value as 8 bits
		data_bin := make([][]string, height)
		for m := range data_bin {
			data_bin[m] = make([]string, nbr_byte)
		}

		// Scan the data and put it into a byte array
		scanner.Scan()

		a := scanner.Bytes()
		x := 0
		y := 0

		for g := 0; g < len(a); g++ {

			// Conversion of hexadecimal characters into 8 bits

			format := fmt.Sprintf("%s%d%s", "%0", 8, "b") // Create the conversion format
			data_hexa[y][x] = int(a[g])

			bin = fmt.Sprintf(format, data_hexa[y][x])

			data_bin[y][x] = bin

			x++
			if x == nbr_byte {
				x = 0
				y = y + 1
			}

		}
		datastring := make([]string, height)

		// Convert our matrix into a bit string
		for i := 0; i < height; i++ {
			for j := 0; j < nbr_byte; j++ {
				datastring[i] = datastring[i] + data_bin[i][j]
			}
		}

		// Create the final matrix which will be bit strings but this time without padding bits
		databit_without_padding := make([][]rune, height)
		for m := range databit_without_padding {
			databit_without_padding[m] = make([]rune, width)
		}
		for i := 0; i < height; i++ {
			bit := []rune(datastring[i])
			for j := 0; j < width; j++ {

				databit_without_padding[i][j] = bit[j]

			}
		}
		// Complete the final data matrix
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {

				if databit_without_padding[i][j] == '0' {
					data[i][j] = false

				} else if databit_without_padding[i][j] == '1' {
					data[i][j] = true

				}

			}
		}
	}
	// Return each element of the struct
	return &PBM{data, width, height, magicNumber}, nil
}

// The Size function returns the values of height and width
func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

// The At function returns the values of data at each position in the matrix
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// The Set function sets the value of the pixel at (x, y)
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save function
func (pbm *PBM) Save(filename string) error {

	// Create the save file named filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	// Write the values of magicNumber, width, and height to the save file
	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// If the magic number is P1
	if pbm.magicNumber == "P1" {
		// Iterate through the data matrix
		for _, i := range pbm.data {
			for _, j := range i {
				if j {
					fmt.Fprint(fileSave, "1 ") // If data[i][j] is true, write 1 to the save file
				} else {
					fmt.Fprint(fileSave, "0 ") // Otherwise, write 0
				}
			}
			fmt.Fprintln(fileSave)
		}
	}

	// If the magicNumber is P4
	if pbm.magicNumber == "P4" {
		// Select the number of bytes present on a line
		number_bytes := math.Round(float64(pbm.width/8) + 1)

		// Calculate the number of padding bits
		padding_bytes := 8 * number_bytes / float64(pbm.width)

		// If there are padding bits, create a matrix to store the data in bits (0 or 1)
		datarune_padding := make([][]rune, pbm.height)
		for m := range datarune_padding {
			datarune_padding[m] = make([]rune, pbm.width+int(padding_bytes))
		}

		// Iterate through our matrix
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < pbm.width; j++ {
				if pbm.data[i][j] { // If pbm.data at position [i][j] is true
					datarune_padding[i][j] = 1 // datarune_padding becomes '1'
				}
				if !pbm.data[i][j] { // Otherwise,
					datarune_padding[i][j] = 0 // datarune_padding becomes '0'
				}
			}
		}
		// If there are padding bits, complete with zeros
		if pbm.width%8 != 0 {
			for i := 0; i < pbm.height; i++ {
				for j := pbm.width; j < pbm.width+int(padding_bytes); j++ {
					datarune_padding[i][j] = 0
				}
			}
		}
		// Create a matrix that will contain the bytes converted to bits
		datastring_padding := make([][]string, pbm.height)
		for m := range datastring_padding {
			datastring_padding[m] = make([]string, int(number_bytes))
		}

		// Complete the matrix with the bytes
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < int(number_bytes); j++ {
				var a string

				for m := 8 * j; m < 8*(1+j); m++ {

					b := fmt.Sprintf("%v", datarune_padding[i][m])

					a = a + b

				}
				datastring_padding[i][j] = a

			}

		}
		// Convert bytes to hexadecimal
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < int(number_bytes); j++ {

				ui, err := strconv.ParseUint(datastring_padding[i][j], 2, 64)
				if err != nil {
					return err
				}
				hexa := fmt.Sprintf("%x", ui)
				if len(hexa)%2 == 0 {

					datastring_padding[i][j] = hexa
				} else {
					datastring_padding[i][j] = "0" + hexa
				}
			}
		}

		// Convert hexadecimal to characters
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < int(number_bytes); j++ {
				decoded, err := hex.DecodeString(datastring_padding[i][j])
				if err != nil {
					fmt.Println("Hexadecimal decoding error:", err)
					return err
				}
				fmt.Fprintf(fileSave, string(decoded)) // And write the result to our save file
			}
		}
	}

	return nil
}

// Function to invert colors
func (pbm *PBM) Invert() {

	// Iterate through data and invert the boolean values associated with data[i][j] by default
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

// Function to horizontally invert the image
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ { // Iterate vertically through half of the matrix
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i] // And swap each pixel
	}
}

// Function to vertically invert the image
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ { // Iterate through our data matrix
		count := pbm.width - 1 // Create our counter to invert the image only once
		for j := 0; j < pbm.width/2; j++ {

			// Use a temporary variable to store our value and then invert it
			valTemp := pbm.data[i][j]
			pbm.data[i][j] = pbm.data[i][count]
			pbm.data[i][count] = valTemp
			count--
		}
	}
}

// Function to choose the magicNumber
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
