package models

import ()

/**
 *  Node
 *  ====
 *  Represents a single node. I can be subtopic or topic.
 *  Each node will have `Name`, `Children`, and ``
 */
type Node struct {
	Name   string `json:"name"`
	IsTree bool   `json:"isTree"`
	Path   string `json:"path"`
}
