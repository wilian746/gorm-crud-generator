package product

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wilian746/go-generator/internal/rules/product"
	"github.com/wilian746/go-generator/pkg/repository/adapter"
	"github.com/wilian746/go-generator/pkg/repository/database"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("Should return instance of handler", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		repository := adapter.NewAdapter(conn)
		assert.NotEmpty(t, NewHandler(repository))
	})
}

func TestHandler_Options(t *testing.T) {
	t.Run("Should return no content when call options", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		r, _ := http.NewRequest(http.MethodOptions, "/product", nil)
		w := httptest.NewRecorder()
		h.Options(w, r)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestHandler_Get(t *testing.T) {
	t.Run("Should return ok when call getAll", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		rules := product.NewRules()
		rules.Migrate(conn, rules.GetMock())
		h := NewHandler(adapter.NewAdapter(conn))
		r, _ := http.NewRequest(http.MethodGet, "/product", nil)
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Should return internal_error without migrate when call getAll", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		r, _ := http.NewRequest(http.MethodGet, "/product", nil)
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Should return bad request when call getOne", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		h := NewHandler(adapter.NewAdapter(conn))
		r, _ := http.NewRequest(http.MethodGet, "/product/123", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", "123")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return notfound when call getOne", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		ID := uuid.New()
		rules := product.NewRules()
		productMock := rules.GetMock()
		productMock.ID = ID
		rules.Migrate(conn, productMock)
		h := NewHandler(adapter.NewAdapter(conn))
		r, _ := http.NewRequest(http.MethodGet, "/product/"+ID.String(), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Should return internal_error without migrate when call getOne", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		ID := uuid.New().String()
		r, _ := http.NewRequest(http.MethodGet, "/product/"+ID, nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", ID)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Should return OK when call getOne", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		ID := uuid.New()
		rules := product.NewRules()
		productMock := rules.GetMock()
		productMock.ID = ID
		rules.Migrate(conn, productMock)
		databaseAdapter := adapter.NewAdapter(conn)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodGet, "/product/"+ID.String(), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Get(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_Post(t *testing.T) {
	t.Run("Should return OK when call post", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		r, _ := http.NewRequest(http.MethodPost, "/product/", bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Post(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Should return bad request when call post; name is empty", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		productMock.Name = ""
		rules.Migrate(conn, productMock)
		r, _ := http.NewRequest(http.MethodPost, "/product/", bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Post(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return internal_error when call post; no exists table", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		r, _ := http.NewRequest(http.MethodPost, "/product/", bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Post(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("Should return bad request when call delete empty uuid", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		r, _ := http.NewRequest(http.MethodDelete, "/product/", nil)
		ctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Delete(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return bad request when call delete wrong uuid", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		r, _ := http.NewRequest(http.MethodDelete, "/product/123", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", "123")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Delete(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return not found when call delete", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		ID := uuid.New()
		r, _ := http.NewRequest(http.MethodDelete, "/product/"+ID.String(), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Delete(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Should return internal_error when call delete", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		h := NewHandler(adapter.NewAdapter(conn))
		ID := uuid.New()
		r, _ := http.NewRequest(http.MethodDelete, "/product/"+ID.String(), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Delete(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	t.Run("Should return NoContent when call delete", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodDelete, "/product/"+productMock.ID.String(), nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", productMock.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Delete(w, r)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestHandler_Put(t *testing.T) {
	t.Run("Should return OK when call put", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/"+productMock.ID.String(), bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", productMock.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
	t.Run("Should return bad request when call put; wrong uuid", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/"+"123", bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", "123")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return bad request when call put; empty uuid", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/", bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return bad request when call put; name is empty", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		productMock.Name = ""
		rules.Migrate(conn, productMock)
		_ = databaseAdapter.Create(nil, productMock, productMock.TableName())
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/"+productMock.ID.String(), bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", productMock.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Should return not found when call put;", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		rules.Migrate(conn, productMock)
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/"+productMock.ID.String(), bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", productMock.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Should return internal_error when call put;", func(t *testing.T) {
		conn := database.GetConnection("sqlite3", ":memory:")
		databaseAdapter := adapter.NewAdapter(conn)
		rules := product.NewRules()
		productMock := rules.GetMock()
		h := NewHandler(databaseAdapter)
		r, _ := http.NewRequest(http.MethodPut, "/product/"+productMock.ID.String(), bytes.NewReader(productMock.Bytes()))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("ID", productMock.ID.String())
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()
		h.Put(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
