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
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("install called")
        name := args[0]

        fmt.Printf("Intalling %q...", name) 


        repoURL := "https://github.com/rgolangh/podman-quadlets"
        tmpDir, err := os.MkdirTemp("", "pq")
        if err != nil {
            fmt.Errorf("eeeeeeeeeeeee \n")
            return err
        }
       fmt.Println("tmp dir name " + tmpDir)
       err = downloadDirectory(repoURL, name, tmpDir)
       if err != nil {
           return err
       }
       
        return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")




}


func downloadDirectory(repoURL, directoryPath, destinationPath string) error {
    fmt.Println("cloning repo")
	// Clone the repository
	_, err := git.PlainClone(destinationPath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

    filesWritten := false
    filepath.Walk(filepath.Join(destinationPath, directoryPath),
        func(path string, info fs.FileInfo, err error) error {
            fmt.Printf("walking the directory %v. workfing on file %v\n", path,  info.Name())
            switch ext := filepath.Ext(info.Name()); ext {
            //TODO need to copy folder strucure if exists. Like if there's a/foo.yaml
            // which the foo.kube points at in Yaml=a/foo.yaml
            case ".container", ".kube", ".volume", ".network", ".image", ".yaml":
                fmt.Printf("handling file %s\n", ext)
               
                configDir, err := os.UserConfigDir()
                 if err != nil {
                    log.Print("failed reading user config dir")
                    log.Fatal(err)
                }
                dest := filepath.Join(configDir, "containers", "systemd")

                bytesRead, err := os.ReadFile(path)
                if err != nil {
                    log.Print("failed reading file ")
                    log.Fatal(err)
                }

                err = os.WriteFile(filepath.Join(dest, info.Name()), bytesRead, 0644)
                if err != nil {
                    log.Fatal(err)
                }
                filesWritten = true
            default:
                fmt.Printf("ignoring %v...\n", ext)
            }
           return nil
    })
    if (filesWritten) {
        fmt.Println("Finisihed writing files")
        fmt.Println("Reloading systemd daemon for the current user")
        cmd := exec.Command("systemctl", "daemon-reload", "--user")
        if err := cmd.Run(); err != nil {
            log.Println("Failed reloading systemctal daemon-reload")
            return err
        }
    } else {
        fmt.Println("Finished without writing files")
    }

	return nil
}
