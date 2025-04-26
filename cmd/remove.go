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
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/rgolangh/pq/pkg/quadlet"
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
		quadletsByName := quadlet.ListQuadlets()
		log.Debugf("installed quadlets %v", quadletsByName)
		if quadlet, ok := quadletsByName[name]; ok && len(quadlet.Files) > 0 {
			log.Debugf("quadlet files %v", quadlet.Files)
			var confirm string
			for _, svc := range quadlet.Files {
				if filepath.Ext(svc.FileName) == ".container" {
					systemd.Stop(strings.Replace(filepath.Base(svc.FileName), ".container", ".service", 1))
				}
			}
			fmt.Printf("Remove quadlet %q from path %s? [y/N] ", quadlet.Name, quadlet.Path)
			fmt.Scanln(&confirm)
			if confirm == "y" {
				os.RemoveAll(quadlet.Path)
				log.Infof("removed %q from path %s\n", quadlet.Name, quadlet.Path)
				systemd.DaemonReload()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
