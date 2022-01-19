package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const BaseURL = "https://api.todoist.com/rest/v1/"
const TasksURL = "tasks?filter=%23Family%20%26%20today"
const ClosePrefixURL = "tasks/"
const CloseSuffixURL = "/close"

/**
{
    "id": 5515474815,
    "assigner": 0,
    "project_id": 2274460838,
    "section_id": 0,
    "order": 8,
    "content": "Test Recurring",
    "description": "",
    "completed": false,
    "label_ids": [],
    "priority": 1,
    "comment_count": 0,
    "creator": 226996,
    "created": "2022-01-19T14:38:38Z",
    "due": {
      "recurring": true,
      "string": "every day",
      "date": "2022-01-19"
    },
    "url": "https://todoist.com/showTask?id=5515474815&sync_id=5515474815"
*/

type Task struct {
	ID          int    `json:"id"`
	Assigner    int    `json:"assigner"`
	ProjectID   int    `json:"project_id"`
	SectionID   int    `json:"section_id"`
	Order       int    `json:"order"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type tasks []Task

type TodoistClient struct {
	APIToken string
	projects map[int]string
}

func NewTodoistClient(api_token string) *TodoistClient {
	tc := &TodoistClient{}
	tc.APIToken = api_token
	tc.projects = make(map[int]string)
	return tc
}

func (tc *TodoistClient) getTasks() (tasks, error) {
	data, err := tc.get(BaseURL + TasksURL)
	if err != nil {
		fmt.Printf("Error %s", err)
		return nil, err
	}

	allTasks := make(tasks, 0)
	uerr := json.Unmarshal(data, &allTasks)
	if uerr != nil {
		fmt.Printf("Error %s", uerr)
		return nil, uerr
	}
	return allTasks, err
}

func (tc *TodoistClient) closeTask(task Task) error {
	_, err := tc.post(fmt.Sprint(BaseURL, ClosePrefixURL, task.ID, CloseSuffixURL), nil)
	return err
}

func (tc *TodoistClient) get(url string) ([]byte, error) {
	req, rerr := http.NewRequest("GET", url, nil)
	if rerr != nil {
		fmt.Printf("Error %s", rerr)
		return nil, rerr
	}
	return tc.do(req)
}

func (tc *TodoistClient) post(url string, body io.Reader) ([]byte, error) {
	req, rerr := http.NewRequest("POST", url, body)
	fmt.Println(url)
	if rerr != nil {
		fmt.Printf("Error %s", rerr)
		return nil, rerr
	}
	return tc.do(req)
}

func (tc *TodoistClient) do(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tc.APIToken)

	c := http.Client{Timeout: time.Duration(15) * time.Second}
	resp, err := c.Do(req)

	if err != nil {
		fmt.Printf("Error %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
