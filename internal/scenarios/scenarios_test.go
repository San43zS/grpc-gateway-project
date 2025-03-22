package scenarios

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/internal/storage/mocksStorage"
	"testing"

	"context"
)

func TestServiceUser_CreateUser(t *testing.T) {
	// Create a mock instance of the QuoteRepository
	storage := mocksStorage.New()
	// Create the service with the mocked repository
	service := New(storage)

	// Define test cases in a table format
	testCases := []struct {
		name        string
		user        *models.User
		expectedUsr *models.User
		expectedErr error
	}{
		{
			name: "CreateUser - Success",
			user: &models.User{
				Name:  "John Doe",
				Email: "2Ia6Y@example.com",
				Age:   21,
			},
			expectedErr: nil,
		},
		{
			name: "CreateUser - Failure",
			user: &models.User{
				Name:  "Yagorka",
				Email: "",
				Age:   21,
			},
			expectedErr: emailEmptyErr,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage.UserMock.
				On("Create", mock.Anything, mock.Anything).
				Return(tc.expectedUsr, tc.expectedErr)

			_, err := service.CreateUser(context.Background(), tc.user)

			// Assert the results using testify's assert package
			assert.Equal(t, tc.expectedErr, err)

			// Assert that all expected repository interactions were called
			storage.UserMock.AssertExpectations(t)
		})
	}
}

//Add more test as needed... I'm so lazy...
