package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"kvstore/internal/errors"
	"kvstore/internal/i18n"
	"kvstore/internal/store"

	"github.com/go-chi/chi/v5"
)

type KVRequest struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
}

type KVResponse struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
	TTL   int64       `json:"TTL"` // 新增字段
}

const (
	defaultPageSize = 10
	defaultPageNum  = 1
)

func RegisterKVRoutes(r chi.Router, s *store.Store) {
	r.Route("/kv", func(r chi.Router) {
		r.Get("/", listKeysHandler(s))
		r.Put("/{key}", setHandler(s))
		r.Get("/{key}", getHandler(s))
		r.Delete("/{key}", deleteHandler(s))
	})
}

func getHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		val, ttl, ok := s.GetWithTTL(key)
		if ok != nil {
			http.Error(w, i18n.Tr(errors.KV_ERR_KEY_NOT_FOUND), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(KVResponse{Key: key, Value: val, TTL: ttl})
	}
}

func setHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		var req KVRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, i18n.Tr(errors.KV_ERR_INVALID_JSON), http.StatusBadRequest)
			return
		}
		if req.TTL < 0 {
			http.Error(w, i18n.Tr(errors.KV_ERR_INVALID_TTL), http.StatusBadRequest)
			return
		}
		err := s.Set(key, req.Value, req.TTL)
		if err != nil {
			http.Error(w, i18n.Tr(errors.KV_ERR_SET_VALUE_FAILED), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(KVResponse{Key: key, Value: req.Value, TTL: req.TTL})
	}
}

func deleteHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if err := s.Delete(key); err != nil {
			http.Error(w, i18n.Tr(errors.KV_ERR_DELETE_FAILED), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func listKeysHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析分页参数
		pageStr := r.URL.Query().Get("page")
		sizeStr := r.URL.Query().Get("size")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = defaultPageNum
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = defaultPageSize
		}

		// 获取所有 keys
		keys, err := s.ListKeys()
		if err != nil {
			http.Error(w, i18n.Tr(errors.KV_ERR_LIST_KEYS_FAILED), http.StatusInternalServerError)
			return
		}

		// 计算分页
		total := len(keys)
		start := (page - 1) * size
		end := start + size
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}

		// 构建分页响应
		resp := map[string]interface{}{
			"data":  keys[start:end],
			"total": total,
			"page":  page,
			"size":  size,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
