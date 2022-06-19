package jira

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"gopkg.in/andygrunwald/go-jira.v1"
)

const DateFormat = "2006-01-02T15:04:05.000-0700"

type Client struct {
	*jira.Client
}

func NewClient(url, accessToken string) (Client, error) {
	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken, TokenType: "bearer"})
	ctx := context.Background()
	client := oauth2.NewClient(ctx, token)

	jiraClient, err := jira.NewClient(client, url)
	if err != nil {
		return Client{}, err
	}

	return Client{jiraClient}, nil
}

// GetIssues returns all Jira issues specified by project and fields created or updated in last n days.
func (c *Client) GetIssues(project string, fields map[string][]string, n uint) ([]jira.Issue, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("project = '%s' AND ", project))
	for key, value := range fields {
		if len(value) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("'%s' IN ('%s') AND", key, strings.Join(value, "','")))
	}
	sb.WriteString(fmt.Sprintf("(created>=-%dd OR updated>=-%dd)", n, n))

	req := sb.String()
	issues, _, err := c.Issue.Search(req, nil)
	return issues, err
}
