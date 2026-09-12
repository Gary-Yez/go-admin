package sys_devtools

import "strings"

func (field Filed) IsNumber() bool {
	return strings.HasPrefix(field.Type, "int") || strings.HasPrefix(field.Type, "uint") || strings.HasPrefix(field.Type, "float")
}

func (field Filed) IsInteger() bool {
	return strings.HasPrefix(field.Type, "int") || strings.HasPrefix(field.Type, "uint")
}

// 数字控件限制到对应类型范围，宽整数额外限制在 JavaScript 安全整数范围内。
func (field Filed) NumberMin() string {
	if strings.HasPrefix(field.Type, "uint") {
		return "0"
	}
	switch field.Type {
	case "int8":
		return "-128"
	case "int16":
		return "-32768"
	case "int32":
		return "-2147483648"
	case "float32":
		return "-3.4028234663852886e38"
	case "float64":
		return "-Number.MAX_VALUE"
	default:
		return "Number.MIN_SAFE_INTEGER"
	}
}

func (field Filed) NumberMax() string {
	switch field.Type {
	case "int8":
		return "127"
	case "int16":
		return "32767"
	case "int32":
		return "2147483647"
	case "uint8":
		return "255"
	case "uint16":
		return "65535"
	case "uint32":
		return "4294967295"
	case "float32":
		return "3.4028234663852886e38"
	case "float64":
		return "Number.MAX_VALUE"
	default:
		return "Number.MAX_SAFE_INTEGER"
	}
}
