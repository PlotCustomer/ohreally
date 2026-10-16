package webshop

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"dvbd.com/webshop/internal/product"
)

func newTestService() *Service {
	return New(Config{ProductRepository: product.NewMemoryRepository(product.SeedCatalog())})
}

func TestHealthReportsWebshop(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	recorder := httptest.NewRecorder()

	newTestService().Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "Webshop") {
		t.Fatalf("health body = %q, want it to mention the Webshop service", recorder.Body.String())
	}
}

func TestListProductsContainsBothSpecializations(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	recorder := httptest.NewRecorder()

	newTestService().Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("products status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var products []product.Product
	if err := json.Unmarshal(recorder.Body.Bytes(), &products); err != nil {
		t.Fatalf("decode products: %v", err)
	}
	if len(products) == 0 {
		t.Fatal("products is empty")
	}

	var fruit, pasta bool
	for _, p := range products {
		switch p.Category {
		case product.CategoryFruitGroenten:
			fruit = true
		case product.CategoryPasta:
			pasta = true
		}
	}
	if !fruit || !pasta {
		t.Fatalf("catalog must contain Fruit & Groenten and Pasta, got fruit=%v pasta=%v", fruit, pasta)
	}
}

func TestSelecteerProductenKeepsSelectionPerKlant(t *testing.T) {
	service := newTestService()
	handler := service.Handler()

	selectRequest := httptest.NewRequest(http.MethodPost, "/api/selection", strings.NewReader(`{"productId":"pasta-spaghetti","quantity":2}`))
	selectRequest.Header.Set("Content-Type", "application/json")
	selectRecorder := httptest.NewRecorder()
	handler.ServeHTTP(selectRecorder, selectRequest)

	if selectRecorder.Code != http.StatusCreated {
		t.Fatalf("select status = %d, want %d", selectRecorder.Code, http.StatusCreated)
	}
	cookies := selectRecorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("select response did not set the Klant cookie")
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/selection", nil)
	for _, cookie := range cookies {
		listRequest.AddCookie(cookie)
	}
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("selection status = %d, want %d", listRecorder.Code, http.StatusOK)
	}
	if !strings.Contains(listRecorder.Body.String(), "pasta-spaghetti") {
		t.Fatalf("selection body = %q, want it to contain the selected product", listRecorder.Body.String())
	}
}

func TestSelecterenFormRedirectsToSelection(t *testing.T) {
	form := url.Values{"productId": {"fruit-appel"}, "quantity": {"1"}}
	request := httptest.NewRequest(http.MethodPost, "/selecteren", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	newTestService().Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("selecteren status = %d, want %d", recorder.Code, http.StatusSeeOther)
	}
	if location := recorder.Header().Get("Location"); location != "/selectie" {
		t.Fatalf("selecteren location = %q, want /selectie", location)
	}
}

func TestSelectionPageRendersSelectedProduct(t *testing.T) {
	service := newTestService()
	handler := service.Handler()

	form := url.Values{"productId": {"pasta-penne"}, "quantity": {"3"}}
	selectRequest := httptest.NewRequest(http.MethodPost, "/selecteren", strings.NewReader(form.Encode()))
	selectRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	selectRecorder := httptest.NewRecorder()
	handler.ServeHTTP(selectRecorder, selectRequest)
	cookies := selectRecorder.Result().Cookies()

	pageRequest := httptest.NewRequest(http.MethodGet, "/selectie", nil)
	for _, cookie := range cookies {
		pageRequest.AddCookie(cookie)
	}
	pageRecorder := httptest.NewRecorder()
	handler.ServeHTTP(pageRecorder, pageRequest)

	if pageRecorder.Code != http.StatusOK {
		t.Fatalf("selection page status = %d, want %d", pageRecorder.Code, http.StatusOK)
	}
	body := pageRecorder.Body.String()
	for _, want := range []string{"Penne", "Pasta", "Totaal"} {
		if !strings.Contains(body, want) {
			t.Fatalf("selection page body does not contain %q", want)
		}
	}
}
