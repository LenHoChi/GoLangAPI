package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"DemoProject2/db"
	"DemoProject2/models"

	"encoding/json"

)

var itemIDKey = "itemID"

func items(router chi.Router) {
	router.Get("/", getAllItems2)
	router.Post("/", createItem)

	router.Route("/{itemId}", func(router chi.Router) {
		router.Use(ItemContext)
		router.Get("/", getItem)
		router.Put("/", updateItem)
		router.Delete("/", deleteItem)
	})
}

func ItemContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "itemId")
		if itemId == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("item ID is required")))
			return
		}
		id, err := strconv.Atoi(itemId)
		fmt.Println("id la: ",id)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("invalid item ID")))
		}
		 ctx := context.WithValue(r.Context(), itemIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func getItem(w http.ResponseWriter, r *http.Request) {
	itemName := r.Context().Value(itemIDKey).(int)
	item, err := dbInstance.GetItemById(itemName)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
func getAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := dbInstance.GetAllItems()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, items); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
func getAllItems2(w http.ResponseWriter, r *http.Request) {
	items, err := db.GetAllItems2(dbInstance)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, items); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
func createItem(w http.ResponseWriter, r *http.Request) {
	item := &models.Item{}
	if err := render.Bind(r, item); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddItem(item); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}


func deleteItem(w http.ResponseWriter, r *http.Request) {
	itemName := r.Context().Value(itemIDKey).(int)
	err := dbInstance.DeleteItem(itemName)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	itemName := r.Context().Value(itemIDKey).(int)
	itemData := models.Item{}
	if err := render.Bind(r, &itemData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	item, err := dbInstance.UpdateItem(itemName, itemData)
	if err != nil {
		if err == db.ErrNoMatch {
			// render.Render(w, r, ErrNotFound)
			// render.Render(w, r, Ok(fmt.Errorf("ok ne!")))
			responseWithJSON(w, http.StatusOK, "Update user successfully")
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &item); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func responseWithJSON(response http.ResponseWriter, statusCode int, data interface{}){
	result, _ := json.Marshal(data)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}
