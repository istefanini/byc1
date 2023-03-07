package controllers

import (
	"byc1/infra"
	"byc1/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// CREATE TABLE persona (id serial PRIMARY KEY, nombre VARCHAR ( 50 ), apellido VARCHAR ( 50 ), direccion VARCHAR ( 50 ), telefono VARCHAR( 50 ));

func Filehandle(c *gin.Context) {

	file, _, err := c.Request.FormFile("myfile")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "No hay archivo")
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "No se pudo abrir el archivo")
		return
	}
	rows, err := f.GetRows("Contactos")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "No se pudo leer la hoja del archivo .xlsx")
		return
	}
	var personas []models.Persona
	for index, row := range rows {
		if index == 0 {
			continue
		}
		if index == 99 {
			log.Println(personas)
		}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(c.Writer, "ID Error Row=%d Value=%s", index+1, row[0])
			return
		}
		personas = append(personas, models.Persona{
			ID:        id,
			Nombre:    row[1],
			Apellidos: row[2],
			Dirección: row[3],
			Teléfono:  row[4],
		})
	}
	log.Println(personas)
	if err := InsertPersonas(c, personas); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, err.Error())
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(c.Writer, "Datos insertados!!!!!!!!!!")
}

func InsertPersonas(c *gin.Context, personas []models.Persona) error {
	ctx := context.Background()
	DBConection := infra.DbLocal

	tx, err := DBConection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	const query = "insert into persona (id,nombre,apellido,direccion,telefono) values($1,$2,$3,$4,$5)"
	for _, p := range personas {
		stm, err := tx.Prepare(query)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stm.Exec(p.ID, p.Nombre, p.Apellidos, p.Dirección, p.Teléfono)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func Filehandlebyc(c *gin.Context) {
	file, _, err := c.Request.FormFile("myfilebyc")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "No hay archivo con la key: myfilebyc",
			"time": time.Now(),
		})
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "No se pudo abrir el archivo",
			"time": time.Now(),
		})
		return
	}
	rows, err := f.GetRows("Hoja1")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "No se pudo leer la hoja del archivo .xlsx",
			"time": time.Now(),
		})
		return
	}
	var filas []models.Row
	for index, row := range rows {
		if index == 0 {
			continue
		}
		if index == 99 {
			log.Println(filas)
		}
		filas = append(filas, models.Row{
			Periodo:             row[0],
			TipoCosto_ID:        row[1],
			UnidadNegocioJDE_ID: row[2],
			TipoItem_ID:         row[3],
			SubtipoItem_ID:      row[4],
			Item_ID:             row[5],
			TipoEpisodio_ID:     row[6],
			Episodio_ID:         row[7],
			Valor:               row[8],
			OperadorAritmetico:  row[9],
			Sitio_ID:            row[10],
		})
	}
	log.Println(filas)
	if err := InsertRows(c, filas); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Datos insertados correctamente",
		"time": time.Now(),
	})
}

func InsertRows(c *gin.Context, filas []models.Row) error {
	ctx := context.Background()
	DBConection := infra.DbPayment

	tx, err := DBConection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//elimino las filas de la bd cuyo "Periodo" sea igual al "Periodo" de las nuevas filas recibidas como parametro
	const selectQuery = "select Periodo from G_Costos_Setup where Periodo = $1"
	var periodo string
	for _, f := range filas {
		err := tx.QueryRow(selectQuery, f.Periodo).Scan(&periodo)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			return err
		}
		if periodo == f.Periodo {
			const deleteQuery = "delete from G_Costos_Setup where Periodo = $1"
			_, err = tx.Exec(deleteQuery, f.Periodo)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	//inserto todas las filas "nuevas" recibidas como parametro en la función
	const insertQuery = "insert into G_Costos_Setup (Periodo,TipoCosto_ID,UnidadNegocioJDE_ID,TipoItem_ID,SubtipoItem_ID,Item_ID,TipoEpisodio_ID,Episodio_ID,Valor,OperadorAritmetico,Sitio_ID) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	for _, f := range filas {
		stm, err := tx.Prepare(insertQuery)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stm.Exec(f.Periodo, f.TipoCosto_ID, f.UnidadNegocioJDE_ID, f.TipoItem_ID, f.SubtipoItem_ID, f.Item_ID, f.TipoEpisodio_ID, f.Episodio_ID, f.Valor, f.OperadorAritmetico, f.Sitio_ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
