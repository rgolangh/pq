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
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/log-go"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a quadlet definition from the remote repository",
	Long: `Inspect a quadlet definition from the remote repository
to help spotting the content details, spot potential collisions or breakage the installation
of the quadlet may cause (like port usages, volume mounts and so on).`,
	Aliases: []string{"show"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing quadlet name")
		}
		quadletName := args[0]
		tmpDir, err := os.MkdirTemp("", "pq")
		if err != nil {
			return err
		}
		log.Infof("Inspect quadlet %q from repo %s", quadletName, repoURL)
		log.Debug("tmp dir name " + tmpDir)
		err = outputQuadlet(repoURL, quadletName, tmpDir, os.Stdout)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	inspectCmd.Flags().StringVarP(
		&repoURL,
		"repo",
		"r",
		"https://github.com/rgolangh/podman-quadlets",
		"The repo url (currently only git support), where the quadlets are stored")

}

func outputQuadlet(repoURL, quadletName, downloadPath string, out io.Writer) error {
	log.Debug("cloning repo")
	// Clone the repository
	_, err := git.PlainClone(downloadPath, &git.CloneOptions{
		Depth: 1,
		URL:   repoURL,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	srcPath := filepath.Join(downloadPath, repoSubdir, quadletName)
	log.Debugf("showing quadlet in path %s \n", srcPath)

	// Read all the entries in the source directory
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}

	// Loop through all the entries
	for _, entry := range entries {
		// If it's a directory ignore for now
		if entry.IsDir() {
			continue
		} else {
			log.Debugf("Reading entry %s \n", entry.Name())
			// write file to out
			f, err := os.ReadFile(filepath.Join(srcPath, entry.Name()))
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "# Source: %s %s/%s\n", repoURL, quadletName, entry.Name())
			fmt.Fprintln(out, string(f))
		}
	}

	return nil
}
