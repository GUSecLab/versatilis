package versatilis

import (
	_ "embed"

	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/semver"
)

//go:embed version.md
var versionString string

var Version *semver.Version

type VersatilisErrNotImplemented struct{}

type Versatilis struct {
	Version *semver.Version
	Log     *log.Logger
}

func (e *VersatilisErrNotImplemented) Error() string {
	return "not yet implemented"
}

func New() *Versatilis {
	var err error
	var v Versatilis
	v.Version, err = semver.NewVersion(versionString)
	if err != nil {
		panic(err)
	}
	v.Log = log.New()
	v.Log.SetLevel(log.InfoLevel)
	return &v
}
