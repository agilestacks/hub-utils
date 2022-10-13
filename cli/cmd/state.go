// Copyright (c) 2022 EPAM Systems, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var Project string
var StateAPILocation string
var Verbose bool
var Out Output

type Output string

const (
	JsonO  Output = "json"
	TableO Output = "table"
)

func (e *Output) String() string {
	return string(*e)
}

func (e *Output) Set(v string) error {
	switch v {
	case "table", "json":
		*e = Output(v)
		return nil
	default:
		return errors.New(`must be one of "table" or "json"`)
	}
}

func (e *Output) Type() string {
	return "string"
}

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "HubCTL remote state viewer and manager",
}

func Execute() {
	cobra.CheckErr(stateCmd.Execute())
}

func init() {
	stateCmd.PersistentFlags().StringVarP(&StateAPILocation, "stateAPILocation", "l", "us-central1", "Location of State API endpoint")
	stateCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose")
	stateCmd.PersistentFlags().VarP(&Out, "output", "o", "Output format. Must be one of [table, json]")
	stateCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "GCP Project ID")
}

func altProjectSources() {
	Project = os.Getenv("GOOGLE_PROJECT")
	if Project != "" {
		return
	}
	cmd := exec.Command("gcloud", "config", "get-value", "core/project")
	stdout, _ := cmd.Output()
	Project = strings.TrimSuffix(string(stdout), "\n")
	if Project == "" {
		if Out == JsonO {
			err := map[string]string{"error": "GCP Project ID is not set"}
			msg, _ := json.Marshal(err)
			fmt.Println(string(msg))
		} else {
			fmt.Println("GCP Project ID is not set. Please do one of the following:")
			fmt.Println("* re-run the command with --project flag")
			fmt.Println("* set GOOGLE_PROJECT env variable")
			fmt.Println("* set the Project ID using `gcloud config set project <project-id>` command")
		}
		os.Exit(1)
	}
}
