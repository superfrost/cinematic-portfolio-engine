package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
	owner  string
	repo   string
}

func NewGitHubClient(ctx context.Context, token, owner, repo string) *GitHubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &GitHubClient{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

type FileInfo struct {
	Name string
	Path string
	SHA  string
}

func (g *GitHubClient) ListFiles(ctx context.Context, path string) ([]FileInfo, error) {
	_, dc, _, err := g.client.Repositories.GetContents(ctx, g.owner, g.repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", path, err)
	}

	var files []FileInfo
	for _, f := range dc {
		if f.GetType() == "file" && strings.HasSuffix(f.GetName(), ".json") {
			files = append(files, FileInfo{
				Name: f.GetName(),
				Path: f.GetPath(),
				SHA:  f.GetSHA(),
			})
		}
	}
	return files, nil
}

type FileContent struct {
	Content string
	SHA     string
}

func (g *GitHubClient) GetFile(ctx context.Context, path string) (*FileContent, error) {
	fc, _, _, err := g.client.Repositories.GetContents(ctx, g.owner, g.repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", path, err)
	}

	content, err := fc.GetContent()
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return &FileContent{Content: content, SHA: fc.GetSHA()}, nil
}

func (g *GitHubClient) CreateFile(ctx context.Context, path, content, msg string) error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(msg),
		Content: []byte(content),
		Branch:  github.Ptr("main"),
	}
	_, _, err := g.client.Repositories.CreateFile(ctx, g.owner, g.repo, path, opts)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	return nil
}

func (g *GitHubClient) UpdateFile(ctx context.Context, path, sha, content, msg string) error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(msg),
		Content: []byte(content),
		SHA:     github.Ptr(sha),
		Branch:  github.Ptr("main"),
	}
	_, _, err := g.client.Repositories.UpdateFile(ctx, g.owner, g.repo, path, opts)
	if err != nil {
		return fmt.Errorf("update %s: %w", path, err)
	}
	return nil
}

func (g *GitHubClient) DeleteFile(ctx context.Context, path, sha, msg string) error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.Ptr(msg),
		SHA:     github.Ptr(sha),
		Branch:  github.Ptr("main"),
	}
	_, _, err := g.client.Repositories.DeleteFile(ctx, g.owner, g.repo, path, opts)
	if err != nil {
		return fmt.Errorf("delete %s: %w", path, err)
	}
	return nil
}

func (g *GitHubClient) RenameFile(ctx context.Context, oldPath, newPath, msg string) error {
	fc, err := g.GetFile(ctx, oldPath)
	if err != nil {
		return fmt.Errorf("rename get old: %w", err)
	}

	if err := g.CreateFile(ctx, newPath, fc.Content, msg); err != nil {
		return fmt.Errorf("rename create new: %w", err)
	}

	if err := g.DeleteFile(ctx, oldPath, fc.SHA, msg); err != nil {
		return err
	}

	return nil
}

type DeployStatus struct {
	Status     string
	Conclusion string
	UpdatedAt  time.Time
	Success    bool
}

func (g *GitHubClient) GetDeployStatus(ctx context.Context) (*DeployStatus, error) {
	runs, _, err := g.client.Actions.ListWorkflowRunsByFileName(ctx, g.owner, g.repo, "deploy.yml", &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return nil, fmt.Errorf("list workflow runs: %w", err)
	}

	if len(runs.WorkflowRuns) == 0 {
		return nil, fmt.Errorf("no deploy runs found")
	}

	run := runs.WorkflowRuns[0]
	ds := &DeployStatus{
		Status:     run.GetStatus(),
		Conclusion: run.GetConclusion(),
		UpdatedAt:  run.GetUpdatedAt().Time,
		Success:    run.GetConclusion() == "success",
	}
	return ds, nil
}
