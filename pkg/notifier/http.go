package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
)

type Http struct {
	client           *http.Client
	url              url.URL
	method           string
	acceptedStatuses []int
}

func NewHttp(client *http.Client, url url.URL, method string, acceptedStatuses []int) *Http {
	return &Http{client: client, url: url, method: method, acceptedStatuses: acceptedStatuses}
}

type request struct {
	Old *aggregator.Rate `json:"old,omitempty"`
	New *aggregator.Rate `json:"new,omitempty"`
	Msg string           `json:"message"`
}

func (h *Http) Notify(ctx context.Context, rateChange aggregator.RateChange) error {
	r := request{
		Old: rateChange.Old,
		New: rateChange.New,
		Msg: rateChange.String(),
	}

	body := bytes.NewBufferString("")
	if err := json.NewEncoder(body).Encode(r); err != nil {
		return err
	}

	httpReq, err := http.NewRequest(h.method, h.url.String(), body)
	if err != nil {
		return err
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

		return fmt.Errorf("wrong HTTP status: %d", resp.StatusCode)
	}

	fmt.Printf("HTTP notification was sent successfully")

	return nil
}
