/*
* Progamme de création des classeurs CSV
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 17/07/2025
* Version : 1.0.0
*
* Description :
*   Ce script permet de créer un classeur et de sauvegarder les données
 */
package writercsv

import (
	"encoding/csv"
	"fmt"
	"github/managerDirectory/functions"
	"github/managerDirectory/models"
	"github/managerDirectory/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func WriteCsv(path_out string, files []string) {
	if len(files) > 0 {
		//today to string
		fileTime := time.Now().Format("20060102")

		// Créer un nouveau fichier ou ouvrir un fichier existant
		dstPath := filepath.Join(path_out, fmt.Sprintf("%v_listes.csv", fileTime))
		file, err := os.Create(dstPath)
		if err != nil {
			log.Fatalf("Erreur lors de la création du fichier : %v", err)
		}
		defer file.Close()

		// Créer un nouvel écrivain CSV qui écrit dans le fichier
		writer := csv.NewWriter(file)
		writer.Comma = ';'
		defer writer.Flush()

		// Définir l'encodage UTF-8
		file.WriteString("\xEF\xBB\xBF") // Ajouter le BOM (Byte Order Mark) pour UTF-8

		// lister les fichiers et dossiers
		filesInfos, _ := utils.GetDirFilesInfo(files, "")
		fmt.Printf("\nL'ensemble des %v fichiers/dossiers compte %v informations", len(files), len(filesInfos))
		time.Sleep(5 * time.Second)

		// Écriture des en-têtes (noms des champs de la struct)
		if err := writer.Write(functions.GetStructFields(filesInfos[0])); err != nil {
			log.Fatal("Erreur écriture en-têtes:", err)
		}

		// Écriture des données
		total := len(files)
		for key, file := range filesInfos {
			fmt.Printf("key : %v sur : %v\n", key+1, total)
			if err := writer.Write(functions.StructToSlice(file)); err != nil {
				log.Fatal("Erreur lors de l'écriture des données:", err)
			}
		}
		log.Println("Fichier CSV créé avec succès")
	} else {
		log.Println("Non Trouver, pas de CSV crée!")
	}
}

func WriteCsvTwo(path_in, path_out string, ext_file []string, unit string) {
	// lister les fichiers et dossiers
	filesInfos := utils.ListDirsFilesCustom(path_in, ext_file, unit)
	fmt.Printf("\nLe repertoire %v compte %v fichiers\n", path_in, len(filesInfos))
	log.Printf("Le repertoire %v compte %v fichiers\n", path_in, len(filesInfos))
	time.Sleep(5 * time.Second)

	if len(filesInfos) > 0 {
		//today to string
		fileTime := time.Now().Format("20060102")

		// Créer un nouveau fichier ou ouvrir un fichier existant
		dstPath := filepath.Join(path_out, fmt.Sprintf("%v_listes.csv", fileTime))
		file, err := os.Create(dstPath)
		if err != nil {
			log.Fatalf("Erreur lors de la création du fichier : %v", err)
		}
		defer file.Close()

		// Créer un nouvel écrivain CSV qui écrit dans le fichier
		writer := csv.NewWriter(file)
		writer.Comma = ';'
		defer writer.Flush()

		// Définir l'encodage UTF-8
		file.WriteString("\xEF\xBB\xBF") // Ajouter le BOM (Byte Order Mark) pour UTF-8

		// Écriture des en-têtes (noms des champs de la struct)
		if err := writer.Write(functions.GetStructFields(filesInfos[0])); err != nil {
			log.Fatal("Erreur écriture en-têtes:", err)
		}

		// Écriture des données
		total := len(filesInfos)
		for key, file := range filesInfos {
			fmt.Printf("key : %v sur : %v\n", key+1, total)
			if err := writer.Write(functions.StructToSlice(file)); err != nil {
				log.Fatal("Erreur lors de l'écriture des données:", err)
			}
		}
		log.Println("Fichier CSV créé avec succès")
	} else {
		log.Println("Non Trouver, pas de CSV crée!")
	}
}

func WriteCsvChannel(path_out string, resultChan <-chan models.FilesStruct, doneChan chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Début Ecriture dans Fichier CSV")
	defer fmt.Println("Fin Ecriture dans Fichier CSV")

	//today to string
	fileTime := time.Now().Format("20060102")

	// Créer un nouveau fichier ou ouvrir un fichier existant
	dstPath := filepath.Join(path_out, fmt.Sprintf("%v_listes.csv", fileTime))
	file, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	// Créer un nouvel écrivain CSV qui écrit dans le fichier
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Définir l'encodage UTF-8
	file.WriteString("\xEF\xBB\xBF") // Ajouter le BOM (Byte Order Mark) pour UTF-8

	// Écriture des en-têtes (noms des champs de la struct)
	if err := writer.Write(functions.GetFieldsNamesStruct()); err != nil {
		log.Fatal("Erreur écriture en-têtes:", err)
	}

	// Recevoir les résultats & Écriture des données
	//total := len(resultChan)
	for fileInfo := range resultChan {
		fmt.Printf("- len : %v - cap : %v -resultChan  : %v\n", len(resultChan), cap(resultChan), fileInfo)
		if err := writer.Write(functions.StructToSlice(fileInfo)); err != nil {
			log.Fatal("Erreur lors de l'écriture des données:", err)
		}
	}
	log.Println("Fichier CSV créé avec succès")

	// Signale que l'écriture dans le fichier CSV est terminée
	doneChan <- true
}

func WriteCsvChannelTwo(path_out string, resultChan <-chan models.FilesStruct, doneChan chan<- bool) {
	fmt.Println("Début Ecriture dans Fichier CSV")
	defer fmt.Println("Fin Ecriture dans Fichier CSV")

	//today to string
	fileTime := time.Now().Format("20060102")

	// Créer un nouveau fichier ou ouvrir un fichier existant
	dstPath := filepath.Join(path_out, fmt.Sprintf("%v_listes.csv", fileTime))
	file, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	// Créer un nouvel écrivain CSV qui écrit dans le fichier
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Définir l'encodage UTF-8
	file.WriteString("\xEF\xBB\xBF") // Ajouter le BOM (Byte Order Mark) pour UTF-8

	// Écriture des en-têtes (noms des champs de la struct)
	if err := writer.Write(functions.GetFieldsNamesStruct()); err != nil {
		log.Fatal("Erreur écriture en-têtes:", err)
	}

	// Utilise select pour recevoir les résultats du canal resultChan
	// Recevoir les résultats & Écriture des données
	for {
		select {
		case fileInfo, ok := <-resultChan:
			if !ok {
				// Le canal est fermé, on sort de la boucle
				doneChan <- true
				return
			}
			fmt.Printf("- len : %v - cap : %v -resultChan  : %v\n", len(resultChan), cap(resultChan), fileInfo)
			if err := writer.Write(functions.StructToSlice(fileInfo)); err != nil {
				log.Fatal("Erreur lors de l'écriture des données:", err)
				return
			}
			log.Println("Fichier CSV créé avec succès")
		}
	}
}

func WriteCsvChannelThree(path_out, prefixe_fichier string, resultChan <-chan models.FilesStruct) {
	//Fermeture
	defer fmt.Println("Fin Ecriture dans Fichier CSV")
	defer log.Println("Fin Ecriture dans Fichier CSV")

	// Ouverture
	fmt.Println("Début Ecriture dans Fichier CSV")
	log.Println("Début Ecriture dans Fichier CSV")

	//today to string
	fileTime := time.Now().Format("20060102")

	// Créer un nouveau fichier ou ouvrir un fichier existant
	var dstPath string
	if len(prefixe_fichier) > 0 {
		dstPath = filepath.Join(path_out, fileTime, fmt.Sprintf("%v_%v_listes.csv", fileTime, prefixe_fichier))
	} else {
		dstPath = filepath.Join(path_out, fileTime, fmt.Sprintf("%v_listes.csv", fileTime))
	}
	file, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	// Créer un nouvel écrivain CSV qui écrit dans le fichier
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Définir l'encodage UTF-8
	file.WriteString("\xEF\xBB\xBF") // Ajouter le BOM (Byte Order Mark) pour UTF-8

	// Écriture des en-têtes (noms des champs de la struct)
	if err := writer.Write(functions.GetFieldsNamesStruct()); err != nil {
		log.Fatal("Erreur écriture en-têtes:", err)
	}

	// Recevoir les résultats & Écriture des données
	for fileInfo := range resultChan {
		//fmt.Printf("- len : %v - cap : %v -resultChan  : %v\n", len(resultChan), cap(resultChan), fileInfo)
		if err := writer.Write(functions.StructToSlice(fileInfo)); err != nil {
			log.Fatal("Erreur lors de l'écriture des données:", err)
		}
	}
	log.Println("Fichier CSV créé avec succès")
}
