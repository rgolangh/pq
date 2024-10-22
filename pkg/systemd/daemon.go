package systemd

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/log-go"
	"gopkg.in/ini.v1"
)

const RESET = "\033[0m"
const Green = "\033[32m"

type UnitStatus struct {
	ActiveState string `ini:"ActiveState"`
	LoadState   string `ini:"LoadState"`
	SubState    string `ini:"SubState"`
	Description string `ini:"Description"`
}

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

func Status(serviceName string) (UnitStatus, error) {
	cmd := exec.Command("systemctl", "show", "--user", serviceName, "--no-page", "--property=ActiveState,SubState,LoadState,Description")
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
