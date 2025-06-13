package database

import "fmt"

func (f Feed) String() string {
	return fmt.Sprintf("database.Feed { ID:%v, Name:%v, Url:%v, UserID:%v, CreatedAt:%v, UpdatedAt:%v }",
		f.ID, f.Name, f.Url, f.UserID, f.CreatedAt, f.UpdatedAt)
}

func (u User) String() string {
	return fmt.Sprintf("database.User { ID:%v, Name:%v, CreatedAt:%v, UpdatedAt:%v }",
		u.ID, u.Name, u.CreatedAt, u.UpdatedAt)
}
