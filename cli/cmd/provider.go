// Copyright (c) 2022 EPAM Systems, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/google"
)

func baseURL() string {
	return fmt.Sprintf("https://%s-%s.cloudfunctions.net/stacks", StateAPILocation, Project)
}

func NewRequest() (*resty.Request, error) {
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	ts, err := google.DefaultTokenSource(ctx, scopes...)
	if err != nil {
		fmt.Println("Failed to create new token source")
		return nil, err
	}
	token, err := ts.Token()
	if err != nil {
		fmt.Println("Failed to get token")
		return nil, err
	}

	request := resty.New().
		SetDebug(Verbose).
		NewRequest().
		SetAuthScheme(token.TokenType).
		// Use identity token to invoke Cloud Function
		// https://cloud.google.com/functions/docs/securing/authenticating#authenticating_developer_testing
		SetAuthToken(token.Extra("id_token").(string))

	return request, nil
}
