package port

import "context"

type CacheRepository interface {
	SaveUserSelectedCategory(ctx context.Context, username string, category string) error
	SaveLoggedInUserToCache(ctx context.Context, username string) ([]string, error)
	GetLoggedInUsersFromCache(ctx context.Context) ([]string, error)
	RemoveLoggedInUserFromCache(ctx context.Context, username string) ([]string, error)
	DeleteLoggedInUsers(ctx context.Context) error
	Close()
}
