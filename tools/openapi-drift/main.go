// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

const manifestPath = "tools/openapi-drift/sources.json"

type manifest struct {
	Sources  []source  `json:"sources"`
	Packages []apiSpec `json:"packages"`
}

type source struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type apiSpec struct {
	Name                    string `json:"name"`
	Local                   string `json:"local"`
	Source                  string `json:"source"`
	SkipPaths               bool   `json:"skipPaths,omitempty"`
	SkipComponentSetChanges bool   `json:"skipComponentSetChanges,omitempty"`
}

type finding struct {
	name    string
	details []string
}

type summaryRow struct {
	name            string
	localVersion    string
	portalVersion   string
	differenceCount int
}

type differences struct {
	items []string
	total int
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}

	var config manifest
	if err := json.Unmarshal(manifestData, &config); err != nil {
		return fmt.Errorf("parse manifest: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	sources := make(map[string]source, len(config.Sources))
	documents := make(map[string]map[string]any, len(config.Sources))
	for _, src := range config.Sources {
		if src.ID == "" || src.URL == "" {
			return fmt.Errorf("source entry must have an id and URL")
		}
		if _, exists := sources[src.ID]; exists {
			return fmt.Errorf("duplicate source id %q", src.ID)
		}
		body, err := download(client, src.URL)
		if err != nil {
			return fmt.Errorf("download %s: %w", src.URL, err)
		}
		doc, err := parseSpec(body)
		if err != nil {
			return fmt.Errorf("parse upstream spec %s: %w", src.ID, err)
		}
		sources[src.ID] = src
		documents[src.ID] = doc
	}

	var findings []finding
	var summary []summaryRow
	for _, pkg := range config.Packages {
		upstream, ok := documents[pkg.Source]
		if !ok {
			return fmt.Errorf("package %q references unknown source %q", pkg.Name, pkg.Source)
		}
		localBytes, err := os.ReadFile(pkg.Local)
		if err != nil {
			return fmt.Errorf("read local spec %s: %w", pkg.Local, err)
		}
		local, err := parseSpec(localBytes)
		if err != nil {
			return fmt.Errorf("parse local spec %s: %w", pkg.Local, err)
		}

		var diff differences
		localVersion := nestedString(local, "info", "version")
		upstreamVersion := nestedString(upstream, "info", "version")
		summary = append(summary, summaryRow{
			name:          pkg.Name,
			localVersion:  localVersion,
			portalVersion: upstreamVersion,
		})
		localComponents := nestedMap(local, "components")
		upstreamComponents := nestedMap(upstream, "components")
		diff.addMapDifferences(
			"components",
			localComponents,
			upstreamComponents,
			pkg.SkipComponentSetChanges,
		)

		if !pkg.SkipPaths {
			diff.addPathDifferences(nestedMap(local, "paths"), nestedMap(upstream, "paths"))
		}
		summary[len(summary)-1].differenceCount = diff.total

		if localVersion != upstreamVersion {
			diff.add(fmt.Sprintf("info.version: local %q; portal %q", localVersion, upstreamVersion))
		}

		if diff.total > 0 {
			details := diff.items
			if diff.total > len(diff.items) {
				details = append(details, fmt.Sprintf("Showing %d of %d differences.", len(diff.items), diff.total))
			}
			findings = append(findings, finding{
				name: pkg.Name,
				details: append(
					details,
					fmt.Sprintf("Local: %s", pkg.Local),
					fmt.Sprintf("Portal: %s", sources[pkg.Source].URL),
				),
			})
		}
	}

	if err := writeFindings(summary, findings); err != nil {
		return err
	}
	if len(findings) > 0 {
		fmt.Printf("Found differences in %d API definition(s).\n", len(findings))
	} else {
		fmt.Println("No structural differences found.")
	}
	return nil
}

func download(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "sacloud-sdk-go-openapi-drift-check")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %s", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	return body, nil
}

func parseSpec(body []byte) (map[string]any, error) {
	var doc map[string]any
	if err := yaml.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	if len(doc) == 0 {
		return nil, fmt.Errorf("empty document")
	}
	if nestedString(doc, "openapi") == "" && nestedString(doc, "swagger") == "" {
		return nil, fmt.Errorf("document has no OpenAPI or Swagger version")
	}
	if nestedString(doc, "info", "version") == "" {
		return nil, fmt.Errorf("document has no info.version")
	}
	return doc, nil
}

func nestedMap(root map[string]any, keys ...string) map[string]any {
	var current any = root
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[key]
	}
	result, _ := current.(map[string]any)
	return result
}

func nestedString(root map[string]any, keys ...string) string {
	var current any = root
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return ""
		}
		current = m[key]
	}
	result, _ := current.(string)
	return result
}

func objectKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (d *differences) add(details ...string) {
	d.total += len(details)
	for _, detail := range details {
		if len(d.items) < 30 {
			d.items = append(d.items, detail)
		}
	}
}

func (d *differences) addSetDifference(label string, upstream, local []string) {
	localSet := make(map[string]struct{}, len(local))
	for _, key := range local {
		localSet[key] = struct{}{}
	}
	upstreamSet := make(map[string]struct{}, len(upstream))
	for _, key := range upstream {
		upstreamSet[key] = struct{}{}
	}
	for _, key := range upstream {
		if _, ok := localSet[key]; !ok {
			d.add(fmt.Sprintf("%s.%s: present in portal, missing locally", label, key))
		}
	}
	for _, key := range local {
		if _, ok := upstreamSet[key]; !ok {
			d.add(fmt.Sprintf("%s.%s: present locally, missing from portal", label, key))
		}
	}
}

func (d *differences) addMapDifferences(label string, local, upstream map[string]any, skipComponentSetChanges bool) {
	keys := union(objectKeys(local), objectKeys(upstream))
	for _, key := range keys {
		localValue, localOK := local[key]
		upstreamValue, upstreamOK := upstream[key]
		localItems, localIsMap := localValue.(map[string]any)
		upstreamItems, upstreamIsMap := upstreamValue.(map[string]any)
		if skipComponentSetChanges {
			if localIsMap && upstreamIsMap {
				for _, item := range intersection(objectKeys(localItems), objectKeys(upstreamItems)) {
					d.compare(label+"."+key+"."+item, localItems[item], upstreamItems[item])
				}
			}
			continue
		}
		if key == "schemas" && localIsMap && upstreamIsMap {
			d.addSetDifference("components.schemas", objectKeys(upstreamItems), objectKeys(localItems))
			for _, schema := range intersection(objectKeys(localItems), objectKeys(upstreamItems)) {
				d.compare(label+".schemas."+schema, localItems[schema], upstreamItems[schema])
			}
			continue
		}
		if !localOK {
			d.add(fmt.Sprintf("%s.%s: present in portal, missing locally", label, key))
			continue
		}
		if !upstreamOK {
			d.add(fmt.Sprintf("%s.%s: present locally, missing from portal", label, key))
			continue
		}
		d.compare(label+"."+key, localValue, upstreamValue)
	}
}

func (d *differences) addPathDifferences(local, upstream map[string]any) {
	for _, route := range union(objectKeys(local), objectKeys(upstream)) {
		localPath, localOK := local[route]
		upstreamPath, upstreamOK := upstream[route]
		if !localOK {
			d.add(fmt.Sprintf("paths.%s: present in portal, missing locally", route))
			continue
		}
		if !upstreamOK {
			d.add(fmt.Sprintf("paths.%s: present locally, missing from portal", route))
			continue
		}
		d.compare("paths."+route, localPath, upstreamPath)
	}
}

func (d *differences) compare(label string, local, upstream any) {
	localMap, localIsMap := local.(map[string]any)
	upstreamMap, upstreamIsMap := upstream.(map[string]any)
	if localIsMap && upstreamIsMap {
		for _, key := range union(objectKeys(localMap), objectKeys(upstreamMap)) {
			localValue, localOK := localMap[key]
			upstreamValue, upstreamOK := upstreamMap[key]
			switch {
			case !localOK:
				d.add(fmt.Sprintf("%s.%s: present in portal, missing locally", label, key))
			case !upstreamOK:
				d.add(fmt.Sprintf("%s.%s: present locally, missing from portal", label, key))
			default:
				d.compare(label+"."+key, localValue, upstreamValue)
			}
		}
		return
	}

	localArray, localIsArray := local.([]any)
	upstreamArray, upstreamIsArray := upstream.([]any)
	if localIsArray && upstreamIsArray {
		if len(localArray) != len(upstreamArray) {
			d.add(fmt.Sprintf("%s: array length local %d; portal %d", label, len(localArray), len(upstreamArray)))
		}
		limit := min(len(localArray), len(upstreamArray))
		for i := range limit {
			d.compare(fmt.Sprintf("%s[%d]", label, i), localArray[i], upstreamArray[i])
		}
		return
	}

	if !jsonEqual(local, upstream) {
		d.add(fmt.Sprintf("%s: value differs (local %s; portal %s)", label, compactJSON(local), compactJSON(upstream)))
	}
}

func jsonEqual(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}

func compactJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	text := string(encoded)
	if len(text) > 100 {
		return text[:97] + "..."
	}
	return text
}

func union(left, right []string) []string {
	values := make(map[string]struct{}, len(left)+len(right))
	for _, value := range left {
		values[value] = struct{}{}
	}
	for _, value := range right {
		values[value] = struct{}{}
	}
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func intersection(left, right []string) []string {
	rightSet := make(map[string]struct{}, len(right))
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	var result []string
	for _, value := range left {
		if _, ok := rightSet[value]; ok {
			result = append(result, value)
		}
	}
	return result
}

func writeFindings(summary []summaryRow, findings []finding) error {
	var output strings.Builder
	versionDifferences := 0
	structuralDifferences := 0
	for _, row := range summary {
		if row.localVersion != row.portalVersion {
			versionDifferences++
		}
		if row.differenceCount > 0 {
			structuralDifferences++
		}
	}

	fmt.Fprintf(
		&output,
		"# OpenAPI comparison summary\n\nCompared APIs: %d; version differences: %d; APIs with structural differences: %d.\n\n",
		len(summary),
		versionDifferences,
		structuralDifferences,
	)
	output.WriteString("| API component | Repository version | Portal version | Status | Structural differences |\n")
	output.WriteString("|---|---:|---:|---|---:|\n")
	for _, row := range summary {
		status := ":white_check_mark:"
		if row.localVersion != row.portalVersion || row.differenceCount > 0 {
			status = ":x:"
		}
		fmt.Fprintf(
			&output,
			"| %s | %s | %s | %s | %d |\n",
			row.name,
			row.localVersion,
			row.portalVersion,
			status,
			row.differenceCount,
		)
	}

	if len(findings) == 0 {
		output.WriteString("\nNo structural differences found.\n")
	} else {
		output.WriteString("\n# OpenAPI structural differences\n")
		for _, item := range findings {
			fmt.Fprintf(&output, "\n## %s\n", item.name)
			for _, detail := range item.details {
				fmt.Fprintf(&output, "- %s\n", detail)
			}
		}
	}

	text := output.String()
	fmt.Print(text)
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		summary, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("open GitHub Actions step summary: %w", err)
		}
		if _, err := io.WriteString(summary, text); err != nil {
			_ = summary.Close()
			return fmt.Errorf("write GitHub Actions step summary: %w", err)
		}
		if err := summary.Close(); err != nil {
			return fmt.Errorf("close GitHub Actions step summary: %w", err)
		}
	}
	return nil
}
