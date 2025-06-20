package database

import "fmt"

func (f Feed) String() string {
	return fmt.Sprintf("database.Feed { ID: %v, Name: %v, Url: %v, UserID: %v, CreatedAt: %v, UpdatedAt: %v }",
		f.ID, f.Name, f.Url, f.UserID, f.CreatedAt, f.UpdatedAt)
}

func (u User) String() string {
	return fmt.Sprintf("database.User { ID: %v, Name: %v, CreatedAt: %v, UpdatedAt: %v }",
		u.ID, u.Name, u.CreatedAt, u.UpdatedAt)
}

func (ff FeedFollow) String() string {
	return fmt.Sprintf("database.FeedFollow { ID: %v, FeedID: %v, UserID: %v, CreatedAt: %v, UpdatedAt: %v }",
		ff.ID, ff.FeedID, ff.UserID, ff.CreatedAt, ff.UpdatedAt)
}

func (ff AddFeedFollowRow) String() string {
	return fmt.Sprintf("database.FeedFollow { ID: %v, FeedID: %v, FeedName: %v, UserID: %v, UserName: %v, CreatedAt: %v, UpdatedAt: %v }",
		ff.ID, ff.FeedID, ff.FeedName, ff.UserID, ff.UserName, ff.CreatedAt, ff.UpdatedAt)
}

func (ff GetFeedFollowsForUserRow) String() string {
	return fmt.Sprintf("database.FeedFollow { ID: %v, FeedID: %v, FeedName: %v, UserID: %v, UserName: %v, CreatedAt: %v, UpdatedAt: %v }",
		ff.ID, ff.FeedID, ff.FeedName, ff.UserID, ff.UserName, ff.CreatedAt, ff.UpdatedAt)
}
