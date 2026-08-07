package plantao

import "time"

type RelatorioFiltro struct {
	DataInicio    time.Time
	DataFim       time.Time
	ColaboradorID string
	Status        *StatusPlantao
}

type RelatorioItem struct {
	ColaboradorID   string
	PlantaoID       string
	Status          StatusPlantao
	DataInicio      time.Time
	DataFim         time.Time
	Data            time.Time
	NomeColaborador string
	ValorTotal      float64
	Valor           float64
	Observacoes     *string
}
