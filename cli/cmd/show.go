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
	"github.com/go-resty/resty/v2"
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
	base := BaseURL(StateAPILocation, Project)
	client := resty.New()
	resp, err := client.R().
		Get(fmt.Sprintf("%s/%s", base, args[0]))
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Println()
		return
	}
	if resp.StatusCode() == 404 && Out != JsonO {
		fmt.Printf("Error: Stack [%s] has not been found", args[0])
		fmt.Println()
		return
	}
	if Out == JsonO {
		var pretty bytes.Buffer
		err = json.Indent(&pretty, resp.Body(), "", "\t")
		if err != nil {
			return
		}
		fmt.Println(pretty.String())
		return
	}
	var state State
	err = json.Unmarshal(resp.Body(), &state)
	if err != nil {
		fmt.Printf("Error: Stack [%s] has not been found", args[0])
		fmt.Println()
		return
	}
	tableFmtState(state)
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
