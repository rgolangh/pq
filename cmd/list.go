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
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

var installDir string

type quadlet struct {
	name string
	path string
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the available quadlets",
	Long:  `List the available quadlets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if installed {
			quadlets := listInstalled()
			for _, q := range quadlets {
				log.Infof(" - %v\n", q.name)
			}
			return nil
		}
		log.Info("Listing quardlets from ", repoURL)
		log.Info("")
		log.Debug("cloning repo ", repoURL)
		workDir, err := os.MkdirTemp("", "pq")
		if err != nil {
			return err
		}

		// Clone the repository
		_, err = git.PlainClone(workDir, false, &git.CloneOptions{
			Depth: 1,
			URL:   repoURL,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %v", err)
		}

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

	configDir, err := os.UserConfigDir()

	if err != nil {
		log.Error("failed reading user config dir")
		log.Fatal(err)
	}
	installDir = filepath.Join(configDir, "containers", "systemd")
}

func listInstalled() []quadlet {
	installed := []quadlet{}
	log.Debugf("about to walk the install dir %s\n", installDir)
	rootWasWalked := false
	filepath.WalkDir(
		installDir,
		func(path string, dirEntry fs.DirEntry, err error) error {
			if !rootWasWalked {
				rootWasWalked = true
				return nil
			}
			log.Debugf("dirEntry %v\n", dirEntry.Name())
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			if len(entries) > 0 {
				installed = append(installed, quadlet{name: dirEntry.Name(), path: path})
			}
			return nil
		})
	return installed
}
