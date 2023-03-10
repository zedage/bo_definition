package bo_definition

import (
	"strings"
	
	"github.com/iancoleman/strcase"
)

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
	Identification   Identification `yaml:"identification,omitempty" json:"identification,omitempty"`
	Alias            string         `yaml:"alias,omitempty" json:"alias,omitempty"`
	MessageVersion   int            `yaml:"messageVersion,omitempty" json:"messageVersion,omitempty"`
	UniqueKey        []string       `yaml:"uniqueKey,omitempty" json:"uniqueKey,omitempty"`
	InternalAliasUCC string         `yaml:"internalAliasUCC,omitempty" json:"internalAliasUCC,omitempty"`
	InternalAliasLCC string         `yaml:"internalAliasLCC,omitempty" json:"internalAliasLCC,omitempty"` 
}

type Identification struct {
	Value string `yaml:"value"`
}

type Definition struct {
	Type          string       `yaml:"type,omitempty" json:"type,omitempty"`
	NotNullFields []string     `yaml:"notNullFields,omitempty" json:"notNullFields,omitempty"`
	Properties    []Properties `yaml:"properties,omitempty" json:"properties,omitempty"`
}

func (d *Definition) fillMissingValuesCascade(path []string) {
	for i, p := range d.Properties {
		d.Properties[i].InternalFieldNameUCC = strcase.ToCamel(strings.ToLower(p.FieldName))
		d.Properties[i].InternalFieldNameLCC = strcase.ToLowerCamel(strings.ToLower(p.FieldName))
		d.Properties[i].InternalStructType = "string"
		d.Properties[i].InternalPbType = "string"
		if p.Type == "string" {
			if p.Format == "date" {
				d.Properties[i].InternalStructType = "*commons.FcsDate"
				d.Properties[i].InternalPbType = "google.protobuf.Timestamp"
			}
		} else if p.Type == "number" {
			if p.Format == "int8" || p.Format == "int16" || p.Format == "int32" || p.Format == "int64" {
				d.Properties[i].InternalStructType = p.Format
				d.Properties[i].InternalPbType = "uint32"
				if p.Format == "int64" {
					d.Properties[i].InternalPbType = "uint64"
				}
			} else if p.Length <= 2 {
				d.Properties[i].InternalStructType = "int8"
				d.Properties[i].InternalPbType = "uint32"
			} else if p.Length <= 4 {
				d.Properties[i].InternalStructType = "int16"
				d.Properties[i].InternalPbType = "uint32"
			} else if p.Length <= 9 {
				d.Properties[i].InternalStructType = "int32"
				d.Properties[i].InternalPbType = "uint32"
			} else if p.Length <= 18 {
				d.Properties[i].InternalStructType = "int64"
				d.Properties[i].InternalPbType = "uint64"
			} else if p.Format == "float" {
				d.Properties[i].InternalStructType = "float"
			} else if p.Format == "double" {
				d.Properties[i].InternalStructType = "double"
			} else {
				d.Properties[i].InternalStructType = "*commons.FcsTechLnr"
			}
		} else if p.Type == "object" {			
			d.Properties[i].InternalStructType = "*" + strings.Join(path, "_") + "_" + p.FieldName
			d.Properties[i].InternalPbType = strcase.ToLowerCamel(strings.ToLower(strings.Join(path, "_"))) + "_" + p.InternalFieldNameUCC
			p.Item.fillMissingValuesCascade(append(path, p.FieldName))
		} else if p.Type == "array" {
			d.Properties[i].InternalStructType = "* []" + strings.Join(path, "_") + "_" + p.FieldName
			d.Properties[i].InternalPbType = strcase.ToLowerCamel(strings.ToLower(strings.Join(path, "_"))) + "_" + p.InternalFieldNameUCC
			p.Item.fillMissingValuesCascade(append(path, p.FieldName))
		}
	}
}

func (d *Definition) FillMissingValues(boName string) {
	d.fillMissingValuesCascade([]string{boName})
}

type Properties struct {
	FieldName            string     `yaml:"fieldName,omitempty" json:"fieldName,omitempty"`
	Type                 string     `yaml:"type,omitempty" json:"type,omitempty"`
	Length               int        `yaml:"length,omitempty" json:"length,omitempty"`
	Format               string     `yaml:"format,omitempty" json:"format,omitempty"`
	Description          string     `yaml:"description,omitempty" json:"description,omitempty"`
	Cid                  bool       `yaml:"cid,omitempty" json:"cid,omitempty"`
	MinLength            int        `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength            int        `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Enum                 []string   `yaml:"enum,omitempty" json:"enum,omitempty"`
	Item                 Definition `yaml:"item,omitempty,omitempty" json:"item,omitempty"`
	InternalStructType   string     `yaml:"internalStructType,omitempty,omitempty" json:"internalStructType,omitempty"`
	InternalPbType       string     `yaml:"internalPbType,omitempty,omitempty" json:"internalPbType,omitempty"`
	InternalFieldNameUCC string     `yaml:"internalFieldNameUCC,omitempty,omitempty" json:"internalFieldNameUCC,omitempty"`
	InternalFieldNameLCC string     `yaml:"internalFieldNameLCC,omitempty,omitempty" json:"internalFieldNameLCC,omitempty"`
}

type Getter struct {
	MethodName string   `yaml:"methodName,omitempty" json:"methodName,omitempty"`
	Interface  []string `yaml:"interface,omitempty" json:"interface,omitempty"`
}

func (b *Message) FillMissingValues() {
	b.InternalAliasUCC = strcase.ToCamel(strings.ToLower(b.Alias))
	b.InternalAliasLCC = strcase.ToLowerCamel(strings.ToLower(b.Alias))
}
