/*
* Progamme de renommage d'extensions de fichiers
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 17/07/2025
* Version : 1.0.0
*
* Description :
*   Ce script permet de rechercher et de supprimer les dossiers / fichiers selon le mot clé en argument
*
 */
package main

import (
	"bufio"
	"fmt"
	"github/managerDirectory/functions"
	"github/managerDirectory/scanner"
	"os"
	"path/filepath"
	"strings"
)

func scanner_files_filter_remove(dir string) {
	//Demande mot ou expression à rechercher
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Veuillez entrer le mot / expression à scanner : ")
	ext_file, _ := reader.ReadString('\n')
	ext_file = strings.TrimSpace(ext_file)

	// Nettoyage plus robuste
	dir = strings.TrimSpace(dir) // Supprime tous les espaces/retours
	dir = filepath.Clean(dir)    // Normalise le chemin

	fmt.Printf("\nParamètre :\n")
	fmt.Printf("- Dossier sélectionné :\n%s\n", dir)
	fmt.Printf("- Mot ou Expression : %s\n", ext_file)

	files := functions.ListFilesFilter(dir, ext_file)
	for _, v := range files {
		scanner.ListDirFilesInfo(v)
		// changerExtension(v, "")
		err := functions.DeleteIfExists(v)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

}

func main() {
	path, err := scanner.GetPath()
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}
	scanner_files_filter_remove(path)
}
