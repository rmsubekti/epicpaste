package handlers

import (
	"epicpaste/api/utils"
	"epicpaste/system/model"
	u "epicpaste/system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Add a new paste to the store
// @Description Test Description
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param   some_id     path    int     true        "Some ID"
// @Success 200 {string} string  "ok"
// @Router /string/{some_id} [get]
func CreatePaste(c *gin.Context) {
	var paste model.Paste
	user, _ := c.Get("user")

	if user == nil { // user should login first
		utils.JSONErr(http.StatusNotFound, c, "Unauthorized")
		return
	}

	if err := c.ShouldBindJSON(&paste); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}
	paste.Paster = user.(model.User)

	if err := paste.Create(); err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, paste)
}

func EditPaste(c *gin.Context) {
	var paste model.Paste
	user, _ := c.Get("user")

	if user == nil { // user should login first
		utils.JSONErr(http.StatusNotFound, c, nil)
		return
	}

	if err := c.ShouldBindJSON(&paste); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}

	paste.ID = c.Param("id")
	paste.CreatedBy = user.(model.User).ID
	if err := paste.Update(); err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, paste)
}

// DeletePaste handler delete paste data owned by current logged in user
func DeletePaste(c *gin.Context) {
	var paste model.Paste
	user, _ := c.Get("user")

	if user == nil { // user should login first
		utils.JSONErr(http.StatusNotFound, c, nil)
		return
	}

	paste.ID = c.Param("id")
	paste.CreatedBy = user.(model.User).ID
	if err := paste.Delete(); err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": "deleted", "paste": paste})
}

func ViewPaste(c *gin.Context) {
	var paste model.Paste
	user, _ := c.Get("user")

	if err := paste.Get(c.Param("id")); err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
		return
	}

	if !*paste.Public || (user != nil && paste.CreatedBy != user.(model.User).ID) {
		utils.JSONErr(http.StatusNotFound, c, nil)
		return
	}
	c.JSON(http.StatusOK, paste)
}

func ListPublicPaste(c *gin.Context) {
	var paginator u.Paginator
	var pastes model.Pastes

	if err := c.Bind(&paginator); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}

	if err := pastes.List(&paginator); err != nil {
		utils.JSONErr(http.StatusInternalServerError, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, paginator)
}
func UserPastes(c *gin.Context) {
	var paginator u.Paginator
	var pastes model.Pastes
	visitor, _ := c.Get("user")
	ownerId := c.Param("userId")

	if err := c.Bind(&paginator); err != nil {
		utils.JSONErr(http.StatusBadRequest, c, err.Error())
		return
	}

	if visitor != nil && visitor.(model.User).ID == ownerId {
		// all pastes by user
		if err := pastes.ListByUser(ownerId, false, &paginator); err != nil {
			utils.JSONErr(http.StatusInternalServerError, c, err.Error())
			return
		}
	} else {
		// list all public
		if err := pastes.ListByUser(ownerId, true, &paginator); err != nil {
			utils.JSONErr(http.StatusInternalServerError, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, paginator)
}
