/*
* Liste des Fonctions usuelles
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 17/07/2025
* Version : 1.0.0
*
* Description : Définition  des fonctions, des méthodes pour les utiliser dans les sous programmes
*
 */
package functions

import (
	"fmt"
	"github/managerDirectory/models"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Récupèration les noms des champs d'une struct pour les en-têtes
func GetStructFields(item interface{}) []string {
	t := reflect.TypeOf(item)
	var fields []string

	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

func GetStructFieldNames(myStruct interface{}) []string {
	// Obtient le type de la struct
	t := reflect.TypeOf(myStruct)

	// Vérifie si c'est une struct
	if t.Kind() != reflect.Struct {
		fmt.Println("L'argument fourni n'est pas une struct")
		return nil
	}

	// Obtient le nombre de champs dans la struct
	numFields := t.NumField()

	// Crée une slice pour stocker les noms des champs
	fieldNames := make([]string, numFields)

	// Parcourt les champs de la struct et stocke leurs noms
	for i := 0; i < numFields; i++ {
		fieldNames[i] = t.Field(i).Name
	}
	return fieldNames
}

// Obtention des noms des champs de la structure
func GetFieldsNamesStruct() []string {
	return GetStructFields(models.FilesStruct{})
}

// Convertir une structure en slice de strings pour le CSV
func StructToSlice(item interface{}) []string {
	v := reflect.ValueOf(item)
	var row []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.String:
			row = append(row, field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			row = append(row, strconv.FormatInt(field.Int(), 10))
		case reflect.Bool:
			row = append(row, strconv.FormatBool(field.Bool()))
		case reflect.Float32, reflect.Float64:
			row = append(row, strconv.FormatFloat(field.Float(), 'f', 2, 64))
		default:
			row = append(row, "")
		}
	}
	return row
}

// Calcul la durée du traitement (pas de fonctionnement)
func TimeAsString(dt float64) string {
	time := dt
	hours := math.Floor(time / 3600)
	minutes := math.Ceil(math.Mod(time, 3600)/60) - 1
	seconds := int(time) % 60
	return fmt.Sprintf("%v:%v:%v", hours, minutes, seconds)
}

// Calcul la durée du traitement
func CalDuration(dt time.Duration) string {
	hours := int64(dt.Hours())
	minutes := int64(dt.Minutes()) % 60
	seconds := int64(dt.Seconds()) % 60
	return fmt.Sprintf("%v H :%v mn :%v seconds", hours, minutes, seconds)
}

func GetCurrentFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "Impossible d'obtenir le chemin du fichier"
	}
	return filename
}

// Gestion Log
func DefinieLog(prefixe string) (*os.File, error) {
	fileTime := time.Now().Format("20060102")
	filePath := GetCurrentFilePath()
	if _, err := os.Stat(filepath.Join(filepath.Dir(filePath), fileTime)); os.IsNotExist(err) {
		err := os.Mkdir(filepath.Join(filepath.Dir(filePath), fileTime), 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
	filelog, err := os.OpenFile(filepath.Join(filepath.Dir(filePath), fileTime, fmt.Sprintf("%v_%vLog.log", fileTime, prefixe)),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return filelog, nil
}

// Vérification de l'existance de fichier ou de repertoire
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Création d'un repertoire et le fichier
func CreateFileWithDirs(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Création d'un repertoire
func CreateDirs(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	return "_", nil
}

// Taille d'un dossier
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Convertisseur taille
func ConvertUnitSize(size_in_bytes int, unit string) int {
	if unit == "KB" {
		return size_in_bytes / 1024
	} else if unit == "MB" {
		return size_in_bytes / (1024 * 1024)
	} else if unit == "GB" {
		return size_in_bytes / (1024 * 1024 * 1024)
	} else {
		return size_in_bytes
	}
}

// Savoir si c'est un dossier ou un fichier
func CheckDirectory(inf bool) string {
	if inf {
		return "Dossier"
	} else {
		return "Fichier"
	}
}

// Calcul la taille du fichier ou dossier
func CalSizeDirFiles(inf bool, path string, infSize int64) int64 {
	if inf {
		calSize, _ := DirSize(path)
		return calSize

	} else {
		calSize := infSize
		return calSize
	}
}

// Obtention des extensions des fichiers
func GetExtensionFiles(inf bool, path string) string {
	if inf {
		return "-"
	} else {
		return strings.Replace(filepath.Ext(path), ".", "", -1)
	}
}

// Obtention du type fichier (Mode)
func GetFileType(mode os.FileMode) string {
	switch {
	case mode.IsDir():
		return "Répertoire"
	case mode&os.ModeSymlink != 0:
		return "Lien symbolique"
	case mode&os.ModeNamedPipe != 0:
		return "Tube nommé"
	case mode&os.ModeSocket != 0:
		return "Socket"
	case mode&os.ModeDevice != 0:
		return "Période bloc"
	case mode&os.ModeCharDevice != 0:
		return "Période caractère"
	default:
		return "Fichier régulier"
	}
}

// Vérifier des caractères spéciaux
func ContainSpecialCharactor(str string) string {
	specialCharacter := []string{"|", "#", "'", "\"", "\\", "%", "?",
		"\n", "<", "Ø", "ð", ">", "ï", "û",
		",", ";", ".", ":", "!", "(", ")", "/", "[",
		"]", "+", "=", "»", "*", "&", "@", "…"}
	for _, sc := range specialCharacter {
		if strings.Contains(str, sc) {
			return "KO"
		}
	}
	return "OK"
}

// Vérifier des caractères spéciaux (ne fonctinnne pas , unicode à revoir)
func ContainsSpecialChars(path string) string {
	regex := regexp.MustCompile(`[,;.:!?()/\\[]+=»*%&@…]`)
	if regex.MatchString(path) {
		return "KO"
	} else {
		return "OK"
	}
}

// Vérifier si le chemin contient des caractères spéciaux (ne fonctinnne pas , unicode à revoir)
func ContainsSpecialCharsPath(path string) string {
	if len(path) >= 2 && path[1] == ':' && path[2] == '\\' {
		path = path[3:]
	}
	regex := regexp.MustCompile(`[\\/:*?"<>|]`)
	if regex.MatchString(path) {
		return "KO"
	} else {
		return "OK"
	}
}

// Vérifier si le chemin contient des caractères accentués
func ContainsAccents(path string) string {
	regex := regexp.MustCompile(`[àâäéèêëîïôöùûüç]`) // Ajoutez d'autres caractères accentués si nécessaire
	if regex.MatchString(path) {
		return "KO"
	} else {
		return "OK"
	}
}

// Vérifier si le chemin contient des caractères accentués
func ContainsWhiteSpace(str string) string {
	if strings.Contains(str, " ") {
		return "KO"
	} else {
		return "OK"
	}
}

// comparaison de deux chaines de caratères
func CompareEF(a string, b string) bool {
	if strings.EqualFold(a, b) {
		return true
	} else {
		return false
	}
}

// Filtrer les extensions ou autres
func Contains(extensions []string, ext string) bool {
	for _, allowedExt := range extensions {
		if strings.ToLower(allowedExt) == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

func ContainsTwo(extensions []string, ext string) bool {
	for _, allowedExt := range extensions {
		if strings.EqualFold(allowedExt, ext) {
			return true
		}
	}
	return false
}

// Scanner les repertoires pour lister uniquement les fichiers
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

// Scanner les repertoires pour lister fichiers /dossiers
func ListDirsFiles(dir string, extension string) []string {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		//extension d'un fichier (pas de dossier)
		if len(extension) != 0 {
			if !d.IsDir() && filepath.Ext(path) == extension {
				files = append(files, path)
			}
		} else {
			// tous les fichiers et les dossiers
			if !d.IsDir() {
				files = append(files, path)
			} else {
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

// Scanner les repertoires pour lister les fichiers selon le mots clé dans le nom
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

// Fonction pour copier un fichier
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// Fonction pour copier récursivement le contenu d'un dossier
func CopyDirContents(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fileInfo, err := entry.Info()
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, fileInfo.Mode()); err != nil {
				return err
			}
			if err := CopyDirContents(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
			if err := os.Chmod(dstPath, fileInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyDirectory(src, dst, keyword string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip files/directories containing the keyword
		if strings.Contains(info.Name(), keyword) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		destPath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			err = os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(path, destPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func DeleteIfExists(path string) error {
	exists, err := Exists(path)
	if err != nil {
		return err
	}
	if exists {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		fmt.Printf("%s est supprimé -- OK\n", path)
	} else {
		fmt.Printf("%s n'existe pas.\n", path)
	}
	return nil
}
