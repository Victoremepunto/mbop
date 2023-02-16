package ocm

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/google/uuid"
	"github.com/redhatinsights/mbop/internal/models"
)

type SDKMock struct{}

func (ocm *SDKMock) InitSdkConnection(ctx context.Context) error {
	return nil
}

func (ocm *SDKMock) GetUsers(u models.UserBody, q models.UserQuery) (models.Users, error) {
	var users models.Users

	if u.Users == nil {
		return users, nil
	}

	if u.Users[0] == "errorTest" {
		return users, fmt.Errorf("internal AMS Error")
	}

	for _, user := range u.Users {
		orgID, err := rand.Int(rand.Reader, big.NewInt(999999-100000))
		if err != nil {
			return users, err
		}

		displayNameNum, err := rand.Int(rand.Reader, big.NewInt(99-0))
		if err != nil {
			return users, err
		}

		users.AddUser(models.User{
			Username:      user,
			ID:            uuid.New().String(),
			Email:         "lub@dub.com",
			FirstName:     "test",
			LastName:      "case",
			AddressString: "https://usersTest.com",
			IsActive:      true,
			IsInternal:    true,
			Locale:        "en_US",
			OrgID:         strconv.Itoa(int(orgID.Int64())),
			DisplayName:   "FedRAMP" + strconv.Itoa(int(displayNameNum.Int64())),
			Type:          "User",
		})
	}

	return users, nil
}

func (ocm *SDKMock) GetOrgAdmin(users []models.User) (models.OrgAdminResponse, error) {
	response := models.OrgAdminResponse{}

	if users[0].ID == "23456" {
		return response, nil
	}

	if users[0].ID == "errorTest" {
		return response, fmt.Errorf("error retrieving Role Bindings")
	}

	for _, user := range users {
		response[user.ID] = models.OrgAdmin{
			ID:         user.ID,
			IsOrgAdmin: true,
		}
	}

	return response, nil
}

func (ocm *SDKMock) CloseSdkConnection() {
	// nil
}
