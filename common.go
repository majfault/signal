package signal

func deleteFirst[E comparable](s []E, e E) []E {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
