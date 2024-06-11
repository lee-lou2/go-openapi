package client

// ScopeRateLimits 스코프별 초당 요청 제한
var ScopeRateLimits = map[string]int{
	ScopeReadMe: 100,
}
