package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	Order          int    `json:"order"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
	ParentID       string `json:"parent_id"`
}

func (c *Client) GetProject(ctx context.Context, projectId string) (*Project, error) {
	// get all projects to check if the project exists or not
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	allProjects := []Project{}
	if _, _, err := c.sendRequest(req, &allProjects); err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"Projects": fmt.Sprintf("%+v", allProjects),
	}).Info("Projects read")

	// Check if project exists in list of active projects
	exists := false
	for _, project := range allProjects {
		if project.ID == projectId {
			exists = true
			break
		}
	}
	res := Project{}
	if !exists {
		return nil, fmt.Errorf("unable to find projects with id: %s", projectId)
	}
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Reading project")
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	_, _, err = c.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project read")
	return &res, nil
}

type CreateProject struct {
	// Required fields
	Name *string `json:"name"`
	// Optional fields
	ParentID   *string `json:"parent_id,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) CreateProject(ctx context.Context, createProject CreateProject) (*Project, error) {
	payload, err := json.Marshal(createProject)
	log.WithFields(log.Fields{
		"payload": string(payload),
	}).Info("Creating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", c.BaseURL), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	_, _, err = c.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project created")
	return &res, nil
}

type UpdateProject struct {
	ID         *string
	Name       *string `json:"name,omitempty"`
	Color      *string `json:"color,omitempty"`
	IsFavorite *bool   `json:"is_favorite,omitempty"`
	ViewStyle  *string `json:"view_style,omitempty"`
}

func (c *Client) UpdateProject(ctx context.Context, updateProject UpdateProject) (*Project, error) {
	payload, err := json.Marshal(updateProject)
	if updateProject.ID == nil {
		return nil, fmt.Errorf("missing project id")
	}
	log.WithFields(log.Fields{
		"payload":   string(payload),
		"projectId": updateProject.ID,
	}).Info("Updating project")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s", c.BaseURL, *updateProject.ID), bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	req = req.WithContext(ctx)
	res := Project{}
	_, _, err = c.sendRequest(req, &res)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Project": fmt.Sprintf("%+v", res),
	}).Debug("Project updated")
	return &res, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectId string) (statusCode int, body string, err error) {
	log.WithFields(log.Fields{
		"projectId": projectId,
	}).Info("Deleting project")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", c.BaseURL, projectId), nil)
	if err != nil {
		return 0, "", err
	}
	req = req.WithContext(ctx)
	statusCode, body, err = c.sendRequest(req, nil)
	if err != nil {
		return -1, "", err
	}
	return statusCode, body, err
}
