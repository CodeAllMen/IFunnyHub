package page

import (
	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/MobileCPX/IFunnyHub/models/content"
)

//单独item页面
func GetEcahItems(position, page int) []models.Item {
	_, data := content.GetAllItemlByPosition(position, page)
	return data
}

//页码
func GetPageNumList(pagenow, pagetotal int) []int {

	var pagelist []int
	if pagetotal < 7 {
		for i := 1; i <= pagetotal; i++ {
			pagelist = append(pagelist, i)
		}
		return pagelist
	} else {
		if pagenow == 1 {
			pagelist = append(pagelist, 1, 2, 3, 0, pagetotal-1, pagetotal)
		} else if pagenow == 2 {
			pagelist = append(pagelist, 1, 2, 3, 4, 0, pagetotal-1, pagetotal)
		} else if pagenow == 3 {
			pagelist = append(pagelist, 1, 2, 3, 4, 5, 0, pagetotal-1, pagetotal)
		} else if pagenow < pagetotal-2 {
			pagelist = append(pagelist, 1, 0, pagenow-1, pagenow, pagenow+1, pagenow+2, 0, pagetotal-1, pagetotal)
		} else {
			pagelist = append(pagelist, 1, 0, pagenow-1)
			for ; pagenow <= pagetotal; pagenow++ {
				pagelist = append(pagelist, pagenow)
			}
		}
		return pagelist
	}

}
