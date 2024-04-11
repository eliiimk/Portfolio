package main

import (
	"bufio"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type HangmanGame struct {
	WordToGuess    string
	DisplayedWord  string
	GuessedLetters []string
	AttemptsLeft   int
	HangmanDrawing string
}

type Score struct {
	Username string
	Score    int
}

type AdminData struct {
	Scores []Score
}

const port = ":8000"

var game HangmanGame

func main() {
	initGame(1)

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/hangman", hangmanHandler)
	http.HandleFunc("/admin", adminHandler)


	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	fmt.Println("(http://localhost:8000) - server start on port", port)
	http.ListenAndServe(":8000", nil)
}

func initGame(difficulty int) {
	var wordsFile string

	switch difficulty {
	case 1:
		wordsFile = "words/facile.txt"
	case 2:
		wordsFile = "words/moyen.txt"
	case 3:
		wordsFile = "words/difficile.txt"
	default:
		fmt.Println("Invalid difficulty level")
		return
	}

	word, err := getRandomWord(wordsFile)
	if err != nil {
		fmt.Println("Error loading word:", err)
		return
	}

	game = HangmanGame{
		WordToGuess:    word,
		DisplayedWord:  strings.Repeat("_ ", len(word)),
		GuessedLetters: []string{},
		AttemptsLeft:   10,
	}
}

func getRandomWord(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if len(words) == 0 {
		return "", fmt.Errorf("no words found in the file")
	}

	return words[randomInt(len(words))], nil
}

func randomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", game)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println("Username:", username)
		initGame(1)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "start", nil)
}
func hangmanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		letter := r.FormValue("letter")
		if strings.Contains(game.WordToGuess, letter) {
			game.DisplayedWord = updateDisplayedWord(game.WordToGuess, game.DisplayedWord, letter)
		} else {
			game.AttemptsLeft--
			game.HangmanDrawing = drawHangman(game.AttemptsLeft)
		}

		game.GuessedLetters = append(game.GuessedLetters, letter)

		if game.DisplayedWord == game.WordToGuess {
			renderTemplate(w, "gagner", game)
			return
		} else if game.AttemptsLeft == 0 {
			renderTemplate(w, "perdue", game)
			return
		}
	}

	renderTemplate(w, "index", game)

	game.HangmanDrawing = drawHangman(game.AttemptsLeft)
}

func updateDisplayedWord(wordToGuess, displayedWord, letter string) string {
	var updatedWord strings.Builder

	for i, char := range wordToGuess {
		if string(char) == letter {
			updatedWord.WriteString(letter)
		} else {
			updatedWord.WriteString(string(displayedWord[i]))
		}
	}

	return updatedWord.String()
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := fmt.Sprintf("templates/%s.html", tmpl)
	if err := template.Must(template.ParseFiles(tmplPath)).Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	scores := getFakeScores()
	adminData := AdminData{
		Scores: scores,
	}
	renderTemplate(w, "admin", adminData)
}

func getFakeScores() []Score {
	return []Score{
		{Username: "Joueur1", Score: 100},
		{Username: "Joueur2", Score: 80},
		{Username: "Joueur3", Score: 70},
		{Username: "Joueur4", Score: 90},
		{Username: "Joueur5", Score: 85},
		{Username: "Joueur6", Score: 95},
		{Username: "Joueur7", Score: 75},
	}
}

func drawHangman(attemptsLeft int) string {
	var hangmanDrawing string

	switch attemptsLeft {
	case 9:
		hangmanDrawing = `
    
       
        
     
        
        
=========`
	case 8:
		hangmanDrawing = `
    
      |  
      |  
      |  
      |  
      |  
=========`
	case 7:
		hangmanDrawing = `
  +---+  
      |  
      |  
      |  
      |  
      |  
=========`
	case 6:
		hangmanDrawing = `
  +---+  
  |   |  
      |  
      |  
      |  
      |  
=========`
	case 5:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
      |  
      |  
      |  
=========`
	case 4:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
  |   |  
      |  
      |  
=========`
	case 3:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
 /|   |  
      |  
      |  
=========`
	case 2:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
 /|\  |  
      |  
      |  
=========`
	case 1:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
 /|\  |  
 /    |  
      |  
=========`
	case 0:
		hangmanDrawing = `
  +---+  
  |   |  
  O   |  
 /|\  |  
 / \  |  
      |  
=========`
	default:
		hangmanDrawing = "Hangman Complete"
	}

	return hangmanDrawing
}
