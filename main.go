package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/google/go-github/v24/github"
	"github.com/intel/tfortools"
	"golang.org/x/oauth2"
)

var (
	token   string
	pattern string
	org     string
	tmpl    string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-f template] [-p pattern] [-o organization] [-t token]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, tfortools.GenerateUsageUndecorated([][]string{}))
	}
	flag.StringVar(&token, "t", "", "github personal token")
	flag.StringVar(&pattern, "p", ".*", "reggex pattern for filtering repos by name")
	flag.StringVar(&org, "o", "devopsfaith", "github org")
	flag.StringVar(&tmpl, "f", defaultTemplate, "template for render the results")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	// t := template.Must(template.New("contributors").Parse(*tmpl))

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

	if err != nil {
		fmt.Println("error:", err.Error())
	}

	re := regexp.MustCompile(pattern)

	contributors := map[string]github.Contributor{}
	for k, v := range repos {
		if !re.MatchString(*v.Name) {
			continue
		}

		fmt.Println("repo", k, *v.Name)
		css, _, err := client.Repositories.ListContributorsStats(ctx, org, *v.Name)
		if err != nil {
			fmt.Printf("error collecting stats for repo %s: %s\n", *v.Name, err.Error())
		}
		for _, cs := range css {
			contributors[*cs.Author.Login] = *cs.Author
		}
	}

	fmt.Println("dumping contributor stats", len(contributors))

	if err := tfortools.OutputToTemplate(os.Stdout, "contributors", tmpl, contributors, nil); err != nil {
		// if err := t.Execute(os.Stdout, contributors); err != nil {
		fmt.Printf("error executing template:", err)
	}
}

const (
	defaultTemplate = `{{range .}}{{.Login}}
{{end}}`
)
