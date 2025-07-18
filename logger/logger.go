/*
* Progamme de gestion des log
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 202/05/2025
* Version : 1.0.0
*
* Description : Personnalisation des logs
*
 */
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
)

func InitLoggers(path_out string) {
	//debut
	fileTime := time.Now().Format("20060102")

	//getion log
	filelog, err := os.OpenFile(filepath.Join(path_out, fmt.Sprintf("%v_log.log", fileTime)),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	ErrorLogger = log.New(filelog, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(filelog, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(filelog, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)

}
