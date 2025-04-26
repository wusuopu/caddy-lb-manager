package helper

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

type JSONParser struct {
	Value *fastjson.Value
}
/*
	动态解析 body 为 json
*/
func (parser *JSONParser) GetJSONBody(c *gin.Context) *fastjson.Value {
	var a fastjson.Arena
	var p fastjson.Parser

	body, ok := c.Get("body")
	if ok {
		parser.Value = body.(*fastjson.Value)
		return parser.Value
	}
	empty := a.NewObject()

	rawBody, _ := c.Get("rawBody")
	if rawBody == nil {
		return empty
	}
	maxLength := 1024 * 1024 * 4		// 最大支持 4M 长度字符，以免内存不足 
	data, _ := rawBody.([]byte)
	ret, err := p.ParseBytes(data[:min(len(data), maxLength)])
	if err != nil {
		ret = empty
	}
	// 将解析结果缓存下来
	c.Set("body", ret)

	parser.Value = ret
	return parser.Value
}
/*
	动态解析 url query 为 json
*/
func (parser *JSONParser) GetJSONQuery(c *gin.Context) *fastjson.Value {
	query, ok := c.Get("query")
	if ok {
		parser.Value = query.(*fastjson.Value)
		return parser.Value
	}

	obj := parser.ParseJSONQuery(c.Request.URL.RawQuery)
	c.Set("query", obj)
	parser.Value = obj
	return parser.Value
}
func (parser *JSONParser) ParseJSONQuery(qs string) *fastjson.Value {
	var a fastjson.Arena

	obj := a.NewObject()
	query, err := url.ParseQuery(qs)
	if err != nil {
		return obj
	}
	for k, v := range query {
		if len(v) == 1 && len(v[0]) == 0 {
			//该字段为空
			continue
		}

		isArray := false
		isNested := false
		if strings.HasSuffix(k, "[]") {
			isArray = true
			k = k[:(len(k)-2)]
		}
		if strings.HasSuffix(k, "]") {
			// 嵌套结构
			isNested = true
		}
		var val *fastjson.Value
		if isArray {
			val = a.NewArray()
			for idx, el := range v {
				val.SetArrayItem(idx, a.NewString(el))
			}
		} else {
			val = a.NewString(v[len(v) - 1])
		}

		if !isNested {
			// 单层结构，直接赋值
			obj.Set(k, val)
			continue
		}

		parent := obj
		keys := strings.Split(k, "[")
		for idx, key := range keys {
			if strings.HasSuffix(key, "]") {
				key = key[:len(key)-1]
			}
			if key == "" {
				continue
			}

			if idx == (len(keys) - 1) {
				parent.Set(key, val)
			} else {
				if !parent.Exists(key) {
					subItem := a.NewObject()
					parent.Set(key, subItem)

				}
				parent = parent.Get(key)
			}
		}
	}

	return obj
}


// lodash 风格的 GET 操作
func (parser *JSONParser) GetJSONInt64(keys string) (int64, error) {
	val := parser.Value
	
	var item *fastjson.Value
	if keys == "" {
		item = val
	} else {
		k := strings.Split(keys, ".")
		if !val.Exists(k...) {
			return 0, fmt.Errorf("key %s not exists", keys)
		}
		item = val.Get(k...)
	}

	if item.Type() == fastjson.TypeNumber {
		return item.Int64()
	}
	if item.Type() == fastjson.TypeString {
		data, err := item.StringBytes()
		if err != nil {
			return 0, err
		}
		return strconv.ParseInt(string(data), 10, 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
func (parser *JSONParser) GetJSONFloat64(keys string) (float64, error) {
	val := parser.Value
	k := strings.Split(keys, ".")
	if !val.Exists(k...) {
		return 0, fmt.Errorf("key %s not exists", keys)
	}
	item := val.Get(k...)
	if item.Type() == fastjson.TypeNumber {
		return item.Float64()
	}
	if item.Type() == fastjson.TypeString {
		data, err := item.StringBytes()
		if err != nil {
			return 0, err
		}
		return strconv.ParseFloat(string(data), 64)
	}
	return 0, fmt.Errorf("key %s in not number", keys)
}
// lodash 风格的 get 方法
func (parser *JSONParser) GetJSONItem(keys string) *fastjson.Value {
	val := parser.Value
	if keys == "" {
		return val
	}
	return val.Get(strings.Split(keys, ".")...)
}
func (parser *JSONParser) GetJSONString(keys string) (string, error) {
	v := parser.GetJSONItem(keys)
	if v == nil {
		return "", fmt.Errorf("key %s not exists", keys)
	}
	if v.Type() != fastjson.TypeString {
		return "", fmt.Errorf("key %s in not string", keys)
	}
	data, err := v.StringBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (parser *JSONParser) GetJSONBool(keys string) (bool, error) {
	v := parser.GetJSONItem(keys)
	if v != nil && v.Type() == fastjson.TypeTrue {
		return true, nil
	}
	if v != nil && v.Type() == fastjson.TypeFalse {
		return false, nil
	}
	return false, fmt.Errorf("key %s in not bool", keys)
}