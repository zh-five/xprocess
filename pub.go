package xprocess

var Debug = false

func println(a ...interface{}) {
	if Debug {
		println(a...)
	}
}
