/*
* Progamme de définition d'un type personnalisé
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 20/05/2025
* Version : 1.0.0
*
* Description : Go /Golang est hautement typé donc il faut forcément définir un type personnalisé en fonction des retours
*
 */

package models

// type personnalisé
type FilesStruct struct {
	Infos           string
	Nom             string
	Chemin          string
	Extension       string
	Permission      string
	Octal           string
	Type            string
	Longueur        int
	Taille          int64
	TailleConvertie int
	Modification    string
	Accent          string
	CaratereSpecial string
	DateTraitement  string
	EspaceVide      string
}

// Obtention des infos par la méthode Receiver value
func (f FilesStruct) NewFilesStructReceiverValue() string {
	return f.Infos
}

// Modification des infos par la méthode Receiver pointer sur l'adresse en memoire
func (f *FilesStruct) NewFilesStructReceiverPointer(Infos, Nom, Chemin, Extension, Permission, Octal, Type, Modification string, Taille int64, TailleConvertie int) {
	f.Infos = Infos
	f.Nom = Nom
	f.Chemin = Chemin
	f.Extension = Extension
	f.Permission = Permission
	f.Octal = Octal
	f.Type = Type
	f.Taille = Taille
	f.TailleConvertie = TailleConvertie
	f.Modification = Modification
}

// Utilisation de la même adresse sans copie pour la création de nouvelles infos
func NewFilesStruct(Infos, Nom, Chemin, Extension, Permission, Octal, Type, Modification string, Taille int64, TailleConvertie int) *FilesStruct {
	return &FilesStruct{
		Infos:           Infos,
		Nom:             Nom,
		Chemin:          Chemin,
		Extension:       Extension,
		Permission:      Permission,
		Octal:           Octal,
		Type:            Type,
		Taille:          Taille,
		TailleConvertie: TailleConvertie,
		Modification:    Modification,
	}
}
