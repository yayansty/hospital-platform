package clients

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"medic-api/models"
)

type AplicareClient struct {
	BaseURL    string
	ConsID     string
	SecretKey  string
	NPPK       string
	HTTPClient *http.Client
}

type AplicareResponse struct {
	MetaData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"metadata"`

	Response interface{} `json:"response"`
}

type AplicareBed struct {
	Kapasitas          int64  `json:"kapasitas"`
	KodeKelas          string `json:"kodekelas"`
	KodeRuang          string `json:"koderuang"`
	LastUpdate         int64  `json:"last_update"`
	LastUpdateText     string `json:"lastupdate"`
	NamaKelas          string `json:"namakelas"`
	NamaRuang          string `json:"namaruang"`
	RowNumber          int64  `json:"rownumber"`
	Stat               string `json:"stat"`
	Tersedia           int64  `json:"tersedia"`
	TersediaPria       int64  `json:"tersediapria"`
	TersediaPriaWanita int64  `json:"tersediapriawanita"`
	TersediaWanita     int64  `json:"tersediawanita"`
}

type AplicareBedReadResponse struct {
	MetaData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"metadata"`

	Response struct {
		List []AplicareBed `json:"list"`
	} `json:"response"`
}

//type AplicareBedReadResponse struct {
//	MetaData struct {
//		Code    int    `json:"code"`
//		Message string `json:"message"`
//	} `json:"metadata"`

//	Response interface{} `json:"response"`
//}

type AplicareBedRequest struct {
	KodeKelas          string `json:"kodekelas"`
	KodeRuang          string `json:"koderuang"`
	NamaRuang          string `json:"namaruang"`
	Kapasitas          string `json:"kapasitas"`
	Tersedia           string `json:"tersedia"`
	TersediaPria       string `json:"tersediapria"`
	TersediaWanita     string `json:"tersediawanita"`
	TersediaPriaWanita string `json:"tersediapriawanita"`
}

func NewAplicareClient() *AplicareClient {
	return &AplicareClient{
		BaseURL:   strings.TrimRight(os.Getenv("BPJS_URL"), "/"),
		ConsID:    os.Getenv("BPJS_CONS_ID"),
		SecretKey: os.Getenv("BPJS_SECRET_KEY"),
		NPPK:      os.Getenv("BPJS_NOPPK"),

		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
					MaxVersion: tls.VersionTLS12,
				},
				ForceAttemptHTTP2:  false,
				DisableCompression: true,
			},
		},
	}
}

func (c *AplicareClient) UpdateBed(room models.RoomAvailability) (*AplicareResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Sesuai dengan implementasi VB.NET:
	// HMAC-SHA256(ConsID + "&" + Timestamp, SecretKey)
	signature := c.generateSignature(timestamp)

	payload := AplicareBedRequest{
		KodeKelas: room.KodeKelas,
		KodeRuang: room.KodeRuang,
		NamaRuang: room.NamaRuang,
		Kapasitas: strconv.FormatInt(room.Kapasitas, 10),
		Tersedia:  strconv.FormatInt(room.Tersedia, 10),

		// Sesuai sistem lama:
		TersediaPria:       "0",
		TersediaWanita:     "0",
		TersediaPriaWanita: "0",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Aplicare payload: %w", err)
	}

	url := fmt.Sprintf(
		"%s/aplicaresws/rest/bed/update/%s",
		c.BaseURL,
		c.NPPK,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aplicare request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-cons-id", c.ConsID)
	req.Header.Set("X-timestamp", timestamp)
	req.Header.Set("X-signature", signature)
	req.Header.Set("User-Agent", "curl/8.0")
	req.Header.Set("Accept", "*/*")

	fmt.Printf(
		"Aplicare request: method=%s url=%s cons_id=%s timestamp=%s signature=%s body=%s\n",
		req.Method,
		req.URL.String(),
		c.ConsID,
		timestamp,
		signature,
		string(jsonData),
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Aplicare: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Aplicare response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"Aplicare returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var result AplicareResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf(
			"failed to decode Aplicare response: %w; response=%s",
			err,
			string(responseBody),
		)
	}

	return &result, nil
}

func (c *AplicareClient) generateSignature(timestamp string) string {
	message := c.ConsID + "&" + timestamp

	mac := hmac.New(
		sha256.New,
		[]byte(c.SecretKey),
	)

	mac.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *AplicareClient) ReadBeds(start, limit int) (*AplicareBedReadResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	signature := c.generateSignature(timestamp)

	url := fmt.Sprintf(
		"%s/aplicaresws/rest/bed/read/%s/%d/%d",
		c.BaseURL,
		c.NPPK,
		start,
		limit,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aplicare request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-cons-id", c.ConsID)
	req.Header.Set("X-timestamp", timestamp)
	req.Header.Set("X-signature", signature)
	req.Header.Set("User-Agent", "curl/8.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Aplicare: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Aplicare response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"Aplicare returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var result AplicareBedReadResponse

	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf(
			"failed to decode Aplicare response: %w; response=%s",
			err,
			string(responseBody),
		)
	}

	return &result, nil
}
