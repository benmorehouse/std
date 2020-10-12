package repo

import (
	"github.com/benmorehouse/std/configs"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepoCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repo Functionality")
}

var _ = Describe("Create Command", func() {
	var (
		passwordConnector, listConnector Connector
		err                              error
		passwordRepo, listRepo           Repo
		currentList                      string
		stuffToDo                        string
	)

	BeforeEach(func() {
		os.Setenv("STD_CONFIG_TESTING", "true")

		err = configs.SetConfigWithUserRoot()
		Expect(err).ShouldNot(HaveOccurred())

		passwordConnector = PasswordConnector()
		listConnector = ListConnector()

		passwordRepo, err = passwordConnector.Connect()
		Expect(err).ShouldNot(HaveOccurred())

		listRepo, err = listConnector.Connect()
		Expect(err).ShouldNot(HaveOccurred())

		currentList = "todo-list"
		stuffToDo = "clean the dishes"
	})

	AfterEach(func() {
		Expect(listConnector.Disconnect()).ShouldNot(HaveOccurred())
		Expect(passwordConnector.Disconnect()).ShouldNot(HaveOccurred())
	})

	Context("Should be able to run a full lifecycle for lists", func() {
		It("and act accordingly", func() {
			// First put
			err = listRepo.Put(currentList, stuffToDo)
			Expect(err).ShouldNot(HaveOccurred())

			// then get
			Expect(listRepo.Get(currentList)).To(Equal(stuffToDo))

			// then list
			Expect(listRepo.List()).To(Equal([]string{currentList}))

			// then delete
			Expect(listRepo.Remove(currentList)).ShouldNot(HaveOccurred())
		})
	})

	Context("Should be able to run a full lifecycle for passwords", func() {
		It("and act accordingly", func() {
			// First put
			err = passwordRepo.Put(currentList, stuffToDo)
			Expect(err).ShouldNot(HaveOccurred())

			// then get
			Expect(passwordRepo.Get(currentList)).To(Equal(stuffToDo))

			// then list
			Expect(passwordRepo.List()).To(BeNil())

			// then delete
			Expect(passwordRepo.Remove(currentList)).ShouldNot(HaveOccurred())
		})
	})
})
