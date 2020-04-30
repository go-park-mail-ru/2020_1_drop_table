package test_test

import (
	"2020_1_drop_table/configs"
	customer "2020_1_drop_table/internal/app/customer/delivery/grpc/client"
	"2020_1_drop_table/internal/app/customer/delivery/grpc/server"
	models2 "2020_1_drop_table/internal/app/customer/models"

	customerMocks "2020_1_drop_table/internal/app/customer/mocks"

	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
)

func TestAdd(t *testing.T) {
	type CheckStructInput struct {
		Ctx         context.Context
		NewCustomer models2.Customer
	}
	type CheckStructOutput struct {
		AddedCustomer models2.Customer
		Err           error
	}
	type AdditionalInfo struct {
		Anything error
	}

	type testCaseStruct struct {
		InputData      CheckStructInput
		OutputData     CheckStructOutput
		AdditionalInfo AdditionalInfo
	}
	emptContext := context.Background()

	test1Cust := models2.Customer{
		CustomerID:   "0",
		CafeID:       0,
		Type:         "",
		Points:       "",
		SurveyResult: "",
	}

	testCases := []testCaseStruct{
		//test Ok
		{
			InputData: CheckStructInput{
				Ctx:         emptContext,
				NewCustomer: test1Cust,
			},
			OutputData: CheckStructOutput{
				AddedCustomer: test1Cust,
				Err:           nil,
			},
		},
	}

	customerMockUsecase := new(customerMocks.Usecase)
	go server.StartCustomerGrpcServer(customerMockUsecase)
	grpcConn, err := grpc.Dial(configs.GRPCCustomerUrl, grpc.WithInsecure())
	assert.Nil(t, err, "no error when start grpc conn")
	custGrpcClient := customer.NewCustomerClient(grpcConn)

	for _, testCase := range testCases {
		customerMockUsecase.On("Add", mock.AnythingOfType("*context.valueCtx"), testCase.InputData.NewCustomer).Return(testCase.OutputData.AddedCustomer, testCase.OutputData.Err)
		res, err := custGrpcClient.Add(testCase.InputData.Ctx, testCase.InputData.NewCustomer)
		assert.Equal(t, testCase.OutputData.Err, err)
		assert.Equal(t, testCase.OutputData.AddedCustomer, res)
	}

}
