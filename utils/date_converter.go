package utils

import (
	"fmt"
	"time"
)

/*
Converte uma data no formato brasileiro (DD/MM/AAAA) para o formato americano (AAAA-MM-DD)
e retorna um ponteiro para time.Time ou um erro caso a conversão falhe.
*/
func ParseBrToUsDate(data *string) (*time.Time, error) {
	if data == nil || *data == "" {
		return nil, fmt.Errorf("erro, data vazia")
	}

	layout := "02/01/2006"

	t, err := time.Parse(layout, *data)

	if err != nil {
		return nil, err
	}

	return &t, nil
} // Fim ParseBrToUsDate

/*
Converte uma data no formato americano (AAAA-MM-DD) para o formato brasileiro (DD/MM/AAAA)
e retorna a string formatada.
*/
func ParseUsToBrDate(t *time.Time) (string, error) {
	if t == nil || t.IsZero() {
		return "", fmt.Errorf("erro, data zero ou nula")
	}

	return t.Format("02/01/2006"), nil
} // Fim ParseUsToBrDate
