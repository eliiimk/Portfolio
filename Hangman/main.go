package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "time"
)

var mots = []string{
    "mot1",
    "mot2", 
    "mot3",
    // Ajoutez d'autres mots ici
}

func Asciiprinter(Word []rune) {
    for hauteur := 2; hauteur <= 10; hauteur++ {
        for _, letter := range Word {
            Showletter(letter, hauteur)
        }
        fmt.Println("")
    }
}

var numero int

func Showletter(letter rune, line int) {
    file, err := os.Open("asciiArt.txt")
    if err != nil {
        fmt.Println("Erreur lors de l'ouverture du fichier :", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    lineCount := 1
    startLine := 0

    numero = int(letter) - 32

    startLine = (numero * 9) + line

    for scanner.Scan() {
        lineCount++
        if lineCount == startLine {
            fmt.Print(scanner.Text())
        }
    }

    if scanner.Err() != nil {
        fmt.Println("Erreur lors de la lecture du fichier :", scanner.Err())
    }
}

func choisirMot() string {
    mots, err := chargerMotsDepuisFichier("Words.txt")
    if err != nil {
        fmt.Printf("Erreur lors de la lecture du fichier de mots : %v\n", err)
        os.Exit(1)
    }

    rand.Seed(time.Now().Unix())
    index := rand.Intn(len(mots))
    return mots[index]
}
func chargerMotsDepuisFichier(nomFichier string) ([]string, error) {
    fichier, err := os.Open(nomFichier)
    if err != nil {
        return nil, err
    }
    defer fichier.Close()

    mots := make([]string, 0)
    scanner := bufio.NewScanner(fichier)
    for scanner.Scan() {
        mot := scanner.Text()
        mots = append(mots, mot)
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return mots, nil
}

func afficherJose(positions []string, position int, afficher bool) {
    if afficher && position < len(positions) {
        fmt.Println(positions[position])
    }
}
func chargerPositionsDepuisFichier(nomFichier string) ([]string, error) {
    fichier, err := os.Open(nomFichier)
    if err != nil {
        return nil, err
    }
    defer fichier.Close()

    positions := make([]string, 0)
    scanner := bufio.NewScanner(fichier)
    positionActuelle := ""
    for scanner.Scan() {
        ligne := scanner.Text()
        if ligne == "=========" {
            positions = append(positions, positionActuelle)
            positionActuelle = ""
        } else {
            positionActuelle += ligne + "\n"
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return positions, nil
}
func main() {
    positions, err := chargerPositionsDepuisFichier("hangman.txt")
    if err != nil {
        fmt.Printf("Erreur lors de la lecture du fichier de positions : %v\n", err)
        os.Exit(1)
    }
    motADeviner := choisirMot()
    longueurMot := len(motADeviner)
    lettresRevelees := make([]rune, longueurMot)
    tentativesRestantes := 10

    for i := range lettresRevelees {
        lettresRevelees[i] = '_'
    }

    fmt.Println("Bienvenue au jeu du Pendu!")
    fmt.Printf("Le mot à deviner contient %d lettres.\n", longueurMot)

    for tentativesRestantes > 0 {

        fmt.Printf("Tentatives restantes : %d\n", tentativesRestantes)
        Asciiprinter(lettresRevelees)
        fmt.Print("Devinez une lettre : ")

        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        lettre := []rune(scanner.Text())[0]

        lettreCorrecte := false
        for i, char := range motADeviner {
            if lettre == char {
                lettresRevelees[i] = lettre
                lettreCorrecte = true
            }
        }

        if !lettreCorrecte {
            fmt.Println("La lettre n'est pas dans le mot.")
            tentativesRestantes--
            afficherJose(positions, 10-tentativesRestantes, true)
        }

        if string(lettresRevelees) == motADeviner {
            fmt.Printf("Félicitations, vous avez deviné le mot : %s\n", motADeviner)
            break
        }
    }

    if string(lettresRevelees) != motADeviner {
        fmt.Printf("Désolé, vous avez épuisé toutes vos tentatives. Le mot était : %s\n", motADeviner)
 
    }
}
