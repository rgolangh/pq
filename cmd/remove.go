/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/log-go"
	"github.com/rgolangh/pq/pkg/systemd"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a quadlet",
	Long:    "Remove a quadlet",
	Aliases: []string{"uninstall", "rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debugf("remove called with args %v", args)
		name := args[0]
		quadlets := listInstalled()
		log.Debugf("installed quadlets %v", quadlets)
		for _, q := range quadlets {
			log.Debugf("installed quadlet name %q", q.name)
			if name == q.name {
				// FIX protect from symlink or going out of the installed dir
				var confirm string
				fmt.Printf("Remove quadlet %q from path %s?[y/n]", q.name, q.path)
				fmt.Scanln(&confirm)
				if confirm == "y" {
					os.RemoveAll(q.path)
					log.Infof("removed %q from path %s\n", q.name, q.path)
					systemd.DaemonReload()
				}
		        return nil
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
