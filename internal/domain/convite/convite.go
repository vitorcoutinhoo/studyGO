package convite

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorConviteNotFound = errors.New("convite não encontrado ou inválido")
	ErrorConviteExpired  = errors.New("convite expirado")
	ErrorConviteUsed     = errors.New("convite já utilizado")
)

type Convite struct {
	Token         uuid.UUID
	IdColaborador uuid.UUID
	ExpiraEm      time.Time
	Usado         bool
	CreatedAt     time.Time
}

func (c *Convite) Validate() error {
	if c.Usado {
		return ErrorConviteUsed
	}

	if time.Now().After(c.ExpiraEm) {
		return ErrorConviteExpired
	}

	return nil
}
