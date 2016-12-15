package main

import "fmt"

func drawLineWithMessage(message string) {
	var line string
	maxLength := 100
	lineLiteral := "-"

	lineLength := maxLength - len(message)
	halfLength := (lineLength - 2) / 2

	for n := 0; n < halfLength; n++ {
		line = line + lineLiteral
	}
	line = line + " "
	line = line + message
	line = line + " "
	for n := len(line); n < maxLength; n++ {
		line = line + lineLiteral
	}
	fmt.Println(line)
}

func contains(envs []Environment, envName string) bool {
	for _, env := range envs {
		if env.Name == envName {
			return true
		}
	}
	return false
}

func decideAction(localEnvs, remoteEnvs []Environment) []Action {
	var actions []Action

	// exists remote and local: updated
	for _, v := range localEnvs {
		if contains(remoteEnvs, v.Name) {
			actions = append(actions, Action{Environment: v, WillBe: "updated", MessageFormat: "\x1b[33m%s\x1b[0m %s=%s"})
		}
	}
	// exists local only: created
	for _, v := range localEnvs {
		if !contains(remoteEnvs, v.Name) {
			actions = append(actions, Action{Environment: v, WillBe: "created", MessageFormat: "\x1b[32m%s\x1b[0m %s=%s"})
		}
	}
	// exists remote only: deleted
	for _, v := range remoteEnvs {
		if !contains(localEnvs, v.Name) {
			actions = append(actions, Action{Environment: v, WillBe: "deleted", MessageFormat: "\x1b[31m%s\x1b[0m %s=%s"})
		}
	}
	return actions
}
