package rbac_test

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/atomic"

	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/rbac"
	"github.com/rancher/opni/pkg/test"
	mock_rbac "github.com/rancher/opni/pkg/test/mock/rbac"
	"github.com/rancher/opni/pkg/util"
)

var testUsers = map[string][]string{
	"user0": {"tenant1", "tenant2"},
	"user1": {"tenant1"},
	"user2": {"tenant2"},
	"user3": {},
}

var _ = Describe("Middleware", Label(test.Unit), func() {
	It("should set tenant IDs for authorized users", func() {
		By("setting up the test controller")
		ctrl := gomock.NewController(GinkgoT())
		mockProvider := mock_rbac.NewMockProvider(ctrl)
		mockProvider.EXPECT().
			SubjectAccess(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, sar *corev1.SubjectAccessRequest) (*corev1.ReferenceList, error) {
				if clusters, ok := testUsers[sar.Subject]; ok {
					items := make([]*corev1.Reference, len(clusters))
					for i, cluster := range clusters {
						items[i] = &corev1.Reference{
							Id: cluster,
						}
					}
					return &corev1.ReferenceList{
						Items: items,
					}, nil
				}
				return nil, errors.New("user not found")
			}).
			AnyTimes()
		defer ctrl.Finish()
		app := fiber.New()
		logger.ConfigureAppLogger(app, "test")

		By("adding test middleware to insert the userID local")
		id := atomic.NewInt32(0)
		app.Use(func(c *fiber.Ctx) error {
			num := id.Load()
			c.Locals(rbac.UserIDKey, fmt.Sprintf("user%d", num))
			id.Store((num + 1) % 5) // includes nonexistent "user4"
			return c.Next()
		})

		By("adding the rbac middleware")
		app.Use(rbac.NewMiddleware(mockProvider, util.NewDelimiterCodec("foo", "|")))

		By("adding test middleware to check the resulting headers")
		app.Use(func(c *fiber.Ctx) error {
			defer GinkgoRecover()
			// This middleware should only be hit if the user is authorized
			userId := c.Locals(rbac.UserIDKey)
			Expect(userId).NotTo(BeNil())
			req := c.Request()
			orgId := string(req.Header.Peek("foo"))
			Expect(orgId).NotTo(BeEmpty())
			tenants := testUsers[userId.(string)]
			Expect(tenants).NotTo(BeEmpty())
			Expect(orgId).To(Equal(strings.Join(tenants, "|")))
			return c.Next()
		})

		By("adding a default 200 handler")
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		By("checking request status codes")
		for i := 0; i < 50; i++ {
			idNum := id.Load() // order is important - this gets modified in the handler
			userId := fmt.Sprintf("user%d", idNum)
			resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
			Expect(err).NotTo(HaveOccurred())
			if idNum > 2 {
				Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized), userId)
			} else {
				Expect(resp.StatusCode).To(Equal(fiber.StatusOK), userId)
			}
		}
	})
	It("should return 401 unauthorized if no user ID key is set", func() {
		By("setting up the test controller")
		ctrl := gomock.NewController(GinkgoT())
		mockProvider := mock_rbac.NewMockProvider(ctrl)
		mockProvider.EXPECT().
			SubjectAccess(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, sar *corev1.SubjectAccessRequest) ([]string, error) {
				defer GinkgoRecover()
				Fail("this should not be called")
				return nil, nil
			}).
			AnyTimes()
		defer ctrl.Finish()
		app := fiber.New()
		logger.ConfigureAppLogger(app, "test")

		By("adding the rbac middleware")
		app.Use(rbac.NewMiddleware(mockProvider, util.NewDelimiterCodec("foo", "|")))

		By("adding a default 200 handler")
		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		By("checking request status codes")
		for i := 0; i < 50; i++ {
			resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		}
	})
})
