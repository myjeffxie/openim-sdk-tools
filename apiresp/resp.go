// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiresp

import (
	"encoding/json"
	"reflect"

	"github.com/myjeffxie/openim-sdk-tools/errs"
	"github.com/myjeffxie/openim-sdk-tools/utils/jsonutil"
)

type ApiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data,omitempty"`
}

func (r *ApiResponse) MarshalJSON() ([]byte, error) {
	type apiResponse ApiResponse
	tmp := (*apiResponse)(r)
	if tmp.Data != nil {
		if format, ok := tmp.Data.(ApiFormat); ok {
			format.ApiFormat()
		}
		if isAllFieldsPrivate(tmp.Data) {
			tmp.Data = json.RawMessage(nil)
		} else {
			data, err := jsonutil.JsonMarshal(tmp.Data)
			if err != nil {
				return nil, err
			}
			tmp.Data = json.RawMessage(data)
		}
	}
	return jsonutil.JsonMarshal(tmp)
}

func isAllFieldsPrivate(v any) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf == nil {
		return false
	}
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return false
	}
	num := typeOf.NumField()
	for i := 0; i < num; i++ {
		c := typeOf.Field(i).Name[0]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func ApiSuccess(data any) *ApiResponse {
	return &ApiResponse{Data: data}
}

func ParseError(err error) *ApiResponse {
	if err == nil {
		return ApiSuccess(nil)
	}
	unwrap := errs.Unwrap(err)
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		resp := ApiResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: codeErr.Detail()}
		if resp.ErrDlt == "" {
			resp.ErrDlt = err.Error()
		}
		return &resp
	}
	return &ApiResponse{ErrCode: errs.ServerInternalError, ErrMsg: err.Error()}
}
