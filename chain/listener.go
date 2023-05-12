package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/brc20Relayer/config"
	"github.com/mapprotocol/brc20Relayer/dao"
	"github.com/mapprotocol/brc20Relayer/startstore"
	"github.com/mapprotocol/brc20Relayer/utils"
	"gorm.io/gorm"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

const (
	RetryLimit    = 5
	RetryInterval = time.Minute * 2
)

const RequestURLFormat = "https://unisat.io/brc20-api-v2/address/%s/brc20/ordi/history?start=%d&limit=%d&type=receive"

const (
	Limit   = 10
	Address = "bc1p32e2t5j5uhe5dm3zn8kevuu59fewk9ut5fjngm3m4agx68hc409qeeelse"
)

type Listener struct {
	Config         *config.Config
	Client         *ethclient.Client
	StartStorePath string
	Stop           chan struct{}
}

type Item struct {
	Type              string `json:"type"`
	Valid             bool   `json:"valid"`
	TxId              string `json:"txid"`
	TxIdx             uint64 `json:"TxIdx"`
	InscriptionNumber uint64 `json:"inscriptionNumber"`
	InscriptionId     string `json:"inscriptionId"`
	From              string `json:"from"`
	To                string `json:"to"`
	Satoshi           uint64 `json:"satoshi"`
	Amount            string `json:"amount"`
	OverallBalance    string `json:"overallBalance"`
	TransferBalance   string `json:"transferBalance"`
	AvailableBalance  string `json:"availableBalance"`
	Height            uint64 `json:"height"`
	Blocktime         uint64 `json:"blocktime"`
}

type Detail []*Item

func (d Detail) Len() int           { return len(d) }
func (d Detail) Less(i, j int) bool { return d[i].Height < d[j].Height }
func (d Detail) Swap(i, j int)      { d[i].Height, d[j].Height = d[j].Height, d[i].Height }

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total  int    `json:"total"`
		Start  int    `json:"start"`
		Detail Detail `json:"detail"`
	} `json:"data"`
}

func (l *Listener) Poll() error {
	var (
		latestStart = l.Config.Start
		latestCount = l.Config.Count
		retry       = RetryLimit
	)

	logger := log.New("func", "Poll")
	logger.Info("Polling...", "start", latestStart)

	for {
		select {
		case <-l.Stop:
			return errors.New("polling terminated")
		default:
			if retry == 0 {
				logger.Error("Polling failed, retries exceeded")
				l.Stop <- struct{}{}
				return nil
			}

			if latestCount == Limit {
				latestStart++
			}

			url, err := getRequestURL(Address, latestStart, Limit)
			if err != nil {
				logger.Error("Failed get request url", "address", Address, "start", latestStart, "limit", Limit, "error", err)
				l.Stop <- struct{}{}
				return nil
			}
			detail, err := requestAndParse(url)
			if err != nil {
				logger.Error("Failed get request and parse response", "url", url, "error", err)
				retry--
				continue
			}

			sort.Sort(detail)
			if latestCount < Limit {
				detail = detail[latestCount:]
			}

			count := len(detail)
			bsChain := make(chan *dao.BRC20, count)
			wg := &sync.WaitGroup{}
			for _, item := range detail {
				wg.Add(1)
				go func(item *Item, bs chan<- *dao.BRC20, wg *sync.WaitGroup) {
					defer wg.Done()

					var (
						err    error
						txHash common.Hash
					)
					// send transaction to MAP relay chain
					for i := 0; i < RetryLimit; i++ {
						txHash, err = MintToken()
						if err != nil {
							logger.Warn("Failed to mint token", "error", err)
							continue
						}
					}
					if err != nil {
						logger.Error("Failed to mint token", "error", err)
						l.Stop <- struct{}{}
					}

					m, _ := json.Marshal(item)
					b := &dao.BRC20{
						Metadata:      string(m),
						Height:        item.Height,
						InscriptionID: item.InscriptionId,
						TxID:          item.TxId,
						TxIdx:         item.TxIdx,
						Account:       "",
						TxHash:        txHash.Hex(),
						CreatedAt:     time.Time{},
						UpdatedAt:     time.Time{},
						DeletedAt:     gorm.DeletedAt{},
					}
					bs <- b
				}(item, bsChain, wg)
			}

			wg.Wait()

			// insert data to database
			bs := make([]*dao.BRC20, 0, count)
			for b := range bsChain {
				bs = append(bs, b)
			}
			if err := dao.NewBRC20().BatchCreate(bs); err != nil {
				logger.Error("Failed to batch create brc20", "bs", utils.JSON(bs))
			}

			processedCount := latestCount + uint64(count)
			if err := startstore.StoreLatestStart(l.StartStorePath, latestStart, processedCount); err != nil {
				logger.Error("Failed to store latest start number", "start", latestStart, "count", count, "error", err)
			}

			retry = RetryLimit
			if processedCount == Limit {
				latestStart++
			} else {
				time.Sleep(RetryInterval)
			}
		}
	}
}

func getRequestURL(address string, start, limit uint64) (string, error) {
	if start < 0 {
		return "", errors.New("start invalid")
	}
	if limit < 1 {
		return "", errors.New("limit invalid")
	}
	if utils.IsEmpty(address) {
		return "", errors.New("address invalid")
	}
	return fmt.Sprintf(RequestURLFormat, address, start, limit), nil
}

func requestAndParse(url string) (detail Detail, err error) {
	resp := &Response{}
	rb, err := utils.Get(url, nil, nil)
	if err != nil {
		return detail, err
	}
	if err := json.Unmarshal(rb, resp); err != nil {
		log.Error("json unmarshal failed", "func", "requestAndParse", "rsp", string(rb), "err", err)
		return detail, err
	}
	// todo judge code ?
	return resp.Data.Detail, nil
}

func MintToken() (common.Hash, error) {
	return common.Hash{}, nil
}
