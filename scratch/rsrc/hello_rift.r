hello_rift {
	sum = (a, b) -> a + b
	
	a = 1 + 20 * 3
	b = 2
	c = sum(a, b) # 3

	d = [1, 2]
	e = [3, 4]
	f = sum(d, e) # [1, 2, 3, 4]

	githubUser = github.user(12345)
}
