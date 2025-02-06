package helloworld

type HWRouter struct {
	test string
}

func (r HWRouter) Path() string {
	return "/helloworld"
}
