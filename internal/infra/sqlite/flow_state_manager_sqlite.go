package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/luizrgf2/universal-flow/internal/core/entities"
	"github.com/luizrgf2/universal-flow/internal/core/types"

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

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = createTables(db); err != nil {
		return nil, err
	}

	return &FlowStateManagerSqlite{db: db}, nil
}

func createTables(db *sql.DB) error {
	flowsTable := `
	CREATE TABLE IF NOT EXISTS flows (
		id TEXT PRIMARY KEY,
		flow_name TEXT,
		status TEXT,
		current_node_id TEXT,
		next_node_id TEXT,
		previous_node_id TEXT,
		previous_nodes_runned TEXT
	);`

	nodesTable := `
	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		flow_id TEXT,
		name TEXT,
		script_path TEXT,
		status TEXT,
		state_input TEXT,
		state_output TEXT,
		error TEXT,
		output_nodes TEXT,
		selected_node TEXT,
		FOREIGN KEY(flow_id) REFERENCES flows(id)
	);`

	if _, err := db.Exec(flowsTable); err != nil {
		return err
	}

	if _, err := db.Exec(nodesTable); err != nil {
		return err
	}

	return nil
}

func (f *FlowStateManagerSqlite) CreateFlow(flow *entities.Flow) error {
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}

	previousNodesRunnedJSON, err := json.Marshal(flow.PreviousNodesRunned)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO flows (id, flow_name, status, current_node_id, next_node_id, previous_node_id, previous_nodes_runned) VALUES (?, ?, ?, ?, ?, ?, ?)",
		flow.ID,
		flow.FlowName,
		flow.Status,
		flow.CurrentNode,
		flow.NextNode,
		flow.PreviousNode,
		string(previousNodesRunnedJSON),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, node := range flow.Nodes {
		outputNodesJSON, err := json.Marshal(node.OutputNodes)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			"INSERT INTO nodes (id, flow_id, name, script_path, status, state_input, state_output, error, output_nodes, selected_node) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			node.ID,
			flow.ID,
			node.Name,
			node.ScriptPath,
			node.Status,
			node.State.Input,
			node.State.Output,
			node.Error,
			string(outputNodesJSON),
			node.SelectedNode,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (f *FlowStateManagerSqlite) UpdateFlow(flow *entities.Flow) error {
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}

	previousNodesRunnedJSON, err := json.Marshal(flow.PreviousNodesRunned)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"UPDATE flows SET flow_name = ?, status = ?, current_node_id = ?, next_node_id = ?, previous_node_id = ?, previous_nodes_runned = ? WHERE id = ?",
		flow.FlowName,
		flow.Status,
		flow.CurrentNode,
		flow.NextNode,
		flow.PreviousNode,
		string(previousNodesRunnedJSON),
		flow.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, node := range flow.Nodes {
		outputNodesJSON, err := json.Marshal(node.OutputNodes)
		if err != nil {
			tx.Rollback()
			return err
		}

		var nodeError *string
		if node.Error != nil {
			nodeError = node.Error
		}

		_, err = tx.Exec(
			"UPDATE nodes SET name = ?, script_path = ?, status = ?, state_input = ?, state_output = ?, error = ?, output_nodes = ?, selected_node = ? WHERE id = ? AND flow_id = ?",
			node.Name,
			node.ScriptPath,
			node.Status,
			node.State.Input,
			node.State.Output,
			nodeError,
			string(outputNodesJSON),
			node.SelectedNode,
			node.ID,
			flow.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (f *FlowStateManagerSqlite) GetFlowState(flowId string) (*entities.Flow, error) {
	var flow entities.Flow
	var previousNodesRunnedJSON string
	var currentNode, nextNode, previousNode sql.NullString

	row := f.db.QueryRow("SELECT id, flow_name, status, current_node_id, next_node_id, previous_node_id, previous_nodes_runned FROM flows WHERE id = ?", flowId)

	err := row.Scan(
		&flow.ID,
		&flow.FlowName,
		&flow.Status,
		&currentNode,
		&nextNode,
		&previousNode,
		&previousNodesRunnedJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flow with id %s not found", flowId)
		}
		return nil, err
	}

	if currentNode.Valid {
		flow.CurrentNode = &currentNode.String
	}
	if nextNode.Valid {
		flow.NextNode = &nextNode.String
	}
	if previousNode.Valid {
		flow.PreviousNode = &previousNode.String
	}

	if err := json.Unmarshal([]byte(previousNodesRunnedJSON), &flow.PreviousNodesRunned); err != nil {
		return nil, err
	}

	rows, err := f.db.Query("SELECT id, name, script_path, status, state_input, state_output, error, output_nodes, selected_node FROM nodes WHERE flow_id = ?", flowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var node entities.Node
		var outputNodesJSON string
		var nodeError, selectedNode sql.NullString
		var status string

		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.ScriptPath,
			&status,
			&node.State.Input,
			&node.State.Output,
			&nodeError,
			&outputNodesJSON,
			&selectedNode,
		)
		if err != nil {
			return nil, err
		}

		node.Status = types.NodeStatus(status)

		if nodeError.Valid {
			node.Error = &nodeError.String
		}
		if selectedNode.Valid {
			node.SelectedNode = &selectedNode.String
		}

		if err := json.Unmarshal([]byte(outputNodesJSON), &node.OutputNodes); err != nil {
			return nil, err
		}

		flow.Nodes = append(flow.Nodes, node)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &flow, nil
}
