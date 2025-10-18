package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rahuldhole/xenv/app/field"
	"github.com/rahuldhole/xenv/app/prompt"
	"github.com/rahuldhole/xenv/app/script"
)

func DetermineOutputFile(formFile string) string {
	dir := filepath.Dir(formFile)
	base := filepath.Base(formFile)

	knownExtensions := []string{".xenv", ".template", ".example"}
	for _, ext := range knownExtensions {
		if strings.HasSuffix(base, ext) {
			base = strings.TrimSuffix(base, ext)
			break
		}
	}

	if !strings.HasPrefix(base, ".") {
		base = "." + base
	}

	return filepath.Join(dir, base)
}

func CheckForDSL(formFile string) bool {
	file, err := os.Open(formFile)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "@text") ||
			strings.Contains(line, "@number") ||
			strings.Contains(line, "@float") ||
			strings.Contains(line, "@date") ||
			strings.Contains(line, "@file") ||
			strings.Contains(line, "@url") ||
			strings.Contains(line, "@password") ||
			strings.Contains(line, "@select") ||
			strings.Contains(line, "@boolean") ||
			strings.Contains(line, "@list") ||
			strings.Contains(line, "@hidden") {
			return true
		}
	}
	return false
}

func ProcessFormFile(formFile, outputFile string, hasDSL bool, mergeMode bool, defaultsMode bool, allScripts bool) ([]string, error) {
	existingOutVars := make(map[string]string)
	existingFileLines := []string{}
	if mergeMode {
		if st, err := os.Stat(outputFile); err == nil && st.Size() > 0 {
			if fOut, err := os.Open(outputFile); err == nil {
				sc := bufio.NewScanner(fOut)
				kvRe := regexp.MustCompile(`^([^=]+)=(.*)$`)
				for sc.Scan() {
					line := sc.Text()
					existingFileLines = append(existingFileLines, line)
					if m := kvRe.FindStringSubmatch(line); len(m) > 2 {
						existingOutVars[m[1]] = m[2]
					}
				}
				fOut.Close()
			}
		}
	}

	file, err := os.Open(formFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keyValueRegex := regexp.MustCompile(`^([^=]+)=(.*)$`)
	var outputLines []string
	var currentFieldInfo *field.Info
	skipMode := false
	hiddenNext := false
	envVars := make(map[string]string)

	file.Seek(0, 0)
	scannerEnv := bufio.NewScanner(file)
	for scannerEnv.Scan() {
		if matches := regexp.MustCompile(`^([^=]+)=(.*)$`).FindStringSubmatch(scannerEnv.Text()); len(matches) > 0 {
			envVars[matches[1]] = matches[2]
		}
	}
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	type tmplEntry struct {
		Comments []string
		Key      string
		Line     string
	}
	templateEntries := make(map[string]*tmplEntry)
	pendingComments := []string{}
	resolvedVars := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			pendingComments = append(pendingComments, line)
			fieldInfo := field.ParseDirective(line)
			if fieldInfo != nil {
				switch fieldInfo.Type {
				case field.SkipField:
					skipMode = true
				case field.HiddenField:
					hiddenNext = true
				default:
					currentFieldInfo = fieldInfo
				}
			}
			outputLines = append(outputLines, line)
			continue
		}

		if matches := keyValueRegex.FindStringSubmatch(line); len(matches) > 0 {
			key := matches[1]
			templateDefault := matches[2]
			currentOutputValue, hasExisting := existingOutVars[key]
			existingValue := templateDefault
			if mergeMode && hasExisting {
				existingValue = currentOutputValue
			}
			conflictSuffix := ""
			if mergeMode && hasExisting && currentOutputValue != templateDefault {
				conflictSuffix = fmt.Sprintf(" (current: %s | template: %s)", currentOutputValue, templateDefault)
			}

			var finalVal string

			if defaultsMode {
				if currentFieldInfo != nil && currentFieldInfo.Script != "" && allScripts {
					finalVal = script.RunAuto(currentFieldInfo.Script, outputLines, envVars, resolvedVars, key)
				} else {
					finalVal = existingValue
				}
				currentFieldInfo = nil
			} else {
				switch {
				case skipMode:
					finalVal = existingValue
				case hiddenNext:
					finalVal = existingValue
					hiddenNext = false
				case currentFieldInfo != nil && currentFieldInfo.Type != field.NoField:
					if currentFieldInfo.Readonly {
						label := currentFieldInfo.Label
						if label == "" {
							label = key
						}
						if conflictSuffix != "" {
							fmt.Printf("%s%s (readonly): %s\n", label, conflictSuffix, existingValue)
						} else {
							fmt.Printf("%s (readonly): %s\n", label, existingValue)
						}
						finalVal = existingValue
						currentFieldInfo = nil
					} else if currentFieldInfo.Script != "" {
						if allScripts {
							finalVal = script.RunAuto(currentFieldInfo.Script, outputLines, envVars, resolvedVars, key)
						} else {
							fmt.Printf("Script for %s%s\n", key, conflictSuffix)
							fmt.Print("Run script? [y/N/v (view)]: ")
							reader := bufio.NewReader(os.Stdin)
							resp, _ := reader.ReadString('\n')
							resp = strings.TrimSpace(strings.ToLower(resp))
							if resp == "view" || resp == "v" {
								fmt.Println("Script content:")
								fmt.Println(currentFieldInfo.Script)
								fmt.Print("Run this script now? [y/N]: ")
								resp, _ = reader.ReadString('\n')
								resp = strings.TrimSpace(strings.ToLower(resp))
							}
							if resp == "y" || resp == "yes" {
								finalVal = script.RunAuto(currentFieldInfo.Script, outputLines, envVars, resolvedVars, key)
							} else {
								fmt.Println("Skipped script.")
								finalVal = existingValue
							}
						}
						currentFieldInfo = nil
					} else {
						label := currentFieldInfo.Label
						if label == "" {
							label = key
						}
						if conflictSuffix != "" {
							label += conflictSuffix
						}
						val, err := prompt.ForValue(currentFieldInfo, label, existingValue)
						if err != nil {
							return nil, err
						}
						finalVal = val
						currentFieldInfo = nil
					}
				case !hasDSL || key == "PLAIN_TEXT_VAR":
					label := key
					if conflictSuffix != "" {
						label += conflictSuffix
					}
					val, err := prompt.Input(label, existingValue)
					if err != nil {
						return nil, err
					}
					finalVal = val
				default:
					finalVal = existingValue
				}
			}

			finalLine := key + "=" + finalVal
			outputLines = append(outputLines, finalLine)
			resolvedVars[key] = finalVal

			if _, ok := templateEntries[key]; !ok {
				templateEntries[key] = &tmplEntry{
					Comments: append([]string{}, pendingComments...),
					Key:      key,
					Line:     finalLine,
				}
			} else {
				templateEntries[key].Line = finalLine
			}
			pendingComments = []string{}
		} else {
			outputLines = append(outputLines, line)
			pendingComments = []string{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !mergeMode {
		return outputLines, nil
	}

	mergedLines := []string{}
	seenKeys := make(map[string]struct{})
	kvRe := regexp.MustCompile(`^([^=]+)=(.*)$`)

	for _, line := range existingFileLines {
		if m := kvRe.FindStringSubmatch(line); len(m) > 2 {
			k := m[1]
			if entry, ok := templateEntries[k]; ok {
				mergedLines = append(mergedLines, entry.Line)
				seenKeys[k] = struct{}{}
				continue
			}
			mergedLines = append(mergedLines, line)
			seenKeys[k] = struct{}{}
		} else {
			mergedLines = append(mergedLines, line)
		}
	}

	for k, entry := range templateEntries {
		if _, ok := seenKeys[k]; ok {
			continue
		}
		for _, c := range entry.Comments {
			mergedLines = append(mergedLines, c)
		}
		mergedLines = append(mergedLines, entry.Line)
	}

	return mergedLines, nil
}
