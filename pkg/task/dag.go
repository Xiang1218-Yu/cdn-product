package task

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCycleDependency  = errors.New("cycle dependency detected")
	ErrTaskNotFound     = errors.New("task not found")
	ErrDependencyNotMet = errors.New("dependency not met")
)

type DAG struct {
	nodes map[string]*Task
	edges map[string][]string
	mu    sync.RWMutex
}

func NewDAG() *DAG {
	return &DAG{
		nodes: make(map[string]*Task),
		edges: make(map[string][]string),
	}
}

func (d *DAG) AddTask(task *Task) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	previousTask, hadTask := d.nodes[task.ID]
	previousEdges, hadEdges := d.edges[task.ID]
	stagedEdges := append([]string(nil), task.DependsOn...)

	d.nodes[task.ID] = task
	d.edges[task.ID] = stagedEdges
	if d.hasCycle() {
		if hadTask {
			d.nodes[task.ID] = previousTask
		} else {
			delete(d.nodes, task.ID)
		}
		if hadEdges {
			d.edges[task.ID] = previousEdges
		} else {
			delete(d.edges, task.ID)
		}
		return ErrCycleDependency
	}

	return nil
}

func (d *DAG) GetTask(id string) (*Task, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	task, exists := d.nodes[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (d *DAG) GetReadyTasks() []*Task {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var readyTasks []*Task
	for taskID, task := range d.nodes {
		if task.Status != StatusPending {
			continue
		}

		if d.areDependenciesMet(taskID) {
			readyTasks = append(readyTasks, task)
		}
	}

	return readyTasks
}

func (d *DAG) areDependenciesMet(taskID string) bool {
	dependencies := d.edges[taskID]
	for _, depID := range dependencies {
		depTask, exists := d.nodes[depID]
		if !exists || depTask.Status != StatusSuccess {
			return false
		}
	}
	return true
}

func (d *DAG) hasCycle() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for nodeID := range d.nodes {
		if d.detectCycle(nodeID, visited, recStack) {
			return true
		}
	}
	return false
}

func (d *DAG) detectCycle(nodeID string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	for _, depID := range d.edges[nodeID] {
		if !visited[depID] {
			if d.detectCycle(depID, visited, recStack) {
				return true
			}
		} else if recStack[depID] {
			return true
		}
	}

	recStack[nodeID] = false
	return false
}

func (d *DAG) UpdateTaskStatus(taskID string, status TaskStatus) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	task, exists := d.nodes[taskID]
	if !exists {
		return ErrTaskNotFound
	}

	task.Status = status
	task.UpdatedAt = time.Now()
	return nil
}
