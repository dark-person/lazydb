package lazydb

import (
	"strconv"
	"strings"
)

// Fill A SQL statement contains '?' with arguments.
// This function behavior similar as prepare statement, however it will not consider the database config.
//
// This function should be consider as a reference when logging/debugging,
// it SHOULD NOT be used in construct real queries.
//
// If args contains unsupported type, panic will occur. This function only supports:
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
func asFilledQuery(sql string, args ...any) string {
	resultStr := sql

	for _, item := range args {
		switch itemVal := item.(type) {
		case string:
			resultStr = strings.Replace(resultStr, "?", `"`+itemVal+`"`, 1)

		case int:
			resultStr = strings.Replace(resultStr, "?", strconv.Itoa(itemVal), 1)
		case int8:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatInt(int64(itemVal), 10), 1)
		case int16:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatInt(int64(itemVal), 10), 1)
		case int32:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatInt(int64(itemVal), 10), 1)
		case int64:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatInt(itemVal, 10), 1)

		case uint:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatUint(uint64(itemVal), 10), 1)
		case uint8:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatUint(uint64(itemVal), 10), 1)
		case uint16:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatUint(uint64(itemVal), 10), 1)
		case uint32:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatUint(uint64(itemVal), 10), 1)
		case uint64:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatUint(itemVal, 10), 1)

		case float64:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatFloat(itemVal, 'G', -1, 64), 1)
		case float32:
			resultStr = strings.Replace(resultStr, "?", strconv.FormatFloat(float64(itemVal), 'G', -1, 32), 1)

		default:
			panic("unsupported type")
		}
	}
	return resultStr
}
