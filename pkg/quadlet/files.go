package quadlet

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/log-go"
	"gopkg.in/ini.v1"
)

var installDir string

type Quadlet struct {
	Name string
	Path string
}
type QuadletFile struct {
	Unit      Unit      `ini:"Unit"`
	Container Container `ini:"Container,omitempty"`
	FileName  string
}

type Container struct {
	Image         string `ini:",omitempty"`
	ContainerName string `ini:",omitempty"`
	Exec          string `ini:",omitempty"`
	Network       string `ini:",omitempty"`
	PublishPort   string `ini:",omitempty"`
	Volume        string `ini:",omitempty"`
}

type Unit struct {
	Description string `ini:",omitempty"`
	Requires    string `ini:",omitempty"`
	Requisite   string `ini:",omitempty"`
	Wants       string `ini:",omitempty"`
	BindsTo     string `ini:",omitempty"`
	PartOf      string `ini:",omitempty"`
	Conflicts   string `ini:",omitempty"`
	Before      string `ini:",omitempty"`
	OnFailure   string `ini:",omitempty"`
	After       string `ini:",omitempty"`
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	installDir = filepath.Join(configDir, "containers", "systemd")
}

func ListQuadletFiles() []QuadletFile {
	// list files in quadlet folder

	// extract unit names

	quadletFiles := []QuadletFile{}
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
			if dirEntry.IsDir() {
				entries, err := os.ReadDir(path)
				if err != nil {
					log.Errorf("failed to read dir %v", err)
					return err
				}
				for _, de := range entries {
					log.Debugf("quadlet file %s", de.Name())
					if strings.HasSuffix(de.Name(), ".container") && !de.IsDir() {
						log.Debugf("found container file %s", de.Name())
						qf, err := newQuadletFile(path, de)
						if err != nil {
							return err
						}
						quadletFiles = append(quadletFiles, qf)
					}
				}
			}
			return nil
		},
	)
	return quadletFiles
}

func newQuadletFile(path string, de fs.DirEntry) (QuadletFile, error) {
	iniFile, err := ini.Load(filepath.Join(path, de.Name()))
	log.Debugf("ini file loaded %+v", iniFile)
	if err != nil {
		return QuadletFile{}, err
	}
	qf := QuadletFile{}
	iniFile.MapTo(&qf)
	qf.FileName = filepath.Join(path, de.Name())
	log.Debugf("container file %+v", qf)
	return qf, nil

}

func ListInstalled() []Quadlet {
	installed := []Quadlet{}
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
				installed = append(installed, Quadlet{Name: dirEntry.Name(), Path: path})
			}
			return nil
		})
	return installed
}
