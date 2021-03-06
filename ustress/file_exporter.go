package ustress

import (
	"fmt"
	"os"

	log "github.com/metrosystems-cpe/ustress/log"
)

// NewFile returns a new file to write data to
func NewFile(filename string) *os.File {
	createDirIfNotExist(HTTPfolder)
	f, err := os.Create(HTTPfolder + filename)
	f, err = os.OpenFile(HTTPfolder+filename, os.O_RDWR|os.O_APPEND, 0766) // For read access.
	if err != nil {
		log.LogWithFields.Errorln(err.Error())
	}
	return f
}

// SaveFileReport it will save the report as json on local storage
func SaveFileReport(r *Report) {
	jsonReport := r.JSON()
	fileWriter := NewFile(fmt.Sprintf("%s.json", r.UUID))
	defer fileWriter.Close()
	fmt.Fprintf(fileWriter, string(jsonReport))
}


// CreateDirIfNotExist the function name says it all
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.LogWithFields.Errorln(err.Error())
		}
	}
}
