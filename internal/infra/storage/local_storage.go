package storage

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"plantao/internal/infra/config"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(cfg *config.Config) *LocalStorage {
	return &LocalStorage{
		BasePath: cfg.FalePath.Path,
	}
}

func (s *LocalStorage) Save(file io.ReadSeeker, fileName string) (string, error) {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	if contentType != "image/jpeg" &&
		contentType != "image/png" &&
		contentType != "image/webp" &&
		contentType != "image/jpg" {
		return "", errors.New("arquivo inválido")
	}

	file.Seek(0, 0)

	ext := filepath.Ext(fileName)

	newName := uuid.New().String() + ext

	err = os.MkdirAll(s.BasePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	path := filepath.Join(s.BasePath, newName)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return "/uploads/" + newName, nil
}

func (s *LocalStorage) Delete(filePath string) error {
	if filePath == "" {
		return nil
	}

	fullPath := filepath.Join(".", filePath)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao deletar foto do colaborador: %w", err)
	}

	return nil
}
