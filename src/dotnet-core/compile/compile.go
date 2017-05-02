package compile

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Compiler struct {
	Stager *libbuildpack.Stager
}

func Run(dc *Compiler) error {
	if err := dc.RestoreCache(); err != nil {
		dc.Stager.Log.Error("Unable to restore buildpack cache: %s", err.Error())
		return err
	}

	if err := dc.ClearNugetCache(); err != nil {
		dc.Stager.Log.Error("Unable to clear NuGet packages cache: %s", err.Error())
		return err
	}

	return nil
}

func (dc *Compiler) RestoreCache() error {
	dc.Stager.Log.BeginStep("Restoring files from buildpack cache")

	depDirs := []string{".dotnet", ".node", "libunwind"}

	for _, name := range depDirs {
		err := os.Rename(filepath.Join(dc.Stager.CacheDir, name), filepath.Join(dc.Stager.BuildDir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func (dc *Compiler) ClearNugetCache() error {
	nugetCacheDir := filepath.Join(dc.Stager.CacheDir, ".nuget")

	nugetCacheExists, err := libbuildpack.FileExists(nugetCacheDir)
	if err != nil {
		return err
	}

	if !nugetCacheExists {
		return nil
	}

	if os.Getenv("CACHE_NUGET_PACKAGES") == "false" {
		dc.Stager.Log.BeginStep("Clearing NuGet packages cache")
		return os.RemoveAll(nugetCacheDir)
	}

	return nil
}
