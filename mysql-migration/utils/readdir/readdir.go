package readdir

import (
	"log"
	"os"
	"path/filepath"
)

type ReadDirResult struct {
	FileName string
	Content  string
}

func Run(folderPath string, skipErrorRead bool) ([]*ReadDirResult, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}
	var result []*ReadDirResult

	// loop every file
	for _, file := range files {
		// check if not dir
		if !file.IsDir() {
			// read file content
			filePath := filepath.Join(folderPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				if skipErrorRead {
					log.Println("Error:", err)
					continue
				}
				return nil, err
			}
			result = append(result, &ReadDirResult{
				FileName: file.Name(),
				Content:  string(content),
			})
		}
	}
	return result, nil
}
