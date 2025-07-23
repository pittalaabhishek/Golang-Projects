package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	twitterAPIBaseURL = "https://api.twitter.com/1.1"
	retweetersFile    = "retweeters.txt"
)

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
}

type BearerTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Tweet struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

var (
	tweetIDStr     string
	pickWinnerFlag bool
	numWinners     int
)

func init() {
	flag.StringVar(&tweetIDStr, "tweetID", "", "The ID of the tweet to check for retweets.")
	flag.BoolVar(&pickWinnerFlag, "pickWinner", false, "Set to true to pick a winner from the list of retweeters.")
	flag.IntVar(&numWinners, "numWinners", 1, "Number of winners to pick (only used with -pickWinner).")
	flag.Parse()

	if tweetIDStr == "" && !pickWinnerFlag {
		fmt.Println("Usage: go run main.go -tweetID <ID> [-pickWinner] [-numWinners <N>]")
		fmt.Println("       Provide a tweet ID to fetch retweeters, or use -pickWinner to select from existing list.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if pickWinnerFlag && numWinners < 1 {
		log.Fatal("Error: -numWinners must be at least 1 when -pickWinner is set.")
	}

	rand.Seed(time.Now().UnixNano())
}

func main() {
	config := Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	if config.ConsumerKey == "" || config.ConsumerSecret == "" {
		log.Println("TWITTER_CONSUMER_KEY and TWITTER_CONSUMER_SECRET environment variables must be set.")
		log.Println("Skipping API calls due to missing credentials.")
	}

	var accessToken string
	if config.ConsumerKey != "" && config.ConsumerSecret != "" {
		var err error
		accessToken, err = authenticate(config)
		if err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}
		log.Println("Successfully authenticated with Twitter API.")
	}

	if tweetIDStr != "" && accessToken != "" {
		log.Printf("Fetching retweeters for tweet ID: %s", tweetIDStr)
		newRetweeters, err := getRetweeters(accessToken, tweetIDStr)
		if err != nil {
			log.Fatalf("Failed to get retweeters: %v", err)
		}

		if len(newRetweeters) > 0 {
			log.Printf("Found %d new retweeters.", len(newRetweeters))
			err = updateRetweetersFile(newRetweeters)
			if err != nil {
				log.Fatalf("Failed to update retweeters file: %v", err)
			}
			log.Println("Retweeters list updated successfully.")
		} else {
			log.Println("No new retweeters found or API returned an empty list.")
		}
	} else if tweetIDStr != "" && accessToken == "" {
		log.Println("Cannot fetch retweeters: Twitter API credentials are not set.")
	}


	if pickWinnerFlag {
		log.Println("Picking winner(s)...")
		allRetweeters, err := readRetweetersFromFile()
		if err != nil {
			log.Fatalf("Failed to read retweeters for winner picking: %v", err)
		}

		if len(allRetweeters) == 0 {
			log.Fatal("No retweeters found in the file to pick a winner from.")
		}

		winners := pickWinners(allRetweeters, numWinners)
		fmt.Println("\n--- WINNER(S) ---")
		for i, winner := range winners {
			fmt.Printf("%d. @%s\n", i+1, winner)
		}
		fmt.Println("-------------------")
	}

	if tweetIDStr == "" && !pickWinnerFlag {
		fmt.Println("Please provide a -tweetID or use -pickWinner.")
		flag.PrintDefaults()
	}
	
	if tweetIDStr != "" && !pickWinnerFlag && (accessToken == "" || len(accessToken) == 0){
		log.Println("To fetch retweeters, please set TWITTER_CONSUMER_KEY and TWITTER_CONSUMER_SECRET environment variables.")
	}
}

func authenticate(config Config) (string, error) {
	authString := url.QueryEscape(config.ConsumerKey) + ":" + url.QueryEscape(config.ConsumerSecret)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token",
		strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Add("Authorization", "Basic "+encodedAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", "29") // Length of "grant_type=client_credentials"

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResponse BearerTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	if tokenResponse.TokenType != "bearer" || tokenResponse.AccessToken == "" {
		return "", fmt.Errorf("invalid token type or missing access token in response")
	}

	return tokenResponse.AccessToken, nil
}

func getRetweeters(accessToken, tweetID string) ([]string, error) {
	log.Println("NOTE: Using simulated retweeter data due to API access restrictions.")
	simulatedRetweeters := []string{
			"userA", "userB", "userC", "userD", "userE",
			"userF", "userG", "userH", "userI", "userJ",
			"userA",
			"userK", "userL", "userM", "userN", "userO",
	}
	uniqueUsernames := make(map[string]struct{})
	for _, name := range simulatedRetweeters {
			uniqueUsernames[strings.ToLower(name)] = struct{}{}
	}

	var retweeters []string
	for username := range uniqueUsernames {
			retweeters = append(retweeters, username)
	}
	return retweeters, nil
}

func readRetweetersFromFile() ([]string, error) {
	file, err := os.OpenFile(retweetersFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open retweeters file for reading: %w", err)
	}
	defer file.Close()

	existingUsers := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		if username != "" {
			existingUsers[username] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading retweeters file: %w", err)
	}

	var users []string
	for user := range existingUsers {
		users = append(users, user)
	}
	return users, nil
}

func updateRetweetersFile(newRetweeters []string) error {
	existingUsers, err := readRetweetersFromFile()
	if err != nil {
		return fmt.Errorf("failed to read existing retweeters: %w", err)
	}

	allUniqueUsers := make(map[string]struct{})
	for _, user := range existingUsers {
		allUniqueUsers[user] = struct{}{}
	}	
	for _, user := range newRetweeters {
		allUniqueUsers[strings.ToLower(user)] = struct{}{}
	}

	var sortedUsers []string
	for user := range allUniqueUsers {
		sortedUsers = append(sortedUsers, user)
	}

	var buf bytes.Buffer
	for _, user := range sortedUsers {
		buf.WriteString(user + "\n")
	}

	if err := ioutil.WriteFile(retweetersFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write retweeters file: %w", err)
	}
	return nil
}

func pickWinners(retweeters []string, count int) []string {
	if count <= 0 {
		return []string{}
	}
	if count >= len(retweeters) {
		return retweeters
	}

	shuffledRetweeters := make([]string, len(retweeters))
	copy(shuffledRetweeters, retweeters)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(shuffledRetweeters), func(i, j int) {
		shuffledRetweeters[i], shuffledRetweeters[j] = shuffledRetweeters[j], shuffledRetweeters[i]
	})

	return shuffledRetweeters[:count]
}