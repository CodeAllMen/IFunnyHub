package page

import (
	"fmt"
	"github.com/MobileCPX/IFunnyHub/models"
	"github.com/MobileCPX/IFunnyHub/models/content"
)

// 首页四种item
func GetItemsIndex() map[string][]models.Item {
	fmt.Println("pos 1")
	_, video := content.GetFiveItemlByPosition(1)
	fmt.Println("pos 2")
	_, game := content.GetFiveItemlByPosition(2)
	fmt.Println("pos 3")
	_, picture := content.GetFiveItemlByPosition(3)
	fmt.Println("pos 4")
	_, ringtone := content.GetFiveItemlByPosition(4)
	return map[string][]models.Item{
		"Video":    video,
		"Picture":  picture,
		"Ringtone": ringtone,
		"Game":     game,
	}
}
