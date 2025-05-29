package ctxdata

import "context"

// GetUId 从上下文中获取用户ID。
// 该函数主要用于从给定的上下文对象中提取用户ID，如果存在并且类型正确则返回该ID，否则返回空字符串。
// 参数:
//
//	ctx context.Context: 传递的上下文对象，用于在不同层级的调用之间传递请求特定的数据。
//
// 返回值:
//
//	string: 如果上下文中包含类型为字符串的用户ID，则返回该ID；否则返回空字符串。
func GetUId(ctx context.Context) string {
	// 尝试从上下文中获取用户ID，并检查其是否为字符串类型。
	if u, ok := ctx.Value(Identity).(string); ok {
		// 如果是字符串类型且存在，则返回用户ID。
		return u
	}

	// 如果用户ID不存在或类型不匹配，则返回空字符串。
	return ""
}
