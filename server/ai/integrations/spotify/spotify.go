package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.spotify.com/v1/"
	defaultAuthURL = "https://accounts.spotify.com/api/token"

	playerPauseURI = "me/player/pause"
	playerPlayURI  = "me/player/play"

	code = ""
)

type Spotify struct {
	credentials Credentials

	baseURL string
	authURL string

	token oAuthResponse
}

func New(creds Credentials, opts ...Option) Spotify {
	s := Spotify{
		credentials: creds,
	}
	s.applyDefaults()
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

func (s *Spotify) applyDefaults() {
	s.baseURL = defaultBaseURL
	s.authURL = defaultAuthURL
}

type Option func(*Spotify)

func WithCustomBaseURL(url string) Option {
	return func(t *Spotify) {
		t.baseURL = url
	}
}

func WithCustomAuthURL(url string) Option {
	return func(t *Spotify) {
		t.authURL = url
	}
}

func (s *Spotify) authorize(ctx context.Context) error {
	limit := time.Now().Add(-10 * time.Minute)
	if s.token.valid(limit) {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code) // TODO: autoryzacja powinna być przemyślana na nowo + potrzeba storage'a -> interakcja na UI powinna być jednorazowa (a jej konieczność wykrywalna tutaj), następnie odświeżanie pownno bazować na refresh tokenie
	data.Set("redirect_uri", "http://localhost:8080/auth")
	encodedData := data.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.authURL, strings.NewReader(encodedData))
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(s.credentials.ClientID, s.credentials.ClientSecret)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authorize with status: %s", response.Status)
	}

	var container oAuthResponse
	if err := json.NewDecoder(response.Body).Decode(&container); err != nil {
		return err
	}

	s.token = container

	s.token.validTo = time.Now().Add(time.Second * time.Duration(s.token.ExpiresIn))
	return nil
}

func (s *Spotify) PausePlayer(ctx context.Context) error {
	if err := s.authorize(ctx); err != nil {
		return err
	}

	address, err := url.JoinPath(s.baseURL, playerPauseURI)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, address, nil)
	if err != nil {
		return err
	}

	bearer := fmt.Sprintf("Bearer %s", s.token.AccessToken)
	request.Header.Set("Authorization", bearer)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logError(response.Body)
		return fmt.Errorf("failed to pause player with status: %s", response.Status)
	}
	return nil
}

func (s *Spotify) PlayPlayer(ctx context.Context) error {
	if err := s.authorize(ctx); err != nil {
		return err
	}

	address, err := url.JoinPath(s.baseURL, playerPlayURI)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, address, nil)
	if err != nil {
		return err
	}

	bearer := fmt.Sprintf("Bearer %s", s.token.AccessToken)
	request.Header.Set("Authorization", bearer)

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logError(response.Body)
		return fmt.Errorf("failed to play player with status: %s", response.Status)
	}
	return nil
}

func logError(body io.Reader) {
	var container = make(map[string]any)
	if err := json.NewDecoder(body).Decode(&container); err != nil {
		panic(err)
	}
	log.Println(container)
}
