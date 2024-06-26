package xprocess

// 检查进程是不是唯一的，不是则kill旧进程
// 为避免误杀和优雅退出, 此处kill不是强杀. 需在 onKill() 方法中自行处理退出操作
func UniqueCheckAndKillOld(flag string, onKill func()) {
	uniqueCheckAndKillOld(flag, onKill)
}

// 检查进程是不是唯一的
func UniqueCheck(flag string) bool {
	return uniqueCheck(flag)
}
