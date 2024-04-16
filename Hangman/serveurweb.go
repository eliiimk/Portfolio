package main

import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
)

var hangmanProcess *exec.Cmd // Stocke le processus du jeu Hangman Go

func handlePlayHangman(w http.ResponseWriter, r *http.Request) {
    if hangmanProcess == nil {
        fmt.Fprintln(w, "Le jeu Hangman n'est pas encore démarré.")
        return
    }

    out, err := hangmanProcess.Output()
    if err != nil {
        fmt.Fprintf(w, "Erreur lors de l'exécution du jeu Hangman : %s\n", err)
        return
    }

    fmt.Fprintf(w, "Résultat du jeu Hangman :\n%s\n", out)
}

func main() {
    // Démarrer le jeu Hangman Go en arrière-plan
    hangmanProcess = exec.Command("go", "run", "main.go")
    if err := hangmanProcess.Start(); err != nil {
        log.Fatalf("Erreur lors du démarrage du jeu Hangman : %s\n", err)
    }
    defer hangmanProcess.Process.Kill() // Arrêter le jeu lorsque le serveur se termine

    http.HandleFunc("/play-hangman", handlePlayHangman)

    fmt.Println("Serveur démarré sur le port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
