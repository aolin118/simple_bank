package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/aolin118/simple_bank/db/mock"
	db "github.com/aolin118/simple_bank/db/sqlc"
	"github.com/aolin118/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)
	// hashedPassword, err := util.HashPassword(password)
	// require.NoError(t, err)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkRespones func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkRespones: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},

		{
			name: "BadRequest",
			body: gin.H{
				"username":  user.Username,
				"test":      password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username:       user.Username,
				// 	HashedPassword: hashedPassword,
				// 	FullName:       user.FullName,
				// 	Email:          user.Email,
				// }
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRespones: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// add more cases
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRespones(t, recorder)
		})

	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	return db.User{
		Username: util.RandomString(6),
		FullName: util.RandomString(6),
		Email:    util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, gotUser, user)
}
