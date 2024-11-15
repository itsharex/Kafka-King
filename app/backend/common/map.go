package common

import (
	"reflect"
)

// StructToMap 将结构体转换为map[string]any
// 支持嵌套结构体、指针类型
func StructToMap(obj any) map[string]any {
	result := make(map[string]any)
	if obj == nil {
		return result
	}

	v := reflect.ValueOf(obj)
	t := v.Type()

	// 如果是指针,获取其底层元素
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	// 如果不是结构体类型则返回空map
	if t.Kind() != reflect.Struct {
		return result
	}

	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 处理指针类型
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}

		// 根据字段类型进行处理
		switch value.Kind() {
		case reflect.Struct:
			// 递归处理嵌套的结构体
			result[field.Name] = StructToMap(value.Interface())
		case reflect.Slice, reflect.Array:
			// 处理切片和数组
			length := value.Len()
			sliceResult := make([]any, length)
			for j := 0; j < length; j++ {
				item := value.Index(j)
				// 如果切片元素是结构体,递归处理
				if item.Kind() == reflect.Struct {
					sliceResult[j] = StructToMap(item.Interface())
				} else {
					sliceResult[j] = item.Interface()
				}
			}
			result[field.Name] = sliceResult
		default:
			// 其他类型直接放入result
			result[field.Name] = value.Interface()
		}
	}

	return result
}
