package create

import (
	"fmt"
	"testing"

	"github.com/benmorehouse/std/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCreateCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Command")
}

var _ = Describe("Create Command", func() {

	var (
		ctrl          *gomock.Controller
		connectorMock *mocks.MockConnector
		boltMock      *mocks.MockRepo
		userMock      *mocks.MockInteractor
		commonErr     error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		connectorMock = mocks.NewMockConnector(ctrl)
		boltMock = mocks.NewMockRepo(ctrl)
		userMock = mocks.NewMockInteractor(ctrl)
		commonErr = fmt.Errorf("Some bad error")
	})

	Context("if not able to connect", func() {
		It("err", func() {
			connectorMock.EXPECT().Connect().Return(nil, commonErr)
			err := process(connectorMock, userMock, nil)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("if name of bucket passed in as argument", func() {
		It("use in lifecycle", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			userMock.EXPECT().RunLifeCycle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			boltMock.EXPECT().Get(gomock.Any()).Return("")
			err := process(connectorMock, userMock, []string{"bucket"})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("if name of bucket gotten from user", func() {
		It("get for lifecycle", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			userMock.EXPECT().Input().Return("bucket")
			boltMock.EXPECT().Get("bucket").Return("")
			boltMock.EXPECT().List().Return([]string{"sauce"})
			userMock.EXPECT().RunLifeCycle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			err := process(connectorMock, userMock, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("if name of bucket passed in as argument", func() {
		It("use in lifecycle", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			userMock.EXPECT().Input().Return("bucket")
			gomock.InOrder(
				boltMock.EXPECT().Get("bucket").Return("stuff"),
				boltMock.EXPECT().Get("bucket").Return(""),
			)
			boltMock.EXPECT().List().Return([]string{"sauce"})
			userMock.EXPECT().RunLifeCycle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			err := process(connectorMock, userMock, []string{"bucket"})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
