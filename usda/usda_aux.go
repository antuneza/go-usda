package usda

import (
	"context"
)

func (c *USDAclient) GetListByType(ctx context.Context, listType string) (USDAList, error) {

	params := &USDAListReqParams{Lt: listType}

	opts := &USDAQueryOptions{Max: 1500, Offset: 0, Sort: "id"}

	return c.GetList(ctx, params, opts)
}

func (c *USDAclient) GetNutrientsReport(ctx context.Context) (USDANutrientReport, error) {

	params := &USDANutrientsReqParams{Fg: []string{"0300"}, NutrientsID: []string{"306", "204"}}

	opts := &USDAQueryOptions{Max: 10, Offset: 0, Sort: "c"}

	return c.GetUSDANutrientReport(ctx, params, opts)
}

func (c *USDAclient) GetFoodNutrientsReport(ctx context.Context, ndbno string) (USDANutrientReport, error) {

	params := &USDANutrientsReqParams{NutrientsID: []string{"306", "204"}, Ndbno: ndbno}

	opts := &USDAQueryOptions{Max: 100, Offset: 0, Sort: "c"}

	return c.GetUSDANutrientReport(ctx, params, opts)
}

func (c *USDAclient) GetBasicFoodReport(ctx context.Context, ndbno string) (USDAFoodsReport, error) {

	params := &USDAFoodsReqParams{Ndbno: []string{ndbno}, Type: "b"}

	return c.GetUSDAFoodsReportV2(ctx, params)
}

func (c *USDAclient) GetFoodNameSearch(ctx context.Context, name string) (USDASearch, error) {

	params := &USDASearchReqParams{Q: name}

	opts := &USDAQueryOptions{Max: 100, Offset: 0, Sort: "n"}

	return c.GetUSDASearch(ctx, params, opts)
}
