package TagCommunityIssues

import (
	"context"
	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

// HTTP Cloud Function.
func TagCommunityIssues(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, nil)
	if err != nil {
		log.Fatal(err)
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("received event ", event)

	switch event := event.(type) {
	case *github.IssuesEvent:
		action := event.GetAction()
		if action == "opened" || action == "edited" || action == "reopened" {
			handleIssuesEvent(event)
		}
	}
}

func handleIssuesEvent(e *github.IssuesEvent) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_OAUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	//check organization membership
	issue := e.Issue
	userName := *issue.User.Login
	targetOrg := os.Getenv("TARGET_GITHUB_ORG")
	isMember, _, err := client.Organizations.IsMember(ctx, targetOrg, userName)

	if err != nil {
		log.Fatal("error", err)
	}

	if !isMember {
		addCommunityLabel(e, client, ctx)
	}

	log.Println("handled issues event")
}

func addCommunityLabel(e *github.IssuesEvent, client *github.Client, ctx context.Context) {
	issue := e.Issue
	currentLabels := make([]string, 0, 0)
	labels := issue.Labels
	for _, l := range labels {
		currentLabels = append(currentLabels, *l.Name)
	}

	label := os.Getenv("COMMUNITY_LABEL")
	if !labelExists(currentLabels, label) {
		owner := *e.Repo.Owner.Login
		repo := *e.Repo.Name
		issueNumber := *issue.Number
		currentLabels = append(currentLabels, label)
		request := github.IssueRequest{}
		request.Labels = &currentLabels
		client.Issues.Edit(ctx, owner, repo, issueNumber, &request)
		log.Printf("added the %v label to issue: %v/%v#%v", label, owner, repo, issueNumber)
	}
}

func labelExists(labels []string, label string) bool {
	for _, n := range labels {
		if label == n {
			return true
		}
	}
	return false
}