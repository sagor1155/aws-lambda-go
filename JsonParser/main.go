package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReturnData struct {
	Status               string
	ContractId           int
	ProductPlanID        int
	OriginalSalesChannel string
	IsREProduct          bool
	IsNewREProduct       bool
}

type ApiResponse struct {
	Status     string
	StatusCode int
	Message    string
	Response   string
}

func main() {
	url := "http://localhost:8080/api/ppid"
	contractId := 335989648
	access_token := ""
	newPPIDList := []int{113372, 113382, 113392, 111872}

	returnData := ReturnData{}
	if returnData = GetProductPlanID(url, contractId, access_token); returnData.Status == "success" {
		if returnData.IsREProduct = IsREProduct(returnData.OriginalSalesChannel); returnData.IsREProduct {
			returnData.IsNewREProduct = IsNewREProduct(returnData.ProductPlanID, newPPIDList)
		}
	}
	fmt.Println("Status:", returnData.Status)
	fmt.Println("ContractId:", returnData.ContractId)
	fmt.Println("ProductPlanID:", returnData.ProductPlanID)
	fmt.Println("OriginalSalesChannel:", returnData.OriginalSalesChannel)
	fmt.Println("IsREProduct:", returnData.IsREProduct)
	fmt.Println("IsNewREProduct:", returnData.IsNewREProduct)
}

func IsNewREProduct(ppid int, newPPIDList []int) bool {
	for _, p := range newPPIDList {
		if p == ppid {
			return true
		}
	}
	return false
}

func IsREProduct(originalSalesChannel string) bool {
	if originalSalesChannel == "RE" {
		return true
	} else {
		return false
	}
}

func GetProductPlanID(url string, contractId int, access_token string) ReturnData {
	returnData := ReturnData{
		Status:               "error",
		ContractId:           contractId,
		IsREProduct:          false,
		IsNewREProduct:       false,
		OriginalSalesChannel: "",
		ProductPlanID:        0,
	}
	if apiResponse := MakeRequest(url); apiResponse.Status == "success" {
		ParseProductPlanID(apiResponse.Response, &returnData)
	}
	return returnData
}

func ParseProductPlanID(response string, returnData *ReturnData) {
	resBytes := []byte(response)
	var jsonRes map[string]interface{}
	_ = json.Unmarshal(resBytes, &jsonRes)

	if jsonRes["data"] != nil {
		arr := jsonRes["data"].([]interface{})

		if len(arr) == 1 {
			arrJson := arr[0].(map[string]interface{})
			if arrJson["contractId"] != nil {
				returnData.ContractId = int(arrJson["contractId"].(float64))
			}
			if arrJson["productPlanId"] != nil {
				returnData.ProductPlanID = int(arrJson["productPlanId"].(float64))
			}
			if arrJson["originalSalesChannel"] != nil {
				returnData.OriginalSalesChannel = arrJson["originalSalesChannel"].(string)
			}
			returnData.Status = "success"
		}
	}
}

func MakeRequest(URL string) ApiResponse {
	apiResponse := ApiResponse{Status: "error", Response: ""}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		apiResponse.Status = "error"
		apiResponse.Message = err.Error()
		fmt.Println("ERROR:", apiResponse.Message)
		return apiResponse
	}
	defer res.Body.Close()

	if apiResponse.StatusCode = res.StatusCode; apiResponse.StatusCode == 200 {
		resBody, _ := ioutil.ReadAll(res.Body)
		response := string(resBody)
		apiResponse.Response = response
		apiResponse.Status = "success"
		apiResponse.Message = "Successfully received response from API!!"
	}
	return apiResponse
}
