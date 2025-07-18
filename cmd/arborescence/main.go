/*
* Progamme de lister les fichiers/dossiers
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 17/07/2025
* Version : 1.0.0
*
* Description : Il permet de lister les dossiers / fichiers tout en calculant des différentes informations : tailles, date, etc..
*
 */

package main

import (
	"context"
	"fmt"
	"github/managerDirectory/functions"
	"github/managerDirectory/models"
	"github/managerDirectory/utils"
	"github/managerDirectory/writercsv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	//Debut
	startTime := time.Now()
	fileTime := time.Now().Format("20060102")

	// Vérification des arguments
	//fmt.Printf("la taille des arguments %v\n", len(os.Args))
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [répertoire] [stockage] <unit> <prefix> <extensions> ")
		fmt.Println("Usage: go run . [répertoire] [stockage] <unit> <prefix> <extensions> ")
		os.Exit(1)
	}

	// Obtention des Arguments
	path_in := os.Args[1]
	path_out := os.Args[2]
	unit := os.Args[3]
	prefixe_fichier := os.Args[4]
	ext_file := os.Args[5]

	// Nettoyage plus robuste
	path_in = strings.TrimSpace(path_in)   // Supprime tous les espaces/retours
	path_in = filepath.Clean(path_in)      // Normalise le chemin
	path_out = strings.TrimSpace(path_out) // Supprime tous les espaces/retours
	path_out = filepath.Clean(path_out)    // Normalise le chemin

	unit = strings.TrimSpace(unit)
	prefixe_fichier = strings.TrimSpace(prefixe_fichier)
	ext_file = strings.TrimSpace(ext_file)

	// Vérification que le répertoire existe
	if _, err := os.Stat(path_in); os.IsNotExist(err) {
		fmt.Printf("Erreur: Le répertoire à scannner %s n'existe pas\n", path_in)
		os.Exit(1)
	}

	// Création du dossier de destination
	if _, err := functions.CreateDirs(path_out); err != nil {
		fmt.Printf("Erreur : lors de la création du répertoire de destination: %v\n", err)
		os.Exit(1)
	}

	// Vérification de l'extension à scanner / ou on vide le tableau
	ext_fileSlice := strings.Split(ext_file, ",")
	var extSlice []string
	var extSliceLog string

	for _, item := range ext_fileSlice {
		if item == "." || functions.CompareEF(item, "go") || item == "" {
			extSlice = []string{}
			extSliceLog = "Toutes les extensions trouvées !"
		} else {
			if !strings.HasPrefix(item, ".") {
				item = "." + item
				extSlice = append(extSlice, item)
				extSliceLog = strings.Join(extSlice, ";")
			} else {
				extSlice = append(extSlice, item)
				extSliceLog = strings.Join(extSlice, ";")
			}
		}
	}

	var nameFileLog string
	if prefixe_fichier == "" || functions.CompareEF(prefixe_fichier, "prefixe_fichier") {
		prefixe_fichier = ""
		nameFileLog = fmt.Sprintf("%v_log.log", fileTime)
	} else {
		nameFileLog = fmt.Sprintf("%v_%v_log.log", fileTime, prefixe_fichier)
	}

	unitList := []string{"KB", "MB", "GB"}
	if functions.CompareEF(unit, "unit") || unit == "" || !functions.ContainsTwo(unitList, unit) {
		unit = ""
	}

	//Getion log
	if _, err := os.Stat(filepath.Join(path_out, fileTime)); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(path_out, fileTime), 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	filelog, err := os.OpenFile(filepath.Join(path_out, fileTime, nameFileLog),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//Console
	fmt.Printf("\nParamètres :\n")
	fmt.Printf("- Répertoire à Scanner ==> %s\n", path_in)
	fmt.Printf("- Répertoire de dépôt CSV ==> %s\n", path_out)
	fmt.Printf("- Unité des tailles (fichiers/dossiers) ==> %s\n", unit)
	fmt.Printf("- Préfixe ==> %s\n", prefixe_fichier)
	fmt.Printf("- Extensions (fichiers) ==> %s\n", extSliceLog)
	fmt.Println("=============================================")

	//Log
	log.SetOutput(filelog)
	log.Println("Début du programme")
	log.Printf("- Répertoire à Scanner ==> %s\n", path_in)
	log.Printf("- Répertoire de dépôt CSV ==> %s\n", path_out)
	log.Printf("- Unité des tailles (fichiers/dossiers) ==> %s\n", unit)
	log.Printf("- Préfixe ==> %s\n", prefixe_fichier)
	log.Printf("- Extensions (fichiers) ==> %s\n", extSliceLog)
	time.Sleep(2 * time.Second)

	// Configuration
	workers := runtime.NumCPU() * 2 // Nombre de workers parallèles
	resultChanSize := 10000         // Taille du buffer résultats
	fileChanSize := 10000           // Taille du buffer fichiers

	//Canaux pour colleter les résultats (channels)
	fileChan := make(chan string, fileChanSize)
	resultChan := make(chan models.FilesStruct, resultChanSize)
	doneChan := make(chan struct{}) // Canal de synchronisation

	// Context pour la gestion d'annulation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Goroutine d'écriture CSV
	go func() {
		writercsv.WriteCsvChannelThree(path_out, prefixe_fichier, resultChan)
		close(doneChan) // Signal de fin d'écriture
	}()

	// Lancer les workers && Goroutine pour calculer et obtenir les infos pour remplir resultChan
	var getInfoWg sync.WaitGroup // Crée un WaitGroup pour les goroutines getDirFilesInfosCustomChannel
	getInfoWg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer getInfoWg.Done()
			utils.GetDirFilesInfosCustomChannelThree(ctx, fileChan, resultChan, unit)
		}()
	}

	//Parcourir l'arborescence et envoyer les fichiers aux workers/ fileChan
	go utils.ListDirsFilesChannelThree(ctx, path_in, extSlice, fileChan)

	// Attend que toutes les goroutines de getDirFilesInfosCustomChannel aient terminé
	getInfoWg.Wait()  // Attendre fin des workers
	close(resultChan) // Fermer le canal des résultats

	// Attendre fin de l'écriture CSV
	<-doneChan

	//Fin
	lastTime := time.Now()
	duree := lastTime.Sub(startTime)
	log.Println("Fin du programme")
	fmt.Printf("\nLa durée du traitement est %v", functions.CalDuration(duree))
	log.Printf("La durée du traitement est %v", functions.CalDuration(duree))
	log.Println("=============================================")
}
