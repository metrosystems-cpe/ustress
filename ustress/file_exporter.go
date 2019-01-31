package ustress

import (
	"os"

	log "git.metrosystems.net/reliability-engineering/ustress/log"
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

// CreateDirIfNotExist the function name says it all
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.LogWithFields.Errorln(err.Error())
		}
	}
}
