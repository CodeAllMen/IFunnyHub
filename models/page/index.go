package page

import (
	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/MobileCPX/IFunnyHub/models/content"
)

//首页四种item
func GetItemsIndex() map[string][]models.Item {
	_, video := content.GetFiveItemlByPosition(1)
	_, game := content.GetFiveItemlByPosition(2)
	_, picture := content.GetFiveItemlByPosition(3)
	_, ringtone := content.GetFiveItemlByPosition(4)
	return map[string][]models.Item{
		"Video":    video,
		"Picture":  picture,
		"Ringtone": ringtone,
		"Game":     game,
	}
}
