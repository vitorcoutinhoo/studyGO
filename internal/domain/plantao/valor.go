package plantao

import "errors"

var (
	ErrNegativeValue = errors.New("Valor não pode ser negativo!")
)

type Dinheiro struct {
	centavos int64
}

func NewDinheiro(centavos float64) (*Dinheiro, error) {
	if centavos < 0 {
		return nil, ErrNegativeValue
	}

	return &Dinheiro{centavos: int64(centavos)}, nil
}


