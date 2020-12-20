package dataservices

import "github.com/eltropycal/models/dbmodels"

func (pc *PostgresClient) GetFoodDetailsByID(id string) (dbmodels.Food, error) {
	query := `SELECT name, description, unit_rate_rupees, unit_time_mins from foods where id=$1`
	rows, err := pc.DB.Query(query, id)
	foodItem := dbmodels.Food{}
	if err != nil {
		return foodItem, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&foodItem.Name, &foodItem.Description, &foodItem.UnitRateRupees, &foodItem.UnitTimeMinutes)
		if err != nil {
			return foodItem, err
		}
		return foodItem, nil
	}
	return foodItem, nil
}
