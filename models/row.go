package models

type Row struct {
	Periodo             string
	TipoCosto_ID        int
	UnidadNegocioJDE_ID string
	TipoItem_ID         int
	SubtipoItem_ID      int
	Item_ID             int
	TipoEpisodio_ID     int
	Episodio_ID         int
	Valor               float64
	OperadorAritmetico  string
	Sitio_ID            int
}
