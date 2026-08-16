package partition

type onlyType struct {
	keywords    []string
	partitioned bool
}

func NewOnlyType(
	keywords []string,
	partitioned bool,
) PatternPartitionType {
	return &onlyType{
		keywords:    keywords,
		partitioned: partitioned,
	}
}

func (o *onlyType) IsPartitioned() bool {
	return o.partitioned
}

func (o *onlyType) Accepts(text string) (bool, string) {
	for _, keyword := range o.keywords {
		if keyword == text {
			return true, text
		}
	}
	return false, text
}
