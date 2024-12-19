package server

import (
	"consumer-payment-service/environment"
	"consumer-payment-service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHttpMount(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	cfg := &environment.Config{
		THIRD_PARTY_SERVICE_BASE_URL: "http://example.domain.com/third-party",
	}

	mockDataStore := mocks.NewMockMongoDBStore(controller)
	mockThirdPartyClient := mocks.NewMockThirdPartyAPIClient(controller)

	router := MountServer(cfg, mockDataStore, mockThirdPartyClient)
	assert.NotNil(t, router)

}
