@main => {	
	for = (s, e, f) -> {
		_for = (accum) -> {
			if accum < e {
				f(accum)
				_for(accum + 1)
			}
		}
		_for(s)
	}
	
	for(0, 100, (i) -> std:println("HAI [", i, "]"))

	while = (c, f) -> {
		if c() {
			f()
			while(c, f)
		}
	}
	
	while(() -> 1 == 0, () -> std:println("HAI"))
}