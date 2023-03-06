package controllers

import (
	"byc1/infra"
	"byc1/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		fmt.Fprintln(c.Writer, "No hay archivo")
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "No se pudo abrir el archivo")
		return
	}
	rows, err := f.GetRows("Hoja1")
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(c.Writer, "No se pudo leer la hoja del archivo .xlsx")
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
		// id, err := strconv.Atoi(row[0])
		// if err != nil {
		// 	c.Writer.WriteHeader(http.StatusBadRequest)
		// 	fmt.Fprintln(c.Writer, "ID Error Row=%d Value=%s", index+1, row[0])
		// 	return
		// }
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
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintln(c.Writer, "Datos insertados!!!!!!!!!!")
}

func InsertRows(c *gin.Context, filas []models.Row) error {
	ctx := context.Background()
	DBConection := infra.DbLocal

	tx, err := DBConection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	const query = "insert into tablaprueba (Periodo,TipoCosto_ID,UnidadNegocioJDE_ID,TipoItem_ID,SubtipoItem_ID,Item_ID,TipoEpisodio_ID,Episodio_ID,Valor,OperadorAritmetico,Sitio_ID) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	for _, f := range filas {
		stm, err := tx.Prepare(query)
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

// var logs []models.Log

// func GetAllLogs(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	tsql := fmt.Sprintf("SELECT TOP (100) [Transaccion_ID],[EventID],[SysFechaC],[Estado],isnull(LogProceso, '') as LogProceso, isnull(MsgFinal, '') as MsgFinal FROM [Interoperabilidad].[dbo].[MQMsgDisparados]")
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	logs := []models.Log{}
// 	for rows.Next() {
// 		var log models.Log
// 		err = rows.Scan(&log.Transaccion_ID, &log.EventID, &log.SysFechaC, &log.Estado, &log.LogProceso, &log.MsgFinal)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		logs = append(logs, log)
// 	}
// 	c.JSON(http.StatusOK, logs)
// }

// func GetLogsbyParams(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	logId := c.Query("Transaccion_ID")
// 	estado := c.Query("Estado")
// 	tsql := fmt.Sprintf(
// 		`SELECT [Transaccion_ID],[EventID],[SysFechaC],[Estado],isnull(LogProceso, '') as LogProceso, isnull(MsgFinal, '') as MsgFinal
// 		FROM [Interoperabilidad].[dbo].[MQMsgDisparados]
// 		where ([Transaccion_ID] IS NULL OR [Transaccion_ID]=` + logId + `)
// 		AND ([Estado] IS NULL OR [Estado]=` + estado + `)`)
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	logs := []models.Log{}
// 	for rows.Next() {
// 		var log models.Log
// 		err = rows.Scan(&log.Transaccion_ID, &log.EventID, &log.SysFechaC, &log.Estado, &log.LogProceso, &log.MsgFinal)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		logs = append(logs, log)
// 	}
// 	c.JSON(http.StatusOK, logs)
// }

// func GetLogsbyID(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	logId := c.Query("Transaccion_ID")
// 	tsql := fmt.Sprintf(
// 		`SELECT [Transaccion_ID],[SysFechaC],[Estado],isnull(LogProceso, '') as LogProceso, isnull(MsgFinal, '') as MsgFinal
// 		FROM [Interoperabilidad].[dbo].[MQMsgDisparados]
// 		where ([Transaccion_ID] IS NULL OR [Transaccion_ID]=` + logId + `)`)
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	logs := []models.Log{}
// 	for rows.Next() {
// 		var log models.Log
// 		err = rows.Scan(&log.Transaccion_ID, &log.SysFechaC, &log.Estado, &log.LogProceso, &log.MsgFinal)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		logs = append(logs, log)
// 	}
// 	c.JSON(http.StatusOK, logs)
// }

// func GetEvents(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	tsql := fmt.Sprintf(
// 		`SELECT [EventID],[Descripcion],[MQConnID]
// 		FROM [Interoperabilidad].[dbo].[MQEvents] `)
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	events := []models.Event{}
// 	for rows.Next() {
// 		var event models.Event
// 		err = rows.Scan(&event.EventID, &event.Descripcion, &event.MQConnID)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		events = append(events, event)
// 	}
// 	c.JSON(http.StatusOK, events)
// }

// func GetLogs(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	estado := c.Query("Estado")
// 	desde := c.Query("Desde")
// 	hasta := c.Query("Hasta")
// 	evento := c.Query("Evento")
// 	tsql := fmt.Sprintf(
// 		`SELECT [Transaccion_ID],[EventID],[SysFechaC],[Estado],isnull(LogProceso, '') as LogProceso, isnull(MsgFinal, '') as MsgFinal
// 		FROM [Interoperabilidad].[dbo].[MQMsgDisparados]
// 		WHERE [Estado] IS NULL OR [Estado] IN (` + estado + `)
// 		AND [EventID] IS NULL OR [EventID] IN ('` + evento + `')
// 		AND [SysFechaC] BETWEEN '` + desde + ` 00:00' AND '` + hasta + ` 23:59'`)
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	logs := []models.Log{}
// 	for rows.Next() {
// 		var log models.Log
// 		err = rows.Scan(&log.Transaccion_ID, &log.EventID, &log.SysFechaC, &log.Estado, &log.LogProceso, &log.MsgFinal)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		logs = append(logs, log)
// 	}
// 	c.JSON(http.StatusOK, logs)
// }

// func ReprocesLogs(c *gin.Context) {
// 	ctx := context.Background()
// 	DBConection := infra.DbPayment
// 	Transacciones := c.Query("Transacciones")
// 	stringSlice := strings.Split(Transacciones, ",")
// 	for i, elem := range stringSlice {
// 		stringSlice[i] = "'" + elem + "'"
// 	}
// 	str := strings.Join(stringSlice, ",")
// 	tsql := `UPDATE [Interoperabilidad].[dbo].[MQMsgDisparados] SET Estado = 1, Intentos = 0
// 	WHERE Transaccion_ID IN (` + str + `);`
// 	rows, err := DBConection.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	c.JSON(http.StatusOK, logs)
// }

// func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(httpStatusCode)
// 	resp := make(map[string]string)
// 	resp["message"] = message
// 	jsonResp, _ := json.Marshal(resp)
// 	w.Write(jsonResp)
// }
