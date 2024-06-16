package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Struct struct {
	Name   string
	Fields []Field
}

func main() {
	// Đọc định nghĩa ASN.1 từ file input
	inputFile := "data.asn1"
	asn1Definition, err := readASN1Definition(inputFile)
	if err != nil {
		fmt.Println("Error reading ASN.1 definition:", err)
		return
	}

	// Phân tích định nghĩa ASN.1 và tạo các struct
	structs := parseASN1Definition(asn1Definition)

	// Ghi các struct Go vào file output
	outputFile := "generated_structs.go"
	err = generateGoCode(outputFile, structs)
	if err != nil {
		fmt.Println("Error generating Go code:", err)
	}
}

func readASN1Definition(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func parseASN1Definition(definition string) []Struct {
	lines := strings.Split(definition, "\n")
	var structs []Struct
	var currentStruct *Struct

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}

		if strings.Contains(line, "::= SEQUENCE") || strings.Contains(line, "::= SET") || strings.Contains(line, "::= CHOICE") {
			if currentStruct != nil {
				structs = append(structs, *currentStruct)
			}
			name := strings.Fields(line)[0]
			currentStruct = &Struct{Name: name}
		} else if strings.Contains(line, "}") {
			if currentStruct != nil {
				structs = append(structs, *currentStruct)
				currentStruct = nil
			}
		} else if currentStruct != nil {
			parts := strings.Fields(line)
			nameData := parts[0]
			typeData := ""
			for i := 1; i < len(parts); i++ {
				typeData += parts[i] + " "
			}
			typeData = typeData[:len(typeData)-1]

			if typeData[len(typeData)-1] == ',' {
				typeData = typeData[:len(typeData)-1]
			}
			typeData = mapASN1TypeToGoType(typeData)

			if len(parts) >= 2 {
				field := Field{
					Name: nameData,
					Type: typeData,
				}
				currentStruct.Fields = append(currentStruct.Fields, field)
			}
		}
	}

	return structs
}

func mapASN1TypeToGoType(asn1Type string) string {
	switch asn1Type {
	case "BOOLEAN":
		return "bool"
	case "INTEGER":
		return "int"
	case "BIT STRING":
		return "asn1.BitString"
	case "OCTET STRING":
		return "[]byte"
	case "NULL":
		return "asn1.RawValue"
	case "UTF8String":
		return "string"
	case "IA5String":
		return "string"
	case "OBJECT IDENTIFIER":
		return "asn1.ObjectIdentifier"
	default:
		return "interface{}"
	}
}

func generateGoCode(filename string, structs []Struct) error {
	const structTemplate = `
package main
// This is auto generate, DO NOT FUCKING EDIT
import (
	"encoding/asn1"
)


{{- range . }}
type {{ .Name }} struct {
	{{- range .Fields }}
	{{ .Name }} {{ .Type }} ` + "`asn1:\"{{ .Tag }}\"`" + `
	{{- end }}
}
{{- end }}
`
	tmpl, err := template.New("structs").Parse(structTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Default().Println(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	err = tmpl.Execute(writer, structs)
	if err != nil {
		return err
	}

	return writer.Flush()
}
