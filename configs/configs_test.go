package configs

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

func TestConfigs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config testing")
}

var _ = Describe("Configs", func() {
	Describe("When writing the original dotfile", func() {
		Context("Should be able to pick up on when one is there", func() {
			It("no err", func() {
				// First write the test file
				rawBytes, err := yaml.Marshal(&config{LogFile: "logFile"})
				Expect(err).ShouldNot(HaveOccurred())
				testConfigName := "firstTestConfig.yml"
				Expect(ioutil.WriteFile(testConfigName, rawBytes, 0744)).ShouldNot(HaveOccurred())

				// then test the basic
				err = loadDotfileConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(STDConf.LogFile).Should(Equal("logFile"))

				// then remove the file at the end
				Expect(os.Remove(testConfigName)).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("when creating the default dotfile config", func() {
		Context("Should match some basic configs", func() {
			It("no err", func() {
				// write default config
				testConfigName := "secondTestConfig.yml"
				Expect(writeDefaultConfig(testConfigName)).ShouldNot(HaveOccurred())

				bytes, err := ioutil.ReadFile(testConfigName)
				Expect(err).ShouldNot(HaveOccurred())

				var testConfig *config
				Expect(yaml.UnmarshalStrict(bytes, &testConfig)).ShouldNot(HaveOccurred())
				Expect(testConfig).ShouldNot(BeNil())
				Expect(testConfig.LogFile).Should(Equal("logFile"))
				Expect(os.Remove(testConfigName)).ShouldNot(HaveOccurred())
			})
		})
	})
})
