package cache

import (
	"fmt"
)

// MakeMap converts []any to map[string]any.
func MakeMap(values ...any) map[string]any {
	if len(values)%2 != 0 {
		return nil
	}
	mp := make(map[string]any, len(values))
	for i := 0; i < len(values); i += 2 {
		field, ok := values[i].(string)
		if field == "" || !ok {
			fmt.Println(222)
			return nil
		}
		fmt.Println("field:", field, "	value:", values[i+1])
		mp[field] = values[i+1]
	}
	return mp
}

func AddFuncName(funcName string) {

}
