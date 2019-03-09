package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
)

func main() {
	token := flag.String("t", "", "github personal token")
	pattern := flag.String("p", ".*", "reggex pattern for filtering repos by name")
	org := flag.String("o", "devopsfaith", "github org")
	tmpl := flag.String("f", "", "template for render the results")

	flag.Parse()

	ctx := context.Background()
	var tc *http.Client

	if *token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	if *tmpl == "" {
		*tmpl = defaultTemplate
	}

	t := template.Must(template.New("contributors").Parse(*tmpl))

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	repos, _, err := client.Repositories.ListByOrg(ctx, *org, opt)

	if err != nil {
		fmt.Println("error:", err.Error())
	}

	re := regexp.MustCompile(*pattern)

	contributors := map[string]github.Contributor{}
	for k, v := range repos {
		if !re.MatchString(*v.Name) {
			continue
		}

		fmt.Println("repo", k, *v.Name)
		css, _, err := client.Repositories.ListContributorsStats(ctx, *org, *v.Name)
		if err != nil {
			fmt.Printf("error collecting stats for repo %s: %s\n", *v.Name, err.Error())
		}
		for _, cs := range css {
			contributors[*cs.Author.Login] = *cs.Author
		}
	}

	fmt.Println("dumping contributor stats", len(contributors))

	if err := t.Execute(os.Stdout, contributors); err != nil {
		fmt.Printf("error executing template:", err)
	}
}

const (
	defaultTemplate = `{{range .}}{{.Login}}
{{end}}`
)
