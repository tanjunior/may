package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	outDir := flag.String("out", "frontend/src", "output directory for generated files")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create out dir: %v\n", err)
		os.Exit(2)
	}

	// Read backend/errors.go
	errsSrc, err := ioutil.ReadFile("backend/errors.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read backend/errors.go: %v\n", err)
		os.Exit(2)
	}
	errs := string(errsSrc)

	// Extract const string assignments (CODE = "VALUE")
	reConstLine := regexp.MustCompile(`([A-Za-z0-9_]+)\s*=\s*"([^"]+)"`)
	consts := reConstLine.FindAllStringSubmatch(errs, -1)
	codes := make([][2]string, 0, len(consts))
	for _, m := range consts {
		// Filter likely error code names (prefix Code)
		if strings.HasPrefix(m[1], "Code") {
			codes = append(codes, [2]string{m[1], m[2]})
		}
	}

	// Extract ErrorMessages map entries by scanning for CodeName: "message"
	messages := map[string]string{}
	for _, c := range codes {
		name := c[0]
		// try plain `CodeName: "msg"`
		rePlain := regexp.MustCompile(fmt.Sprintf(`%s\s*:\s*"([^"]+)"`, regexp.QuoteMeta(name)))
		if m := rePlain.FindStringSubmatch(errs); len(m) >= 2 {
			messages[name] = m[1]
			continue
		}
		// try bracketed `[CodeName]: "msg"`
		reBracket := regexp.MustCompile(fmt.Sprintf(`\[%s\]\s*:\s*"([^"]+)"`, regexp.QuoteMeta(name)))
		if m := reBracket.FindStringSubmatch(errs); len(m) >= 2 {
			messages[name] = m[1]
			continue
		}
	}

	// Generate errorCodes.ts
	var b strings.Builder
	b.WriteString("// GENERATED FROM Go constants in backend/errors.go\n")
	b.WriteString("// Keep in sync with backend; used by frontend for error-code checks and messages.\n\n")
	for _, c := range codes {
		name := c[0]
		val := c[1]
		b.WriteString(fmt.Sprintf("export const %s = \"%s\";\n", name, val))
	}
	b.WriteString("\nexport const ErrorMessages: Record<string, string> = {\n")
	for _, c := range codes {
		name := c[0]
		if msg, ok := messages[name]; ok {
			b.WriteString(fmt.Sprintf("  [%s]: \"%s\",\n", name, msg))
		}
	}
	b.WriteString("};\n\n")
	b.WriteString("export default {\n")
	for _, c := range codes {
		b.WriteString(fmt.Sprintf("  %s,\n", c[0]))
	}
	b.WriteString("  ErrorMessages,\n};\n")

	err = ioutil.WriteFile(filepath.Join(*outDir, "errorCodes.ts"), []byte(b.String()), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write errorCodes.ts: %v\n", err)
		os.Exit(2)
	}

	// Read backend/model.go to infer Product fields
	modelSrc, err := ioutil.ReadFile("backend/model.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read backend/model.go: %v\n", err)
		os.Exit(2)
	}
	prodFields := parseProductFields(string(modelSrc))

	// Generate apiTypes.ts
	var t strings.Builder
	t.WriteString("// GENERATED: API response types for frontend\n\n")
	t.WriteString("export interface APIError {\n  code: string;\n  message: string;\n  details?: any;\n}\n\n")
	t.WriteString("export interface ErrorEnvelope {\n  success: false;\n  status: number;\n  error: APIError;\n}\n\n")
	t.WriteString("export interface SuccessEnvelope<T> {\n  success: true;\n  status: number;\n  data: T;\n  meta?: Record<string, any>;\n}\n\n")
	// Product
	t.WriteString("export interface Product {\n")
	for _, f := range prodFields {
		t.WriteString(fmt.Sprintf("  %s: %s;\n", f.Name, f.TSType))
	}
	t.WriteString("}\n\n")
	t.WriteString("export type ProductListResponse = SuccessEnvelope<Product[]>;\n")
	t.WriteString("export type ProductResponse = SuccessEnvelope<Product>;\n")
	t.WriteString("export type APIFailure = ErrorEnvelope;\n")

	err = ioutil.WriteFile(filepath.Join(*outDir, "apiTypes.ts"), []byte(t.String()), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write apiTypes.ts: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stderr, "generated files in %s\n", *outDir)
}

type Field struct {
	Name   string
	TSType string
}

// parseProductFields does a small heuristic parse for `type Product struct` fields
// and returns a slice of Field suitable for TypeScript generation. It also
// injects the embedded gorm.Model fields as typical fields.
func parseProductFields(src string) []Field {
	// default fields from gorm.Model
	fields := []Field{
		{Name: "id", TSType: "number"},
		{Name: "code", TSType: "string"},
		{Name: "price", TSType: "number"},
		{Name: "createdAt", TSType: "string"},
		{Name: "updatedAt", TSType: "string"},
		{Name: "deletedAt", TSType: "string | null"},
	}

	// try to parse explicit fields in the struct (e.g., Code string, Price uint)
	scanner := bufio.NewScanner(strings.NewReader(src))
	inStruct := false
	reField := regexp.MustCompile(`^\s*([A-Za-z0-9_]+)\s+([A-Za-z0-9_\.\[\]]+)`)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "type Product struct") {
			inStruct = true
			continue
		}
		if inStruct {
			if strings.HasPrefix(strings.TrimSpace(line), "}") {
				break
			}
			m := reField.FindStringSubmatch(line)
			if len(m) >= 3 {
				name := m[1]
				typ := m[2]
				// skip gorm.Model embedded
				if name == "gorm.Model" || name == "Model" {
					continue
				}
				// map Go types to TS
				ts := "any"
				switch typ {
				case "string":
					ts = "string"
				case "uint", "int", "uint32", "uint64", "int32", "int64":
					ts = "number"
				default:
					if strings.HasPrefix(typ, "[]") {
						ts = "any[]"
					} else {
						ts = "any"
					}
				}
				// replace or append field
				replaced := false
				for i := range fields {
					if strings.EqualFold(fields[i].Name, name) {
						fields[i].TSType = ts
						replaced = true
						break
					}
				}
				if !replaced {
					fields = append(fields, Field{Name: strings.ToLower(name[:1]) + name[1:], TSType: ts})
				}
			}
		}
	}
	return fields
}
