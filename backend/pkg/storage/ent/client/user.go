package client

import (
	"context"

	"github.com/redhat-appstudio/quality-studio/api/apis/user/v1alpha1"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/predicate"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db/users"
)

func (d *Database) CreateUser(u v1alpha1.User) error {
	userAlreadyExists := d.client.Users.Query().
		Where(users.UserEmail(u.UserEmail)).
		ExistX(context.TODO())
	if userAlreadyExists {
		_, err := d.client.Users.Update().
			Where(predicate.Users(users.UserEmail(u.UserEmail))).
			SetConfig(u.Config).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to update user: %w", err)
		}
	} else {
		_, err := d.client.Users.Create().
			SetUserName(u.UserName).
			SetUserEmail(u.UserEmail).
			SetConfig(u.Config).
			Save(context.TODO())
		if err != nil {
			return convertDBError("failed to create user: %w", err)
		}
	}

	return nil
}

func (d *Database) GetUser(userEmail string) (*db.Users, error) {
	user, err := d.client.Users.Query().
		Where(users.UserEmail(userEmail)).Only(context.TODO())

	if err != nil {
		return nil, convertDBError("get user: %w", err)
	}

	return user, nil
}

func (d *Database) ListAllUsers() ([]*db.Users, error) {
	users, err := d.client.Users.Query().All(context.Background())

	if err != nil {
		return nil, convertDBError("failed to return users status: %w", err)
	}

	return users, nil
}
