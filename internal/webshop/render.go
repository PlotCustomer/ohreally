package webshop

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	"dvbd.com/webshop/internal/product"
	"dvbd.com/webshop/internal/selection"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func (s *Service) render(c echo.Context, status int, name string, data any) error {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(status)
	_, err := buf.WriteTo(c.Response())
	return err
}

type catalogView struct {
	SelectionCount int
	CategoryGroups []categoryGroup
}

type categoryGroup struct {
	Category string
	Products []productView
}

type productView struct {
	ID          string
	Name        string
	Category    string
	Description string
	Price       string
	Unit        string
	Stock       int
}

type selectionView struct {
	Count   int
	Total   string
	Entries []selectionEntry
}

type selectionEntry struct {
	Product   productView
	Quantity  int
	UnitPrice string
	Subtotal  string
}

func newCatalogView(products []product.Product, count int) catalogView {
	grouped := make(map[product.Category][]productView)
	for _, p := range products {
		grouped[p.Category] = append(grouped[p.Category], newProductView(p))
	}

	view := catalogView{SelectionCount: count}
	// Show the Product specializations in model order, then any other category.
	for _, category := range product.Categories {
		if items, ok := grouped[category]; ok {
			view.CategoryGroups = append(view.CategoryGroups, categoryGroup{Category: string(category), Products: items})
			delete(grouped, category)
		}
	}
	for _, p := range products {
		if items, ok := grouped[p.Category]; ok {
			view.CategoryGroups = append(view.CategoryGroups, categoryGroup{Category: string(p.Category), Products: items})
			delete(grouped, p.Category)
		}
	}
	return view
}

func newProductView(p product.Product) productView {
	return productView{
		ID:          p.ID,
		Name:        p.Name,
		Category:    string(p.Category),
		Description: p.Description,
		Price:       formatMoney(p.PriceCents, p.Currency),
		Unit:        p.Unit,
		Stock:       p.Stock,
	}
}

func newSelectionView(lines []selection.Line) selectionView {
	view := selectionView{Count: len(lines)}
	var total int64
	for _, line := range lines {
		subtotal := line.Product.PriceCents * int64(line.Quantity)
		total += subtotal
		view.Entries = append(view.Entries, selectionEntry{
			Product:   newProductView(line.Product),
			Quantity:  line.Quantity,
			UnitPrice: formatMoney(line.Product.PriceCents, line.Product.Currency),
			Subtotal:  formatMoney(subtotal, line.Product.Currency),
		})
	}
	view.Total = formatMoney(total, "EUR")
	return view
}

func formatMoney(cents int64, currency string) string {
	symbol := "€"
	if currency != "" && currency != "EUR" {
		symbol = currency + " "
	}
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%s%d,%02d", sign, symbol, cents/100, cents%100)
}
