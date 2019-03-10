package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v24/github"
	"github.com/intel/tfortools"
	"golang.org/x/oauth2"
)

var (
	token   string
	pattern string
	orgs    string
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
	flag.StringVar(&orgs, "o", "devopsfaith", "comma separated list of github orgs")
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

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	contributorsMap := map[string]github.Contributor{}
	for _, org := range strings.Split(orgs, ",") {
		repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

		if err != nil {
			fmt.Println("error:", err.Error())
		}

		re := regexp.MustCompile(pattern)

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
				if cs.Author.Contributions == nil {
					cs.Author.Contributions = cs.Total
				}
				if c, ok := contributorsMap[*cs.Author.Login]; ok {
					*cs.Author.Contributions += *c.Contributions
				}
				contributorsMap[*cs.Author.Login] = *cs.Author
			}
		}
	}

	contributors := []Contributor{}
	for _, c := range contributorsMap {
		contributor := Contributor{}
		if c.Login != nil {
			contributor.Login = *c.Login
		}
		if c.ID != nil {
			contributor.ID = *c.ID
		}
		if c.AvatarURL != nil {
			contributor.AvatarURL = *c.AvatarURL
		}
		if c.GravatarID != nil {
			contributor.GravatarID = *c.GravatarID
		}
		if c.URL != nil {
			contributor.URL = *c.URL
		}
		if c.HTMLURL != nil {
			contributor.HTMLURL = *c.HTMLURL
		}
		if c.FollowersURL != nil {
			contributor.FollowersURL = *c.FollowersURL
		}
		if c.FollowingURL != nil {
			contributor.FollowingURL = *c.FollowingURL
		}
		if c.GistsURL != nil {
			contributor.GistsURL = *c.GistsURL
		}
		if c.StarredURL != nil {
			contributor.StarredURL = *c.StarredURL
		}
		if c.SubscriptionsURL != nil {
			contributor.SubscriptionsURL = *c.SubscriptionsURL
		}
		if c.OrganizationsURL != nil {
			contributor.OrganizationsURL = *c.OrganizationsURL
		}
		if c.ReposURL != nil {
			contributor.ReposURL = *c.ReposURL
		}
		if c.EventsURL != nil {
			contributor.EventsURL = *c.EventsURL
		}
		if c.ReceivedEventsURL != nil {
			contributor.ReceivedEventsURL = *c.ReceivedEventsURL
		}
		if c.Type != nil {
			contributor.Type = *c.Type
		}
		if c.SiteAdmin != nil {
			contributor.SiteAdmin = *c.SiteAdmin
		}
		if c.Contributions != nil {
			contributor.Contributions = *c.Contributions
		}
		contributors = append(contributors, contributor)
	}

	fmt.Println("dumping contributor stats", len(contributors))

	if err := tfortools.OutputToTemplate(os.Stdout, "contributors", tmpl, contributors, nil); err != nil {
		fmt.Printf("error executing template:", err)
	}
}

type Contributor struct {
	Login             string `json:"login,omitempty"`
	ID                int64  `json:"id,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	GravatarID        string `json:"gravatar_id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	Type              string `json:"type,omitempty"`
	SiteAdmin         bool   `json:"site_admin,omitempty"`
	Contributions     int    `json:"contributions,omitempty"`
}

const (
	defaultTemplate = `{{range .}}{{.Login}}
{{end}}`
)
