package builder

// generate go client boilerplate which can perform the request

import (
	"fmt"
	"text/template"
)

type goVars struct {
	Method           string
	CreateBodyReader string
	BodyReader       string
	Query            map[string]string
}

// Go returns go snippet calling the url
func (b *Builder) Go(req *Request) error {

	const goTemplate = `
package main

import (
    "context"
    "net/http"
    "strings"
)

func main() {

    ctx := context.Background()
    {{ .CreateBodyReader }}
    req, err := http.NewRequestWithContext(
        ctx,
        {{ .Method }},
        {{ .BodyReader }},
    )
    if err != nil {
        log.Fatal(err)
    }
    {{ range $name, $value := .Query }}
    req.Header.Add("{{ $name }}", "{{ $value }}")
    {{ end }}
    tr := &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: false,
    }
    client := &http.Client{Transport: tr}

    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("resp=%+v", resp)

}`

	bodyReader := "nil"
	createBodyReader := "\n"
	if req.Method == "POST" || req.Method == "PUT" {
		createBodyReader = fmt.Sprintf("body := strings.NewReader(`%s`)", req.Body)
		bodyReader = "body"
	}

	foo := goVars{
		Method:           req.Method,
		BodyReader:       bodyReader,
		CreateBodyReader: createBodyReader,
		Query:            req.Query,
	}
	// Create a new template and parse the letter into it.
	t, err := template.New("go").Parse(goTemplate)

	if err != nil {
		return err
	}

	err = t.Execute(b.w, foo)
	if err != nil {
		return err
	}

	return nil
}
