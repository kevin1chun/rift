math => {
	sum = (a, b) -> a + b
}

main => {
	a = 10
	b = 20
	c = math:sum(a, b)
	std:println(a, " + ", b, " = ", c)
}