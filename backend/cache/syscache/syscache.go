package syscache

//数据库需要一些meta信息来组织表结构，这些信息会频繁访问
//为了减少这一部分的开销，需要设置系统信息的一些高速缓存，包括配置信息
//需要同时提供脏数据更新的功能
type SysCache struct {
}
