/*
* Progamme de renommage d'extensions de fichiers
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 17/07/2025
* Version : 1.0.0
*
* Description :
*   Ce script permet de rechercher et renommer massivement
*   les extensions de fichiers de manière interactive.
*
 */

package main

import (
	"fmt"
	"github/managerDirectory/extension"
	"github/managerDirectory/functions"
	"github/managerDirectory/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	//Debut
	startTime := time.Now()

	//Gestion log
	filelog, _ := functions.DefinieLog("changeExtension")
	log.SetOutput(filelog)
	log.Println("Début du programme")

	path, err := extension.GetPath("Veuillez entrer le chemin du dossier : ")
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}

	ScanExt, err := extension.GetExtentionFromAsk("Veuillez entrer l'extension à scanner : ")
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}

	NewExt, err := extension.GetExtentionFromAsk("Veuillez entrer la nouvelle extension : ")
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}

	if path == "" || ScanExt == "" || NewExt == "" {
		fmt.Println("Merci de renseigner les paramètres demandées")
		log.Println("Merci de renseigner les paramètres demandées")
		os.Exit(1)
	}
	// Nettoyage plus robuste
	path = strings.TrimSpace(path) // Supprime tous les espaces/retours
	path = filepath.Clean(path)    // Normalise le chemin

	fmt.Printf("\nParamètre :\n")
	fmt.Printf("- Dossier sélectionné :%s\n", path)
	fmt.Printf("- Extension à Scanner : %s\n", ScanExt)
	fmt.Printf("- Extension à utiliser: %s\n", NewExt)

	log.Printf("Paramètre :")
	log.Printf("- Dossier sélectionné :%s\n", path)
	log.Printf("- Extension à Scanner : %s\n", ScanExt)
	log.Printf("- Extension à utiliser: %s\n", NewExt)

	files := functions.ListDirsFiles(path, ScanExt) //  ".mp4" ,".part"
	for _, path := range files {
		// fmt.Println(path)
		utils.ListFilesInfo(path)
		//extension.changerExtension(path, NewExt)
	}

	//Fin
	lastTime := time.Now()
	duree := lastTime.Sub(startTime)
	log.Println("Fin du programme")
	fmt.Printf("\nLa durée du traitement est %v", functions.CalDuration(duree))
	log.Printf("La durée du traitement est %v", functions.CalDuration(duree))
	log.Println("=============================================")
}
