package deref

func String(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func Bool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
