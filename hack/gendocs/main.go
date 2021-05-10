/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"go.searchlight.dev/icinga-operator/pkg/cmds"
	"go.searchlight.dev/icinga-operator/plugins/hyperalert"

	"github.com/spf13/cobra/doc"
	"gomodules.xyz/runtime"
	"k8s.io/klog/v2"
)

const (
	version = "8.0.0"
)

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	genHostfactsDocs()
	genHyperalertDocs()
	genSearchlightDocs()
}

func genHostfactsDocs() {
	var (
		tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Hostfacts
description: Searchlight Hostfacts Reference
menu:
  product_searchlight_{{ .Version }}:
    identifier: hostfacts-cli
    name: Hostfacts
    parent: reference
    weight: 10
menu_name: product_searchlight_{{ .Version }}
---
`))

		_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  product_searchlight_{{ .Version }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: hostfacts-cli
{{- if .RootCmd }}
    weight: 0
{{ end }}
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_{{ .Version }}
{{- if .RootCmd }}
url: /products/searchlight/{{ .Version }}/reference/hostfacts/
aliases:
  - /products/searchlight/{{ .Version }}/reference/hostfacts/hostfacts/
{{ end }}
---
`))
	)
	rootCmd := cmds.NewCmdHostfacts()
	dir := runtime.GOPath() + "/src/go.searchlight.dev/icinga-operator/docs/reference/hostfacts"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		klog.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		klog.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		filename = filepath.Base(filename)
		base := strings.TrimSuffix(filename, path.Ext(filename))
		name := strings.Title(strings.Replace(base, "_", " ", -1))
		parts := strings.Split(name, " ")
		if len(parts) > 1 {
			name = strings.Join(parts[1:], " ")
		}
		data := struct {
			ID      string
			Name    string
			Version string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			name,
			version,
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			klog.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/hostfacts/" + name
	}
	doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		klog.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{ Version string }{version})
	if err != nil {
		klog.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		klog.Fatalln(err)
	}
}

func genHyperalertDocs() {
	var (
		tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Hyperalert
description: Searchlight Hyperalert Reference
menu:
  product_searchlight_{{ .Version }}:
    identifier: hyperalert-cli
    name: Hyperalert
    parent: reference
    weight: 20
menu_name: product_searchlight_{{ .Version }}
---
`))

		_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  product_searchlight_{{ .Version }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: hyperalert-cli
{{- if .RootCmd }}
    weight: 0
{{ end }}
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_{{ .Version }}
{{- if .RootCmd }}
url: /products/searchlight/{{ .Version }}/reference/hyperalert/
aliases:
  - /products/searchlight/{{ .Version }}/reference/hyperalert/hyperalert/
{{ end }}
---
`))
	)
	rootCmd := hyperalert.NewCmd()
	dir := runtime.GOPath() + "/src/go.searchlight.dev/icinga-operator/docs/reference/hyperalert"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		klog.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		klog.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		filename = filepath.Base(filename)
		base := strings.TrimSuffix(filename, path.Ext(filename))
		name := strings.Title(strings.Replace(base, "_", " ", -1))
		parts := strings.Split(name, " ")
		if len(parts) > 1 {
			name = strings.Join(parts[1:], " ")
		}
		data := struct {
			ID      string
			Name    string
			Version string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			name,
			version,
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			klog.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/hyperalert/" + name
	}
	doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		klog.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{ Version string }{version})
	if err != nil {
		klog.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		klog.Fatalln(err)
	}
}

func genSearchlightDocs() {
	var (
		tplFrontMatter = template.Must(template.New("index").Parse(`---
title: Hyperalert
description: Searchlight CLI Reference
menu:
  product_searchlight_{{ .Version }}:
    identifier: searchlight-cli
    name: Searchlight
    parent: reference
    weight: 20
menu_name: product_searchlight_{{ .Version }}
---
`))

		_ = template.Must(tplFrontMatter.New("cmd").Parse(`---
title: {{ .Name }}
menu:
  product_searchlight_{{ .Version }}:
    identifier: {{ .ID }}
    name: {{ .Name }}
    parent: searchlight-cli
{{- if .RootCmd }}
    weight: 0
{{ end }}
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_{{ .Version }}
{{- if .RootCmd }}
url: /products/searchlight/{{ .Version }}/reference/searchlight/
aliases:
  - /products/searchlight/{{ .Version }}/reference/searchlight/searchlight/
{{ end }}
---
`))
	)
	rootCmd := cmds.NewCmdSearchlight()
	dir := runtime.GOPath() + "/src/go.searchlight.dev/icinga-operator/docs/reference/searchlight"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.RemoveAll(dir)
	if err != nil {
		klog.Fatalln(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		klog.Fatalln(err)
	}

	filePrepender := func(filename string) string {
		filename = filepath.Base(filename)
		base := strings.TrimSuffix(filename, path.Ext(filename))
		name := strings.Title(strings.Replace(base, "_", " ", -1))
		parts := strings.Split(name, " ")
		if len(parts) > 1 {
			name = strings.Join(parts[1:], " ")
		}
		data := struct {
			ID      string
			Name    string
			Version string
			RootCmd bool
		}{
			strings.Replace(base, "_", "-", -1),
			name,
			version,
			!strings.ContainsRune(base, '_'),
		}
		var buf bytes.Buffer
		if err := tplFrontMatter.ExecuteTemplate(&buf, "cmd", data); err != nil {
			klog.Fatalln(err)
		}
		return buf.String()
	}

	linkHandler := func(name string) string {
		return "/docs/reference/searchlight/" + name
	}
	doc.GenMarkdownTreeCustom(rootCmd, dir, filePrepender, linkHandler)

	index := filepath.Join(dir, "_index.md")
	f, err := os.OpenFile(index, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		klog.Fatalln(err)
	}
	err = tplFrontMatter.ExecuteTemplate(f, "index", struct{ Version string }{version})
	if err != nil {
		klog.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		klog.Fatalln(err)
	}
}
