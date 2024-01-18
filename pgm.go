package Netpbm

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Definition of the PGM struct
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// Function to read the file
func ReadPGM(filename string) (*PGM, error) {

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
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// Create the scope variable that will scan the width line entirely, then independently retrieve width and height
	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	// Retrieve the max value
	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// Create the data matrix that will contain each element of our file
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	// If the magicNumber is equal to P2, iterate over the data matrix to retrieve all values
	if magicNumber == "P2" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()           // Scan each line
			byteCase := strings.Fields(line) // Retrieve each element separated by a space

			if len(byteCase) < width {
				break
			}

			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j]) // Finally, retrieve the value
				data[i][j] = uint8(value)             // Add it to the matrix, converting it to the correct unit
			}
		}
	}

	// If the magicNumber is equal to P5:
	if magicNumber == "P5" {
		var bin string

		// Create a matrix to store characters as a string of 8 bits
		databin := make([][]string, height)
		for m := range databin {
			databin[m] = make([]string, width)
		}

		// Retrieve the data into a byte array
		scanner.Scan()
		a := scanner.Bytes()
		x := 0
		y := 0

		for g := 0; g < len(a); g++ {

			// Convert the different characters into a string of 8 bits
			format := fmt.Sprintf("%s%d%s", "%0", 8, "b")

			bin = fmt.Sprintf(format, a[g])

			// Complete the matrix
			databin[y][x] = bin

			x++
			if x == width {
				x = 0
				y = y + 1
			}

		}

		// Convert the bit strings to uint8
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				result, err := strconv.ParseInt(databin[i][j], 2, 0)

				if err != nil {
					fmt.Println("Error:", err)
				}

				data[i][j] = uint8(result) // Fill data with our final converted values

			}
		}
	}

	// Return each element of the struct
	return &PGM{data, width, height, magicNumber, max}, nil

}

// Function Size returns the values of height and width
func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

// Function At returns the values of data at each position in the matrix
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Function Set sets the value of the pixel at (x, y)
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save function
func (pgm *PGM) Save(filename string) error {

	// Create the save file named filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}
	defer fileSave.Close() // Once the function is done, close the file

	// If the magicNumber is P2
	if pgm.magicNumber == "P2" {
		// Write the values of magicNumber, width, height, and max to the save file
		fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

		// Iterate over the data matrix
		for i := range pgm.data {
			for j := range pgm.data[i] {
				fmt.Fprintf(fileSave, "%d ", pgm.data[i][j]) // Write each value of data to its correct position in the save file
			}
			fmt.Fprintln(fileSave)
		}
	}

	if pgm.magicNumber == "P5" {

		// Write to FileSave (with the correct format) the values of magicNumber, width, height, and max
		fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

		// Retrieve the values of each pixel as a string of 8 bits
		datastring_bin := make([][]string, pgm.height)
		for m := range datastring_bin {
			datastring_bin[m] = make([]string, pgm.width)
		}
		// Convert the values of the pixels to a string of 8 bits
		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.height; j++ {
				datastring_bin[i][j] = strconv.FormatInt(int64(pgm.data[i][j]), 2)

			}
		}

		// Convert the values of each pixel to an 8-bit version in hexadecimal form
		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {

				ui, err := strconv.ParseUint(datastring_bin[i][j], 2, 64)
				if err != nil {
					return err
				}
				hexa := fmt.Sprintf("%x", ui)
				if len(hexa)%2 == 0 {

					datastring_bin[i][j] = hexa
				} else {
					datastring_bin[i][j] = "0" + hexa
				}
			}
		}

		// Decode hexadecimal into a readable character
		for i := 0; i < pgm.height; i++ {
			for j := 0; j < pgm.width; j++ {
				decoded, err := hex.DecodeString(datastring_bin[i][j])
				if err != nil {
					fmt.Println("Hexadecimal decoding error:", err)
					return err
				}
				fmt.Fprintf(fileSave, string(decoded)) // Finally, write our result to the save file
			}
		}
	}

	return nil
}

// Function to invert colors
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] { // Iterate over the matrix
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j] // Subtract max from data to get the opposite value
		}
	}
}

// Function to horizontally flip the image
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ { // Iterate vertically over half of the matrix
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i] // Swap each pixel
	}
}

// Function to vertically flip the image
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ { // Iterate over our data matrix
		count := pgm.width - 1 // Create our counter to flip the image only once
		for j := 0; j < pgm.width/2; j++ {

			// Use a temporary variable to store our value and then swap
			valTemp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][count]
			pgm.data[i][count] = valTemp
			count--
		}
	}
}

// Function to set the magicNumber
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Function to change the max color value
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	newMax := float64(maxValue) / float64(pgm.max) // newMax is our Multiplicator
	pgm.max = int(maxValue)                        // pgm.max becomes our new max value

	// We run through the matrix
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			// We change the max
			pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * float64(newMax))
		}
	}
}

// Function to rotate the image 90° clockwise
func (pgm *PGM) Rotate90CW() {

	// Iterate over the original matrix to perform the rotation
	rotateData := make([][]uint8, pgm.width)
	for i := range rotateData {
		rotateData[i] = make([]uint8, pgm.height)
	}

	// Iterate over the original matrix to perform the rotation
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			d := pgm.height - j - 1
			rotateData[i][d] = pgm.data[j][i] // Store rotation data in the new matrix
		}
	}

	// Update the dimensions and data of the image with the rotateData matrix
	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotateData
}

// ToPBM converts a PGM image to a PBM image.
func (pgm *PGM) ToPBM() *PBM {

	// Recreate the PBM matrix
	pbm := &PBM{
		data:        make([][]bool, pgm.height),
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	// Iterate over the PGM matrix
	for y := 0; y < pgm.height; y++ {
		pbm.data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			// Convert each pixel value to a boolean value based on the limit
			pbm.data[y][x] = pgm.data[y][x] < uint8(pgm.max/2)
		}
	}
	return pbm
}
