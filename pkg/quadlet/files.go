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
	Name  string
	Path  string
	Files []QuadletFile
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

// ListQuadlets list all quadlet files from the default quadlet dir
// returns a map of quadle name to quadlet files
func ListQuadlets() map[string]Quadlet {
	log.Debugf("about to walk the install dir %s\n", installDir)
	quadletsByName := make(map[string]Quadlet)
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
				quadlet := Quadlet{
					Name:  dirEntry.Name(),
					Path:  path,
					Files: []QuadletFile{},
				}
				for _, de := range entries {
					if !de.IsDir() {
						log.Debugf("quadlet file %s", de.Name())
						i := strings.LastIndex(de.Name(), ".")
						if i < 0 {
							continue
						}
						switch de.Name()[i:] {
						case "container", "pod", "kube":
							log.Debug("container file")
						case "volume":
							log.Debug("container file")
						case "network":
							log.Debug("network file")
						default:
							log.Debug("some other file")
						}

						qf, err := newQuadletFile(path, de)
						if err != nil {
							return err
						}
						log.Debugf("converted quadlet file %+v", qf)
						quadlet.Files = append(quadlet.Files, qf)
					}
				}
				quadletsByName[dirEntry.Name()] = quadlet
			}
			return nil
		},
	)
	return quadletsByName
}

func newQuadletFile(path string, de fs.DirEntry) (QuadletFile, error) {
	iniFile, err := ini.Load(filepath.Join(path, de.Name()))
	log.Debugf("ini file loaded")
	if err != nil {
		return QuadletFile{}, err
	}
	qf := QuadletFile{}
	iniFile.MapTo(&qf)
	qf.FileName = filepath.Join(path, de.Name())
	log.Debugf("quadlet file %+v", qf)
	return qf, nil

}
