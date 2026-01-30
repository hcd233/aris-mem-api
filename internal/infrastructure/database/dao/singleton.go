package dao

var (
	userDAOSingleton       *UserDAO
	todoItemDAOSingleton   *TodoItemDAO
	dialogDAOSingleton     *DialogDAO
	tagDAOSingleton        *TagDAO
	articleDAOSingleton    *ArticleDAO
	articleTagDAOSingleton *ArticleTagDAO
	actionDAOSingleton     *ActionDAO
)

func init() {
	userDAOSingleton = &UserDAO{}
	todoItemDAOSingleton = &TodoItemDAO{}
	dialogDAOSingleton = &DialogDAO{}
	tagDAOSingleton = &TagDAO{}
	articleDAOSingleton = &ArticleDAO{}
	articleTagDAOSingleton = &ArticleTagDAO{}
	actionDAOSingleton = &ActionDAO{}
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

// GetDialogDAO 获取对话DAO
//
//	return *DialogDAO
//	author centonhuang
//	update 2025-11-14 16:07:18
func GetDialogDAO() *DialogDAO {
	return dialogDAOSingleton
}

// GetTagDAO 获取标签DAO
//
//	return *TagDAO
//	author centonhuang
//	update 2026-01-29 10:00:00
func GetTagDAO() *TagDAO {
	return tagDAOSingleton
}

// GetArticleDAO 获取文章DAO
//
//	return *ArticleDAO
//	author centonhuang
//	update 2026-01-29 10:00:00
func GetArticleDAO() *ArticleDAO {
	return articleDAOSingleton
}

// GetArticleTagDAO 获取文章标签关联DAO
//
//	return *ArticleTagDAO
//	author centonhuang
//	update 2026-01-29 10:00:00
func GetArticleTagDAO() *ArticleTagDAO {
	return articleTagDAOSingleton
}

// GetActionDAO get user action DAO
//
//	return *ActionDAO
//	author centonhuang
//	update 2026-01-30 21:00:00
func GetActionDAO() *ActionDAO {
	return actionDAOSingleton
}
