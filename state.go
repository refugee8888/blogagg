package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/refugee8888/blogagg/internal/config"
	"github.com/refugee8888/blogagg/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	NameToHandler map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("username is required")
	}
	getting, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("unknown username")
		os.Exit(1)
	}
	if err := s.cfg.SetUser(getting.Name); err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("username is required")
	}
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{uuid.New(), time.Now(), time.Now(), cmd.Arguments[0]})
	if err != nil {
		os.Exit(1)
		return err
	}

	if err := s.cfg.SetUser(newUser.Name); err != nil {
		return err
	}
	fmt.Println("user created")
	fmt.Println(newUser)

	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := s.cfg.SetUser(""); err != nil {
		return err
	}
	fmt.Println("table has been reset")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	getter, err := s.db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}

	for _, v := range getter {
		if s.cfg.CurrentUserName == v.Name {
			fmt.Println("* " + v.Name + " (current)")
		} else {
			fmt.Println("* " + v.Name)
		}
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}
	if _, err := s.db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed %s fetched: %w", nextFeed.Name, err)
	}

	fetchedFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %s: %w", nextFeed.Name, err)
	}

	for _, item := range fetchedFeed.Channel.Item {
		publishedAt, ok := parseDate(item.PubDate)
		if !ok {
			log.Printf("couldn't parse date %q, using now", item.PubDate)
			publishedAt = time.Now()
		}

		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      uuid.NullUUID{UUID: nextFeed.ID, Valid: true},
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue // already exists (ON CONFLICT DO NOTHING)
			}
			log.Printf("couldn't create post %q: %v", item.Link, err)
			continue
		}
	}

	log.Printf("feed %s collected, %d items", nextFeed.Name, len(fetchedFeed.Channel.Item))
	return nil
}

func handlerFeed(s *state, cmd command) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("time between requests is required")
	}
	time_between_reqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()
	for ; ; <-ticker.C {

		fmt.Printf("Collecting feeds every %v", time_between_reqs)
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("%w", err)
		}

	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) == 1 {
		n, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("invalid limit %q: %w", cmd.Arguments[0], err)
		}
		limit = n
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	for _, p := range posts {
		fmt.Printf("%s from %s\n", p.PublishedAt, p.FeedName)
		fmt.Printf("--- %s ---\n", p.Title)
		fmt.Printf("    %v\n", p.Description.String)
		fmt.Printf("Link: %s\n", p.Url)
		fmt.Println("=====================================")
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if ok := len(cmd.Arguments) < 2; ok {
		return fmt.Errorf("both feed name and url are required")
	}

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{uuid.New(), time.Now(), time.Now(), cmd.Arguments[0], cmd.Arguments[1], user.ID})
	if err != nil {
		return err
	}

	if _, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		uuid.New(), time.Now(),
		time.Now(),
		user.ID,
		newFeed.ID,
	}); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(newFeed.Name, newFeed.ID, newFeed.Url)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	getter, err := s.db.GetFeeds(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}

	for _, v := range getter {
		fmt.Println(v.Name, v.Url, v.UserName)
	}

	return nil
}

func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("feed url is required")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}
	var curFeed uuid.UUID
	for _, v := range feeds {
		if v.Url == cmd.Arguments[0] {
			curFeed = v.ID
		}
	}
	fmt.Println(user.ID, curFeed)
	getter, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		uuid.New(), time.Now(),
		time.Now(),
		user.ID,
		curFeed,
	})
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(getter.FeedName, getter.UserName)

	return nil
}

func handlerUserFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, v := range feeds {
		fmt.Println(v.Follows, v.Username, v.Feedname)
	}
	return nil
}

func handlerFeedUnfollow(s *state, cmd command, user database.User) error {
	if ok := len(cmd.Arguments) == 0; ok {
		return fmt.Errorf("feed url is required")
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{cmd.Arguments[0], user.ID})
	if err != nil {
		return fmt.Errorf("%w", err)
	} else {
		fmt.Println("feed unfollowed")
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	v, ok := c.NameToHandler[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return v(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if name == "" {
		return fmt.Errorf("command name can't be empty")
	}
	c.NameToHandler[name] = f
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func parseDate(s string) (time.Time, bool) {
	layouts := []string{time.RFC1123Z, time.RFC1123, time.RFC3339, "2006-01-02"}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
