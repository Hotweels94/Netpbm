package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Création de la struct PGM
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// Création de la fonction Point
type Point struct {
	X, Y int
}

// Création de la fonction Pixel
type Pixel struct {
	R, G, B uint8
}

// Fonction de lecture du fichier
func ReadPPM(filename string) (*PPM, error) {

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
	data := make([][]Pixel, height)
	for i := range data {
		data[i] = make([]Pixel, width)
	}

	// Si le magicNumber est égal a P3 on parcours data case par case pour récupérer toutes les valeurs du tableau
	if magicNumber == "P3" {
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()              // On scan chaque ligne
			byteCaseRGB := strings.Fields(line) // On récupère chaque élément qui sont séparés par un espace
			for j := 0; j < width; j++ {
				r, _ := strconv.Atoi(byteCaseRGB[j*3])           // Rouge pour chaque premier élement de chaque bloc de 3 nombres
				g, _ := strconv.Atoi(byteCaseRGB[j*3+1])         // Bleu pour chaque second élement de chaque bloc de 3 nombres
				b, _ := strconv.Atoi(byteCaseRGB[j*3+2])         // Vert pour chaque troisiéme élement de chaque bloc de 3 nombres
				data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)} // Et data [i][j] est un rassemblement des 3 couleurs (r, g, b)
			}
		}
	}

	if magicNumber == "P6" {

	}

	// On retourne chaque élément de la struct
	return &PPM{data, width, height, magicNumber, max}, nil
}

// La fonction Size retourne les valeurs de height et width
func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

// La fonction At retourne les valeurs de data a chaque position de la matrice
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// La fonction Set définit la valeur du pixel en (x, y)
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Fonction de sauvegarde
func (ppm *PPM) Save(filename string) error {

	// On crée le fichier de sauvegarde nommé filename
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}

	// On écrit les valeurs de magicNumber, width, height et max dans le fichier de sauvegarde
	fmt.Fprintf(fileSave, "%s\n%d %d\n %d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	// On parcours la matrice data
	for i := range ppm.data {
		for j := range ppm.data[i] {
			// et on écrit chaque valeurs de data a sa bonne position dans le fichier de sauvegarde couleur par couleur
			fmt.Fprintf(fileSave, "%d %d %d ", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
		}
		fmt.Fprintln(fileSave)
	}
	return nil
}

// Fonction pour inverser les couleurs
func (ppm *PPM) Invert() {
	for i := range ppm.data {
		for j := range ppm.data[i] { // On parcours la matrice

			// On soustrait a data la valeur max pour avoir la valeur opposé et ce pour chaque couleur (pour r, pour g, pour b)
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

// Fonction pour inverser l'image horizontallement
func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ { // On parcours verticalement la moitié de la matrice
		ppm.data[i], ppm.data[ppm.height-i-1] = ppm.data[ppm.height-i-1], ppm.data[i] // Et on intervertit chaque pixel
	}
}

// Fonction pour inverser l'image verticalement
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ { // On parcours notre matrice data
		count := ppm.width - 1 // Création de notre compteur pour inverser l'image une seule fois
		for j := 0; j < ppm.width/2; j++ {

			// Utilisation d'une variable temporaire pour stocker notre valeur puis inversement
			valTemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][count]
			ppm.data[i][count] = valTemp
			count--
		}
	}
}

// Fonction pour choisir le magicNumber
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// Fonction pour changer de valeur de couleur max
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if maxValue <= 255 || maxValue >= 1 {
		newMax := float64(maxValue) / float64(ppm.max)
		ppm.max = int(maxValue)
		for i := 0; i < ppm.height; i++ {
			for j := 0; j < ppm.width; j++ {
				ppm.data[i][j].R = uint8(math.Round(float64(ppm.data[i][j].R) * float64(newMax)))
				ppm.data[i][j].G = uint8(math.Round(float64(ppm.data[i][j].G) * float64(newMax)))
				ppm.data[i][j].B = uint8(math.Round(float64(ppm.data[i][j].B) * float64(newMax)))
			}
		}
	}
}

// Fonction pour faire faire un rotation a 90° dans le sens des aiguilles d'une montre a notre image
func (ppm *PPM) Rotate90CW() {

	// Création d'une nouvelle matrice rotateData pour stocker les données de rotation
	rotateData := make([][]Pixel, ppm.width)
	for i := range rotateData {
		rotateData[i] = make([]Pixel, ppm.height)
	}

	// Parcours de la matrice d'origine pour effectuer la rotation
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			d := ppm.height - j - 1
			rotateData[i][d] = ppm.data[j][i] // Stockage des données de rotation dans la nouvelle matrice
		}
	}

	// On mets à jour les dimensions et les données de l'image avec la matrice rotateData
	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotateData
}

// Fonction pour convertir de PPM a PGM
func (ppm *PPM) ToPGM() *PGM {

	//Création de pgm reprenant le pointeur de la struct PGM avec les memes valurs pour width, height, magicNumber et max
	pgm := &PGM{
		magicNumber: "P2",
		width:       ppm.width,
		height:      ppm.height,
		max:         ppm.max,
	}

	// Recréation de pgm.data
	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	// On parcourt la matrice et on prend la moyenne des 3 couleurs pour avoir une valeur de gris pour pgm.data
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = (ppm.data[i][j].R + ppm.data[i][j].G + ppm.data[i][j].B) / 3
		}
	}

	return pgm
}

// Fonction pour convertir de PPM a PBM
func (ppm *PPM) ToPBM() *PBM {

	// Création de pbm reprenant le pointeur de la struct PBM avec les memes valurs pour width, height, magicNumber
	pbm := &PBM{
		magicNumber: "P1",
		width:       ppm.width,
		height:      ppm.height,
	}

	// Recréation de pbm.data
	data := make([][]bool, ppm.height)
	for i := range data {
		data[i] = make([]bool, ppm.width)
	}

	// Création de lim qui est ma valeur qui détermine si mon pixel sera blanc ou noir (inférieur a lim c'est blanc et au dessus noir)
	lim := uint8(ppm.max / 2)

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			// Convertir chaque pixel en noir ou blanc en fonction de la limite
			pbm.data[i][j] = ppm.data[i][j].R > lim || ppm.data[i][j].G > lim || ppm.data[i][j].B > lim
		}
	}
	return pbm
}

// Pour la fonction DrawLine nous allons utiliser l'Algorithme de Bresenham.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// calcul du pas horizontal de notre droite
	deltaX := p2.X - p1.X
	if deltaX < 0 { // On utilisera les valeurs absolue de deltaX pour se diriger dans la bonne direction
		deltaX = -deltaX
	}

	// calcul du pas vertical de notre droite
	deltaY := p2.Y - p1.Y
	if deltaY < 0 { // On utilisera les valeurs absolue de deltaY pour se diriger dans la bonne direction
		deltaY = -deltaY
	}

	// Ici on calcul signX qui nous permet de savoir si on trace notre droite de gauche a droite ou l'inverse ( comme un coefficient directeur).
	signX := -1 // Dans ce cas de droite a gauche.
	if p1.X < p2.X {
		signX = 1 // Dans ce cas de gauche a droite.
	}

	// Ici on calcul signY qui nous permet de savoir si on trace notre droite de bas en haut ou l'inverse.
	signY := -1 // Dans ce cas de haut en bas. (C'est l'inverse de ce que l'on a en Math)
	if p1.Y < p2.Y {
		signY = 1 // Dans ce cas de bas en haut.
	}

	err := deltaX - deltaY

	// Création de la boucle de dessin de notre droite.
	for {

		// On vérifie que le pixel est dans les limites de l'image
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
			// Ici on colorie le pixel
			ppm.Set(p1.X, p1.Y, color)
		}

		// Si la ligne a fini de se tracer on stop la boucle de dessin.
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		// Cette Variable est utilisé pour savoir a quelle moment on va devoir avancer en direction de Y (verticalement)
		err2 := 2 * err

		// Si err2 est supérieur a l'opposé de deltaY on doit avancer dans la direction X
		if err2 > -deltaY {
			err -= deltaY // On compense le fait que l'on a avancé dans la direction X
			p1.X += signX // Et on fait le déplacement
		}

		// Si err2 est inférieur a deltaX on doit avancer dans la direction Y
		if err2 < deltaX {
			err += deltaX // On compense le fait que l'on a avancé dans la direction Y
			p1.Y += signY // Et on fait le déplacement
		}

		// Enfin on revérifie que le point est bien dans les limites
		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
	}
}

// Fonction pour dessiner un rectangle vide
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {

	// On vérifie que les points soit au bon endroits
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

	// On créer les 3 coins du rectangle (+ p1 dans la fonction soit 4 points en tout)
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	// On les relie tous de telle sorte a faire une boucle
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// Fonction pour dessiner un rectangle plain
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {

	// On parcours ligne par ligne notre matrice
	for i := 0; i < height; i++ {

		// point1 et point2 sont de part et d’autre de la longueur du rectangle
		point1 := Point{p1.X, p1.Y + i}
		point2 := Point{p1.X + width, p1.Y + i}
		ppm.DrawLine(point1, point2, color) // Et on les relie
	}
}

// Fonction qui dessine un cercle plein
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	//Création des variables nécéssaires
	x := radius - 1
	y := 0
	dx := 1
	dy := 1
	err := dx - (radius * 2)

	// Tant que le cercle n'es pas entièrement dessiné
	for x > y {
		if x == radius-1 && y == 0 { // Top point
			ppm.Set(center.X+x, center.Y+y+1, color)
		} else {
			ppm.Set(center.X+x, center.Y+y, color)
		}

		if x != y { // Avoid overwriting the top point with the right point
			ppm.Set(center.X+y, center.Y+x, color)
		}

		if x == 0 && y == radius-1 { // Right point
			ppm.Set(center.X-y, center.Y+x-1, color)
		} else {
			ppm.Set(center.X-y, center.Y+x, color)
		}

		ppm.Set(center.X-x, center.Y+y, color)

		if x == -radius+1 && y == 0 { // Bottom point
			ppm.Set(center.X-x, center.Y-y-1, color)
		} else {
			ppm.Set(center.X-x, center.Y-y, color)
		}

		if x != y { // Avoid overwriting the bottom point with the left point
			ppm.Set(center.X-y, center.Y-x, color)
		}

		if x == 0 && y == -radius+1 { // Left point
			ppm.Set(center.X+y, center.Y-x+1, color)
		} else {
			ppm.Set(center.X+y, center.Y-x, color)
		}

		ppm.Set(center.X+x, center.Y-y, color)

		// On adapte err suivant la direction du dessin et on compense l'avancement (Basé aussi sur l'algorithme de Bresenham)
		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (radius * 2)
		}
	}
}

// Fonction qui dessine un cercle rempli
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {

	for i := center.X - radius; i <= center.X+radius; i++ {
		for j := center.Y - radius; j <= center.Y+radius; j++ {
			// Calcul de la distance entre le point parcourant le tour du cercle et le centre du cercle (Formule de distance entre 2 points dans un cercle)
			distance := math.Sqrt((float64(i-center.X) * float64(i-center.X)) + (float64(j-center.Y) * float64(j-center.Y)))
			// Si la distance est inférieure ou égale au rayon du cercle on colorie le pixel
			if distance <= float64(radius) {
				ppm.Set(i, j, color)
			}
		}
	}
}

// Fonction qui dessine un triangle vide
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// Il suffit de relier nos 3 points
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// Fonction qui dessine un triangle plein
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {

}

// Fonction qui dessine un polygone vide
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// On parcours tout les points un par un pour les relier entre eux
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	// On trace la dernière droite du dernier au premier point
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

func main() {
	ppm, _ := ReadPPM("test.ppm")
	color := Pixel{R: 255, G: 0, B: 0}
	point1 := Point{X: 1, Y: 3}
	point2 := Point{X: 4, Y: 2}
	point3 := Point{X: 8, Y: 3}
	ppm.DrawTriangle(point1, point2, point3, color)
	ppm.Save("testsave.ppm")
}
