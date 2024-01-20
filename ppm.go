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

// Creating the PPM struct
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// Creating the Point Struct
type Point struct {
	X, Y int
}

// Creating the Pixel Struct
type Pixel struct {
	R, G, B uint8
}

// Function for reading the file
func ReadPPM(filename string) (*PPM, error) {

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

	// Check for comments and skip them if present
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// Create the scope variable to scan the width line entirely, then retrieve width and height independently
	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	// Get the max value
	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// Create the data matrix to store each element of the file
	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
	}

	// If the magicNumber is P3, iterate through data cell by cell to retrieve all values from the array
	if magicNumber == "P3" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()              // We scan each line
			byteCaseRGB := strings.Fields(line) // We retrieve each element separated by a space
			for j := 0; j < width; j++ {
				r, _ := strconv.Atoi(byteCaseRGB[j*3])           // Red for each first element of each block of 3 numbers
				g, _ := strconv.Atoi(byteCaseRGB[j*3+1])         // Blue for each second element in each block of 3 numbers
				b, _ := strconv.Atoi(byteCaseRGB[j*3+2])         // Green for each third element in each block of 3 numbers
				data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)} // And data [i][j] is a collection of the 3 colors (r, g, b)
			}
		}
	}

	// If magic Number is P6
	if magicNumber == "P6" {
		var bin string

		// Create a matrix to store characters as a string of 8 bits
		databin := make([][]string, height)
		for m := range databin {
			databin[m] = make([]string, 3*width)
		}

		// Retrieve the data into a byte array
		scanner.Scan()
		byte_array := scanner.Bytes()
		x := 0
		y := 0

		for g := 0; g < len(byte_array); g++ {

			// Convert the different characters into a string of 8 bits
			format := fmt.Sprintf("%s%d%s", "%0", 8, "b")

			bin = fmt.Sprintf(format, byte_array[g])

			// Complete the matrix
			databin[y][x] = bin

			x++
			if x == 3*width { // Multiplacate by 3 because one pixel is equal to 3 value because we are in RGB
				x = 0
				y = y + 1
			}
		}

		// Convert the bit strings to uint8 but this time for each color
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				r, _ := strconv.ParseInt(databin[i][j*3], 2, 0)
				g, _ := strconv.ParseInt(databin[i][j*3+1], 2, 0)
				b, _ := strconv.ParseInt(databin[i][j*3+2], 2, 0)

				// And we stock it in data for each color (depend of start value 0, 1, 2 to know which color is it)
				data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
			}
		}
	}

	// Return each element of the struct
	return &PPM{data, width, height, magicNumber, max}, nil
}

// The Size function returns the values of height and width
func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

// The At function returns the values of data at each position in the matrix
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// The Set function sets the value of the pixel at (x, y)
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save function
func (ppm *PPM) Save(filename string) error {

	// Create the save file named filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}
	defer fileSave.Close()

	// If magicNumber is P3
	if ppm.magicNumber == "P3" {
		// Write the values of magicNumber, width, height, and max to the save file
		fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
		// Iterate through the data matrix
		for i := range ppm.data {
			for j := range ppm.data[i] {
				// Write each value of data to its correct position in the save file, color by color
				fmt.Fprintf(fileSave, "%d %d %d ", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
			}
			fmt.Fprintln(fileSave)
		}
	}

	// If magicNumber is P6
	if ppm.magicNumber == "P6" {

		// Write to FileSave (with the correct format) the values of magicNumber, width, height, and max
		fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

		// Retrieve the values of each pixel as a string of 8 bits
		datastring_bin := make([][]string, ppm.height)
		for m := range datastring_bin {
			datastring_bin[m] = make([]string, ppm.width)
		}
		// Convert the values of the pixels to a string of 8 bits for each color
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.height; j++ {
				r := int64(ppm.data[i][j].R)
				g := int64(ppm.data[i][j].G)
				b := int64(ppm.data[i][j].B)
				datastring_bin[i][j] = strconv.FormatInt(r, 2) + strconv.FormatInt(g, 2) + strconv.FormatInt(b, 2)
			}
		}

		// Convert the values of each pixel to an 8-bit version in hexadecimal form
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {

				ui, _ := strconv.ParseUint(datastring_bin[i][j], 2, 64)
				hexa := fmt.Sprintf("%x", ui)
				if len(hexa)%2 == 0 {

					datastring_bin[i][j] = hexa
				} else {
					datastring_bin[i][j] = "0" + hexa
				}
			}
		}

		// Decode hexadecimal into a readable character
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				decoded, _ := hex.DecodeString(datastring_bin[i][j])
				fmt.Fprintf(fileSave, string(decoded)) // Finally, write our result to the save file
			}
		}
	}
	return nil
}

// Function to invert colors
func (ppm *PPM) Invert() {
	for i := range ppm.data {
		for j := range ppm.data[i] { // Browse the matrix

			// subtract a data from the max value to obtain the opposite value for each color (for r, for g, for b)
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

// Function to horizontally flip the image
func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ { // We vertically traverse half of the matrix
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i] // And invert each pixel
	}
}

// Function to vertically flip the image
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ { // Run through our data matrix
		count := ppm.width - 1 // Creation of our counter to invert the image once only
		for j := 0; j < ppm.width/2; j++ {

			// Use a temporary variable to store our value and vice versa
			valTemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][count]
			ppm.data[i][count] = valTemp
			count--
		}
	}
}

// Function to set the magicNumber
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// Function to change the max color value
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if maxValue <= 255 || maxValue >= 1 {
		newMax := float64(maxValue) / float64(ppm.max) // newMax is our Multiplicator
		ppm.max = int(maxValue)                        // ppm.max becomes our new max value

		// We run through the matrix
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				// We change the max value for each color
				ppm.data[i][j].R = uint8(math.Round(float64(ppm.data[i][j].R) * float64(newMax)))
				ppm.data[i][j].G = uint8(math.Round(float64(ppm.data[i][j].G) * float64(newMax)))
				ppm.data[i][j].B = uint8(math.Round(float64(ppm.data[i][j].B) * float64(newMax)))
			}
		}
	}
}

// Function to rotate the image 90 degrees clockwise
func (ppm *PPM) Rotate90CW() {

	// Create a new rotateData matrix to store rotation data
	rotateData := make([][]Pixel, ppm.width)
	for i := range rotateData {
		rotateData[i] = make([]Pixel, ppm.height)
	}

	// Traverse the original matrix to perform the rotation
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			d := ppm.height - j - 1
			rotateData[i][d] = ppm.data[j][i] // Store rotation data in the new matrix
		}
	}

	// Image dimensions and data are updated using the rotateData matrix
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotateData
}

// Function to convert from PPM to PGM
func (ppm *PPM) ToPGM() *PGM {

	//Creation of pgm using the struct PGM pointer with the same values for width, height, magicNumber and max
	pgm := &PGM{
		magicNumber: "P2",
		width:       ppm.width,
		height:      ppm.height,
		max:         ppm.max,
	}

	// Recreate pgm.data
	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	// Scan the matrix and average the 3 colors to obtain a gray value for pgm.data
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			grey := (float64(ppm.data[i][j].R) + float64(ppm.data[i][j].G) + float64(ppm.data[i][j].B)) / 3
			pgm.data[i][j] = uint8(grey)
		}
	}

	return pgm
}

// Function to convert from PPM to PBM
func (ppm *PPM) ToPBM() *PBM {

	// Creation of pbm using the struct PBM pointer with the same values for width, height and magicNumber
	pbm := &PBM{
		magicNumber: "P1",
		width:       ppm.width,
		height:      ppm.height,
	}

	// Recreate pbm.data
	data := make([][]bool, ppm.height)
	for i := range data {
		data[i] = make([]bool, ppm.width)
	}

	// Assign data to pbm.data
	pbm.data = data

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			// Convert each pixel to black or white according to the limit
			lim := (uint16(ppm.data[i][j].R) + uint16(ppm.data[i][j].G) + uint16(ppm.data[i][j].B)) / 3
			pbm.data[i][j] = lim < uint16(ppm.max/2)
		}
	}

	return pbm
}

// For the DrawLine function, we'll use Bresenham's Algorithm.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// calculate the horizontal pitch of our line
	deltaX := p2.X - p1.X
	if deltaX < 0 { // We'll use deltaX absolute values to point us in the right direction
		deltaX = -deltaX
	}

	// calculate the vertical pitch of our line
	deltaY := p2.Y - p1.Y
	if deltaY < 0 { // We'll use the absolute values of deltaY to point us in the right direction.
		deltaY = -deltaY
	}

	// Here we calculate signX, which lets us know whether we're drawing our line from left to right or vice versa (like a directing coefficient).
	signX := -1 // In this case from right to left.
	if p1.X < p2.X {
		signX = 1 // In this case, left to right.
	}

	// Here we calculate signY, which lets us know whether we're drawing our line from bottom to top or vice versa.
	signY := -1 // In this case from top to bottom (the opposite of what we normally have in Math)
	if p1.Y < p2.Y {
		signY = 1 // In this case from bottom to top.
	}

	err := deltaX - deltaY

	// Create the drawing loop for our right.
	for {

		// Check that the pixel is within the image boundaries
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			// Here we color the pixel
			ppm.Set(p1.X, p1.Y, color)
		}

		// If the line has finished drawing, stop the drawing loop.
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		// This Variable is used to determine when we need to move in the direction of Y (vertically).
		err2 := 2 * err

		// If err2 is greater than the opposite of deltaY, we must move in the X direction.
		if err2 > -deltaY {
			err -= deltaY // We compensate for the fact that we've moved in direction X
			p1.X += signX // And it move
		}

		// If err2 is less than deltaX, we must move in direction Y
		if err2 < deltaX {
			err += deltaX // We compensate for the fact that we've moved in direction Y
			p1.Y += signY // And it move
		}

		// Finally, we double-check that the point is within the limits
		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
	}
}

// Function to draw an empty rectangle
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {

	// Check that the points are not out of bounds
	if p1.X < 0 {
		p1.X = 0
	}
	if p1.Y < 0 {
		p1.Y = 0
	}

	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
	}

	// Create the 3 corners of the rectangle (+ p1 in the function, either 4 points in all)
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	// We link them all so as to make a loop
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// Function to draw a filled rectangle
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {

	// We go through our matrix line by line
	for i := 0; i < height+1; i++ {

		// point1 and point2 are on either side of the length of the rectangle
		point1 := Point{p1.X, p1.Y + i}
		point2 := Point{p1.X + width, p1.Y + i}
		ppm.DrawLine(point1, point2, color) // And we link them
	}
}

// Function that draws a empty circle
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	// We go through the matrix
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			deltaX := float64(i) - float64(center.X)
			deltaY := float64(j) - float64(center.Y)
			distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				ppm.Set(i, j, color)
			}
		}
	}
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

// Function that draws a filled circle
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {

	ppm.DrawCircle(center, radius, color)
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var number_points int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				number_points += 1
				positions = append(positions, j)
			}
		}
		if number_points > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color

			}
		}
	}
}

// Function that draws an empty triangle
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {

	// We just need to link our 3 points between them
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// Function that draws an filled triangle
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {

	// We call our function to draw an empty Triangle
	ppm.DrawTriangle(p1, p2, p3, color)

	for i := 0; i < ppm.height; i++ { // We run through the matrix
		var positions []int   // We initialize the positions to color
		var number_points int // and how many points to color

		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color { // if data equal a part of our a empty triangle (The ribs of the triangle)
				number_points += 1 // number of points to color in the triangle increase
				positions = append(positions, j)
			}
		}

		if number_points == 2 { // If on the line there is ONLY 2 points colored (left and right part of our triangle)
			for k := positions[0] + 1; k < positions[1]; k++ { // k run through the line inside the triangle
				ppm.data[i][k] = color // we color the triangle
			}
		}
	}
}

// Function that draws an empty polygon
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// We go through all the points one by one to link them together
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	// Draw the last straight line from the last point to the first point (close the loop)
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {

	// We call our function to draw an empty Polygon
	ppm.DrawPolygon(points, color)

	for i := 0; i < ppm.height; i++ { // We run through the matrix
		var positions []int   // We initialize the positions to color
		var number_points int // and how many points to color
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color { // if data equal a part of our a empty polygon (The ribs of the polygon)
				number_points += 1 // number of points to color in the polygon increase
				positions = append(positions, j)
			}
		}
		if number_points == 2 { // If on the line there is ONLY 2 points colored (left and right part of our polygon)
			for k := positions[0] + 1; k < positions[1]; k++ { // k run through the line inside the polygon
				ppm.data[i][k] = color // we color the polygon

			}
		}
	}
}
