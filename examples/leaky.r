@main => {
	do_it = (f) -> {
		_do_it = (i) -> {
			std:println("i = ", i, ", f = ", f)
		}
		_do_it("real i")
	}

	# This leaks into the `do_it` func!
	f = "not f"
	# ...but this does not
	i = "not i"

	do_it("real f")
}