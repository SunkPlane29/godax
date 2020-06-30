package godax

import (
	"net/http"
	"reflect"
	"testing"
)

const (
	baseRestURL = "https://test.pro.coinbase"
	baseWsURL   = "wss://test.ws-feed.pro.coinbase.com"
	key         = "super_secret_key_123_abc"
	secret      = "MTIzYWJjU3VwZXJTZWNyZXRTZWNyZXQ="
	passphrase  = "1q2w3e4r"
)

func TestClient_ListAccounts(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	defaultFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	tests := []struct {
		name    string
		fields  fields
		mock    HTTPClient
		want    []ListAccount
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to ListAccounts with no results",
			fields:  defaultFields(),
			want:    []ListAccount{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to ListAccounts with one account",
			fields: defaultFields(),
			want: []ListAccount{{
				ID:        "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
				Currency:  "BTC",
				Balance:   "10.01",
				Available: "15.449977",
				Hold:      "wat",
			}},
			wantRaw: `[{
                "id": "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
                "currency": "BTC",
                "balance": "10.01",
                "available": "15.449977",
                "hold": "wat"
            }]`,
		},
		{
			name:   "when a successful call is made to ListAccounts with many accounts",
			fields: defaultFields(),
			want: []ListAccount{{
				ID:        "766b7a10-06bb-4b1d-a4b3-679d025352ad",
				Currency:  "BTC",
				Balance:   "0.00000000000",
				Available: "123.456789",
				Hold:      "0.101",
			}, {
				ID:        "543c3da9-a71d-4a4c-b6e7-edc132ff704e",
				Currency:  "ETH",
				Balance:   "1.300006",
				Available: "9000.67685938262624",
				Hold:      "0.101",
			}, {
				ID:        "dcbe61c2-a1bd-444c-b41a-3c6b2363afd6",
				Currency:  "BAT",
				Balance:   "9999.677773333",
				Available: "9000.67685938262624",
				Hold:      "0.101",
			}},
			wantRaw: `[{
                "id": "766b7a10-06bb-4b1d-a4b3-679d025352ad",
                "currency": "BTC",
                "balance": "0.00000000000",
                "available": "123.456789",
                "hold": "0.101"
            },{
                "id": "543c3da9-a71d-4a4c-b6e7-edc132ff704e",
                "currency": "ETH",
                "balance": "1.300006",
                "available": "9000.67685938262624",
                "hold": "0.101"
            },{
                "id": "dcbe61c2-a1bd-444c-b41a-3c6b2363afd6",
                "currency": "BAT",
                "balance": "9999.677773333",
                "available": "9000.67685938262624",
                "hold": "0.101"
            }]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.ListAccounts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccount(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	defaultFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	type args struct {
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Account
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccount and no account is found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    Account{},
			wantRaw: `{}`,
		},
		{
			name:   "when a successful call is made to GetAccount and an account is found",
			fields: defaultFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: Account{
				ID:        "a1b2c3d4",
				Balance:   "1.100",
				Holds:     "0.100",
				Available: "101.56",
				Currency:  "USD",
			},
			wantRaw: `{
                "id": "a1b2c3d4",
                "balance": "1.100",
                "holds": "0.100",
                "available": "101.56",
                "currency": "USD"
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.GetAccount(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccountHistory(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	defaultFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	type args struct {
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []AccountActivity
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccountHistory and no history is found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    []AccountActivity{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to GetAccountHistory and one history is found",
			fields: defaultFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: []AccountActivity{{
				ID:        "100",
				CreatedAt: "2014-11-07T08:19:27.028459Z",
				Amount:    "0.001",
				Balance:   "239.669",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "d50ec984-77a8-460a-b958-66f114b0de9b",
					TradeID:   "74",
					ProductID: "BTC-USD",
				},
			}},
			wantRaw: `[{
                "id": "100",
                "created_at": "2014-11-07T08:19:27.028459Z",
                "amount": "0.001",
                "balance": "239.669",
                "type": "fee",
                "details": {
                    "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
                    "trade_id": "74",
                    "product_id": "BTC-USD"
                }
            }]`,
		},
		{
			name:   "when a successful call is made to GetAccountHistory and multiple histories are found",
			fields: defaultFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: []AccountActivity{{
				ID:        "100",
				CreatedAt: "2014-11-07T08:19:27.028459Z",
				Amount:    "0.001",
				Balance:   "239.669",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "d50ec984-77a8-460a-b958-66f114b0de9b",
					TradeID:   "74",
					ProductID: "BTC-USD",
				},
			}, {
				ID:        "80",
				CreatedAt: "2015-12-04T08:19:27.028459Z",
				Amount:    "0.011",
				Balance:   "4059.212345",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "8b9258f8-811b-429b-810d-71fede464b29",
					TradeID:   "99",
					ProductID: "BTC-ETH",
				},
			}},
			wantRaw: `[{
                "id": "100",
                "created_at": "2014-11-07T08:19:27.028459Z",
                "amount": "0.001",
                "balance": "239.669",
                "type": "fee",
                "details": {
                    "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
                    "trade_id": "74",
                    "product_id": "BTC-USD"
                }
            },{
                "id": "80",
                "created_at": "2015-12-04T08:19:27.028459Z",
                "amount": "0.011",
                "balance": "4059.212345",
                "type": "fee",
                "details": {
                    "order_id": "8b9258f8-811b-429b-810d-71fede464b29",
                    "trade_id": "99",
                    "product_id": "BTC-ETH"
                }
            }]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.GetAccountHistory(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccountHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccountHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccountHolds(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	defaultFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	type args struct {
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []AccountHold
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccountHolds and no holds are found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    []AccountHold{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to GetAccountHolds and one hold is found",
			fields: defaultFields(),
			args:   args{accountID: "1q2w3e4r"},
			want: []AccountHold{{
				ID:        "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
				AccountID: "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
				CreatedAt: "2014-11-06T10:34:47.123456Z",
				UpdatedAt: "2014-11-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}},
			wantRaw: `[{
                "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
                "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
                "created_at": "2014-11-06T10:34:47.123456Z",
                "updated_at": "2014-11-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
            }]`,
		},
		{
			name:   "when a successful call is made to GetAccountHolds and multiple holds are found",
			fields: defaultFields(),
			args:   args{accountID: "1q2w3e4r"},
			want: []AccountHold{{
				ID:        "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
				AccountID: "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
				CreatedAt: "2014-11-06T10:34:47.123456Z",
				UpdatedAt: "2014-11-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}, {
				ID:        "3d58f10b-3d9a-4d38-bb51-c8800f5ad4ca",
				AccountID: "b6f8fee0-f47f-481a-98ee-08d397681edb",
				CreatedAt: "2015-10-06T10:34:47.123456Z",
				UpdatedAt: "2015-10-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}},
			wantRaw: `[{
                "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
                "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
                "created_at": "2014-11-06T10:34:47.123456Z",
                "updated_at": "2014-11-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
            },
            {
                "id": "3d58f10b-3d9a-4d38-bb51-c8800f5ad4ca",
                "account_id": "b6f8fee0-f47f-481a-98ee-08d397681edb",
                "created_at": "2015-10-06T10:34:47.123456Z",
                "updated_at": "2015-10-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
            }]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.GetAccountHolds(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccountHolds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccountHolds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func validateHeaders(t *testing.T, client *Client) {
	compareHeader(t, client, "CB-ACCESS-KEY", key)
	compareHeader(t, client, "CB-ACCESS-PASSPHRASE", passphrase)
	compareHeader(t, client, "User-Agent", userAgent)
	validateHeaderPresent(t, client, "CB-ACCESS-SIGN")
	validateHeaderPresent(t, client, "CB-ACCESS-TIMESTAMP")
}

func compareHeader(t *testing.T, c *Client, wantHeader string, wantContent string) {
	if c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader) != wantContent {
		t.Errorf(
			"%s header should be %s, was '%s'\n",
			wantHeader,
			wantContent,
			c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader),
		)
	}
}

func validateHeaderPresent(t *testing.T, c *Client, wantHeader string) {
	if c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader) == "" {
		t.Errorf("%s header should not be empty\n", wantHeader)
	}
}