package main
import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
)

func handlePlayHangman(w http.ResponseWriter, r *http.Request) {
    // Exécutez votre jeu Hangman Go
    cmd := exec.Command("go", "run", "serveurweb.go")
    out, err := cmd.Output()
    if err != nil {
        fmt.Fprintf(w, "Erreur lors de l'exécution du jeu Hangman : %s\n", err)
        return
    }

    fmt.Fprintf(w, "Résultat du jeu Hangman :\n%s\n", out)
}

func main() {
    http.HandleFunc("/play-hangman", handlePlayHangman)

    fmt.Println("Serveur démarré sur le port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
