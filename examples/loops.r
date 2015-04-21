@main => {	
	for = (s, e, f) -> {
		_for = (accum, next) -> {
			if accum < e {
				f(accum)
				next(accum + 1, next)
			}
		}
		_for(s, _for)
	}
	
	# for(0, 100, (i) -> std:println("HAI [", i, "]"))

	while = (c, f) -> {
		_while = (next) -> {
			if c() {
				f()
				next(next)
			}
		}
		_while(_while)
	}
	

	i = 10
	while(() -> {
		i = i - 1 
		i > 0
	}, () -> std:println("HAI [", i, "]"))
}