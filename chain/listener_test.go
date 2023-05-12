package chain

import (
	"fmt"
	"github.com/mapprotocol/brc20Relayer/utils"
	"reflect"
	"testing"
)

func Test_requestAndParse(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name       string
		args       args
		wantDetail Detail
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDetail, err := requestAndParse(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("requestAndParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("requestAndParse() gotDetail = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestName(t *testing.T) {
	nums := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for _, num := range nums {
		n := num
		utils.Go(func() {
			fmt.Println("============================== num: ", n)
		})
	}
	select {}
}
