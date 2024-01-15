package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Création de la struct PGM
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// Fonction de lecture du fichier
func ReadPGM(filename string) (*PGM, error) {

	// on ouvre le fichier
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
		if scanner.Text()[0] == '#' {
			continue
		}
		break
	}

	// On crée la variable scope qui va scanner la ligne de width entiérement puis on récupére independamment width et height
	scope := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(scope[0])
	height, _ := strconv.Atoi(scope[1])

	// on récupère le max
	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// On crée la matrice data qui va contenir chaque élément de notre fichier
	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	// Si le magicNumber est égal a P2 on parcours data case par case pour récupérer toutes les valeurs du tableau
	if magicNumber == "P2" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()           // On scan chaque ligne
			byteCase := strings.Fields(line) // On récupère chaque élément qui sont séparés par un espace

			if len(byteCase) < width {
				break
			}

			for j := 0; j < width; j++ {
				value, _ := strconv.Atoi(byteCase[j]) // Enfin on récupére la valeur
				data[i][j] = uint8(value)             // Et on l'ajoute au tableau en la convertissant dans la bonne unité
			}
		}
	}

	// Si le magicNumber est égal a P2 on parcours data case par case pour récupérer toutes les valeurs du tableau
	if magicNumber == "P5" {

	}

	// On retourne chaque élément de la struct
	return &PGM{data, width, height, magicNumber, max}, nil

}

// La fonction Size retourne les valeurs de height et width
func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

// La fonction At retourne les valeurs de data a chaque position de la matrice
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// La fonction Set définit la valeur du pixel en (x, y)
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Fonction de sauvegarde
func (pgm *PGM) Save(filename string) error {

	// On crée le fichier de sauvegarde nommé filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	// On écrit les valeurs de magicNumber, width, height et max dans le fichier de sauvegarde
	fmt.Fprintf(fileSave, "%s\n%d %d\n %d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// On parcours la matrice data
	for i := range pgm.data {
		for j := range pgm.data[i] {
			fmt.Fprintf(fileSave, "%d ", pgm.data[i][j]) // et on écrit chaque valeurs de data a sa bonne position dans le fichier de sauvegarde
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

// Fonction pour inverser les couleurs
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] { // On parcours la matrice
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j] // On soustrait a data la valeur max pour avoir la valeur opposé
		}
	}
}

// Fonction pour inverser l'image horizontallement
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ { // On parcours verticalement la moitié de la matrice
		pgm.data[i], pgm.data[pgm.height-i-1] = pgm.data[pgm.height-i-1], pgm.data[i] // Et on intervertit chaque pixel
	}
}

// Fonction pour inverser l'image verticalement
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ { // On parcours notre matrice data
		count := pgm.width - 1 // Création de notre compteur pour inverser l'image une seule fois
		for j := 0; j < pgm.width/2; j++ {

			// Utilisation d'une variable temporaire pour stocker notre valeur puis inversement
			valTemp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][count]
			pgm.data[i][count] = valTemp
			count--
		}
	}
}

// Fonction pour choisir le magicNumber
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Fonction pour changer de valeur de couleur max
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	if maxValue >= 1 && maxValue <= 255 { // On vérifie que maxValue soit compris entre 1 et 255 inclus
		pgm.max = int(maxValue)

		for i := 0; i < pgm.height; i++ { // On parcours la matrice
			for j := 0; j < pgm.width; j++ {
				pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 255)) // Chaque valeur de data est modifié selon le coefficient multiplicateur de maxValue
			}
		}
	}
}

// Fonction pour faire faire un rotation a 90° dans le sens des aiguilles d'une montre a notre image
func (pgm *PGM) Rotate90CW() {

	// Création d'une nouvelle matrice rotateData pour stocker les données de rotation
	rotateData := make([][]uint8, pgm.width)
	for i := range rotateData {
		rotateData[i] = make([]uint8, pgm.height)
	}

	// Parcours de la matrice d'origine pour effectuer la rotation
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			d := pgm.height - j - 1
			rotateData[i][d] = pgm.data[j][i] // Stockage des données de rotation dans la nouvelle matrice
		}
	}

	// On mets à jour les dimensions et les données de l'image avec la matrice rotateData
	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotateData
}

func (pgm *PGM) ToPBM() *PBM {

	// Création de pbm reprenant le pointeur de la struct PBM avec les memes valurs pour width, height, magicNumber
	pbm := &PBM{
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	// Recréation de pbm.data
	pbm.data = make([][]bool, pgm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, pgm.width)
	}

	// Création de lim qui est ma valeur qui détermine si mon pixel sera blanc ou noir (inférieur a lim c'est blanc et au dessus noir)
	lim := uint8(pgm.max / 2)

	// On parcourt notre tableau et on lui ajoute la valeur correspondante a la bonne couleur
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pbm.data[i][j] = pgm.data[i][j] > lim
		}
	}

	return pbm
}

/* func main() {
	pgm, _ := ReadPGM("test.pgm")
	pbm := pgm.ToPBM()
	pbm.Save("output.pbm")
} */
