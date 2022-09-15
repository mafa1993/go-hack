package scanner

// 插件的约束，插件实现对一个域名和port 的扫描，返回result
type Checker interface {
	Check(host string, port uint64) *Result
}

type Result struct {
	Vulnerable bool
	Details    string
}
