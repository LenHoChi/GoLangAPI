package db

import (
	"database/sql"
	"DemoProject2/models"
	"fmt"
)

func (db Database) GetAllItems() (*models.ItemList, error) {
	list := &models.ItemList{}

	rows, err := db.Conn.Query("SELECT * FROM items")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Name, &item.Description)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}
func GetAllItems2(db Database) (*models.ItemList, error) {
	list := &models.ItemList{}

	rows, err := db.Conn.Query("SELECT * FROM items")
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Name, &item.Description)
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
	//row := db.Conn.QueryRow(query, itemName)
	// err := row.Scan(&item.Name, &item.Description)
	err := db.Conn.QueryRow(query, itemName).Scan(&item.Name, &item.Description)
	if err!= nil{
		if err == sql.ErrNoRows {
			fmt.Println("vao day")
			return item, ErrNoMatch
		}
		return item, err
	}
	return item, nil

	// switch err := row.Scan(&item.Name, &item.Description); err {
	// case sql.ErrNoRows:
	// 	return item, ErrNoMatch
	// default:
	// 	return item, err
	//}
}
func(db Database) Len(id string) bool {
	return true
}
func (db Database) DeleteItem(itemName int) error {
	query := `DELETE FROM items WHERE name = $1;`
	_, err := db.Conn.Exec(query, itemName)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateItem(itemName int, itemData models.Item) (models.Item, error) {
	item := models.Item{}
	query := `UPDATE items SET description=$1 WHERE name=$2`
	err := db.Conn.QueryRow(query, itemData.Description, itemName).Scan(&item.Name, &item.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return item, ErrNoMatch
		}
		return item, err
	}
	return item, nil

	// db.Conn.QueryRow(query, itemData.Description, itemName)
	// return item, nil
}
