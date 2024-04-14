package handler

import (
	"avito-banner/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Servicer interface {
	GetUserBanner(ctx context.Context, token string, tagId, featureId int, useLastRevision bool) (domain.ResponseUserBanner, error)
	GetBanner(ctx context.Context, tagId, featureId, limit, offset int) ([]domain.ResponseBanner, error)
	PostBanner(ctx context.Context, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) (int, error)
	PatchBanner(ctx context.Context, bannerId int, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) error
	DeleteBanner(ctx context.Context, bannerId int) error
}

func GetUserBanner(service Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if len(token) == 0 {
			http.Error(w, "no token", http.StatusBadRequest)
			return
		}

		tagIdStr := r.URL.Query().Get("tag_id")
		if len(tagIdStr) == 0 {
			http.Error(w, "no tag id", http.StatusBadRequest)
			return
		}

		featureIdStr := r.URL.Query().Get("feature_id")
		if len(featureIdStr) == 0 {
			http.Error(w, "no feature id", http.StatusBadRequest)
			return
		}

		useLastRevisionStr := r.URL.Query().Get("use_last_revision")
		if len(useLastRevisionStr) == 0 {
			useLastRevisionStr = "false"
		}

		tagId, errParse := strconv.Atoi(tagIdStr)
		if errParse != nil {
			http.Error(w, "tag_id must be integer", http.StatusBadRequest)
			return
		}

		featureId, errParse := strconv.Atoi(featureIdStr)
		if errParse != nil {
			http.Error(w, "feature_id must be integer", http.StatusBadRequest)
			return
		}

		useLastRevision, errParse := strconv.ParseBool(useLastRevisionStr)
		if errParse != nil {
			http.Error(w, "useLastRevision must be boolean", http.StatusBadRequest)
			return
		}

		response, errServ := service.GetUserBanner(context.Background(), token, tagId, featureId, useLastRevision)
		if errServ != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errServ.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if _, err := w.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	}
}

func GetBanner(service Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token != "admin_token" {
			http.Error(w, "token should be admin_token", http.StatusBadRequest)
			return
		}

		tagIdStr := r.URL.Query().Get("tag_id")
		featureIdStr := r.URL.Query().Get("feature_id")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		tagId, errParse := strconv.Atoi(tagIdStr)
		if errParse != nil {
			tagId = -1
		}
		featureId, errParse := strconv.Atoi(featureIdStr)
		if errParse != nil {
			featureId = -1
		}
		limit, errParse := strconv.Atoi(limitStr)
		if errParse != nil {
			limit = -1
		}
		offset, errParse := strconv.Atoi(offsetStr)
		if errParse != nil {
			offset = -1
		}

		response, errServ := service.GetBanner(context.Background(), tagId, featureId, limit, offset)
		if errServ != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errServ.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if _, err := w.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	}
}

func PostBanner(service Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token != "admin_token" {
			http.Error(w, "token should be admin_token", http.StatusBadRequest)
			return
		}

		// Проверяем Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "error: not application/json in Content-Type", http.StatusBadRequest)
			return
		}

		// Читаем тело запроса
		var PostRequestBody domain.PostRequestBody
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&PostRequestBody); err != nil {
			http.Error(w, "error with decoding content", http.StatusBadRequest)
			return
		}

		response, errServ := service.PostBanner(context.Background(),
			PostRequestBody.TagIds,
			PostRequestBody.FeatureID,
			PostRequestBody.Content,
			PostRequestBody.IsActive)

		if errServ != nil {
			//TODO какие заголовки?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errServ.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if _, err := w.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

	}
}

func PatchBannerId(service Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Преобразование bannerId в int64
		bannerId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		token := r.Header.Get("token")
		if token != "admin_token" {
			http.Error(w, "token should be admin_token", http.StatusBadRequest)
			return
		}

		// Проверяем Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "error: not application/json in Content-Type", http.StatusBadRequest)
			return
		}

		// Читаем тело запроса
		var PostRequestBody domain.PostRequestBody
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&PostRequestBody); err != nil {
			http.Error(w, "error with decoding content", http.StatusBadRequest)
			return
		}

		errServ := service.PatchBanner(context.Background(),
			bannerId,
			PostRequestBody.TagIds,
			PostRequestBody.FeatureID,
			PostRequestBody.Content,
			PostRequestBody.IsActive)

		if errServ != nil {
			//TODO какие заголовки?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errServ.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteBannerId(service Servicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Преобразование bannerId в int64
		bannerId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		token := r.Header.Get("token")
		if token != "admin_token" {
			http.Error(w, "token should be admin_token", http.StatusBadRequest)
			return
		}

		errServ := service.DeleteBanner(context.Background(), bannerId)

		if errServ != nil {
			//TODO какие заголовки?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errServ.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
