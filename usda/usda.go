/*
Package usda implements an API client library for the USDA Food Composition Databases
https://ndb.nal.usda.gov/ndb/doc/index
*/
package usda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
)

const entryPointURL string = "api.nal.usda.gov/ndb/"

type USDAclient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(apiKey string, httpClient *http.Client) *USDAclient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	bURL, err := url.Parse("https://" + apiKey + "@" + entryPointURL)

	if err != nil {
		// Some problem with the API Key. Security breach Point
		panic(err)
	}
	return &USDAclient{baseURL: bURL, httpClient: httpClient}
}

func (c *USDAclient) newRequest(pathAPI string, data interface{}) (*http.Request, error) {

	//finalPathURL := c.baseURL.ResolveReference(&url.URL{Path: pathAPI})
	finalPathURL := c.baseURL.String() + pathAPI

	bStream := bytes.NewBuffer([]byte{}) // Use new()
	err := json.NewEncoder(bStream).Encode(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, finalPathURL, bStream)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, err
}

func (c *USDAclient) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
		fmt.Println(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

type USDAQueryOptions struct {
	Max    int    `url:"max,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Sort   string `url:"sort,omitempty"`
}

func addQueryOptions(pathAPI string, opts *USDAQueryOptions) (string, error) {

	if opts == nil {
		return pathAPI, nil
	}

	v, err := query.Values(opts)
	if err != nil {
		return pathAPI, err
	}

	u, err := url.Parse(pathAPI)
	if err != nil {
		return pathAPI, err
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}

type USDAListReqParams struct {
	Lt string `json:"lt,omitempty"`
}

type USDAList struct {
	List struct {
		Lt    string `json:"lt"`
		Start int    `json:"start"`
		End   int    `json:"end"`
		Total int    `json:"total"`
		Sr    string `json:"sr"`
		Sort  string `json:"sort"`
		Item  []struct {
			Offset int    `json:"offset"`
			ID     string `json:"id"`
			Name   string `json:"name"`
		} `json:"item"`
	} `json:"list"`
}

func (c *USDAclient) GetList(ctx context.Context, params *USDAListReqParams, opts *USDAQueryOptions) (USDAList, error) {

	var ul USDAList

	endQuery, err := addQueryOptions("list", opts)
	if err != nil {
		return ul, err
	}

	req, err := c.newRequest(endQuery, params)
	if err != nil {
		return ul, err
	}

	_, err = c.do(ctx, req, &ul)

	return ul, err
}

type USDANutrientsReqParams struct {
	Fg          []string `json:"fg,omitempty"`
	Ndbno       string   `json:"ndbno,omitempty"`
	NutrientsID []string `json:"nutrients"`
	Subset      string   `json:"subset,omitempty"`
}

type USDANutrientReport struct {
	Report struct {
		Sr     string `json:"sr"`
		Groups []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"groups"`
		Subset string `json:"subset"`
		End    int    `json:"end"`
		Start  int    `json:"start"`
		Total  int    `json:"total"`
		Foods  []struct {
			Ndbno     string  `json:"ndbno"`
			Name      string  `json:"name"`
			Weight    float64 `json:"weight"`
			Measure   string  `json:"measure"`
			Nutrients []struct {
				NutrientID int     `json:"nutrient_id"`
				Nutrient   string  `json:"nutrient"`
				Unit       string  `json:"unit"`
				Value      string  `json:"value"`
				Gm         float64 `json:"gm"`
			} `json:"nutrients"`
		} `json:"foods"`
	} `json:"report"`
}

func (c *USDAclient) GetUSDANutrientReport(ctx context.Context, params *USDANutrientsReqParams, opts *USDAQueryOptions) (USDANutrientReport, error) {

	var unr USDANutrientReport

	endQuery, err := addQueryOptions("nutrients", opts)
	if err != nil {
		return unr, err
	}

	req, err := c.newRequest(endQuery, params)
	if err != nil {
		return unr, err
	}

	_, err = c.do(ctx, req, &unr)

	return unr, err
}

type USDAFoodsReqParams struct {
	Ndbno []string `json:"ndbno"`
	Type  string   `json:"type,omitempty"`
}

type USDAFoodsReport struct {
	Foods []struct {
		Food struct {
			Sr   string `json:"sr"`
			Type string `json:"type"`
			Desc struct {
				Ndbno string  `json:"ndbno"`
				Name  string  `json:"name"`
				Sd    string  `json:"sd"`
				Fg    string  `json:"fg"`
				Sn    string  `json:"sn"`
				Cn    string  `json:"cn"`
				Manu  string  `json:"manu"`
				Nf    float64 `json:"nf"`
				Cf    float64 `json:"cf"`
				Ff    float64 `json:"ff"`
				Pf    float64 `json:"pf"`
				R     string  `json:"r"`
				Rd    string  `json:"rd"`
				Ds    string  `json:"ds"`
				Ru    string  `json:"ru"`
			} `json:"desc"`
			Nutrients []struct {
				NutrientID interface{} `json:"nutrient_id"` //Issue
				Name       string      `json:"name"`
				Group      string      `json:"group"`
				Unit       string      `json:"unit"`
				Value      interface{} `json:"value"` //Issue
				Derivation string      `json:"derivation"`
				Sourcecode interface{} `json:"sourcecode"`
				Dp         interface{} `json:"dp"` //Issue
				Se         string      `json:"se"`
				Measures   []struct {
					Label string      `json:"label"`
					Eqv   float64     `json:"eqv"`
					Eunit string      `json:"eunit"`
					Qty   float64     `json:"qty"`
					Value interface{} `json:"value"` //Issue
				} `json:"measures"`
			} `json:"nutrients"`
			Sources []struct {
				ID      int    `json:"id"`
				Title   string `json:"title"`
				Authors string `json:"authors"`
				Vol     string `json:"vol"`
				Iss     string `json:"iss"`
				Year    string `json:"year"`
			} `json:"sources"`
			Footnotes []interface{} `json:"footnotes"`
			Langual   []interface{} `json:"langual"`
		} `json:"food,omitempty"`
		Error string `json:"error,omitempty"`
	} `json:"foods"`
	Count    int     `json:"count"`
	Notfound int     `json:"notfound"`
	API      float64 `json:"api"`
}

func (c *USDAclient) GetUSDAFoodsReportV2(ctx context.Context, params *USDAFoodsReqParams) (USDAFoodsReport, error) {

	var ufr USDAFoodsReport

	req, err := c.newRequest("V2/reports", params)
	if err != nil {
		return ufr, err
	}

	_, err = c.do(ctx, req, &ufr)

	return ufr, err
}

type USDASearchReqParams struct {
	Q  string `json:"q,omitempty"`
	Ds string `json:"ds,omitempty"`
	Fg string `json:"fg,omitempty"`
}

type USDASearch struct {
	List struct {
		Q     string `json:"q"`
		Sr    string `json:"sr"`
		Ds    string `json:"ds"`
		Start int    `json:"start"`
		End   int    `json:"end"`
		Total int    `json:"total"`
		Group string `json:"group"`
		Sort  string `json:"sort"`
		Item  []struct {
			Offset int    `json:"offset"`
			Group  string `json:"group"`
			Name   string `json:"name"`
			Ndbno  string `json:"ndbno"`
			Ds     string `json:"ds"`
			Manu   string `json:"manu"`
		} `json:"item"`
	} `json:"list"`
}

func (c *USDAclient) GetUSDASearch(ctx context.Context, params *USDASearchReqParams, opts *USDAQueryOptions) (USDASearch, error) {

	var us USDASearch

	endQuery, err := addQueryOptions("search", opts)
	if err != nil {
		return us, err
	}

	req, err := c.newRequest(endQuery, params)
	if err != nil {
		return us, err
	}

	_, err = c.do(ctx, req, &us)

	return us, err
}
