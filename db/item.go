package db

import (
	"database/sql"
	"DemoProject2/models"
)

func (db Database) GetAllItems() (*models.ItemList, error) {
	list := &models.ItemList{}

	rows, err := db.Conn.Query("SELECT * FROM items")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Description, &item.Name)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}

func (db Database) AddItem(item *models.Item) error {
	query := `INSERT INTO items (name, description) VALUES ($1, $2)`
	db.Conn.QueryRow(query, item.Name, item.Description)
	return nil
}

func (db Database) GetItemById(itemName int) (models.Item, error) {
	item := models.Item{}

	query := `SELECT * FROM items WHERE name = $1;`
	row := db.Conn.QueryRow(query, itemName)
	switch err := row.Scan(&item.Name, &item.Description); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}
func(db Database) Len(id string) bool {
	return true
}
// func (db Database) DeleteItem(itemId int) error {
// 	query := `DELETE FROM items WHERE id = $1;`
// 	_, err := db.Conn.Exec(query, itemId)
// 	switch err {
// 	case sql.ErrNoRows:
// 		return ErrNoMatch
// 	default:
// 		return err
// 	}
// }

// func (db Database) UpdateItem(itemId int, itemData models.Item) (models.Item, error) {
// 	item := models.Item{}
// 	query := `UPDATE items SET name=$1, description=$2 WHERE id=$3 RETURNING id, name, description, created_at;`
// 	err := db.Conn.QueryRow(query, itemData.Name, itemData.Description, itemId).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return item, ErrNoMatch
// 		}
// 		return item, err
// 	}

// 	return item, nil
// }
