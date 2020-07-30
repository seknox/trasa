package idps

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// DOTokenSource required for digital ocean api
type DOTokenSource struct {
	AccessToken string
}

// Token required for digital ocean api
func (t *DOTokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// GetDoDroplets fetches droplet details from Digital ocean api
func GetDoDroplets(token []byte) ([]godo.Droplet, error) {

	drops := make([]godo.Droplet, 0)
	tokenSource := &DOTokenSource{
		AccessToken: string(token),
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	drops, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		return drops, err
	}

	return drops, nil
}
