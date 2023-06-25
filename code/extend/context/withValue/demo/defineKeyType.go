// demo 包定义了一个 User 类型用于存储在上下文中
package demo

import "context"

// User 用于定义存储在Context中的值的类型
type User struct{}

// key 在包内定义的一个不可导出的类型.其目的在于为了防止在同一个上下文内,
// 包外的代码使用字面量相同的键覆盖包内的键.单独定义类型,可以确保即使包外
// 使用了相同的字面量,也不会导致键冲突,因为类型不同
type key int

// userKey 是 demo.User 类型的值所对应的键.该变量也是不可导出的
// 客户端使用 demo.NewContext 和 demo.FromContext 函数来存取上下文中的值
// 而非直接使用该键来存取上下文中的值
var userKey key = 1

// NewContext 返回一个新的包含了 user 的 Context
func NewContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// FromContext 从 ctx 中返回之前存储的 *User 类型的值,若之前未存储过,则返回的标量为false
func FromContext(ctx context.Context) (*User, bool) {
	value := ctx.Value(userKey)
	user, ok := value.(*User)
	return user, ok
}
