package delete

import (
	"fmt"
	"testing"

	"github.com/benmorehouse/std/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDeleteCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delete Command")
}

var _ = Describe("Delete Command", func() {

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
			boltMock.EXPECT().Get(gomock.Any()).Return("welcome to std")
			boltMock.EXPECT().Remove(gomock.Any()).Return(nil)
			err := process(connectorMock, userMock, []string{"bucket"})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("if name of bucket gotten from user", func() {
		It("get for lifecycle", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			userMock.EXPECT().Input().Return("bucket")
			boltMock.EXPECT().Get("bucket").Return("welcome to std")
			boltMock.EXPECT().List().Return([]string{"sauce"})
			userMock.EXPECT().RunLifeCycle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			boltMock.EXPECT().Remove(gomock.Any()).Return(nil)
			err := process(connectorMock, userMock, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("if name of bucket passed in as argument", func() {
		It("use in lifecycle", func() {
			connectorMock.EXPECT().Connect().Return(boltMock, nil)
			userMock.EXPECT().Input().Return("bucket")
			gomock.InOrder(
				boltMock.EXPECT().Get("bucket").Return(""),
				boltMock.EXPECT().Get("bucket").Return("oh this has something"),
			)
			boltMock.EXPECT().List().Return([]string{"sauce"})
			userMock.EXPECT().RunLifeCycle(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			boltMock.EXPECT().Remove(gomock.Any()).Return(nil)
			err := process(connectorMock, userMock, []string{"bucket"})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
