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
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/rgolangh/pq/pkg/quadlet"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listServicesCmd = &cobra.Command{
	Use:   "list-services",
	Short: "List the istalled services from quadlets",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := quadlet.ListQuadletFiles()
		for _, s := range services {
			if strings.HasSuffix(s.FileName, ".container") {
				svc := strings.Replace(filepath.Base(s.FileName), ".container", ".service", 1)
				log.Infof(" - %s", svc)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listServicesCmd)
}
