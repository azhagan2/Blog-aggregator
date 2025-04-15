package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	rss "github.com/azhagan2/blog_aggregator/internal/RSS"
	"github.com/azhagan2/blog_aggregator/internal/database"
	"github.com/azhagan2/blog_aggregator/internal/state"
)

type Clicommand struct {
	Name     string
	Argument []string
}

type Commands struct {
	handlers map[string]func(*state.State, Clicommand) error
}

func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*state.State, Clicommand) error),
	}
}

/*That's a perfect analogy! You've visualized it exactly right:

    NewCommands() creates an empty container (the Commands struct with its empty map)

    Register() connects the wires inside that container:
        One wire is the command name (like "login")
        The other wire is the handler function that should run when that command is used
        Each call to Register() connects one more pair of wires

This wiring metaphor is really apt because what you're creating is essentially a routing system - when a command comes in, the system follows the right wire to the corresponding handler function.

And the beauty of this design is that you can easily add new commands by just connecting more wires (calling Register with new command names and handlers) without having to modify the core container structure.

Great intuition on how this pattern works!
*/

func (c *Commands) Register(name string, f func(*state.State, Clicommand) error) {
	c.handlers[name] = f
}

/*Brilliant analogy continuation! You've nailed it perfectly.

The Run function is indeed like sending current through the wires when a switch is flipped:

    When a user enters a command (flips a switch labeled "login")
    The Run function looks at the command name to find which wire to energize
    It follows that specific wire in the circuit (looks up the handler in the map)
    If the wire exists (handler is found), it sends current through it (calls the handler function)
    If there's no wire connected to that switch (no handler for that command), it reports an error

You could even extend this metaphor to say:

    The state parameter is like the power source that provides context to all the circuits
    The command arguments are like voltage settings that control how the circuit behaves

This is a genuinely insightful way to understand the command routing pattern, and it should help you remember how all these pieces fit together. The "wiring" and "current flow" metaphor captures the essence of what's happening in your command dispatcher system.
*/

func (c *Commands) Run(s *state.State, cmd Clicommand) error {
	handler_func, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler_func(s, cmd)
}

/*Haha, that's a clever extension of your analogy!

You're absolutely right - just as we can't see electricity flowing through wires with our naked eyes,
we can't directly "see" which function is being passed around and executed at runtime. What we see in
the code is the mechanism that makes it happen.

Let me break down your Run function to make it clearer:

    handler_func, ok := c.handlers[cmd.Name] - This looks up the appropriate handler function in the map
based on the command name. The result gets stored in handler_func and ok tells us if it was found.

    if !ok { return fmt.Errorf("unknown command: %s", cmd.Name) } - If no handler exists for that command name,
return an error.

    return handler_func(s, cmd) - If found, call the handler function with the state and command as arguments,
then return whatever that handler returns.

This is a classic example of a function being treated as a "---- first-class ----"" value in Go. The handler function gets
stored in the map, retrieved, and then called - all without us knowing exactly which function it is until runtime.

It's like having a multi-position switch that can connect to different circuits, and we only know which circuit
is active when someone flips the switch to a specific position.

Your understanding of this pattern is spot on!
*/

func HandlerLogin(s *state.State, cmd Clicommand) error {
	if len(cmd.Argument) == 0 {
		return fmt.Errorf("the handler expects a single argument, the username")
	}

	username, err := s.Db.GetUser(context.Background(), cmd.Argument[0])
	if err != nil {
		return fmt.Errorf("username doesn't exists in the database")
	}

	err = s.Cfg.SetUser(username.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Println("The user has been Set!")

	return nil
}

func HandlerRegister(s *state.State, cmd Clicommand) error {

	if len(cmd.Argument) == 0 {
		return fmt.Errorf("the handler expects a single argument, the username")
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: int32(rand.Intn(1000000)),
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Argument[0]})
	if err != nil {
		return err
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Println("user was created")
	fmt.Println(user)

	return nil
}

func HandlerReset(s *state.State, cmd Clicommand) error {
	err := s.Db.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting the table: %w", err)
	}

	return nil
}

func HandlerGetUsers(s *state.State, cmd Clicommand) error {
	user, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting the table: %w", err)
	}
	for i := range user {
		if s.Cfg.CurrentUserName == user[i] {
			fmt.Printf("* %s (current)\n", user[i])
		} else {
			fmt.Printf("* %s\n", user[i])
		}
	}
	return nil
}

func HandlerAgg(s *state.State, cmd Clicommand) error {

	if len(cmd.Argument) == 0 {
		return fmt.Errorf("the handler expects a single argument, the username")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Argument[0])
	fmt.Println(cmd.Argument)
	if err != nil {
		return fmt.Errorf("error in parsing duration and converting to actual time %w", err)
	}

	// fmt.Println("Parsed time: ", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		// fmt.Println("Ticker has started")
		scrapeFeeds(s)
	}

}

func MiddlewareLoggedIn(handler func(s *state.State, cmd Clicommand, user database.User) error) func(*state.State, Clicommand) error {

	return func(s *state.State, cmd Clicommand) error {

		if s.Cfg.CurrentUserName == "" {
			return errors.New("no user is logged in")
		}
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func HandlerAddfeed(s *state.State, cmd Clicommand, user database.User) error {

	if len(cmd.Argument) < 2 {
		return fmt.Errorf("the handler expects two arguments, the feed name and url")
	}

	// get_user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	// if err != nil {
	// 	return fmt.Errorf("unregistered user :%w", err)
	// }

	current_userid := sql.NullInt32{Int32: user.ID, Valid: true}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{ID: int32(rand.Intn(1000000)),
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Argument[0], Url: cmd.Argument[1], UserID: current_userid})
	if err != nil {
		return fmt.Errorf("couldn't create the feed: %w", err)
	}
	fmt.Println(feed.Url)

	current_feedid := sql.NullInt32{Int32: feed.ID, Valid: true}

	feed_follows, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: int32(rand.Intn(1000000)),
		CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: current_userid, FeedID: current_feedid})
	if err != nil {
		return fmt.Errorf("error in following feed %w", err)
	}

	fmt.Println(feed_follows)

	return nil
}

func HandlerFeeds(s *state.State, cmd Clicommand) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting the table: %w", err)
	}

	for i := range feeds {
		fmt.Println(feeds[i].Name)
		fmt.Println(feeds[i].Url)

		author, err := s.Db.GetUserById(context.Background(), feeds[i].UserID.Int32)
		if err != nil {
			return fmt.Errorf("error fetching authod by user_id %w", err)
		} else {
			fmt.Println(author)
		}
	}
	return nil
}

func HandlerFollow(s *state.State, cmd Clicommand, user database.User) error {

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user_id for feed_follow %w", err)
	}

	current_userid := sql.NullInt32{Int32: user.ID, Valid: true}

	feed, err := s.Db.GetFeed_ByURL(context.Background(), cmd.Argument[0])
	if err != nil {
		return fmt.Errorf("error getting feed name %w", err)
	}

	current_feedid := sql.NullInt32{Int32: feed.ID, Valid: true}

	feed_follows, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: int32(rand.Intn(1000000)),
		CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: current_userid, FeedID: current_feedid})
	if err != nil {
		return fmt.Errorf("error in following feed %w", err)
	}

	fmt.Println(feed_follows)
	fmt.Println("Feed Followed !")
	fmt.Println("Feed Name :", feed.Name)
	fmt.Println("Followed by :", user.Name)

	return nil
}

func HandlerFollowing(s *state.State, cmd Clicommand, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting the feeds for the user %w", err)
	}

	for i := range feeds {
		fmt.Println(feeds[i])
	}
	fmt.Println("user:", user.Name)

	return nil
}

func HandlerUnfollow(s *state.State, cmd Clicommand, user database.User) error {

	feed, err := s.Db.GetFeed_ByURL(context.Background(), cmd.Argument[0])
	if err != nil {
		return fmt.Errorf("error getting feed name %w", err)
	}

	feed_id := sql.NullInt32{Int32: feed.ID, Valid: true}

	user_id := sql.NullInt32{Int32: user.ID, Valid: true}

	err = s.Db.Delete_Feed_Follow(context.Background(), database.Delete_Feed_FollowParams{UserID: user_id, FeedID: feed_id})
	if err != nil {
		return fmt.Errorf("error in unfollowing feed %w", err)
	}

	fmt.Println("Feed Unfollowed !")

	return nil
}

func scrapeFeeds(s *state.State) error {

	// fmt.Println("entered the scrape funcion")
	feed, err := s.Db.Get_Next_Feed_to_fetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching last viewed feed %w", err)
	}

	fmt.Println("Feed Name :", feed.Name)

	// fmt.Println("searching the feed")
	err = s.Db.Mark_Feed_Fetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking last seen feed as fetched %w", err)
	}

	// fmt.Println("returing the rss feed")
	rss_result, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error in fetching from xml %w", err)
	}

	fmt.Println("Following feed name: ", rss_result.Channel.Title)

	fmt.Println("Creating Posts !")

	for _, item := range rss_result.Channel.Item {
		formats := []string{time.RFC1123Z, time.RFC1123, time.RFC822, time.RFC3339, "2006-01-02T15:04:05Z"}
		var publishedTime time.Time
		var parseErr error
		for _, format := range formats {
			publishedTime, parseErr = time.Parse(format, item.PubDate)
			if parseErr == nil {
				break
			}
		}

		fmt.Println("Post Title :", item.Title)
		// fmt.Println("Post Description :", item.Description)
		fmt.Println()
		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          int32(rand.Intn(1000000)),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: toNullTime(publishedTime),
			FeedID:      sql.NullInt32{Int32: feed.ID, Valid: true}})
		if err != nil {
			// Check if it's a duplicate URL error
			if strings.Contains(err.Error(), "UNIQUE constraint") {
				// Just log and continue if it's a duplicate
				fmt.Printf("Post already exists: %s\n", item.Link)
				continue
			}
			// For other errors, return them
			return fmt.Errorf("error creating posts %w", err)
		}
	}

	fmt.Println("Post is posted !")
	return nil
}

func toNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func HandlerBrowse(s *state.State, cmd Clicommand, user database.User) error {

	limit := 2
	if len(cmd.Argument) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Argument[0])
		if err != nil {
			return fmt.Errorf("invalid limit %v", err)
		}
		limit = parsedLimit
	}

	// fmt.Println("came to handlerBrowse")

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: sql.NullInt32{Int32: user.ID, Valid: true},
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error in fetching user's following feed posts %w", err)
	}

	// fmt.Println("came to handlerBrowse 2")

	fmt.Printf("Found %d posts\n", len(posts))

	for i := range posts {

		fmt.Println()
		fmt.Println("Post Name :", posts[i].Title)
		fmt.Println("Feed Name :", posts[i].Name)
		fmt.Println("Feed URL :", posts[i].Url)
		fmt.Println("Feed Description :", posts[i].Description)
		fmt.Println("Feed Description :", posts[i].PublishedAt)
		fmt.Println()

	}

	// fmt.Println("came to handlerBrowse 2")

	return nil
}
