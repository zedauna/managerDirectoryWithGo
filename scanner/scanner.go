package scanner

import (
	"bufio"
	"fmt"
	"github/managerDirectory/functions"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadPath() (string, error) {
	var dossierPath string
	// Obtenir le chemin à parcourir
	if len(os.Args) != 0 {
		dossierPath = os.Args[1]
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
		}
		return "", fmt.Errorf("erreur d'accès : %v", err)
	}

	if !info.IsDir() {
		fmt.Printf("'%s' n'est pas un dossier valide. Veuillez réessayer.\n", dossierPath)
	}
	// Conversion en chemin absolu
	absPath, err := filepath.Abs(dossierPath)
	if err != nil {
		return "", fmt.Errorf("erreur de conversion en chemin absolu : %v", err)
	}
	return absPath, nil
}

func GetPath() (string, error) {
	// Obtenir le chemin à parcourir
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Veuillez entrer le chemin du dossier : ")
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

func ReadGetPath() (string, error) {
	var path string
	path, err := ReadPath()
	if path == "" {
		path, err = GetPath()
	}
	return path, err
}

func ListFiles(dir string, extension string) []string {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if len(extension) != 0 {
			if !d.IsDir() && filepath.Ext(path) == extension {
				files = append(files, path)
			}
		} else {
			if !d.IsDir() {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func ListFilesFilter(dir string, wildcard string) []string {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != dir {
			// Vérifier si le nom du dossier contient le mot-clé
			if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(wildcard)) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				files = append(files, absPath)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

// erreur dans la fonction
func ListFilesFilterRemove(dir string, wildcard string) ([]string, []string) {
	files := []string{}
	files_echec := []string{}

	// Vérification que le répertoire existe
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Erreur: Le répertoire %s n'existe pas\n", dir)
		os.Exit(1)
	}

	//Message de confirmation
	fmt.Printf("Recherche des dossiers contenant '%s' dans %s...\n", wildcard, dir)
	fmt.Println("Appuyez sur 'Y' puis Entrée pour confirmer la suppression")
	fmt.Println("Toute autre entrée annulera l'opération")

	// Demande de confirmation
	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToUpper(confirm) != "Y" {
		fmt.Println("Opération annulée")
		os.Exit(0)
	}

	// Parcours récursif du répertoire
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != dir {
			// Vérifier si le nom du dossier contient le mot-clé
			if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(wildcard)) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				exists, err := functions.Exists(absPath)
				if err != nil {
					return err
				}
				// verifier si le fichier est toujours disponible
				if exists {
					fmt.Printf("Suppression de %s... ", absPath)
					err = os.RemoveAll(absPath)
					if err != nil {
						fmt.Printf("ÉCHEC: %v\n", err)
						files_echec = append(files_echec, absPath)
					} else {
						fmt.Println("OK")
						files = append(files, absPath)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Erreur lors de la suppression: %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Opération terminée")
	return files, files_echec
}

func ListFilesFilterCopy(dir string, destDir string, wildcard string) error {
	// Vérification que le répertoire existe
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Erreur: Le répertoire source %s n'existe pas\n", dir)
		os.Exit(1)
	}

	// Parcours récursif du répertoire
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Ne traiter que les dossiers (sauf le dossier racine)
		if info.IsDir() && path != dir {
			// Vérifier si le nom du dossier contient le mot-clé / on saut
			if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(wildcard)) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), strings.ToLower(wildcard)) {
				return nil
			}

			// Vérifier si le nom du dossier ne contient pas le mot-clé
			if info.IsDir() {
				relPath, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}
				destPath := filepath.Join(destDir, relPath)
				fmt.Printf("Copie de %s vers %s... ", path, destPath)

				// Création du dossier de destination
				if err = os.MkdirAll(destPath, info.Mode()); err != nil {
					fmt.Printf("ÉCHEC (création dossier): %v\n", err)
					return nil
				}

				// Copie du contenu du dossier
				if err := functions.CopyDirContents(path, destPath); err != nil {
					fmt.Printf("ÉCHEC (copie contenu): %v\n", err)
				} else {
					fmt.Println("OK")
				}

			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Erreur lors de la copie: %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("Opération terminée")
	return nil
}

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
	fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
}

func ListDirFilesInfo(filePath string) {
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
		fmt.Printf("- Chemin : %s\n", filepath.ToSlash(filePath)) // Standardise les slashs
		fmt.Printf("- Taille : %d octets\n", fileInfo.Size())
		fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	}

	if !fileInfo.IsDir() {
		fmt.Printf("\nInfos fichier :\n")
		fmt.Printf("- Chemin : %s\n", filepath.ToSlash(filePath)) // Standardise les slashs
		fmt.Printf("- Taille : %d octets\n", fileInfo.Size())
		fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	}
}
