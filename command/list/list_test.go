package list

import (
	"fmt"
	"testing"

	"github.com/benmorehouse/std/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestListCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "List Command")
}

var _ = Describe("Create Command", func() {

	var (
		ctrl          *gomock.Controller
		connectorMock *mocks.MockConnector
		boltMock      *mocks.MockRepo
		commonErr     error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		connectorMock = mocks.NewMockConnector(ctrl)
		boltMock = mocks.NewMockRepo(ctrl)
		commonErr = fmt.Errorf("Some bad error")
	})

	Context("if not able to connect", func() {
		It("err", func() {
			connectorMock.EXPECT().Connect().Return(nil, commonErr)
			err := process(connectorMock, nil)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("show lists", func() {
		It("no error", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			boltMock.EXPECT().List().Return([]string{"bucket"})
			connectorMock.EXPECT().Disconnect().Return(nil)
			err := process(connectorMock, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
