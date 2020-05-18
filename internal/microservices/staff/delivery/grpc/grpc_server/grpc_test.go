package staff

import (
	staff "2020_1_drop_table/internal/microservices/staff/delivery/grpc/client"
	staffMocks "2020_1_drop_table/internal/microservices/staff/mocks"
	"2020_1_drop_table/internal/microservices/staff/models"
	"context"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	type CheckStructInput struct {
		Ctx     context.Context
		StaffId int
	}

	type CheckStructOutput struct {
		Staff models.SafeStaff
		Err   error
	}
	type AdditionalInfo struct {
		Anything error
	}

	type testCaseStruct struct {
		InputData      CheckStructInput
		OutputData     CheckStructOutput
		AdditionalInfo AdditionalInfo
	}

	session := sessions.Session{Values: map[interface{}]interface{}{"userID": 228}}
	ctx := context.WithValue(context.Background(), "session", &session)
	emptContext := ctx

	test1St := models.SafeStaff{
		StaffID:  228,
		Name:     "",
		Email:    "",
		EditedAt: time.Time{},
		Photo:    "",
		IsOwner:  false,
		CafeId:   0,
		Position: "",
	}

	testCases := []testCaseStruct{
		//test Ok
		{
			InputData: CheckStructInput{
				Ctx:     emptContext,
				StaffId: 228,
			},
			OutputData: CheckStructOutput{
				Staff: test1St,
				Err:   nil,
			},
		},
	}

	staffMockUsecase := new(staffMocks.Usecase)
	urlForTests := "localhost:8091"
	go StartStaffGrpcServer(staffMockUsecase, urlForTests)
	grpcConn, err := grpc.Dial(urlForTests, grpc.WithInsecure())
	assert.Nil(t, err, "no error when start grpc conn")
	custGrpcClient := staff.NewStaffClient(grpcConn)

	for _, testCase := range testCases {
		staffMockUsecase.On("GetByID", mock.AnythingOfType("*context.valueCtx"), testCase.InputData.StaffId).Return(testCase.OutputData.Staff, testCase.OutputData.Err)
		res, err := custGrpcClient.GetById(testCase.InputData.Ctx, testCase.InputData.StaffId)
		assert.Equal(t, testCase.OutputData.Err, err)
		assert.Equal(t, testCase.OutputData.Staff, res)
	}

}
