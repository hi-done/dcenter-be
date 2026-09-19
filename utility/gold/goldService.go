package gold

import (
	"context"
	"dcenter/internal/dao"
	"dcenter/internal/model"
	"dcenter/internal/model/do"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// goldApiURL is the data API behind https://www.huilvbiao.com/gold.
	goldApiURL = "https://www.huilvbiao.com/api/gold_indexApi"
	// goldPageURL is used as the Referer header.
	goldPageURL = "https://www.huilvbiao.com/gold"
	// goldFetchInterval is how often gold prices are fetched.
	goldFetchInterval = time.Minute
	// goldFetchTimeout is the HTTP timeout for a single fetch.
	goldFetchTimeout = 15 * time.Second
)

// goldProducts maps the remote hq_str variable to a stable symbol, display name and unit.
var goldProducts = []struct {
	Var    string
	Symbol string
	Name   string
	Unit   string
}{
	{"hq_str_gds_AUTD", "AUTD", "黄金延期", "元/克"},
	{"hq_str_hf_GC", "GC", "纽约黄金", "美元/盎司"},
	{"hq_str_hf_XAU", "XAU", "伦敦金", "美元/盎司"},
}

// hqStrPattern matches lines like: var hq_str_gds_AUTD="928.00,0,...,黄金延期";
var hqStrPattern = regexp.MustCompile(`var\s+(hq_str_\w+)\s*=\s*"([^"]*)"`)

type goldService struct {
	ctx       context.Context
	price     float64
	unit      string
	priceTime time.Time
}

var (
	GS *goldService
	mu sync.Mutex
)

func NewGoldService() *goldService {
	if GS != nil {
		return GS
	}
	mu.Lock()
	defer mu.Unlock()
	if GS == nil {
		GS = &goldService{
			ctx: context.Background(),
		}
		GS.start()
	}
	return GS
}

func (s *goldService) start() {
	go func() {
		s.fetchGoldPrices()

		ticker := time.NewTicker(goldFetchInterval)
		defer ticker.Stop()
		for range ticker.C {
			s.fetchGoldPrices()
		}
	}()
}

func (s *goldService) fetchGoldPrices() {
	records, err := s.fetchFromHuilvbiao()
	if err != nil {
		g.Log().Errorf(s.ctx, "[GoldPrice] Fetch failed: %v", err)
		return
	}

	// 实现今天的金价更新到gold_price数据表中，当天的同一类别金价只保留一条
	for _, record := range records {
		if err = s.saveGoldPrice(record); err != nil {
			g.Log().Errorf(s.ctx, "[GoldPrice] Save %s quote failed: %v", record.Symbol, err)
		}
	}

	g.Log().Infof(s.ctx, "[GoldPrice] Updated %d records, first quote time: %s",
		len(records), records[0].QuoteTime.Format("2006-01-02 15:04:05"))
}

// saveGoldPrice keeps only one row per symbol per quote date: the row of the
// same symbol and quote date is updated if it exists, otherwise a new row is inserted.
func (s *goldService) saveGoldPrice(record model.GoldPrice) error {
	columns := dao.GoldPrice.Columns()
	quoteTime := gtime.NewFromTime(record.QuoteTime)

	data := do.GoldPrice{
		Symbol:     record.Symbol,
		Name:       record.Name,
		Price:      record.Price,
		OpenPrice:  record.OpenPrice,
		HighPrice:  record.HighPrice,
		LowPrice:   record.LowPrice,
		PrevSettle: record.PrevSettle,
		Unit:       record.Unit,
		QuoteTime:  quoteTime,
	}

	existing, err := dao.GoldPrice.Ctx(s.ctx).
		Fields(columns.Id).
		Where(columns.Symbol, record.Symbol).
		WhereGTE(columns.QuoteTime, quoteTime.StartOfDay()).
		WhereLTE(columns.QuoteTime, quoteTime.EndOfDay()).
		One()
	if err != nil {
		return err
	}

	if existing.IsEmpty() {
		_, err = dao.GoldPrice.Ctx(s.ctx).Data(data).Insert()
		return err
	}

	_, err = dao.GoldPrice.Ctx(s.ctx).
		Data(data).
		Where(columns.Id, existing[columns.Id].Int64()).
		Update()
	return err
}

func (s *goldService) fetchFromHuilvbiao() ([]model.GoldPrice, error) {
	url := fmt.Sprintf("%s?t=%d", goldApiURL, time.Now().UnixMilli())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", goldPageURL)

	client := &http.Client{Timeout: goldFetchTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	return s.parseGoldIndex(string(body))
}

func (s *goldService) parseGoldIndex(body string) ([]model.GoldPrice, error) {
	quotes := make(map[string]string)
	for _, m := range hqStrPattern.FindAllStringSubmatch(body, -1) {
		quotes[m[1]] = m[2]
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("no hq_str data found in response")
	}

	var records []model.GoldPrice
	for _, p := range goldProducts {
		fields := strings.Split(quotes[p.Var], ",")
		if len(fields) < 13 {
			g.Log().Debugf(s.ctx, "[GoldPrice] Quote for %s missing or malformed, skipped", p.Var)
			continue
		}

		price, err := parsePrice(fields[0])
		if err != nil {
			g.Log().Debugf(s.ctx, "[GoldPrice] Invalid price for %s: %q", p.Var, fields[0])
			continue
		}

		priceTime := parseQuoteTime(fields[12], fields[6])

		if p.Symbol == "AUTD" {
			s.price = price
			s.unit = p.Unit
			s.priceTime = priceTime
		}

		records = append(records, model.GoldPrice{
			Symbol:     p.Symbol,
			Name:       p.Name,
			Price:      price,
			OpenPrice:  parsePriceOrZero(fields[8]),
			HighPrice:  parsePriceOrZero(fields[4]),
			LowPrice:   parsePriceOrZero(fields[5]),
			PrevSettle: parsePriceOrZero(fields[7]),
			Unit:       p.Unit,
			QuoteTime:  priceTime,
		})
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no valid gold price parsed")
	}
	return records, nil
}

func (s *goldService) GetGoldPrice() (float64, string, time.Time) {
	return s.price, s.unit, s.priceTime
}

// parsePrice converts a quote field into a float64.
func parsePrice(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// parsePriceOrZero is a lenient variant of parsePrice used for optional fields.
func parsePriceOrZero(s string) float64 {
	v, err := parsePrice(s)
	if err != nil {
		return 0
	}
	return v
}

// parseQuoteTime combines the date and clock fields into a local time.
func parseQuoteTime(date, clock string) time.Time {
	date = strings.TrimSpace(date)
	clock = strings.TrimSpace(clock)
	if date == "" || clock == "" {
		return time.Now()
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+clock, time.Local)
	if err != nil {
		return time.Now()
	}
	return t
}
