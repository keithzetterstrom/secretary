package docs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/keithzetterstrom/secretary/internal/repository/models"
	"github.com/keithzetterstrom/secretary/utils/logger"
)

type Config struct {
	SpreadsheetId        string `yaml:"spreadsheet_id"`
	TokenFileName        string `yaml:"token_file_name"`
	ClientSecretFileName string `yaml:"client_secret_file_name"`
}

type Client struct {
	srv           *sheets.Service
	spreadsheetId string
	log           logger.Logger
}

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) (*http.Client, error) {
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok), nil
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.Wrap(err, "failed to read authorization code")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve token from web")
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return errors.Wrap(err, "failed to cache OAuth token")
	}
	json.NewEncoder(f).Encode(token)
	return nil
}

func New(cfg Config, log logger.Logger) (*Client, error) {
	ctx := context.Background()
	b, err := os.ReadFile(cfg.ClientSecretFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read client secret file")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse client secret file to config")
	}

	client, err := getClient(config, cfg.TokenFileName)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve Docs client")
	}

	return &Client{
		srv:           srv,
		spreadsheetId: cfg.SpreadsheetId,
		log:           log,
	}, nil
}

func (c *Client) Get(ctx context.Context) ([][]interface{}, error) {
	resp, err := c.srv.Spreadsheets.Values.Get(c.spreadsheetId, "A1:C").Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve data from sheet")
	}

	return resp.Values, nil
}

func (c *Client) Set(value models.UserRegistration) error {
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data: []*sheets.ValueRange{
			{
				Range:  "A1:C1",
				Values: [][]interface{}{{value.Email, value.Name, value.TgUserName}},
			},
		},
	}

	_, err := c.srv.Spreadsheets.Values.BatchUpdate(c.spreadsheetId, rb).Do()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve data from sheet")
	}

	return nil
}

func (c *Client) Append(ctx context.Context, value models.UserRegistration) error {
	range2 := "A1"
	valueInputOption := "RAW"
	insertDataOption := "INSERT_ROWS"

	vr := &sheets.ValueRange{
		Values: [][]interface{}{{value.Email, value.Name, value.TgUserName}},
	}

	_, err := c.srv.Spreadsheets.Values.
		Append(c.spreadsheetId, range2, vr).
		ValueInputOption(valueInputOption).
		InsertDataOption(insertDataOption).
		Context(ctx).IncludeValuesInResponse(true).Do()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve data from sheet")
	}

	return nil
}

func (c *Client) Update(ctx context.Context, value models.UserRegistration) error {
	range2 := "A1:C1"
	valueInputOption := "RAW"

	vr := &sheets.ValueRange{
		Values: [][]interface{}{{value.Email, value.Name, value.TgUserName}},
	}

	_, err := c.srv.Spreadsheets.Values.
		Update(c.spreadsheetId, range2, vr).
		ValueInputOption(valueInputOption).
		Context(ctx).Do()
	if err != nil {
		return errors.Wrap(err, "failed to retrieve data from sheet")
	}

	return nil
}

func (c *Client) UpdateValue(ctx context.Context, rng string, value interface{}) error {
	spreadsheetId := "16PmkwFTd42Cx0LaFTQUKg8pQBziKM-SLCOyUnlL3JvU"

	valueInputOption := "RAW"

	vr := &sheets.ValueRange{
		Values: [][]interface{}{{value}},
	}

	_, err := c.srv.Spreadsheets.Values.
		Update(spreadsheetId, rng, vr).
		ValueInputOption(valueInputOption).
		Context(ctx).Do()
	if err != nil {
		return errors.Wrap(err, "failed to update value")
	}

	return nil
}
