package main

import "testing"

func TestDrawLineWithMessage(t *testing.T) {
	drawLineWithMessage("t")
	drawLineWithMessage("te")
	drawLineWithMessage("tes")
	drawLineWithMessage("test")
	drawLineWithMessage("testtesttest")
	drawLineWithMessage("testtesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttesttest")
}

func TestContains(t *testing.T) {

	envs := []Environment{
		Environment{Name: "env1", Value: "val1"},
		Environment{Name: "env2", Value: "val2"},
	}
	if contains(envs, "env1") != true {
		t.Errorf("should contain")
	}
	if contains(envs, "env3") != false {
		t.Errorf("should not contain")
	}
}

func TestDecideAction(t *testing.T) {
	envs1 := []Environment{
		Environment{Name: "env1", Value: "val1"},
		Environment{Name: "env2", Value: "val2"},
		Environment{Name: "env3", Value: "val3"},
		Environment{Name: "env6", Value: "val6"},
	}
	envs2 := []Environment{
		Environment{Name: "env1", Value: "val1"},
		Environment{Name: "env2", Value: "val2"},
		Environment{Name: "env4", Value: "val4"},
		Environment{Name: "env5", Value: "val5"},
		Environment{Name: "env6", Value: "val6"},
	}

	actions := decideAction(envs1, envs2)
	for _, v := range actions {
		switch v.Environment.Name {
		case "env1", "env2", "env6":
			if v.WillBe != "updated" {
				t.Errorf("should be updated")
			}
		case "env3":
			if v.WillBe != "created" {
				t.Errorf("should be created")
			}
		case "env4", "env5":
			if v.WillBe != "deleted" {
				t.Errorf("should be deleted")
			}
		}
	}
}
