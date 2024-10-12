package model

import (
	"errors"
)

const (
	OPENAI        = "openai"
	AZURE_OPENAI  = "azure-openai"
	BAICHUAN      = "baichuan"
	MOONSHOT      = "moonshot"
	DEEPSEEK      = "deepseek"
	VOLCANOENGINE = "volcanoengine"
)

type Model struct {
	Driver        string
	OpenAI        OpenAI
	AzureOpenAI   AzureOpenAI
	Baichuan      Baichuan
	Moonshot      Moonshot
	DeepSeek      DeepSeek
	VolcanoEngine VolcanoEngine
}

func NewModel(cfg Model) (Provider, error) {
	if cfg.Driver == OPENAI {
		if o := newOpenAI(cfg.OpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewOpenAI failed")
		}
	} else if cfg.Driver == AZURE_OPENAI {
		if o := newAzureOpenAI(cfg.AzureOpenAI); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewAzureOpenAI failed")
		}
	} else if cfg.Driver == BAICHUAN {
		if o := newBaichuan(cfg.Baichuan); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewBaichuan failed")
		}
	} else if cfg.Driver == MOONSHOT {
		if o := newMoonshot(cfg.Moonshot); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewMoonshot failed")
		}
	} else if cfg.Driver == DEEPSEEK {
		if o := newDeepSeek(cfg.DeepSeek); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewDeepSeek failed")
		}
	} else if cfg.Driver == VOLCANOENGINE {
		if o := newVolcanoEngine(cfg.VolcanoEngine); o != nil {
			return o, nil
		} else {
			return nil, errors.New("NewVolcanoEngine failed")
		}
	}

	return nil, errors.New("model driver not found")
}

func ModelAvailable(driver, modelType, modelName string) bool {
	switch driver {
	case OPENAI, AZURE_OPENAI:
		switch modelType {
		case "llm":
			for _, v := range OPENAI_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range OPENAI_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	case BAICHUAN:
		switch modelType {
		case "llm":
			for _, v := range BAICHUAN_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range BAICHUAN_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	case MOONSHOT:
		switch modelType {
		case "llm":
			for _, v := range MOONSHOT_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range MOONSHOT_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
		return false
	case DEEPSEEK:
		switch modelType {
		case "llm":
			for _, v := range DEEPSEEK_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range DEEPSEEK_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
	case VOLCANOENGINE:
		switch modelType {
		case "llm":
			for _, v := range VOLCANOENGINE_LLM_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		case "embedding":
			for _, v := range VOLCANOENGINE_EMBEDDING_SUPPORTS {
				if v == modelName {
					return true
				}
			}
		}
	}
	return false
}
