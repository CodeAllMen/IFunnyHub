package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/utils"

	"github.com/NewTrident/iFunnyHub/app/models"
	"github.com/NewTrident/iFunnyHub/config/auth"
	"github.com/NewTrident/iFunnyHub/db"
)

// GetDB get DB from request
func GetDB(req *http.Request) *gorm.DB {
	if db := utils.GetDBFromRequest(req); db != nil {
		return db
	}
	return db.DB
}

// GetCurrentGuestUser get current user from request
func GetCurrentGuestUser(req *http.Request) *models.User {
	if currentUser, ok := auth.GuestAuth.GetCurrentUser(req).(*models.User); ok {
		return currentUser
	}
	return nil
}

type Pagenation struct {
	Total   int
	PerPage int
	Page    int
}

func (p Pagenation) Pages() int {
	return p.Total / p.PerPage
}

func (p Pagenation) HasPrev() bool {
	return p.Page > 1
}

func (p Pagenation) PrevNUM() int {
	return p.Page - 1
}

func (p Pagenation) HasNext() bool {
	return p.Page < p.Pages()
}

func (p Pagenation) NextNUM() int {
	return p.Page + 1
}

func (p Pagenation) IterPages(leftEdge, leftCurrent, rightCurrent, rightEdge int) []int {
	pages := []int{}
	last := 0
	for num := 1; num <= p.Pages()+1; num++ {
		if num <= leftEdge || (num > p.Page-leftCurrent-1 && num < p.Page+rightCurrent) || num > p.Pages()-rightEdge {
			if last+1 != num {
				pages = append(pages, 0)
				// continue
			}
			pages = append(pages, num)
			last = num
		}
	}
	return pages
}

func AddURLQuery(u string, key, value string) string {

	if strings.HasSuffix(u, "/") {
		return fmt.Sprintf("%s?%s=%s", u, key, value)
	} else if strings.Contains(u, "?") {
		return fmt.Sprintf("%s&%s=%s", u, key, value)
	} else if !strings.Contains(u, "/") && !strings.Contains(u, "?") {
		return fmt.Sprintf("%s/?%s=%s", u, key, value)
	}
	return fmt.Sprintf("%s?%s=%s", u, key, value)
}
