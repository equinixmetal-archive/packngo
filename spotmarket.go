package packngo

const spotMarketBasePath = "/market/spot/prices"

// SpotMarketService expooses Spot Market methods
type SpotMarketService interface {
	Prices() (PriceMap, *Response, error)
}

// SpotMarketServiceOp implements SpotMarketService
type SpotMarketServiceOp struct {
	client *Client
}

// PriceMap is a map of [facility][type]-> float Price
type PriceMap map[string]map[string]float64

// Prices gets currnt PriceMap from the API
func (s *SpotMarketServiceOp) Prices() (PriceMap, *Response, error) {

	type spotPrice struct {
		Price float64 `json:"price"`
	}

	type deepPriceMap map[string]map[string]spotPrice

	type marketRoot struct {
		SMPs deepPriceMap `json:"spot_market_prices"`
	}

	req, err := s.client.NewRequest("GET", spotMarketBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(marketRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	prices := make(PriceMap)
	for k, v := range root.SMPs {
		prices[k] = map[string]float64{}
		for kk, vv := range v {
			prices[k][kk] = vv.Price
		}
	}
	return prices, resp, err
}
