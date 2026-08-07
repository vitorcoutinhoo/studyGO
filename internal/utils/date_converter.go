package utils

import (
	"fmt"
	"time"
)

// ParseBrToUsDate converte data no formato DD/MM/AAAA ou AAAA-MM-DD para time.Time no fuso horário informado.
// Se loc for nil, usa time.UTC.
func ParseBrToUsDate(data *string, loc *time.Location) (*time.Time, error) {
	if data == nil || *data == "" {
		return nil, fmt.Errorf("erro, data vazia")
	}

	if loc == nil {
		loc = time.UTC
	}

	layouts := []string{"02/01/2006", "2006-01-02"}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, *data, loc)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("data inválida, use DD/MM/AAAA ou AAAA-MM-DD")
	}

	return &t, nil
} // Fim ParseBrToUsDate

// ParseUsToBrDate converte time.Time para o formato brasileiro DD/MM/AAAA no fuso horário informado.
// Se loc for nil, usa o fuso horário já presente no valor.
func ParseUsToBrDate(t *time.Time, loc *time.Location) (string, error) {
	if t == nil || t.IsZero() {
		return "", fmt.Errorf("erro, data zero ou nula")
	}

	if loc != nil {
		converted := t.In(loc)
		return converted.Format("02/01/2006"), nil
	}

	return t.Format("02/01/2006"), nil
} // Fim ParseUsToBrDate

// ToLocation converte um time.Time para o fuso horário especificado.
func ToLocation(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		return t
	}
	return t.In(loc)
} // Fim ToLocation

// BrasilLocation retorna o fuso horário de Brasília (America/Sao_Paulo).
func BrasilLocation() *time.Location {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.UTC
	}
	return loc
} // Fim BrasilLocation
