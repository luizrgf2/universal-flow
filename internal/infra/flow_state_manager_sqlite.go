package infra

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/luizrgf2/universal-flow/internal/core/entities"
	_ "github.com/mattn/go-sqlite3"
)

type FlowStateManagerSqlite struct {
	db *sql.DB
}

func NewFlowStateManagerSqlite(dbPath string) (*FlowStateManagerSqlite, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS flow_states (
		id TEXT PRIMARY KEY,
		data TEXT NOT NULL
	);
	`
	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return &FlowStateManagerSqlite{
		db: db,
	}, nil
}

func (m *FlowStateManagerSqlite) GetFlowState(flowId string) (*entities.Flow, error) {
	query := `SELECT data FROM flow_states WHERE id = ?`
	row := m.db.QueryRow(query, flowId)

	var data string
	err := row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flow not found")
		}
		return nil, err
	}

	var flow entities.Flow
	err = json.Unmarshal([]byte(data), &flow)
	if err != nil {
		return nil, err
	}

	return &flow, nil
}

func (m *FlowStateManagerSqlite) CreateFlow(flow *entities.Flow) error {
	data, err := json.Marshal(flow)
	if err != nil {
		return err
	}

	query := `INSERT INTO flow_states (id, data) VALUES (?, ?)`
	_, err = m.db.Exec(query, flow.ID, string(data))
	return err
}

func (m *FlowStateManagerSqlite) UpdateFlow(flow *entities.Flow) error {
	data, err := json.Marshal(flow)
	if err != nil {
		return err
	}

	query := `UPDATE flow_states SET data = ? WHERE id = ?`
	_, err = m.db.Exec(query, string(data), flow.ID)
	return err
}
