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

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

func init() {
	stateCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show <Stack ID>",
	Short: "Show details of a stack",
	Run:   show,
	Args:  cobra.ExactValidArgs(1),
}

func show(cmd *cobra.Command, args []string) {
	if Project == "" {
		altProjectSources()
	}

	req, err := NewRequest()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp, err := req.
		Get(fmt.Sprintf("%s/%s", baseURL(), args[0]))
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Println()
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

		var state State
		err = json.Unmarshal(resp.Body(), &state)
		if err != nil {
			fmt.Printf("Failed unmarshal json body: %s", err)
			return
		}
		tableFmtState(state)
		return
	}
	if resp.StatusCode() == 404 {
		fmt.Printf("Error: State \"%s\" not found", args[0])
		return
	}
	fmt.Printf("Error: %s", resp.Status())
	return
}

func tableFmtState(state State) {
	fmt.Printf("Showing details of [%s]", state.ID)
	fmt.Println()
	tbl := table.New("", "")
	tbl.AddRow("Stack ID", state.ID)
	tbl.AddRow("Stack Name", state.Name)
	tbl.AddRow("Latest Status", state.Status)
	tbl.AddRow("Last Updated", timeago.English.Format(state.LatestOP.Timestamp))
	tbl.AddRow("Initiator", state.LatestOP.Initiator)
	tbl.AddRow("State File Location", state.StateLocation.Uri)
	tbl.Print()
	fmt.Println()
	tblComponents := table.New("COMPONENT NAME", "STATUS")
	headerFmt := color.New().SprintfFunc()
	tblComponents.WithHeaderFormatter(headerFmt)
	for _, widget := range state.Components {
		tblComponents.AddRow(widget.Name, widget.Status)
	}
	tblComponents.Print()
}
