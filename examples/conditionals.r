@main => {
	is_even = (n) -> {
		mod = n % 2
		mod == 0
	}

	n = 10
	if is_even(n) {
		std:println(n, " is even")
	} else {
		std:println(n, " is odd")
	}

	account_balance = 10000
	withdraw_amount = 1000
	post_balance = account_balance - withdraw_amount

	msg = if post_balance >= 0 {
		"Sufficient funds"
	} else {
		"Not enough funds"
	}

	std:println(msg)
}