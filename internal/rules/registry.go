package rules

var registry []Rule

func Register(rule Rule) {
	registry = append(registry, rule)
}

func Rules() []Rule {
	return registry
}