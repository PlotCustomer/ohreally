package webshop

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"dvbd.com/webshop/internal/klant"
	"dvbd.com/webshop/internal/product"
)

// health reports that the Webshop application service is available.
func (s *Service) health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":    "ok",
		"component": "DVBD.com",
		"service":   "Webshop",
	})
}

// listProducts returns the Product DataObject as JSON.
func (s *Service) listProducts(c echo.Context) error {
	products, err := s.catalog(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, products)
}

// getProduct returns a single Product DataObject as JSON.
func (s *Service) getProduct(c echo.Context) error {
	p, err := s.config.ProductRepository.Get(c.Request().Context(), c.Param("id"))
	if errors.Is(err, product.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, p)
}

// getSelection returns the selection of the Klant as JSON.
func (s *Service) getSelection(c echo.Context) error {
	products, err := s.catalog(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, s.selection.Lines(klant.ID(c), products))
}

// addSelection selects a Product for the Klant (Selecteer producten).
func (s *Service) addSelection(c echo.Context) error {
	var body struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if body.ProductID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "productId is required")
	}
	if _, err := s.config.ProductRepository.Get(c.Request().Context(), body.ProductID); err != nil {
		if errors.Is(err, product.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	s.selection.Add(klant.ID(c), body.ProductID, body.Quantity)
	return c.JSON(http.StatusCreated, map[string]any{
		"productId": body.ProductID,
		"count":     s.selection.Count(klant.ID(c)),
	})
}

// deleteSelection removes a Product from the selection of the Klant.
func (s *Service) deleteSelection(c echo.Context) error {
	s.selection.Remove(klant.ID(c), c.Param("id"))
	return c.NoContent(http.StatusNoContent)
}

// home renders the Webshop catalog grouped by BusinessObject specialization.
func (s *Service) home(c echo.Context) error {
	products, err := s.catalog(c)
	if err != nil {
		return err
	}
	return s.render(c, http.StatusOK, "index.html", newCatalogView(products, s.selection.Count(klant.ID(c))))
}

// showSelection renders the selection of the Klant.
func (s *Service) showSelection(c echo.Context) error {
	products, err := s.catalog(c)
	if err != nil {
		return err
	}
	lines := s.selection.Lines(klant.ID(c), products)
	return s.render(c, http.StatusOK, "selection.html", newSelectionView(lines))
}

// selectProduct handles the Selecteer producten form.
func (s *Service) selectProduct(c echo.Context) error {
	productID := c.FormValue("productId")
	quantity, _ := strconv.Atoi(c.FormValue("quantity"))
	if productID != "" {
		s.selection.Add(klant.ID(c), productID, quantity)
	}
	return c.Redirect(http.StatusSeeOther, "/selectie")
}

// removeSelectedProduct handles removing a product from the selection.
func (s *Service) removeSelectedProduct(c echo.Context) error {
	s.selection.Remove(klant.ID(c), c.FormValue("productId"))
	return c.Redirect(http.StatusSeeOther, "/selectie")
}

func (s *Service) catalog(c echo.Context) ([]product.Product, error) {
	products, err := s.config.ProductRepository.List(c.Request().Context())
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return products, nil
}
