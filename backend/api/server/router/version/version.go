package version

import (
	"context"
	"net/http"
	"runtime"

	"github.com/konflux-ci/quality-studio/pkg/utils/httputils"
)

var (
	ServerVersion = "1.1.0"
	// apiMaturity is the level of maturity the Server has achieved at this version, eg. planning, pre-alpha, alpha, beta, stable, mature, inactive, or deprecated.
	APIMaturity = "v1alpha1"
	// gitCommit is a constant representing the source version that
	// generated this build. It should be set during build via -ldflags.
	GitCommit string
	// buildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	BuildDate string
)

// Define the server version information
type Version struct {
	ServerVersion string `json:"serverAPIVersion"`
	APIMaturity   string `json:"apiMaturity"`
	GitCommit     string `json:"gitCommit"`
	BuildDate     string `json:"buildDate"`
	GoOs          string `json:"goOs"`
	GoArch        string `json:"goArch"`
	GoVersion     string `json:"goVersion"`
}

// version godoc
// @Summary API Server info
// @Description returns the Server information as a JSON
// @Tags Server API
// @Accept json
// @Produce json
// @Router /server/info [get]
// @Success 200 {object} Version
func (s *versionRouter) getVersion(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return httputils.WriteJSON(w, http.StatusOK, Version{
		ServerVersion: ServerVersion,
		APIMaturity:   APIMaturity,
		GitCommit:     GitCommit,
		BuildDate:     BuildDate,
		GoOs:          runtime.GOOS,
		GoArch:        runtime.GOARCH,
		GoVersion:     runtime.Version(),
	})
}
