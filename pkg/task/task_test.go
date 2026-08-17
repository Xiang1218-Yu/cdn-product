package task

import "testing"

func TestTaskMetadataAndDependencies(t *testing.T) {
	task := NewTask("daily", "echo ok", "* * * * *")
	task.SetParam("region", "ap-northeast")
	task.AddDependency("prepare")
	if task.Params["region"] != "ap-northeast" || len(task.DependsOn) != 1 || task.DependsOn[0] != "prepare" {
		t.Fatalf("task metadata was not retained: %#v", task)
	}
	dag := NewDAG()
	parent := NewTask("prepare", "true", "")
	parent.ID = "prepare"
	if err := dag.AddTask(parent); err != nil {
		t.Fatal(err)
	}
	if err := dag.AddTask(task); err != nil {
		t.Fatal(err)
	}
	if got := dag.GetReadyTasks(); len(got) != 1 || got[0].ID != "prepare" {
		t.Fatalf("unexpected ready tasks: %#v", got)
	}
	if err := dag.UpdateTaskStatus("prepare", StatusSuccess); err != nil {
		t.Fatal(err)
	}
	if got := dag.GetReadyTasks(); len(got) != 1 || got[0].ID != task.ID {
		t.Fatalf("dependent task was not released: %#v", got)
	}
}
