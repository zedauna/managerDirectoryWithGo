package extension

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

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
