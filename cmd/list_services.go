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
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/rgolangh/pq/pkg/quadlet"
	"github.com/rgolangh/pq/pkg/systemd"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listServicesCmd = &cobra.Command{
	Use:   "list-services",
	Short: "List the systemd services from installed quadlets",
	RunE: func(cmd *cobra.Command, args []string) error {
		quadletsByName := quadlet.ListQuadlets()
		log.Debugf("quadlets by name %+v", quadletsByName)
		for _, quadlet := range quadletsByName {
			for _, qf := range quadlet.Files {
				if strings.HasSuffix(qf.FileName, ".container") {
					// .container file is generated as .service file
					svc := strings.Replace(filepath.Base(qf.FileName), ".container", ".service", 1)
					unitStatus, err := systemd.Status(svc)
					if err != nil {
						return err
					}
					color := ""
					if unitStatus.ActiveState == "active" {
						color = systemd.Green
					}
					log.Infof("%s - %s %s (%s)",
						quadlet.Name,
						svc,
						fmt.Sprintf("%s%s%s", color, unitStatus.ActiveState, systemd.RESET),
						fmt.Sprintf("%s%s%s", color, unitStatus.SubState, systemd.RESET),
					)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listServicesCmd)
}
