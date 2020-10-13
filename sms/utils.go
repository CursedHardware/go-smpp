package sms

func blocks(n, block byte) byte {
	if n%block == 0 {
		return n / block
	}
	return n/block + 1
}

func unblocks(n, block byte) (x byte) {
	x = n * block
	if n%block != 0 {
		x -= 1
	}
	return x
}
