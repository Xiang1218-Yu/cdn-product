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

	d.nodes[task.ID] = task
	if _, exists := d.edges[task.ID]; !exists {
		d.edges[task.ID] = []string{}
	}

	if len(task.DependsOn) > 0 {
		for _, depID := range task.DependsOn {
			d.edges[task.ID] = append(d.edges[task.ID], depID)
		}
	}

	if d.hasCycle() {
		delete(d.nodes, task.ID)
		delete(d.edges, task.ID)
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

// ClaimReadyTask atomically reserves a pending task before the scheduler starts execution.
func (d *DAG) ClaimReadyTask(taskID string) (*Task, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	task, exists := d.nodes[taskID]
	if !exists || task.Status != StatusPending || !d.areDependenciesMet(taskID) {
		return nil, false
	}

	task.Status = StatusRunning
	task.UpdatedAt = time.Now()
	return task, true
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
