package upload

import (
	"advanced-web.hcmus/util"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	workingDir, _ = os.Getwd()
	mediaDir      = fmt.Sprintf("%s/public/system", workingDir)
	TemporaryDir  = mediaDir + "/temporary"
)

func (temp Temporary) pre() {
	// check temporary dir is exists, if not make it
	if _, err := os.Stat(TemporaryDir); os.IsNotExist(err) {
		if err := os.MkdirAll(mediaDir, 0644); err != nil {
			panic(err)
		}
	}
}

func (temp Temporary) Save(fileName string, fileContent []byte) string {
	// file name have extension: ex.mp3
	fileFullName := filepath.Base(fileName)
	fileExt := filepath.Ext(fileName)
	fileNameForHash := fmt.Sprintf("%s_%d", fileFullName, time.Now().Unix()) // for name is always unique
	hashFileName := fmt.Sprintf("%s%s", util.HexSha256String([]byte(fileNameForHash)), fileExt)
	util.CreateFolderV2(TemporaryDir)
	filePath := TemporaryDir + "/" + hashFileName
	err := ioutil.WriteFile(filePath, fileContent, 0644)
	util.CheckErr(err)

	return filePath
}

func (temp Temporary) SaveV2(fileName string, fileContent []byte, folderPath string) string {
	// file name have extension: ex.mp3
	fileFullName := filepath.Base(fileName)
	fileExt := filepath.Ext(fileName)
	fileNameForHash := fmt.Sprintf("%v_%v", fileFullName, time.Now().Unix()) // for name is always unique
	hashFileName := fmt.Sprintf("%v%v", util.HexSha256String([]byte(fileNameForHash)), fileExt)
	finalFolderPath := mediaDir + "/" + folderPath
	util.CreateFolderV2(finalFolderPath)
	filePath := finalFolderPath + "/" + hashFileName
	err := ioutil.WriteFile(filePath, fileContent, 0644)
	util.CheckErr(err)

	return fmt.Sprintf("/system/%v/%v", folderPath, hashFileName)
}

type Temporary struct {
}

func WithTemporary() Temporary {
	temp := Temporary{}
	temp.pre()
	return temp
}
