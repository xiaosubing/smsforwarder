package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code   int      `json:"code"`
	Result string   `json:"result"`
	Data   DataInfo `json:"data"`
}

type DataInfo struct {
	Phone  string `json:"phone"`
	Sender string `json:"sender"`
	Sms    string `json:"sms"`
	Code   string `json:"code"`
}

// GetMessage  801
func GetMessageInfo(writer http.ResponseWriter, request *http.Request) {
	var getMessageVerify string = ""
	query := request.URL.Query()
	phone := query.Get("phone")
	verifyCode := query.Get("verify")

	if len(phone) == 0 {
		ret := Response{
			Code:   404,
			Result: "invalid phone number",
		}
		if err := json.NewEncoder(writer).Encode(ret); err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}

	}

	if len(getMessageVerify) != 0 {
		if verifyCode != getMessageVerify {
			ret := Response{
				Code:   404,
				Result: "invalid verify code",
			}
			if err := json.NewEncoder(writer).Encode(ret); err != nil {
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}

	// pass
	infos := QueryData(phone)
	var jsonInfo DataInfo
	err := json.Unmarshal([]byte(infos), &jsonInfo)
	if err != nil {
		ret := Response{
			Code:   404,
			Result: "invalid data",
		}
		if err := json.NewEncoder(writer).Encode(ret); err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	var ret = Response{
		Code:   200,
		Result: "sucess",
		Data: DataInfo{
			Phone:  jsonInfo.Phone,
			Sender: jsonInfo.Sender,
			Sms:    jsonInfo.Sms,
			Code:   jsonInfo.Code,
		},
	}
	if err := json.NewEncoder(writer).Encode(ret); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)

	}
}
