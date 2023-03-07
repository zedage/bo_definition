package bo_definition

import "fmt"

type BoRoot struct {
	BoModel BoModel `yaml:"boModel,omitempty" json:"boModel"`
}

type BoModel struct {
	YamlVersion string     `yaml:"yamlVersion,omitempty" json:"yamlVersion,omitempty"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Origin      string     `yaml:"origin,omitempty" json:"origin,omitempty"`
	Message     Message    `yaml:"message,omitempty" json:"message,omitempty"`
	Definition  Definition `yaml:"definition,omitempty" json:"definition,omitempty"`
	Getter      []Getter   `yaml:"getter,omitempty" json:"getter,omitempty"`
}

type Message struct {
	Identification Identification `yaml:"identification,omitempty" json:"identification,omitempty"`
	Alias          string         `yaml:"alias,omitempty" json:"alias,omitempty"`
	MessageVersion int            `yaml:"messageVersion,omitempty" json:"messageVersion,omitempty"`
	UniqueKey      []string       `yaml:"uniqueKey,omitempty" json:"uniqueKey,omitempty"`
}

type Identification struct {
	Value string `yaml:"value"`
}

type Definition struct {
	Type          string       `yaml:"type,omitempty" json:"type,omitempty"`
	NotNullFields []string     `yaml:"notNullFields,omitempty" json:"notNullFields,omitempty"`
	Properties    []Properties `yaml:"properties,omitempty" json:"properties,omitempty"`
}

type Properties struct {
	FieldName   string     `yaml:"fieldName,omitempty" json:"fieldName,omitempty"`
	Type        string     `yaml:"type,omitempty" json:"type,omitempty"`
	Length      int        `yaml:"length,omitempty" json:"length,omitempty"`
	Format      string     `yaml:"format,omitempty" json:"format,omitempty"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Cid         bool       `yaml:"cid,omitempty" json:"cid,omitempty"`
	MinLength   int        `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength   int        `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Enum        []string   `yaml:"enum,omitempty" json:"enum,omitempty"`
	Item        Definition `yaml:"item,omitempty,omitempty" json:"item,omitempty"`
}

type Getter struct {
	MethodName string   `yaml:"methodName,omitempty" json:"methodName,omitempty"`
	Interface  []string `yaml:"interface,omitempty" json:"interface,omitempty"`
}

func (p *Properties) GetGoType() string {
	if p.Type == "string" {
		if p.Format == "date" {
			return "*commons.FcsDate"
		}
	} else if p.Type == "number" {
		if p.Format == "int8" || p.Format == "int16" || p.Format == "int32" || p.Format == "int64" {
			return p.Format
		} else if p.Format == "float" {
			return "float"
		} else if p.Format == "double" {
			return "double"
		}
		return "*commons.FcsTechLnr"
	}
	return "string"
}

func (p *Properties) GetBsonAnnotation() string {
	return fmt.Sprintf("`bson:\"%s,omitempty\"`", p.FieldName)
}
