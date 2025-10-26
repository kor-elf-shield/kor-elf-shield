package firewall

type metadataNaming struct {
	TableName        string `mapstructure:"table_name"`
	ChainInputName   string `mapstructure:"chain_input_name"`
	ChainOutputName  string `mapstructure:"chain_output_name"`
	ChainForwardName string `mapstructure:"chain_forward_name"`
}

func defaultMetadataNaming() metadataNaming {
	return metadataNaming{
		TableName:        "shield",
		ChainInputName:   "input",
		ChainOutputName:  "output",
		ChainForwardName: "forward",
	}
}
