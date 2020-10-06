// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"log"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/names"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

func init() {
	countCmd.AddCommand(countVersionsCmd)
}

var countVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Count the number of versions of specified APIs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.TODO()
		client, err := connection.NewClient(ctx)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		name := args[0]
		if m := names.ApiRegexp().FindStringSubmatch(name); m != nil {
			err := core.ListAPIs(ctx, client, m, countFilter, func(api *rpc.Api) {
				count := 0
				request := &rpc.ListVersionsRequest{
					Parent: api.Name,
				}
				it := client.ListVersions(ctx, request)
				for {
					_, err := it.Next()
					if err == iterator.Done {
						break
					} else if err == nil {
						count++
					}
				}
				log.Printf("%d\t%s", count, api.Name)
			})
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
		}
	},
}