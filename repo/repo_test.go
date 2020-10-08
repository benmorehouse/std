package repo

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRepoCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repo Functionality")
}

var _ = BeforeSuite(func() {
	os.Setenv("STD_CONFIG_TESTING", "true")
})

var _ = Describe("Create Command", func() {
	var (
		connector   Connector
		err         error
		repo        Repo
		currentList string
		stuffToDo   string
	)

	BeforeEach(func() {
		connector = DefaultConnector()
		repo, err = connector.Connect()
		Expect(err).ShouldNot(HaveOccurred())
		currentList = "todo-list"
		stuffToDo = "clean the dishes"
	})

	AfterEach(func() {
		Expect(connector.Disconnect()).ShouldNot(HaveOccurred())
	})

	Context("Should be able to run a full lifecycle", func() {
		It("and act accordingly", func() {
			// First put
			err = repo.Put(currentList, stuffToDo)
			Expect(err).ShouldNot(HaveOccurred())

			// then get
			Expect(repo.Get(currentList)).To(Equal(stuffToDo))

			// then list
			Expect(repo.List()).To(Equal([]string{currentList}))

			// then delete
			Expect(repo.Remove(currentList)).ShouldNot(HaveOccurred())
		})
	})
})
