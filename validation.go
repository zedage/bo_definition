package bo_definition

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type SimpleError struct {
	Text string
}

func (m SimpleError) Error() string {
	return m.Text
}

type ValidationError struct {
	Pfx       string
	FieldName string
	FileName  string
	List      []string
}

func (m *ValidationError) Error() string {
	return fmt.Sprintf("%s.%s is not set in %s", m.Pfx, m.FieldName, m.FileName)
}

func checkIfTrue(flag bool, pfx string, fieldName string, fileName string) error {
	if flag {
		return &ValidationError{Pfx: pfx, FieldName: fieldName, FileName: fileName}
	}
	return nil
}

func checkIfStringNotInList(value string, list []string, pfx string, fieldName string, fileName string) error {
	for _, v := range list {
		if v == value {
			return nil
		}
	}
	return &ValidationError{Pfx: pfx, FieldName: fieldName, FileName: fileName, List: list}
}

func (b *BoRoot) Validate(fileName string) error {
	err := b.BoModel.validate("boModel", fileName)
	if err != nil {
		return err
	}
	def := &b.BoModel.Definition
	err = ensureFieldsExist(b.BoModel.Message.UniqueKey, def, fileName, "boModel.message.uniqueKey", true)
	if err != nil {
		return err
	}
	return nil
}

func (b *BoModel) validate(pfx string, fileName string) error {
	err := checkIfTrue(b.YamlVersion == "", pfx, "yamlVersion", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(b.Origin == "", pfx, "origin", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(b.Description == "", pfx, "description", fileName)
	if err != nil {
		return err
	}
	err = b.Message.validate(pfx+".message", fileName)
	if err != nil {
		return err
	}
	err = b.Definition.validate(pfx+".definition", fileName)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) validate(pfx string, fileName string) error {
	err := checkIfTrue(m.Alias == "", pfx, "alias", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(m.MessageVersion == 0, pfx, "messageVersion", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(len(m.UniqueKey) == 0, pfx, "uniqueKey", fileName)
	if err != nil {
		return err
	}
	err = m.Identification.validate(pfx+".identification", fileName)
	if err != nil {
		return err
	}
	m.InternalAliasUCC = strcase.ToCamel(strings.ToLower(m.Alias))
	m.InternalAliasLCC = strcase.ToLowerCamel(strings.ToLower(m.Alias))

	return nil
}

func (i *Identification) validate(pfx string, fileName string) error {
	err := checkIfTrue(i.Value == "", pfx, "value", fileName)
	if err != nil {
		return err
	}
	return nil
}

func (d *Definition) validate(pfx string, fileName string) error {
	err := checkIfTrue(d.Type == "", pfx, "type", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(len(d.Properties) == 0, pfx, "properties", fileName)
	if err != nil {
		return err
	}
	err = checkIfStringNotInList(d.Type, []string{"object", "array"}, pfx, "type", fileName)
	if err != nil {
		return err
	}
	for i, prop := range d.Properties {
		err = prop.validate(fmt.Sprintf("%s.properties[%v]", pfx, i), fileName)
		if err != nil {
			return err
		}
	}
	err = ensureFieldsExist(d.NotNullFields, d, fileName, pfx+".notNullFields", false)
	if err != nil {
		return err
	}
	return nil
}

func (p *Properties) validate(pfx string, fileName string) error {
	err := checkIfTrue(p.FieldName == "", pfx, "fieldName", fileName)
	if err != nil {
		return err
	}
	err = checkIfTrue(p.Description == "", pfx, "description", fileName)
	if err != nil {
		return err
	}
	err = checkIfStringNotInList(p.Type, []string{"number", "string", "object", "array", "float"}, pfx, "type", fileName)
	if err != nil {
		return err
	}
	switch p.Type {
	case "number":
		err = p.validateNumber(pfx, fileName)
		if err != nil {
			return err
		}
	case "string":
		err = p.validateString(pfx, fileName)
		if err != nil {
			return err
		}
	case "array":
		err = p.Item.validate(pfx+".item", fileName)
		if err != nil {
			return err
		}
	case "object":
		err = p.Item.validate(pfx+".item", fileName)
		if err != nil {
			return err
		}
	}
	p.InternalFieldNameUCC = strcase.ToCamel(strings.ToLower(p.FieldName))
	p.InternalFieldNameLCC = strcase.ToLowerCamel(strings.ToLower(p.FieldName))
	return nil
}

func (p *Properties) validateNumber(pfx string, fileName string) error {
	if p.Format != "" && p.Length > 0 {
		return &SimpleError{Text: fmt.Sprintf("%s for number fields, only 'format' or 'length' can be set, not both", pfx)}
	}
	if p.Format == "" && p.Length == 0 {
		return &SimpleError{Text: fmt.Sprintf("%s for number fields, either 'format' or 'length' must be set", pfx)}
	}
	if p.Format != "" {
		err := checkIfStringNotInList(p.Format, []string{"int8", "int16", "int32", "int64", "int128", "float", "double"}, pfx, "format", fileName)
		if err != nil {
			return err
		}
	} else {
		if p.Length < 3 {
			p.Format = "int8"
		} else if p.Length < 5 {
			p.Format = "int16"
		} else if p.Length < 10 {
			p.Format = "int32"
		} else if p.Length < 19 {
			p.Format = "int64"
		} else {
			p.Format = "int128"
		}
	}
	return nil
}

func (p *Properties) validateString(pfx string, fileName string) error {
	if p.Format != "" {
		err := checkIfStringNotInList(p.Format, []string{"date", "timestamp"}, pfx, "format", fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureFieldsExist(fieldNames []string, def *Definition, fileName string, fieldPrefix string, searchInTree bool) error {
	if !fieldsExistInDefinition(fieldNames, def, []string{}, searchInTree) {
		return &SimpleError{Text: fmt.Sprintf("Not all of the given fields %s in %s exist in the boModel %s", fieldNames, fieldPrefix, fileName)}
	}
	return nil
}

func fieldsExistInDefinition(fieldNames []string, def *Definition, fieldPath []string, searchInTree bool) bool {
	for _, fieldNameTmp := range fieldNames {
		fieldNameList := []string{fieldNameTmp}
		if strings.HasPrefix(fieldNameTmp, "oneOf(") {
			fieldNameTmp = strings.Replace(fieldNameTmp, "oneOf(", "", -1)
			fieldNameTmp = strings.Replace(fieldNameTmp, ")", "", -1)
			fieldNameList = strings.Split(fieldNameTmp, ",")
		}
		for _, fieldName := range fieldNameList {
			flag := false
			for _, fDef := range def.Properties {
				fDefName := fDef.FieldName
				if len(fieldPath) > 0 {
					fDefName = fmt.Sprintf("%s.%s", strings.Join(fieldPath[:], "."), fDef.FieldName)
				}
				if fieldName == fDefName {
					flag = true
					break
				}
				if searchInTree {
					if fDef.Type == "array" || fDef.Type == "object" {
						if fieldsExistInDefinition(fieldNames, &fDef.Item, append(fieldPath, fDef.FieldName), searchInTree) {
							flag = true
							break
						}
					}
				}
			}
			if !flag {
				return false
			}
		}
	}
	return true
}
