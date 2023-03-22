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

		tipocostoint := row[1]
		tipocosto, err := strconv.Atoi(tipocostoint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		tipoitemint := row[3]
		tipoitem, err := strconv.Atoi(tipoitemint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		subtipoitemint := row[4]
		subtipoitem, err := strconv.Atoi(subtipoitemint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		itemint := row[5]
		item, err := strconv.Atoi(itemint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		tipoepisodioint := row[6]
		tipoepisodio, err := strconv.Atoi(tipoepisodioint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		episodioint := row[7]
		episodio, err := strconv.Atoi(episodioint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		sitioint := row[10]
		sitio, err := strconv.Atoi(sitioint)
		if err != nil {
			fmt.Println("Error de conversion")
		}

		valorfloat := row[8]
		valor, err := strconv.ParseFloat(valorfloat, 64)
		if err == nil {
			fmt.Println(valor)
		}

		filas = append(filas, models.Row{
			Periodo:             row[0],
			TipoCosto_ID:        tipocosto,
			UnidadNegocioJDE_ID: row[2],
			TipoItem_ID:         tipoitem,
			SubtipoItem_ID:      subtipoitem,
			Item_ID:             item,
			TipoEpisodio_ID:     tipoepisodio,
			Episodio_ID:         episodio,
			Valor:               valor,
			OperadorAritmetico:  row[9],
			Sitio_ID:            sitio,
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
	// DBConection := infra.DbLocal

	tx, err := DBConection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//elimino las filas de la bd cuyo "Periodo" sea igual al "Periodo" de las nuevas filas recibidas como parametro

	const selectQuery = "select Periodo from G_Costos_Setup where Periodo = $1"
	// const selectQuery = "select Periodo from tablaprueba where Periodo = $1"
	var periodo string
	for _, f := range filas {
		err := tx.QueryRow(selectQuery, f.Periodo).Scan(&periodo)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			return err
		}
		if periodo == f.Periodo {
			const deleteQuery = "delete from G_Costos_Setup where Periodo = $1"
			// const deleteQuery = "delete from tablaprueba where Periodo = $1"
			_, err = tx.Exec(deleteQuery, f.Periodo)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	//inserto todas las filas "nuevas" recibidas como parametro en la funciónn

	stm, err := tx.Prepare("insert into G_Costos_Setup (Periodo, TipoCosto_ID, UnidadNegocioJDE_ID, TipoItem_ID, " +
		"SubtipoItem_ID, Item_ID, TipoEpisodio_ID, Episodio_id, Valor, OperadorAritmetico, Sitio_ID) " +
		"values(@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11)")
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, f := range filas {
		_ = stm.QueryRowContext(
			ctx,
			sql.Named("p1", f.Periodo),
			sql.Named("p2", f.TipoCosto_ID),
			sql.Named("p3", f.UnidadNegocioJDE_ID),
			sql.Named("p4", f.TipoItem_ID),
			sql.Named("p5", f.SubtipoItem_ID),
			sql.Named("p6", f.Item_ID),
			sql.Named("p7", f.TipoEpisodio_ID),
			sql.Named("p8", f.Episodio_ID),
			sql.Named("p9", f.Valor),
			sql.Named("p10", f.OperadorAritmetico),
			sql.Named("p11", f.Sitio_ID),
		)

	}

	// Ejecutar SP
	// const spQuery = "EXEC [dbo].[CALCULAR_COSTOS_FLENI_Etiquetado_XSubtipoIF] @PERIODO = $1, @SITIO = $2"
	// stmt, err := tx.Prepare(spQuery)
	// if err != nil {
	// tx.Rollback()
	// return err
	// }
	// _, err = stmt.Exec(filas[0].Periodo, filas[0].Sitio_ID)
	// if err != nil {
	// tx.Rollback()
	// return err
	// }

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
