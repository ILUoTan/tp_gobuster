package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Fonction pour lire un fichier ligne par ligne
func lireLignes(chemin string) ([]string, error) {
	fichier, err := os.Open(chemin)
	if err != nil {
		return nil, err
	}
	defer fichier.Close()

	var lignes []string
	scanner := bufio.NewScanner(fichier)
	for scanner.Scan() {
		lignes = append(lignes, strings.TrimSpace(scanner.Text()))
	}
	return lignes, scanner.Err()
}

// Fonction qui envoie un requête HTTP pour tester un chemin donné
func envoyerRequete(cible string, chemin string) (int, error) {
	url := fmt.Sprintf("http://%s/%s", cible, chemin)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// Fonction principale du scan
func scanner(cible string, chemins []string, travailleurs int, silencieux bool) {
	var wg sync.WaitGroup
	travail := make(chan string)
	resultats := make(chan string)

	debut := time.Now()

	// Création des travailleur
	for i := 0; i < travailleurs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chemin := range travail {
				statut, err := envoyerRequete(cible, chemin)
				if err != nil {
					continue
				}
				if statut == 200 || !silencieux {
					resultats <- fmt.Sprintf("/%s\t%d", chemin, statut)
				}
			}
		}()
	}

	// Envoi des chemins aux travailleurs
	go func() {
		for _, chemin := range chemins {
			travail <- chemin
		}
		close(travail)
	}()

	// Fermeture du canal résultats après la fin des travailleurs
	go func() {
		wg.Wait()
		close(resultats)
	}()

	fmt.Println("Starting scan...")
	for resultat := range resultats {
		fmt.Println(resultat)
	}

	duree := time.Since(debut)
	fmt.Printf("Scan done in %.6fs\n", duree.Seconds())
}

func main() {
	dictionnaire := flag.String("d", "", "Path to dictionary file")
	silencieux := flag.Bool("q", false, "Quiet mode, only show HTTP 200 results")
	cible := flag.String("t", "", "Target to enumerate (including port)")
	travailleurs := flag.Int("w", 1, "Number of workers to run")
	flag.Parse()

	if *dictionnaire == "" || *cible == "" {
		fmt.Println("Usage of mygb:")
		flag.PrintDefaults()
		return
	}

	chemins, err := lireLignes(*dictionnaire)
	if err != nil {
		fmt.Printf("Échec de la lecture du dictionnaire : %v\n", err)
		return
	}

	fmt.Println("Démarrage de MyGB")
	fmt.Println("--")
	fmt.Printf("Target: http://%s\nList: %s\nWorkers: %d\n", *cible, *dictionnaire, *travailleurs)
	fmt.Println("--")

	scanner(*cible, chemins, *travailleurs, *silencieux)
}
