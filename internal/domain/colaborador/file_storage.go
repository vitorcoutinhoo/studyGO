package colaborador

import "io"

type FileStorage interface {
	Save(file io.ReadSeeker, fileName string) (string, error)
	Delete(filePath string) error
}
