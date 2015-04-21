@main => {
	f = std:open("/tmp/test.txt")

	std:write(f, "Line 1\n")
	std:write(f, "Line 2\n")
	std:write(f, "Line 3\n")

	std:close(f)
}