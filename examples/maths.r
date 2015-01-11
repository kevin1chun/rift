@main => {
	a = 10
	b = 20
	sum = (a, b) -> a + b
	c = sum(a, b)
	std:println(a, " + ", b, " = ", c)
}