@main => {
	mul = (a, b) -> a * b
	twice = (n) -> mul(2, n)

	std:println(twice(10))
}