package server
//todo check what a problem
//import (
//	"2020_1_drop_table/internal/app/cafe/delivery/grpc/client"
//	cafeMock "2020_1_drop_table/internal/app/cafe/mocks"
//	"2020_1_drop_table/internal/app/cafe/models"
//	"context"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"google.golang.org/grpc"
//	"testing"
//	"time"
//)
//
//func TestGetById(t *testing.T) {
//	type CheckStructInput struct {
//		Ctx    context.Context
//		CafeID int
//	}
//	type CheckStructOutput struct {
//		Cafe models.Cafe
//		Err  error
//	}
//	type AdditionalInfo struct {
//		Anything error
//	}
//
//	type testCaseStruct struct {
//		InputData      CheckStructInput
//		OutputData     CheckStructOutput
//		AdditionalInfo AdditionalInfo
//	}
//	emptContext := context.Background()
//
//	test1Cafe := models.Cafe{
//		CafeID:      228,
//		CafeName:    "",
//		Address:     "",
//		Description: "",
//		StaffID:     0,
//		OpenTime:    time.Time{},
//		CloseTime:   time.Time{},
//		Photo:       "",
//	}
//
//	testCases := []testCaseStruct{
//		//test Ok
//		{
//			InputData: CheckStructInput{
//				Ctx:    emptContext,
//				CafeID: test1Cafe.CafeID,
//			},
//			OutputData: CheckStructOutput{
//				Cafe: test1Cafe,
//				Err:  nil,
//			},
//		},
//	}
//
//	cafeUsecase := new(cafeMock.Usecase)
//	urlForTests := "localhost:8093"
//	go StartCafeGrpcServer(cafeUsecase, urlForTests)
//	grpcConn, err := grpc.Dial(urlForTests, grpc.WithInsecure())
//	assert.Nil(t, err, "no error when start grpc conn")
//	cafeGrpcClient := client.NewCafeClient(grpcConn)
//	for _, testCase := range testCases {
//		cafeUsecase.On("GetByID", mock.AnythingOfType("*context.valueCtx"), testCase.InputData.CafeID).Return(testCase.OutputData.Cafe, nil)
//		resCafe, err := cafeGrpcClient.GetByID(testCase.InputData.Ctx, testCase.InputData.CafeID)
//		assert.Equal(t, testCase.OutputData.Err, err)
//		assert.Equal(t, testCase.OutputData.Cafe, resCafe)
//	}
//
//}
