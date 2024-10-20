package systemd

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/log-go"
)

func DaemonReload() error {
	var confirm string
	fmt.Printf("Reload systemd daemon?[y/N]")
	fmt.Scanln(&confirm)
	if confirm == "y" {
		log.Info("Reloading systemd daemon for the current user")
		cmd := exec.Command("systemctl", "daemon-reload", "--user")
		out, err := cmd.Output()
		if err != nil {
			log.Error("Failed reloading systemctl daemon-reload")
			return err
		}
		log.Debug(out)
	}
	return nil
}

func Start(serviceName string) error {
	log.Infof("Starting service %s for current user", serviceName)
	cmd := exec.Command("systemctl", "start", "--user", serviceName)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to start service %s with error: %v", serviceName, err)
		return err
	}
	log.Debug(out)
	return nil
}

func Stop(serviceName string) error {
	log.Infof("Stopping service %s for current user", serviceName)
	cmd := exec.Command("systemctl", "stop", "--user", serviceName)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to stop service %s with error: %v", serviceName, err)
		return err
	}
	log.Debug(out)
	return nil
}
