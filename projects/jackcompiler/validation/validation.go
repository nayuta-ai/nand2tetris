package validation

type Rule func(s string) bool

func Match(val string) Rule {
	return func(s string) bool {
		return val == s
	}
}

func OneOf(options ...string) Rule {
	return func(s string) bool {
		for _, o := range options {
			if o == s {
				return true
			}
		}
		return false
	}
}
