# github-contributors
small tool to collect contributors from github repos and organizations

## Installation

```
go get -u github.com/kpacha/github-contributors
go install github.com/kpacha/github-contributors
```

## Run

```
$ github-contributors -h
Usage: github-contributors [-f template] [-p pattern] [-o organization] [-t token]
  -f string
    	template for render the results (default "{{range .}}{{.Login}}\n{{end}}")
  -o string
    	comma separated list of github orgs (default "devopsfaith")
  -p string
    	reggex pattern for filtering repos by name (default ".*")
  -t string
    	github personal token

[]struct {
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

```

## Output manipulation

This project includes the tfortools lib, so there are tons of functions and macros available:

```
github-contributors -o "devopsfaith" -p "krakend" -f "{{table (cols (sort . "Contributions" "dsc") "Login" "Contributions")}}"
```
