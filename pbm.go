package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Création de la struct PBM
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// Fonction de lecture du fichier
func ReadPBM(filename string) (*PBM, error) {

	// On ouvre le fichier
	file, error := os.Open(filename)
	if error != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier")
		return nil, error
	}
	defer file.Close() // Une fois que le fonction est terminée, on ferme le fichier

	// On crée le scanner
	scanner := bufio.NewScanner(file)

	// on récupère le magicNumber
	scanner.Scan()
	magicNumber := scanner.Text()

	// On vérifie si il y a un commentaire, si oui on le saute
	for scanner.Scan() {
		if len(scanner.Text()) > 0 && scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// On crée la variable scope qui va scanner la ligne de width entiérement puis on récupére independamment width et height
	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	// On crée la matrice data qui va contenir chaque élément de notre fichier
	data := make(([][]bool), height)
	for i := range data {
		data[i] = make(([]bool), width)
	}

	// Si le magicNumber est égal a P1 on parcours data case par case pour récupérer toutes les valeurs du tableau
	if magicNumber == "P1" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()           // On scan chaque ligne
			byteCase := strings.Fields(line) // On récupère chaque élément qui sont séparés par un espace
			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j]) // Enfin on récupére la valeur
				if value == 1 {
					data[i][j] = true // Si value = 1, data sera = true
				} else {
					data[i][j] = false // Sinon ce sera false
				}
			}
		}
	}

	if magicNumber == "P4" {
		ascii := 0
		n := 0
		for n <= width {
			ascii += 1
			n += 8

		}

		data3 := make([][]int, height)
		for g := range data3 {
			data3[g] = make([]int, ascii)
		}

		chars := make([][]rune, height)
		for g := range chars {
			chars[g] = make([]rune, ascii)
		}

		var bin string

		datarune := make([][]string, height)
		for m := range datarune {
			datarune[m] = make([]string, ascii)
		}

		scanner.Scan()
		a := scanner.Bytes()
		x := 0
		y := 0

		for g := 0; g < len(a); g++ {
			format := fmt.Sprintf("%s%d%s", "%0", 8, "b")
			data3[y][x] = int(a[g])

			bin = fmt.Sprintf(format, data3[y][x])

			datarune[y][x] = bin

			x++
			if x == ascii {
				x = 0
				y = y + 1
			}

		}
		datastring := make([]string, height)

		for i := 0; i < height; i++ {
			for j := 0; j < ascii; j++ {
				datastring[i] = datastring[i] + datarune[i][j]
			}
		}

		datarune_padding := make([][]rune, height)
		for m := range datarune_padding {
			datarune_padding[m] = make([]rune, width)
		}
		for i := 0; i < height; i++ {
			l := []rune(datastring[i])
			for j := 0; j < width; j++ {

				datarune_padding[i][j] = l[j]

			}
		}

		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {

				if datarune_padding[i][j] == '0' {
					data[i][j] = false

				} else if datarune_padding[i][j] == '1' {
					data[i][j] = true

				}

			}
		}
	}

	// On retourne chaque élément de la struct
	return &PBM{data, width, height, magicNumber}, nil
}

// La fonction Size retourne les valeurs de height et width
func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

// La fonction At retourne les valeurs de data a chaque position de la matrice
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// La fonction Set définit la valeur du pixel en (x, y)
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Fonction de sauvegarde
func (pbm *PBM) Save(filename string) error {

	// On crée le fichier de sauvegarde nommé filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	// On écrit les valeurs de magicNumber, width et height dans le fichier de sauvegarde
	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// On parcours la matrice data
	for _, i := range pbm.data {
		for _, j := range i {
			if j {
				fmt.Fprint(fileSave, "1 ") // Si data[i][j est égal a true on écrit 1 dans le fichier de sauvergarde
			} else {
				fmt.Fprint(fileSave, "0 ") // Sinon on écrit 0
			}
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

// Fonction pour inverser les couleurs
func (pbm *PBM) Invert() {

	// On parcours data et on inverse les valeurs booléenes associés de base a data[i][j]
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

// Fonction pour inverser l'image horizontallement
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ { // On parcours verticalement la moitié de la matrice
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i] // Et on intervertit chaque pixel
	}
}

// Fonction pour inverser l'image verticalement
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ { // On parcours notre matrice data
		count := pbm.width - 1 // Création de notre compteur pour inverser l'image une seule fois
		for j := 0; j < pbm.width/2; j++ {

			// Utilisation d'une variable temporaire pour stocker notre valeur puis inversement
			valTemp := pbm.data[i][j]
			pbm.data[i][j] = pbm.data[i][count]
			pbm.data[i][count] = valTemp
			count--
		}
	}
}

// Fonction pour choisir le magicNumber
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

/* func main() {
	ReadPBM("test.pbm")
} */
