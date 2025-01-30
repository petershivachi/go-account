package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/petershivachi/go_transact/db/mock"
	db "github.com/petershivachi/go_transact/db/sqlc"
	"github.com/petershivachi/go_transact/util"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testcases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{name: "OK", accountID: account.ID, buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
		}, checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusOK, recorder.Code)
			requireMatchBodyAccount(t, recorder.Body, account)
		},
		},
		{name: "Not Found", accountID: account.ID, buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
		}, checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusNotFound, recorder.Code)
		},
		},
		{name: "Internal Server Error", accountID: account.ID, buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
		}, checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
		},
		{name: "Invalid ID", accountID: 0, buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any())
		}, checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			//check response
			require.Equal(t, http.StatusBadRequest, recorder.Code)
		},
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			//build stubs
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

func requireMatchBodyAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
