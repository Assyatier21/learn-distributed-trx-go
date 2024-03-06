package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/assyatier21/learn-distributed-trx/config"
	"github.com/assyatier21/learn-distributed-trx/driver"
	"github.com/google/uuid"
	"github.com/gosom/gosql2pc"
	"github.com/labstack/echo/v4"
)

type InsertPondRequest struct {
	UUID          string `json:"uuid"`
	CPUUID        string `json:"cp_uuid"`
	Name          string `json:"name"`
	FeederBarcode string `json:"feeder_barcode"`
	LeadID        string `json:"lead_id"`
}

func main() {
	cfg := config.Load()

	cultivationDB, err := driver.NewPostgreSQL(cfg.DBFirstConfig)
	if err != nil {
		panic(err)
	}
	defer cultivationDB.Close()

	neptuneDB, err := driver.NewPostgreSQL(cfg.DBSecondConfig)
	if err != nil {
		panic(err)
	}
	defer neptuneDB.Close()

	e := echo.New()

	e.POST("/pond", func(c echo.Context) error {
		var payload InsertPondRequest
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"success": false, "message": "bad request"})
		}

		payload.UUID = uuid.NewString()
		if payload.CPUUID == "" {
			payload.CPUUID = payload.UUID
		}

		cultivationParticipant := gosql2pc.NewParticipant(cultivationDB, func(ctx context.Context, tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO ponds (uuid, name, lead_id) VALUES ($1, $2, $3)", payload.UUID, payload.Name, payload.LeadID)
			return err
		})

		neptuneParticipant := gosql2pc.NewParticipant(neptuneDB, func(ctx context.Context, tx *sql.Tx) error {
			_, err := tx.ExecContext(ctx, "INSERT INTO ponds (uuid, name, cp_uuid, feeder_barcode) VALUES ($1, $2, $3, $4)",
				payload.UUID, payload.Name, payload.CPUUID, payload.FeederBarcode)
			return err
		})

		params := gosql2pc.Params{
			LogFn: func(format string, args ...any) {
				fmt.Println(format, args)
			},
			Participants: []gosql2pc.Participant{cultivationParticipant, neptuneParticipant},
		}

		if err := gosql2pc.Do(context.Background(), params); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": err.Error()})
		}

		// participants := []*Participant{cultivationParticipant, neptuneParticipant}

		// if err := coordinator.ExecuteDistributedTransaction(ctx, payload, participants); err != nil {
		// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"success": false, "message": err.Error()})
		// }

		return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "message": "success"})
	})

	e.Logger.Fatal(e.Start(":8000"))
}
