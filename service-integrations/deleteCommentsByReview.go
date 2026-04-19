package serviceIntegrations

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func DeleteCommentsByReview(reviewID int64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	baseURL := "http://comments-servise:8080/comments/delete"
	params := url.Values{}
	params.Add("review_id", strconv.FormatInt(reviewID, 10))
	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, nil)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SERVICE_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error executing request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("comment-service returned error: %d", resp.StatusCode)
		return
	}

	log.Println("request to comment-service succeeded")
}
