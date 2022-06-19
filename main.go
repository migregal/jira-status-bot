package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"statusbot/config"
	"statusbot/jira"
	"statusbot/mail"
)

type Issue struct {
	URL     string
	Key     string
	Summary string
	Status  string
	Comment string
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	issues, err := getIssues(cfg)
	if err != nil {
		log.Fatal(err)
	}

	body := getMailBody(issues)

	client := mail.Mailer{
		Server:   cfg.GetMailServer(),
		User:     cfg.GetMailSender(),
		Password: cfg.GetMailSenderPassword(),
	}

	client.SendEmail(getSubject(), body, cfg.GetMailReciever())
}

func getIssues(cfg config.Config) (map[string][]Issue, error) {
	jiraClient, err := jira.NewClient(cfg.GetBaseURL(), cfg.GetAccessToken())
	if err != nil {
		return nil, err
	}

	projects, err := cfg.GetProjects()
	if err != nil {
		return nil, err
	}

	res := map[string][]Issue{}
	for key, value := range projects {
		issues, err := jiraClient.GetIssues(key, value, cfg.GetDays())
		if err != nil {
			return nil, err
		}

		res[key] = []Issue{}
		for _, issue := range issues {
			comment := ""
			data, _, err := jiraClient.Issue.Get(issue.ID, nil)
			if err == nil {
				if data.Fields.Comments == nil || len(data.Fields.Comments.Comments) == 0 {
					continue
				}

				lastComment := data.Fields.Comments.Comments[len(data.Fields.Comments.Comments)-1]

				lastUpdated, err := time.Parse(jira.DateFormat, lastComment.Updated)
				if err != nil {
					continue
				}

				if lastUpdated.Before(time.Now().Add(-7 * 24 * time.Hour)) {
					continue
				}

				comment = lastComment.Body
			}

			issueURL := fmt.Sprintf("%s/browse/%s", cfg.GetBaseURL(), issue.Key)

			res[key] = append(res[key], Issue{
				URL:     issueURL,
				Key:     issue.Key,
				Summary: issue.Fields.Summary,
				Status:  issue.Fields.Summary,
				Comment: comment,
			})
		}
	}

	return res, nil
}

func getMailBody(issues map[string][]Issue) string {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	for key, value := range issues {
		sb.WriteString(fmt.Sprintf("<h3>%s</h3>", key))
		for _, issue := range value {
			sb.WriteString(
				fmt.Sprintf("<p><a href=\"%s\" data-cke-saved-href=\"%s\">%s</a> (%s) — %s %s\r\n\r\n</p>",
					issue.URL, issue.URL, issue.Key, issue.Summary, issue.Status, issue.Comment))
		}
	}
	sb.WriteString("</body></html>")
	return sb.String()
}

func getSubject() string {
	date := time.Now()
	_, week := date.UTC().ISOWeek()
	return fmt.Sprintf("%s, weekly %d", date.Format("02.01.2006"), week)
}
