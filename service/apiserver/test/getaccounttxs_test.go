package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bnb-chain/zkbas/service/apiserver/internal/types"
)

func (s *ApiServerSuite) TestGetAccountTxs() {
	type args struct {
		by     string
		value  string
		offset int
		limit  int
	}

	type testcase struct {
		name     string
		args     args
		httpCode int
	}

	tests := []testcase{
		{"not found by index", args{"account_index", "99999999", 0, 10}, 200},
		{"not found by name", args{"account_name", "fakeaccount.legend", 0, 10}, 200},
		{"not found by pk", args{"account_pk", "fake8470d33c59a5cbf5e10df426eb97c2773ab890c3364f4162ba782a56ca998", 0, 10}, 200},
		{"invalid by", args{"invalidby", "fake8470d33c59a5cbf5e10df426eb97c2773ab890c3364f4162ba782a56ca998", 0, 10}, 400},
	}

	statusCode, txs := GetTxs(s, 0, 100)
	if statusCode == http.StatusOK && len(txs.Txs) > 0 {
		_, account := GetAccount(s, "name", txs.Txs[len(txs.Txs)-1].AccountName)
		tests = append(tests, []testcase{
			{"found by index", args{"account_index", strconv.Itoa(int(account.Index)), 0, 10}, 200},
			{"found by name", args{"account_name", account.Name, 0, 10}, 200},
			{"found by pk", args{"account_pk", account.Pk, 0, 10}, 200},
		}...)
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			httpCode, result := GetAccountTxs(s, tt.args.by, tt.args.value, tt.args.offset, tt.args.limit)
			assert.Equal(t, tt.httpCode, httpCode)
			if httpCode == http.StatusOK {
				if tt.args.offset < int(result.Total) {
					assert.True(t, len(result.Txs) > 0)
					assert.NotNil(t, result.Txs[0].BlockHeight)
					assert.NotNil(t, result.Txs[0].Hash)
					assert.NotNil(t, result.Txs[0].Type)
					assert.NotNil(t, result.Txs[0].StateRoot)
					assert.NotNil(t, result.Txs[0].Info)
					assert.NotNil(t, result.Txs[0].Status)
				}
				fmt.Printf("result: %+v \n", result)
			}
		})
	}

}

func GetAccountTxs(s *ApiServerSuite, by, value string, offset, limit int) (int, *types.Txs) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/accountTxs?by=%s&value=%s&offset=%d&limit=%d", s.url, by, value, offset, limit))
	assert.NoError(s.T(), err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(s.T(), err)

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil
	}
	result := types.Txs{}
	//nolint:errcheck
	json.Unmarshal(body, &result)
	return resp.StatusCode, &result
}
