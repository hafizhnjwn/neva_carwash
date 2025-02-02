package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nevacarwash.com/main/middleware"
	"nevacarwash.com/main/repositories"
	"nevacarwash.com/main/services"
)

type VehicleHandler struct {
	service *services.VehicleService
}

func NewVehicleHandler(service *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "create.html", nil)
		return
	}

	var vehicle repositories.CreateVehicleRequest
	if err := c.ShouldBind(&vehicle); err != nil {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	claims := middleware.JwtClaims(c)
	if claims["id"] == nil {
		c.HTML(http.StatusUnauthorized, "create.html", gin.H{
			"Error": "Unauthorized",
		})
		return
	}
	idFloat, ok := claims["id"].(float64)
	if !ok {
		c.HTML(http.StatusUnauthorized, "create.html", gin.H{
			"Error": "Unauthorized",
		})
		return
	}
	vehicle.UID = fmt.Sprintf("%d", uint(idFloat))
	vehicleID, err := h.service.CreateVehicle(&vehicle)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "create.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/vehicles/%s", vehicleID))
}

func (h *VehicleHandler) GetVehiclesByUsername(c *gin.Context) {
	username := c.Param("username")

	// If username empty, get username from jwt claims
	if username == "" {
		claims := middleware.JwtClaims(c)
		if claims == nil {
			c.HTML(http.StatusUnauthorized, "mylist.html", gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		if usernameClaim, ok := claims["username"].(string); ok {
			username = usernameClaim
		} else {
			c.HTML(http.StatusUnauthorized, "mylist.html", gin.H{
				"Error": "Unauthorized",
			})
			return
		}
	}

	vehicles, err := h.service.GetVehiclesByUsername(username)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "mylist.html", gin.H{
		"vehicles": vehicles,
	})
}

func (h *VehicleHandler) GetVehiclesByProcess(c *gin.Context) {
	process := []string{"Menunggu", "Proses", "Selesai"}

	// Call the service to get the vehicles grouped by status
	groupedVehicles, err := h.service.GetVehiclesByProcess(process)
	if err != nil {
		// Handle error by showing it on the page
		c.HTML(http.StatusInternalServerError, "list.html", gin.H{"error": err.Error()})
		return
	}

	// Pass the grouped vehicles to the template
	c.HTML(http.StatusOK, "list.html", gin.H{
		"groupedVehicles": groupedVehicles,
	})
}

func (h *VehicleHandler) GetVehicleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Redirect(http.StatusSeeOther, "/vehicles")
		return
	}
	vehicle, err := h.service.GetVehicleByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/vehicles")
		return
	}
	var isLoggedIn bool
	claims := middleware.JwtClaims(c)
	if claims != nil {
		isLoggedIn = true
	}
	c.HTML(http.StatusOK, "viewvehicle.html", gin.H{
		"Name":             vehicle.Name,
		"Package":          vehicle.Package,
		"Username":         vehicle.User.Username,
		"Process":          vehicle.Process,
		"Contact":          vehicle.Contact,
		"Plat":             vehicle.Plat,
		"Date":             vehicle.Date,
		"EnterTime":        vehicle.EnterTime,
		"ID":               vehicle.ID,
		"IsOwner":          isLoggedIn,
		"EstimatedTime":    vehicle.EstimatedTime,
		"FinishTime":       vehicle.FinishTime,
		"ShowDeleteButton": false,
	})
}

func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	id := c.Param("id")

	// Show edit form for GET requests
	if c.Request.Method == http.MethodGet {
		vehicle, err := h.service.GetVehicleByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Check if user owns this vehicle
		claims := middleware.JwtClaims(c)
		if claims == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		if idFloat, ok := claims["id"].(float64); ok {
			currentUserID := uint(idFloat)
			if currentUserID != vehicle.UserID {
				c.HTML(http.StatusForbidden, "edit.html", gin.H{
					"Error": "Not authorized to edit this vehicle",
				})
				return
			}
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"ID":      vehicle.ID,
			"Name":    vehicle.Name,
			"Package": vehicle.Package,
			"Contact": vehicle.Contact,
			"Process": vehicle.Process,
			"Plat":    vehicle.Plat,
		})
		return
	}

	// Handle PUT request to update vehicle
	var updatedVehicle repositories.CreateVehicleRequest
	if err := c.ShouldBind(&updatedVehicle); err != nil {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"Error":   err.Error(),
			"Name":    updatedVehicle.Name,
			"Package": updatedVehicle.Package,
			"Contact": updatedVehicle.Contact,
			"Process": updatedVehicle.Process,
			"Plat":    updatedVehicle.Plat,
		})
		return
	}

	if err := h.service.UpdateVehicle(id, updatedVehicle); err != nil {
		c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
			"Error":   err.Error(),
			"ID":      id,
			"Name":    updatedVehicle.Name,
			"Package": updatedVehicle.Package,
			"Contact": updatedVehicle.Contact,
			"Process": updatedVehicle.Process,
			"Plat":    updatedVehicle.Plat,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/vehicles/%s", id))
}

func (h *VehicleHandler) DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	// Show delete confirmation for GET requests
	if c.Request.Method == http.MethodGet {
		vehicle, err := h.service.GetVehicleByID(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Check if user owns this vehicle
		claims := middleware.JwtClaims(c)
		if claims == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		if idFloat, ok := claims["id"].(float64); ok {
			currentUserID := uint(idFloat)
			if currentUserID != vehicle.UserID {
				c.HTML(http.StatusForbidden, "mylist.html", gin.H{
					"Error": "Not authorized to delete this vehicle",
				})
				return
			}
		}

		c.Redirect(http.StatusOK, "/vehicles/listed")
		return
	}

	// Handle DELETE request
	if err := h.service.DeleteVehicle(id); err != nil {
		c.HTML(http.StatusInternalServerError, "mylist.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicles/listed")
}
