package compile_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"dotnet-core/compile"

	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/ansicleaner"
)

var _ = Describe("Compile", func() {
	var (
		buildDir string
		cacheDir string
		err      error
		compiler compile.Compiler
		stager   *libbuildpack.Stager
		logger   libbuildpack.Logger
		buffer   *bytes.Buffer
	)

	BeforeEach(func() {
		buildDir, err = ioutil.TempDir("", "dotnet-core-buildpack.build")
		Expect(err).To(BeNil())

		cacheDir, err = ioutil.TempDir("", "dotnet-core-buildpack.cache")
		Expect(err).To(BeNil())

		buffer = new(bytes.Buffer)

		logger = libbuildpack.NewLogger()
		logger.SetOutput(ansicleaner.New(buffer))
	})

	JustBeforeEach(func() {
		stager = &libbuildpack.Stager{
			BuildDir: buildDir,
			CacheDir: cacheDir,
			Log:      logger,
		}

		compiler = compile.Compiler{
			Stager: stager,
		}
	})

	AfterEach(func() {
		err = os.RemoveAll(buildDir)
		Expect(err).To(BeNil())

		err = os.RemoveAll(cacheDir)
		Expect(err).To(BeNil())
	})

	Describe("RestoreCache", func() {
		var depNames []string

		BeforeEach(func() {
			depNames = []string{".dotnet", ".node", "libunwind"}

			for _, name := range depNames {
				err = os.MkdirAll(filepath.Join(cacheDir, name), 0755)
				Expect(err).To(BeNil())

			}

			err = os.MkdirAll(filepath.Join(cacheDir, "something-else"), 0755)
			Expect(err).To(BeNil())
		})

		It("logs that it is resoring the cache", func() {
			err = compiler.RestoreCache()
			Expect(err).To(BeNil())

			Expect(buffer.String()).To(ContainSubstring("-----> Restoring files from buildpack cache"))
		})

		It("moves the correct dependencies to <buildDir>", func() {
			err = compiler.RestoreCache()
			Expect(err).To(BeNil())

			for _, name := range depNames {
				Expect(filepath.Join(buildDir, name)).To(BeADirectory())
			}

			Expect(filepath.Join(buildDir, "something-else")).NotTo(BeADirectory())
		})
	})

	Describe("ClearNugetCache", func() {
		Context("nuget cache does not exist", func() {
			It("Logs nothing", func() {
				err = compiler.ClearNugetCache()
				Expect(err).To(BeNil())

				Expect(buffer.String()).To(Equal(""))
			})
		})

		Context("nuget cache exists", func() {
			BeforeEach(func() {
				err = os.MkdirAll(filepath.Join(cacheDir, ".nuget"), 0755)
				Expect(err).To(BeNil())
			})

			Context("CACHE_NUGET_PACKAGES is false", func() {
				var oldCacheNugetPackages string

				BeforeEach(func() {
					oldCacheNugetPackages = os.Getenv("CACHE_NUGET_PACKAGES")
					err = os.Setenv("CACHE_NUGET_PACKAGES", "false")
					Expect(err).To(BeNil())
				})

				AfterEach(func() {
					err = os.Setenv("CACHE_NUGET_PACKAGES", "false")
					Expect(err).To(BeNil())
				})

				It("logs a message that the nuget cache is being cleared", func() {
					err = compiler.ClearNugetCache()
					Expect(err).To(BeNil())

					Expect(buffer.String()).To(ContainSubstring("-----> Clearing NuGet packages cache"))
				})

				It("clears the nuget package cache", func() {
					err = compiler.ClearNugetCache()
					Expect(err).To(BeNil())

					Expect(filepath.Join(cacheDir, ".nuget")).NotTo(BeADirectory())
				})
			})
		})
	})
})