// Copyright (c) 2022 EPAM Systems, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

var stackFilter []string

func init() {
	lsCmd.Flags().StringSliceVar(&stackFilter, "filter", []string{}, "Filter by name, status or initiator. Example: --filter \"name=GKE,status=incomplete\"")
	stateCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all stacks within a project",
	Run:   ls,
}

func ls(cmd *cobra.Command, args []string) {
	if Project == "" {
		altProjectSources()
	}

	filterQuery := make(map[string]string, 0)
	for _, param := range stackFilter {
		vals := strings.Split(param, "=")
		if len(vals) == 2 {
			key := vals[0]
			if !strings.HasPrefix(key, "latestOperation") {
				key = strings.ToLower(vals[0])
				if key == "initiator" {
					key = "latestOperation.initiator"
				}
			}
			filterQuery[key] = vals[1]
		}
	}

	req, err := NewRequest()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp, err := req.
		SetQueryParams(filterQuery).
		Get(baseURL())
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if resp.IsSuccess() {
		if Out == JsonO {
			var pretty bytes.Buffer
			err = json.Indent(&pretty, resp.Body(), "", "\t")
			if err != nil {
				fmt.Printf("Failed indent json body: %s", err)
				return
			}
			fmt.Println(pretty.String())
			return
		}
		var states []State
		json.Unmarshal(resp.Body(), &states)
		limit := 20
		if len(states) < limit {
			limit = len(states)
		}
		if len(states) > 0 {
			tableFmtStates(states[:limit])
			return
		}
		fmt.Println("Nothing has been found")
		return
	}
	fmt.Printf("Error: %s", resp.Status())
	return
}

func tableFmtStates(states []State) {
	fmt.Printf("Listing Stacks in [%s] GCP project", Project)
	fmt.Println()
	fmt.Println()
	headerFmt := color.New().SprintfFunc()
	tbl := table.New("ID", "NAME", "STATUS", "UPDATED", "INITIATOR")
	tbl.WithHeaderFormatter(headerFmt)
	for _, widget := range states {
		tbl.AddRow(widget.ID, widget.Name, widget.Status, timeago.English.Format(widget.LatestOP.Timestamp), widget.LatestOP.Initiator)
	}
	tbl.Print()
}
