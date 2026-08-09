package partition

type exceptType struct {
	keywords    []string
	partitioned bool
}

func NewExceptType(
	keywords []string,
	partitioned bool,
) PatternPartitionType {
	return &exceptType{
		keywords:    keywords,
		partitioned: partitioned,
	}
}

func (o *exceptType) IsPartitioned() bool {
	return o.partitioned
}

func (o *exceptType) Accepts(text string) (bool, string) {
	for _, keyword := range o.keywords {
		if keyword == text {
			return false, text
		}
	}
	return true, text
}
