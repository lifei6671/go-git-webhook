package gojson

import (
	"encoding/json"
	"strconv"
)

type JsonObject struct {
	dataMap interface{}

}

func DeserializeObject(value string) *JsonObject {
	var dataMap map[string]interface{}
	jsonObject := &JsonObject{}

	if err := json.Unmarshal([]byte(value),&dataMap); err == nil {
		jsonObject.dataMap = dataMap
	}

	return jsonObject
}

func (p *JsonObject) IsValid() bool {
	return p.dataMap != nil
}

func (p *JsonObject) GetJsonObject(key string) *JsonObject {
	jsonObject := &JsonObject{}
	var dataMap map[string]interface{}

	if v,ok := p.dataMap.(map[string]interface{}) ; ok {
		if v, ok := v[key]; ok {
			if dataMap, ok = v.(map[string]interface{}); ok {
				jsonObject.dataMap = dataMap
			}
		}
	}
	return jsonObject
}

func (p *JsonObject) ToString() string {
	if m, ok := p.dataMap.(string); ok {
		return m
	}
	if m, ok := p.dataMap.(float64); ok {
		return strconv.FormatFloat(m, 'f', -1, 64)
	}
	if m, ok := p.dataMap.(bool); ok {
		return strconv.FormatBool(m)
	}
	return ""
}