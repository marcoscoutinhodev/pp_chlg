package keycloak

import (
	"context"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type IdentityManager struct {
	baseUrl         string
	realm           string
	clientID        string
	clientSecret    string
	shopkeeperGroup string
	customerGroup   string
}

func NewIdentityManager() *IdentityManager {
	return &IdentityManager{
		baseUrl:         os.Getenv("KC_BASE_URL"),
		realm:           os.Getenv("KC_REALM"),
		clientID:        os.Getenv("KC_CLIENT_ID"),
		clientSecret:    os.Getenv("KC_CLIENT_SECRET"),
		shopkeeperGroup: os.Getenv("KC_SHOPKEEPER_GROUP"),
		customerGroup:   os.Getenv("KC_CUSTOMER_GROUP"),
	}
}

func (im IdentityManager) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)
	jwt, err := client.LoginClient(ctx, im.clientID, im.clientSecret, im.realm)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (im IdentityManager) GetGroupID(ctx context.Context, group string) (groupID string, err error) {
	clientToken, err := im.loginClient(ctx)
	if err != nil {
		return
	}

	client := gocloak.NewClient(im.baseUrl)
	groups, err := client.GetGroups(ctx, clientToken.AccessToken, im.realm, gocloak.GetGroupsParams{})
	if err != nil {
		return
	}

	for _, g := range groups {
		if *g.Name == group {
			groupID = *g.ID
			break
		}
	}

	return
}

func (im IdentityManager) CreateUser(ctx context.Context, user entity.User, groupID string) (*gocloak.User, error) {
	clientToken, err := im.loginClient(ctx)
	if err != nil {
		return nil, err
	}

	kcUser := &gocloak.User{
		Enabled:       gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		Email:         &user.Email,
		Username:      &user.Email,
		FirstName:     &user.FirstName,
		LastName:      &user.LastName,
		Attributes: &map[string][]string{
			"TaxpayeerIdentification": {user.TaxpayeerIdentification},
		},
	}

	client := gocloak.NewClient(im.baseUrl)
	userID, err := client.CreateUser(ctx, clientToken.AccessToken, im.realm, *kcUser)
	if err != nil {
		return nil, err
	}

	if err := client.AddUserToGroup(ctx, clientToken.AccessToken, im.realm, userID, groupID); err != nil {
		return nil, err
	}

	kcUser, err = client.GetUserByID(ctx, clientToken.AccessToken, im.realm, userID)
	if err != nil {
		return nil, err
	}

	return kcUser, err
}

func (im IdentityManager) AuthenticateUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)
	userJWT, err := client.Login(ctx, im.clientID, im.clientSecret, im.realm, username, password)
	if err != nil {
		return nil, err
	}

	return userJWT, nil
}
