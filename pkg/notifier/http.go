package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
)

type Http struct {
	client           *http.Client
	url              url.URL
	method           string
	acceptedStatuses []int
	extraHeaders     map[string]string
}

func NewHttp(
	client *http.Client,
	url url.URL,
	method string,
	acceptedStatuses []int,
	extraHeaders map[string]string,
) *Http {
	return &Http{
		client:           client,
		url:              url,
		method:           method,
		acceptedStatuses: acceptedStatuses,
		extraHeaders:     extraHeaders,
	}
}

type request struct {
	State      string       `json:"state"`
	Attributes requestAttrs `json:"attributes"`
}

type requestAttrs struct {
	UnitOfMeasurement string `json:"unit_of_measurement"`
	FriendlyName      string `json:"friendly_name"`
}

func (h *Http) Notify(ctx context.Context, rateChange aggregator.RateChange) error {
	fmt.Println("Sending HTTP notification")
	r := request{
		State: fmt.Sprintf("%.3f", rateChange.New.Rate),
		Attributes: requestAttrs{
			UnitOfMeasurement: rateChange.New.To,
			FriendlyName:      fmt.Sprintf("%s to %s", rateChange.New.From, rateChange.New.To),
		},
	}

	body := bytes.NewBufferString("")
	if err := json.NewEncoder(body).Encode(r); err != nil {
		return err
	}

	httpReq, err := http.NewRequest(h.method, h.url.String(), body)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for header, value := range h.extraHeaders {
		httpReq.Header.Set(header, value)
	}

	resp, err := h.client.Do(httpReq.WithContext(ctx))
	if err != nil {
		return err
	}

	if len(h.acceptedStatuses) > 0 {
		for _, status := range h.acceptedStatuses {
			if resp.StatusCode == status {
				return nil
			}
		}

		bod, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		return fmt.Errorf("wrong HTTP status: %d\n\nContent: %s\n\n", resp.StatusCode, bod)
	}

	fmt.Println("HTTP notification was sent successfully")

	return nil
}
