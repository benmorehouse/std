package configs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
	userOS "os/user"
)

func TestConfigs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config testing")
}

var _ = Describe("Configs", func() {

	var (
		usr            *userOS.User
		err            error
		rootDir        string
		testConfigName string
	)

	BeforeEach(func() {
		usr, err = userOS.Current()
		Expect(err).ShouldNot(HaveOccurred())
		rootDir = usr.HomeDir
		testConfigName = filepath.Join(rootDir, ".std/.testing/logging/std.log")
		os.Setenv(testingConfigEnv, "true")
		err = setConfigWithUserRoot()
		Expect(err).ShouldNot(HaveOccurred())
		err = makePaths()
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("When writing the original dotfile", func() {
		Context("Should be able to pick up on when one is there", func() {
			It("no err", func() {
				rawBytes, err := yaml.Marshal(&config{LogFile: "logFile"})
				Expect(err).ShouldNot(HaveOccurred())

				err = ioutil.WriteFile(testConfigName, rawBytes, 0744)
				Expect(err).ShouldNot(HaveOccurred())

				err = loadDotfileConfig()
				Expect(err).ShouldNot(HaveOccurred())

				Expect(STDConf.LogFile).Should(Equal(testConfigName))

				err = os.Remove(testConfigName)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("when creating the default dotfile config", func() {
		Context("Should match some basic configs", func() {
			It("no err", func() {
				// write default config
				err = writeDefaultConfig(testConfigName)
				Expect(err).ShouldNot(HaveOccurred())

				bytes, err := ioutil.ReadFile(testConfigName)
				Expect(err).ShouldNot(HaveOccurred())

				var testConfig *config
				Expect(yaml.UnmarshalStrict(bytes, &testConfig)).ShouldNot(HaveOccurred())
				Expect(testConfig).ShouldNot(BeNil())
				Expect(testConfig.LogFile).Should(Equal(testConfigName))

				err = os.Remove(testConfigName)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
