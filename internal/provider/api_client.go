package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type APIClient struct {
	client  *http.Client
	headers map[string]string
	baseURL string
}

func (client *APIClient) Post(ctx context.Context, url string, body []byte, desiredHTTPStatus string) ([]byte, error) {
	fullURL := client.baseURL + url
	tflog.Trace(ctx, fmt.Sprintf("my full URL is %s", fullURL))
	// Create a POST request with the specified URL and body.
	req, err := http.NewRequest(
		http.MethodPost,
		fullURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	for key, value := range client.headers {
		req.Header.Set(key, value)
	}

	// Execute the request using the client's Do method.
	httpResp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body:", err)
	}
	tflog.Trace(ctx, fmt.Sprintf("Response Status: %s", httpResp.Status))
	tflog.Trace(ctx, fmt.Sprintf("Response Body: %s", bodyBytes))
	if httpResp.Status != desiredHTTPStatus {
		return nil, fmt.Errorf("The HTTP response code was %s but I wanted %s.", httpResp.Status, desiredHTTPStatus)
	}
	return bodyBytes, nil
}

func (client *APIClient) LookupExactlyOne(ctx context.Context, url string, body []byte) (json.RawMessage, error) {
	respBytes, err := client.Post(ctx, url, body, "200 OK")
	if err != nil {
		return nil, fmt.Errorf("error from HTTP client: %v", err)
	}

	var itemsReturned []json.RawMessage
	err = json.Unmarshal(respBytes, &itemsReturned)
	if err != nil {
		return nil, fmt.Errorf("I couldn't unmarshal the outer list because %v", err)
	}
	if len(itemsReturned) != 1 {
		return nil, fmt.Errorf("I wanted exactly one org but it returned %v", len(itemsReturned))
	}
	return itemsReturned[0], nil
}
