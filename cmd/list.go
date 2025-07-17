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
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/go-git/go-git/v6"
	"github.com/rgolangh/pq/pkg/quadlet"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the available quadlets",
	Long:  `List the available quadlets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if installed {
			quadlets := quadlet.ListQuadlets()
			for _, quadlet := range quadlets {
				log.Infof("- %v\n", quadlet.Name)
			}
			return nil
		}
		log.Info("Listing quadlets from ", repoURL)
		log.Info("")
		log.Debug("cloning repo ", repoURL)
		workDir, err := os.MkdirTemp("", "pq")
		if err != nil {
			return err
		}

		// Clone the repository
		_, err = git.PlainClone(workDir, &git.CloneOptions{
			Depth: 1,
			URL:   repoURL,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %v", err)
		}

		workDir = path.Join(workDir, repoSubdir)
		filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				rel, err := filepath.Rel(workDir, path)
				if err != nil {
					return err
				}
				if rel == "." {
					// depth 0
					return nil
				}
				if strings.Count(rel, string(filepath.Separator)) >= 1 {
					// scanning only 1st level directories
					log.Debug("skipping dir ", rel)
					return fs.SkipDir
				}
				if rel[0] == '.' || info.Name()[0] == '.' {
					// skip hidden dirs
					log.Debug("skipping ", info.Name())
					return fs.SkipDir
				}
				log.Info("- ", info.Name())
			}
			return nil
		})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(
		&repoURL,
		"repo",
		"r",
		"https://github.com/rgolangh/podman-quadlets",
		"The repo url (currently only git support), where the quadlets are stored")
	listCmd.Flags().BoolVar(
		&installed,
		"installed",
		false,
		"list only the installed quadlets")
}
