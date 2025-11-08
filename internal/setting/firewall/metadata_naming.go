package firewall

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"

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

func (m metadataNaming) Validate() error {
	if err := validate.Name(m.TableName, "table_name"); err != nil {
		return err
	}
	if err := validate.Name(m.ChainInputName, "chain_input_name"); err != nil {
		return err
	}
	if err := validate.Name(m.ChainOutputName, "chain_output_name"); err != nil {
		return err
	}
	if err := validate.Name(m.ChainForwardName, "chain_forward_name"); err != nil {
		return err
	}

	return nil
}
