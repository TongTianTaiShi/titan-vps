package aliyun

import (
	"fmt"

	"github.com/LMF709268224/titan-vps/api/types"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// newBssopen /**
func newBssopen(keyID, keySecret string) (_result *bssopenapi20171214.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(keyID),
		AccessKeySecret: tea.String(keySecret),
	}

	config.Endpoint = tea.String("business.aliyuncs.com")
	_result = &bssopenapi20171214.Client{}
	_result, _err = bssopenapi20171214.NewClient(config)
	return _result, _err
}

// DescribeInstanceBill crate an instance
func DescribeInstanceBill(keyID, keySecret string) (*types.CreateInstanceResponse, error) {
	var out *types.CreateInstanceResponse

	client, err := newBssopen(keyID, keySecret)
	if err != nil {
		return out, err
	}

	describeInstanceBillRequest := &bssopenapi20171214.DescribeInstanceBillRequest{
		BillingCycle: tea.String("2023-08"),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.DescribeInstanceBillWithOptions(describeInstanceBillRequest, runtime)
		if err != nil {
			return err
		}

		fmt.Println("result :", result)

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}

// QueryProductList crate an instance
func QueryProductList(keyID, keySecret string) (*types.CreateInstanceResponse, error) {
	var out *types.CreateInstanceResponse

	client, err := newBssopen(keyID, keySecret)
	if err != nil {
		return out, err
	}

	describeInstanceBillRequest := &bssopenapi20171214.QueryProductListRequest{
		PageNum: tea.Int32(1),
	}

	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		result, err := client.QueryProductListWithOptions(describeInstanceBillRequest, runtime)
		if err != nil {
			return err
		}

		fmt.Println("result :", result)

		return nil
	}()

	if tryErr != nil {
		errors := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errors = _t
		} else {
			errors.Message = tea.String(tryErr.Error())
		}
		return out, errors
	}
	return out, nil
}
