package password

import (
	"fmt"
	"testing"

	"github.com/benmorehouse/std/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPassswordCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Password Command")
}

var _ = Describe("Password Command", func() {

	var (
		ctrl          *gomock.Controller
		vaultMock     *mocks.MockRepo
		userMock      *mocks.MockInteractor
		connectorMock *mocks.MockConnector
		commonErr     error
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		vaultMock = mocks.NewMockRepo(ctrl)
		connectorMock = mocks.NewMockConnector(ctrl)
		userMock = mocks.NewMockInteractor(ctrl)
		commonErr = fmt.Errorf("Some bad error")
	})

	Describe("Total vault password process", func() {
		Context("when putting set", func() {
			It("put", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("key"),
					userMock.EXPECT().Input().Return("value"),
				)
				connectorMock.EXPECT().Connect().Return(vaultMock, nil)

				vaultMock.EXPECT().Put("key", "value").Return(nil)
				err := process(connectorMock, userMock, "put")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when getting set", func() {
			It("get", func() {
				userMock.EXPECT().Input().Return("password")
				vaultMock.EXPECT().Get("password").Return("sugar and spice make everything nice")
				connectorMock.EXPECT().Connect().Return(vaultMock, nil)

				err := process(connectorMock, userMock, "get")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("When putting a password", func() {
		Context("if first try they make a valid key and value", func() {
			It("err", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("key"),
					userMock.EXPECT().Input().Return("value"),
				)

				vaultMock.EXPECT().Put("key", "value").Return(nil)

				Expect(putPassword(vaultMock, userMock)).ShouldNot(HaveOccurred())

			})
		})

		Context("if first try they make a valid key and value", func() {
			It("err", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return(""),
					userMock.EXPECT().Input().Return("key"),
					userMock.EXPECT().Input().Return("value"),
				)

				vaultMock.EXPECT().Put("key", "value").Return(nil)

				Expect(putPassword(vaultMock, userMock)).ShouldNot(HaveOccurred())

			})
		})

		Context("if they don't give a key at first", func() {
			It("err", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("key"),
					userMock.EXPECT().Input().Return(""),
					userMock.EXPECT().Input().Return("value"),
				)

				vaultMock.EXPECT().Put("key", "value").Return(nil)

				Expect(putPassword(vaultMock, userMock)).ShouldNot(HaveOccurred())
			})
		})

		Context("if the put fails", func() {
			It("err", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("key"),
					userMock.EXPECT().Input().Return("value"),
				)

				vaultMock.EXPECT().Put("key", "value").Return(commonErr)

				Expect(putPassword(vaultMock, userMock)).Should(HaveOccurred())
			})
		})
	})

	Describe("When getting a password", func() {
		Context("and getting key right on first time", func() {
			It("no err", func() {
				userMock.EXPECT().Input().Return("password")
				vaultMock.EXPECT().Get("password").Return("sugar and spice make everything nice")
				Expect(getPassword(vaultMock, userMock)).ShouldNot(HaveOccurred())
			})
		})

		Context("and they mispelled something", func() {
			It("no err", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("password"),
					userMock.EXPECT().Input().Return("no-this-password"),
				)

				gomock.InOrder(
					vaultMock.EXPECT().Get("password").Return(""),
					vaultMock.EXPECT().Get("no-this-password").Return("sugar and spice make everything nice"),
				)

				Expect(getPassword(vaultMock, userMock)).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("When removing a password", func() {
		Context("and getting a valid password first try", func() {
			It("remove that password", func() {
				userMock.EXPECT().Input().Return("chase-bank")
				vaultMock.EXPECT().Get("chase-bank").Return("some-dope-password")
				vaultMock.EXPECT().Remove("chase-bank").Return(nil)
				Expect(removePassword(vaultMock, userMock)).ShouldNot(HaveOccurred())
			})
		})

		Context("and giving valid password second try", func() {
			It("remove that password", func() {
				gomock.InOrder(
					userMock.EXPECT().Input().Return("chase-bank-wrong"),
					userMock.EXPECT().Input().Return("chase-bank"),
				)

				gomock.InOrder(
					vaultMock.EXPECT().Get("chase-bank-wrong").Return(""),
					vaultMock.EXPECT().Get("chase-bank").Return("some-dope-password"),
				)
				vaultMock.EXPECT().Remove("chase-bank").Return(nil)
				Expect(removePassword(vaultMock, userMock)).ShouldNot(HaveOccurred())
			})
		})

	})

})
