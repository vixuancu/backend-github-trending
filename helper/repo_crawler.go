package helper

import (
	"backend-github-trending/handle_error"
	"backend-github-trending/model"
	"backend-github-trending/repository"
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/gommon/log"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// CrawlRepo crawls GitHub trending repositories
func CrawlRepo(githubRepo repository.GithubRepo) {
	c := colly.NewCollector()

	// Set realistic headers
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*github.com*",
		Parallelism: 1,
		Delay:       2 * time.Second,
	})

	repos := make([]model.GithubRepo, 0, 30)

	// Parse trending repositories
	c.OnHTML(`article.Box-row`, func(e *colly.HTMLElement) {
		var repo model.GithubRepo

		// Get repository name
		repoName := e.ChildText("h2.h3 a")
		repoName = strings.TrimSpace(strings.ReplaceAll(repoName, "\n", ""))
		repoName = strings.ReplaceAll(repoName, " ", "")
		repo.Name = repoName

		// Get description
		repo.Description = strings.TrimSpace(e.ChildText("p.color-fg-muted"))

		// Get language color
		bgColor := e.ChildAttr(".repo-language-color", "style")
		re := regexp.MustCompile("#[a-zA-Z0-9_]+")
		match := re.FindStringSubmatch(bgColor)
		if len(match) > 0 {
			repo.Color = match[0]
		}

		// Get URL
		repoURL := e.ChildAttr("h2.h3 a", "href")
		if repoURL != "" {
			repo.Url = "https://github.com" + repoURL
		}

		// Get language
		repo.Lang = strings.TrimSpace(e.ChildText("span[itemprop=programmingLanguage]"))

		// Get stars and forks
		e.ForEach("a[href*='/stargazers']", func(index int, el *colly.HTMLElement) {
			repo.Stars = strings.TrimSpace(el.Text)
		})

		e.ForEach("a[href*='/forks']", func(index int, el *colly.HTMLElement) {
			repo.Fork = strings.TrimSpace(el.Text)
		})

		// Get today's stars
		todayStarsText := e.ChildText("span.float-sm-right")
		if strings.Contains(todayStarsText, "stars today") {
			repo.StarsToday = strings.TrimSpace(todayStarsText)
		}

		// Get contributors
		var buildBy []string
		e.ForEach("a[data-hovercard-type='user'] img", func(index int, el *colly.HTMLElement) {
			avatarURL := el.Attr("src")
			if avatarURL != "" {
				buildBy = append(buildBy, avatarURL)
			}
		})
		repo.BuildBy = strings.Join(buildBy, ",")

		if repo.Name != "" {
			repos = append(repos, repo)
		}
	})

	// Process scraped data
	c.OnScraped(func(r *colly.Response) {
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		log.Info(fmt.Sprintf("Scraped %d repositories", len(repos)))

		for _, repo := range repos {
			queue.Submit(&SaveRepoJob{
				repo:      repo,
				repoStore: githubRepo,
			})
		}
	})

	// Visit GitHub trending page
	err := c.Visit("https://github.com/trending")
	if err != nil {
		log.Error("Failed to visit GitHub trending:", err)
	}
}

// SaveRepoJob defines a job to save repository
type SaveRepoJob struct {
	repo      model.GithubRepo
	repoStore repository.GithubRepo
}

// Process implements Job interface
func (j *SaveRepoJob) Process() {
	ctx := context.Background()
	_, err := j.repoStore.SaveRepo(ctx, j.repo)
	if err != nil && err != handle_error.RepoConflict {
		log.Error("Failed to save repository:", err)
	}
}

// SetupCrawler configures and starts the crawler
func SetupCrawler(githubRepo repository.GithubRepo) {
	// Run crawler immediately
	go CrawlRepo(githubRepo)

	// Set up periodic crawling (every 3 hours)
	ticker := time.NewTicker(3 * time.Hour)
	go func() {
		for range ticker.C {
			log.Info("Running scheduled GitHub trending crawler")
			CrawlRepo(githubRepo)
		}
	}()
}
