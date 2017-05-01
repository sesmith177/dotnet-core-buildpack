package compile_test

import (
	"bytes"
	"io/ioutil"
	"os"

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
		It("exists", func() {
			err = compiler.RestoreCache()
			Expect(err).To(BeNil())
		})
	})
})
