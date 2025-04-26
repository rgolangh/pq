package systemd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Masterminds/log-go"
	"gopkg.in/ini.v1"
)

const (
	RESET = "\033[0m"
	Green = "\033[32m"
)

var UserFlag string

func init() {
	if os.Getuid() != 0 {
		UserFlag = "--user"
	}

}

type UnitStatus struct {
	ActiveState string `ini:"ActiveState"`
	LoadState   string `ini:"LoadState"`
	SubState    string `ini:"SubState"`
	Description string `ini:"Description"`
}

func DaemonReload() error {
	var confirm string
	fmt.Printf("Reload systemd daemon? [y/N] ")
	fmt.Scanln(&confirm)
	if confirm == "y" {
		log.Infof("Reloading systemd daemon for the current user(uid: %d)", os.Getuid())
		args := []string{"daemon-reload"}
		if UserFlag != "" {
			args = append(args, UserFlag)
		}
		cmd := exec.Command("systemctl", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s %w", out, err)
		}
	}
	return nil
}

func Status(serviceName string) (UnitStatus, error) {
	args := []string{"show", serviceName, "--no-page", "--property=ActiveState,SubState,LoadState,Description"}
	if UserFlag != "" {
		args = append(args, UserFlag)
	}
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to get the status service %s with error: %v", serviceName, err)
		return UnitStatus{}, err
	}
	log.Debugf("status output %s", out)

	iniFile, err := ini.Load(out)
	if err != nil {
		return UnitStatus{}, err
	}

	us := UnitStatus{}
	if err := iniFile.Section("").MapTo(&us); err != nil {
		return UnitStatus{}, err
	}
	log.Debug("status was successful", us)
	return us, nil
}

func Start(serviceName string) error {
	log.Infof("Starting service %s for current user", serviceName)
	args := []string{"start", serviceName}
	if UserFlag != "" {
		args = append(args, UserFlag)
	}
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to start service %s with error: %s %s", serviceName, err, out)
	}
	log.Debug(out)
	return nil
}

func Stop(serviceName string) error {
	log.Infof("Stopping service %s for current user", serviceName)
	args := []string{"stop", serviceName}
	if UserFlag != "" {
		args = append(args, UserFlag)
	}
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to stop service %s with error: %v", serviceName, err)
		return err
	}
	log.Debug(out)
	return nil
}

func Journal(serviceName string) error {
	if serviceName == "" {
		return fmt.Errorf("Unit name cannot be empty")
	}
	args := []string{"--unit", serviceName, "-p", "info", "--boot", "-n", "10", "--output", "cat"}
	if UserFlag != "" {
		args = append(args, UserFlag)
	}
	cmd := exec.Command("journalctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Error(string(out))
	return nil
}
