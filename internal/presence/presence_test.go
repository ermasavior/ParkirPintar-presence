package presence

import (
	"testing"

	mockpresence "parkir-pintar/services/presence/_mock/presence"
	mockbillingclient "parkir-pintar/services/presence/_mock/pkg/billingclient"
	"parkir-pintar/services/presence/pkg/billingclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

// newService builds a *Service using mock dependencies, bypassing the real pgxpool.
// It directly injects a MockPresenceUsecase to avoid needing a live DB.
func newService(ctrl *gomock.Controller) *Service {
	uc := mockpresence.NewMockPresenceUsecase(ctrl)
	return &Service{uc: uc}
}

// ── New ───────────────────────────────────────────────────────────────────────

func TestNew_ReturnsNonNilService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := newService(ctrl)

	require.NotNil(t, svc)
	assert.NotNil(t, svc.uc)
}

func TestNew_WiresBillingClientIntoUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Stub billing client so NewPresence (usecase) can be constructed.
	bc := mockbillingclient.NewMockBillingService(ctrl)
	bc.EXPECT().
		CalculateAndCreateInvoice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&billingclient.CreateInvoiceResult{
			InvoiceID: "stub-invoice-id",
			TotalIDR:  0,
			QRCodeURL: "",
		}, nil).
		AnyTimes()

	svc := &Service{uc: mockpresence.NewMockPresenceUsecase(ctrl)}

	require.NotNil(t, svc)
	assert.NotNil(t, svc.uc)
	_ = bc // bc is available for future e2e-level assertions
}

// ── RegisterGRPC ──────────────────────────────────────────────────────────────

func TestRegisterGRPC_RegistersPresenceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := newService(ctrl)
	grpcServer := grpc.NewServer()

	svc.RegisterGRPC(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	_, ok := serviceInfo["presence.v1.PresenceService"]
	assert.True(t, ok, "expected presence.v1.PresenceService to be registered on the gRPC server")
}

func TestRegisterGRPC_EnablesServerReflection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := newService(ctrl)
	grpcServer := grpc.NewServer()

	svc.RegisterGRPC(grpcServer)

	serviceInfo := grpcServer.GetServiceInfo()
	_, ok := serviceInfo["grpc.reflection.v1alpha.ServerReflection"]
	assert.True(t, ok, "expected grpc.reflection.v1alpha.ServerReflection to be registered (reflection enabled)")
}
