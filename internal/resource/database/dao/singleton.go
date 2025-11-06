package dao

var (
	userDAOSingleton     *UserDAO
	todoItemDAOSingleton *TodoItemDAO
)

func init() {
	userDAOSingleton = &UserDAO{}
	todoItemDAOSingleton = &TodoItemDAO{}
}

// GetUserDAO 获取用户DAO
//
//	return *baseDAO
//	author centonhuang
//	update 2024-10-17 04:59:37
func GetUserDAO() *UserDAO {
	return userDAOSingleton
}

// GetTodoItemDAO 获取待办事项DAO
//
//	return *TodoItemDAO
//	author centonhuang
//	update 2025-11-07 01:12:23
func GetTodoItemDAO() *TodoItemDAO {
	return todoItemDAOSingleton
}
