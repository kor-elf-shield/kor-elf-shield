package firewall

type metadataNaming struct {
	TableName string `mapstructure:"table_name"`
}

func defaultMetadataNaming() metadataNaming {
	return metadataNaming{
		TableName: "shield",
	}
}
