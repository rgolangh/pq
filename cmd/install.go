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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/rgolangh/pq/pkg/quadlet"
	"github.com/rgolangh/pq/pkg/systemd"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

var (
	repoURL               string
	installed             bool
	dryRun                bool
	noSystemdDaemonReload bool
	installDir            string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a quadlet from a quadlet repo",
	Long: `Download the quadlet folder by NAME and copy
it into the $HOME/.config/containers/systemd/
Files which are not supported should be cleared from the directory
All quadlet repos should have a directory structure where every quadlet is a top level directory and all the 
.service , .network files are inside.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debugf("install called with args %v", args)
		quadletName := args[0]

		tmpDir, err := os.MkdirTemp("", "pq")
		if err != nil {
			return err
		}
		log.Infof("Installing quadlet %q", quadletName)
		log.Debug("tmp dir name " + tmpDir)
		err = downloadDirectory(repoURL, quadletName, tmpDir)
		if err != nil {
			return err
		}

		if !noSystemdDaemonReload {
			systemd.DaemonReload()
		}

		quadletsByName := quadlet.ListQuadlets()
		if q, ok := quadletsByName[quadletName]; ok {
			for _, f := range q.Files {
				if filepath.Ext(f.FileName) == ".container" {
					err := systemd.Start(strings.Replace(filepath.Base(f.FileName), ".container", ".service", 1))
					if err != nil {
						return err
					}

				}

			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(
		&repoURL,
		"repo",
		"r",
		"https://github.com/rgolangh/podman-quadlets",
		"The repo url (currently only git support), where the quadlets are stored")
	installCmd.Flags().BoolVarP(
		&dryRun,
		"dry-run",
		"",
		false,
		"Don't install, just output the generated quadlet for dubugging",
	)
	installCmd.Flags().BoolVarP(
		&noSystemdDaemonReload,
		"no-systemd-daemon-reload",
		"",
		false,
		"No systemd daemon reloading after installing. Useful for controlling when to reload the deamon",
	)
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	installDir = filepath.Join(configDir, "containers", "systemd")

}

func downloadDirectory(repoURL, quadletName, downloadPath string) error {
	log.Debug("cloning repo")
	// Clone the repository
	_, err := git.PlainClone(downloadPath, false, &git.CloneOptions{
		Depth: 1,
		URL:   repoURL,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	if dryRun {
		noSystemdDaemonReload = true
		log.Debug("Install Dry Run")
		cmd := exec.Command("/usr/lib/systemd/system-generators/podman-system-generator", "--user", "--dryrun")
		cmd.Env = append(cmd.Env, "QUADLET_UNIT_DIRS="+filepath.Join(downloadPath, quadletName))
		out, err := cmd.Output()
		if err != nil {
			log.Error(err)
			return err
		}
		fmt.Fprint(os.Stdout, string(out))
		return nil
	}
	err = copyDir(filepath.Join(downloadPath, quadletName), filepath.Join(installDir, quadletName))
	if err != nil {
		log.Errorf("Error copying the directory %v\n", err)
		return err
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	log.Debugf("copying file from %v to %v", src, dst)
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

// copyDir recursively copies a directory from src to dst
func copyDir(src string, dst string) error {
	// Get properties of the source directory
	log.Debugf("copying from %v to %v\n", src, dst)
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Read all the entries in the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Loop through all the entries
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// If it's a directory, recurse
		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Otherwise, copy the file
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
