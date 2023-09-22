package tgpt

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func ParseFCArgumentsToMap(data string) (map[string]string, error) {
	data = sanitizeJSON(data)

	// data = strings.ReplaceAll(data, "\n", "")
	// logx.Infof("parse fc arguments: %s", data)

	var params map[string]interface{}
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %s with %s", err.Error(), data)
	}

	var result = make(map[string]string)
	for k, v := range params {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}

func sanitizeJSON(input string) string {
	// 步骤1: 去除所有 JSON key 之前和之后的 \r, \n, \t
	re1 := regexp.MustCompile(`[\r\n\t]+\s*\"`)
	cleaned := re1.ReplaceAllString(input, " \"")

	// 步骤2: 去除所有 JSON value 之后的 \r, \n, \t
	re2 := regexp.MustCompile(`\"[\r\n\t]+\s*`)
	cleaned = re2.ReplaceAllString(cleaned, "\" ")

	// 步骤3: 去除 JSON 对象或数组结束之前的 \r, \n, \t
	re3 := regexp.MustCompile(`[\r\n\t]+\s*[\}\]]`)
	cleaned = re3.ReplaceAllString(cleaned, " }")

	// 步骤4: 将 value 中的 \r, \n, \t 转换为 \\r, \\n, \\t
	re4 := regexp.MustCompile(`[\r\n\t]`)
	escaped := re4.ReplaceAllStringFunc(cleaned, func(match string) string {
		switch match {
		case "\r":
			return "\\r"
		case "\n":
			return "\\n"
		case "\t":
			return "\\t"
		default:
			return match
		}
	})

	return escaped
}
