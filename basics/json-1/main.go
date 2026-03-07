package main

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	TransactionId     string  `json:"tx_id"`
	Amount            float64 `json:"amount,string"`
	Description       string  `json:"description,omitempty"`
	internalReference string  `json:"-"`
	IsRefunded        *bool   `json:"isRefunded"`
}

// BoolPtr Sürekli yeni değişken upholster adresini almak yerine helper function yazalım 1
func BoolPtr(boolean bool) *bool {
	return &boolean
}
func main() {
	tx0 := Transaction{
		TransactionId:     "21",
		Amount:            210,
		Description:       "",
		internalReference: "44*sad-da",
		IsRefunded:        BoolPtr(false), // burana neden true false diyemiyorum ? :
		/*,
		IsRefunded alanını *bool (boolean tutan bir bellek adresi / pointer) olarak tanımladın.
		Go'da true veya false gibi ham değerler (literaller) bellekte sabit bir adrese sahip değildir.
		Bu yüzden doğrudan adreslerini alamazsın; yani &true veya &false yazamazsın.
		Bu durum string veya integer literaller için de geçerlidir (örn: &"test" veya &42 yazamazsın).
		*/
	}

	marshal, err := json.MarshalIndent(tx0, "", "	")
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
	// out : 	{"tx_id":"21","amount":"210","isRefunded":null}

	incomingJSON := []byte(`{
		"tx_id": "99",
		"amount": "500",
		"description": "Iade islemi",
		"isRefunded": true,
		"bilinmeyen_alan": "bu alan structta yok"
	}`)

	var newTx Transaction

	err = json.Unmarshal(incomingJSON, &newTx) // Jsonu newTx'e cevirecek
	if err != nil {
		fmt.Printf("Unmarshalling err %v", err)
		return
	} //

	fmt.Printf("newTx %+v \n   boolean val := %v", newTx, *newTx.IsRefunded)
}
