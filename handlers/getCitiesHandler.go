package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"review-service/repository"
	"strconv"
	"time"
)

func (h *Handler) GetCitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nocache := r.URL.Query().Get("nocache") == "true"
	ctx := r.Context()

	if !nocache {
		cached, err := h.redis.Get(ctx, "cities").Result()

		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(cached)); err == nil {
				return
			}
			log.Printf("failed to write cache response: %v", err)
		}
		log.Println("No cache, going to DB")
	}

	cities, err := repository.GetAllCities(h.db)
	if err != nil {
		log.Println("get cities error", err.Error())
		http.Error(w, "get cities error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(cities)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.redis.Set(ctx, "cities", data, 24*time.Hour)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetCityByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cityID, err := strconv.ParseInt(r.URL.Query().Get("city_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect city_id", http.StatusBadRequest)
		return
	}

	city, err := repository.GetCityByID(h.db, cityID)
	if err != nil {
		log.Println("get city error", err.Error())
		http.Error(w, "get city error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(city)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
