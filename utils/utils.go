package utils

import (
	"bufio"
	"context"
	"fmt"
	"github/managerDirectory/functions"
	"github/managerDirectory/models"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Afficher les infos d'un fichier
func ListFilesInfo(filePath string) {
	// Nettoyage plus robuste
	filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
	filePath = filepath.Clean(filePath)    // Normalise le chemin
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
		} else {
			fmt.Printf("Erreur d'accès : %v\n", err)
		}
		return
	}
	if fileInfo.IsDir() {
		fmt.Printf("Erreur : '%s' est un dossier, pas un fichier.\n", filePath)
		return
	}

	fmt.Printf("\nInfos fichier :\n")
	fmt.Printf("- Chemin : %s\n", filepath.ToSlash(filePath)) // Standardise les slashs
	fmt.Printf("- Taille : %d octets\n", fileInfo.Size())
	fmt.Printf("- Taille Convertie : %d octets\n", functions.ConvertUnitSize(int(fileInfo.Size()), "MB"))
	fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
}

// Afficher les infos d'un fichier / dossier
func ListDirFilesInfo(filePath string, unit string) {
	// Nettoyage plus robuste
	filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
	filePath = filepath.Clean(filePath)    // Normalise le chemin
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
		} else {
			fmt.Printf("Erreur d'accès : %v\n", err)
		}
		return
	}

	if fileInfo.IsDir() {
		fmt.Printf("\nInfos Dossier :\n")
		fmt.Printf("- Nom Dossier : %s\n", fileInfo.Name())
	} else {
		fmt.Printf("\nInfos fichier :\n")
		fmt.Printf("- Nom Fichier : %s\n", fileInfo.Name())
	}

	fmt.Printf("- Chemin : %s\n", filepath.ToSlash(filePath)) // Standardise les slashs
	fmt.Printf("- Permission : %s\n", fileInfo.Mode())
	fmt.Printf("- Taille : %d octets\n", fileInfo.Size())
	if len(unit) == 0 {
		fmt.Printf("- Taille Convertie : %d %s\n", functions.ConvertUnitSize(int(fileInfo.Size()), "MB"), "MB")
	} else {
		fmt.Printf("- Taille Convertie : %d %s\n", functions.ConvertUnitSize(int(fileInfo.Size()), unit), unit)
	}

	fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
}

// Poser une question à l'utilisateur pour avoir le Path
func GetPath(text string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%v", text)
		dossierPath, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("erreur de lecture : %v", err)
		}
		// Nettoyage du chemin
		dossierPath = strings.TrimSpace(dossierPath)
		dossierPath = strings.TrimRight(dossierPath, "\r\n")     // Gestion des retours chariot Windows
		dossierPath = strings.ReplaceAll(dossierPath, `\`, "\\") // Normalise les séparateurs
		dossierPath = filepath.Clean(dossierPath)

		// Vérification que c'est bien un dossier
		info, err := os.Stat(dossierPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Le dossier '%s' n'existe pas. Veuillez réessayer.\n", dossierPath)
				continue
			}
			return "", fmt.Errorf("erreur d'accès : %v", err)
		}

		if !info.IsDir() {
			fmt.Printf("'%s' n'est pas un dossier valide. Veuillez réessayer.\n", dossierPath)
			continue
		}
		// Conversion en chemin absolu
		absPath, err := filepath.Abs(dossierPath)
		if err != nil {
			return "", fmt.Errorf("erreur de conversion en chemin absolu : %v", err)
		}
		return absPath, nil
	}
}

// Poser une question à l'utilisateur pour avoir l'extension du ficher
func GetExtentionFromAsk(text string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%v", text)
	ext_file, _ := reader.ReadString('\n')
	ext_file = strings.TrimSpace(ext_file)

	if !strings.HasPrefix(ext_file, ".") {
		ext_file = "." + ext_file
	}

	if ext_file == "." {
		ext_file = ""
	}
	return ext_file, nil
}

// Changer les extension d'un fichier
func ChangerExtension(path string, nouvelleExt string) (string, error) {
	var newPath string
	// Séparation de l'extension
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]

	// Formatage de la nouvelle extension
	if nouvelleExt != "" && !strings.HasPrefix(nouvelleExt, ".") {
		nouvelleExt = "." + nouvelleExt
	}

	// Renommage avec gestion des attributs Windows
	if nouvelleExt != ext && nouvelleExt != "" {
		newPath := base + nouvelleExt
		err := syscall.Rename(path, newPath)
		if err != nil {
			return "", fmt.Errorf("échec du renommage Windows: %v", err)
		}
	} else {
		newPath := base
		err := syscall.Rename(path, newPath)
		if err != nil {
			return "", fmt.Errorf("échec du renommage Windows: %v", err)
		}
	}

	fmt.Printf("\nChangement Extension du Fichier :\n")
	fmt.Printf("- Chemin : %s\n", path)
	fmt.Printf("- Sans extension : %s\n", base)
	fmt.Printf("- Ancienne extension : %s\n", ext)
	fmt.Printf("- Nouvelle extension : %s\n", nouvelleExt)
	// fmt.Printf("- Resultat : %s\n", newPath)
	return newPath, nil
}

// Récuperation des informations en tableau selon un type personnalisé
func GetDirFilesInfo(files []string, unit string) ([]models.FilesStruct, error) {
	var filesInfos []models.FilesStruct
	for _, filePath := range files {
		// Nettoyage plus robuste
		filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
		filePath = filepath.Clean(filePath)    // Normalise le chemin
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
			} else {
				fmt.Printf("Erreur d'accès : %v\n", err)
			}
		} else {
			filesInfos = append(filesInfos, models.FilesStruct{
				Infos:           functions.CheckDirectory(fileInfo.IsDir()),
				Nom:             fileInfo.Name(),
				Chemin:          filepath.ToSlash(filePath),
				Extension:       functions.GetExtensionFiles(fileInfo.IsDir(), filePath),
				Permission:      fileInfo.Mode().String(),
				Octal:           fmt.Sprintf("%04o", fileInfo.Mode().Perm()),
				Type:            functions.GetFileType(fileInfo.Mode()),
				Taille:          functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size()),
				TailleConvertie: functions.ConvertUnitSize(int(functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())), unit),
				Modification:    fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}

	return filesInfos, nil

}

// Récuperation des informations en tableau clé-valeur selon un type personnalisé
func GetMapDirFilesInfo(files []string, unit string) (map[string]*models.FilesStruct, error) {
	filesInfos := make(map[string]*models.FilesStruct)
	for _, filePath := range files {
		// Nettoyage plus robuste
		filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
		filePath = filepath.Clean(filePath)    // Normalise le chemin
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
			} else {
				fmt.Printf("Erreur d'accès : %v\n", err)
			}
		} else {
			Infos := functions.CheckDirectory(fileInfo.IsDir())
			Nom := fileInfo.Name()
			Chemin := filepath.ToSlash(filePath)
			Extension := functions.GetExtensionFiles(fileInfo.IsDir(), filePath)
			Permission := fileInfo.Mode().String()
			Octal := fmt.Sprintf("%04o", fileInfo.Mode().Perm())
			Type := functions.GetFileType(fileInfo.Mode())
			Taille := functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())
			TailleConvertie := functions.ConvertUnitSize(int(functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())), unit)
			Modification := fileInfo.ModTime().Format("2006-01-02 15:04:05")

			filesInfos[filePath] = models.NewFilesStruct(Infos, Nom, Chemin, Extension, Permission, Octal, Type, Modification, Taille, TailleConvertie)
		}
	}

	return filesInfos, nil

}

// Récuperation des informations type personnalisé
func GetDirFilesInfosCustom(filePath string, unit string) (*models.FilesStruct, error) {
	var filesInfos *models.FilesStruct
	// Nettoyage plus robuste
	filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
	filePath = filepath.Clean(filePath)    // Normalise le chemin
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
		} else {
			fmt.Printf("Erreur d'accès : %v\n", err)
		}
	} else {
		filesInfos = &models.FilesStruct{
			Infos:           functions.CheckDirectory(fileInfo.IsDir()),
			Nom:             fileInfo.Name(),
			Chemin:          filepath.ToSlash(filePath),
			Extension:       functions.GetExtensionFiles(fileInfo.IsDir(), filePath),
			Permission:      fileInfo.Mode().String(),
			Octal:           fmt.Sprintf("%04o", fileInfo.Mode().Perm()),
			Type:            functions.GetFileType(fileInfo.Mode()),
			Taille:          functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size()),
			TailleConvertie: functions.ConvertUnitSize(int(functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())), unit),
			Modification:    fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		}
		return filesInfos, nil
	}
	return filesInfos, nil
}

// Scanner les repertoires pour lister fichiers /dossiers avec les informations en tableau selon type personnalisé
func ListDirsFilesCustom(dir string, extension []string, unit string) []models.FilesStruct {
	var files []models.FilesStruct
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//extension d'un fichier (pas de dossier)
		if len(extension) > 0 {
			if !d.IsDir() && functions.Contains(extension, filepath.Ext(path)) {
				info, _ := GetDirFilesInfosCustom(path, unit)
				files = append(files, *info)
			}
		} else {
			// tous les fichiers et les dossiers
			if !d.IsDir() {
				info, _ := GetDirFilesInfosCustom(path, unit)
				files = append(files, *info)
			} else {
				info, _ := GetDirFilesInfosCustom(path, unit)
				files = append(files, *info)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

// Récuperation des informations type personnalisé par loop sur le channel fileChan
func GetDirFilesInfosCustomChannel(filePath string, resultChan chan<- models.FilesStruct, unit string) {
	// fmt.Println("Début Récuperation information des fichiers /dossiers")
	// defer fmt.Println("Fin Récuperation information des fichiers /dossiers")

	// Nettoyage plus robuste
	filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
	filePath = filepath.Clean(filePath)    // Normalise le chemin
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
		} else {
			fmt.Printf("Erreur d'accès : %v\n", err)
		}
	} else {
		// Envoyer les resultats dans le canal pour collecter les résultats (resultChan)
		resultChan <- models.FilesStruct{
			Infos:           functions.CheckDirectory(fileInfo.IsDir()),
			Nom:             fileInfo.Name(),
			Chemin:          filepath.ToSlash(filePath),
			Extension:       functions.GetExtensionFiles(fileInfo.IsDir(), filePath),
			Permission:      fileInfo.Mode().String(),
			Octal:           fmt.Sprintf("%04o", fileInfo.Mode().Perm()),
			Type:            functions.GetFileType(fileInfo.Mode()),
			Taille:          functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size()),
			TailleConvertie: functions.ConvertUnitSize(int(functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())), unit),
			Modification:    fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		}
	}
	//close(resultChan)
}

// Scanner les repertoires pour lister fichiers /dossiers avec les informations en tableau selon type personnalisé
func ListDirsFilesChannel(dir string, extension []string, fileChan chan<- string) {
	fmt.Println("Début Scanning : Repertoire")
	defer fmt.Println("Fin Scanning : Repertoire")
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//extension d'un fichier (pas de dossier)
		if len(extension) > 0 {
			if !d.IsDir() && functions.Contains(extension, filepath.Ext(path)) {
				fileChan <- path
			}
		} else {
			// tous les fichiers et les dossiers
			if !d.IsDir() {
				fileChan <- path
			} else {
				fileChan <- path
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	// Ferme le canal fileChan après avoir terminé l'envoi des chemins de fichiers
	//close(fileChan)
}

// Récuperation des informations type personnalisé par loop sur le channel fileChan
func GetDirFilesInfosCustomChannelThree(ctx context.Context, fileChan <-chan string, resultChan chan<- models.FilesStruct, unit string) {
	for filePath := range fileChan {
		select {
		case <-ctx.Done():
			return //Annulation
		default:
			// Nettoyage plus robuste
			filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
			filePath = filepath.Clean(filePath)    // Normalise le chemin
			filePathName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filepath.Base(filePath)))
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Erreur : le fichier '%s' n'existe plus.\n", filePath)
					//log.Printf("Erreur : le fichier '%s' n'existe plus.\n", filePath)
				} else {
					fmt.Printf("Erreur d'accès : %v\n", err)
					//log.Printf("Erreur d'accès : %v\n", err)
				}
			} else {
				//fmt.Printf("%v--\n%v\n", filePathName, fileInfo.Name())
				//fmt.Printf("Trouvé : le fichier '%s' existe.\n", filePath)
				// Envoyer les resultats dans le canal pour collecter les résultats (resultChan)
				data := models.FilesStruct{
					Infos:           functions.CheckDirectory(fileInfo.IsDir()),
					Nom:             fileInfo.Name(),
					Chemin:          filepath.ToSlash(filePath),
					Extension:       functions.GetExtensionFiles(fileInfo.IsDir(), filePath),
					Permission:      fileInfo.Mode().String(),
					Octal:           fmt.Sprintf("%04o", fileInfo.Mode().Perm()),
					Type:            functions.GetFileType(fileInfo.Mode()),
					Longueur:        len(filePath),
					Taille:          functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size()),
					TailleConvertie: functions.ConvertUnitSize(int(functions.CalSizeDirFiles(fileInfo.IsDir(), filePath, fileInfo.Size())), unit),
					Modification:    fileInfo.ModTime().Format("2006-01-02 15:04:05"),
					Accent:          functions.ContainsAccents(filePath),
					CaratereSpecial: functions.ContainSpecialCharactor(filePathName),
					DateTraitement:  time.Now().Format("20060102"),
					EspaceVide:      functions.ContainsWhiteSpace(filePath),
				}
				select {
				case resultChan <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// Scanner les repertoires pour lister fichiers /dossiers avec les informations en tableau selon type personnalisé
func ListDirsFilesChannelThree(ctx context.Context, dir string, extension []string, fileChan chan<- string) {
	// Ferme le canal fileChan après avoir terminé l'envoi des chemins de fichiers
	defer close(fileChan) // dernier message

	//Fermeture
	defer fmt.Println("Fin Calcul, Récupération des informations sur fichiers/dossiers")
	defer log.Println("Fin Calcul, Récupération des informations sur fichiers/dossiers")
	defer fmt.Println("Fin Scanning ==> Répertoire")
	defer log.Println("Fin Scanning ==> Répertoire") // premier message

	//Ouverture
	fmt.Println("Début Scanning ==> Répertoire")
	log.Println("Début Scanning ==> Répertoire")
	fmt.Println("Début Calcul, Récupération des informations sur fichiers/dossiers")
	log.Println("Début Calcul, Récupération des informations sur fichiers/dossiers")

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			//return err
			return nil // Ignorer les erreurs mais continuer
		}
		if path == dir {
			return nil
		}
		//extension d'un fichier (pas de dossier)
		if len(extension) > 0 {
			if !d.IsDir() && functions.Contains(extension, filepath.Ext(path)) {
				//fileChan <- path
				select {
				case fileChan <- path:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		} else {
			// tous les fichiers et les dossiers
			if !d.IsDir() {
				//fileChan <- path
				select {
				case fileChan <- path:
				case <-ctx.Done():
					return ctx.Err()
				}
			} else {
				//fileChan <- path
				select {
				case fileChan <- path:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Erreur de parcours: %v\n", err)
		log.Printf("Erreur de parcours: %v\n", err)
		//cancel() // Annuler le contexte en cas d'erreur
	}

}
