package main

import (
	"fmt"
	"io"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// 1. Sadece POST isteklerini kabul edelim ki okunacak bir Body (veri) gelsin.
	if r.Method != http.MethodPost {
		http.Error(w, "Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	// 2. RAM'de sadece 1 KB'lık (1024 bayt) sabit bir alan ayırıyoruz.
	// 1 GB veri de gelse, sistem sadece bu buffer'ı tekrar tekrar kullanacak.
	buffer := make([]byte, 1024)
	totalBytesRead := 0
	chunkCount := 0

	fmt.Println("--- Yeni Istek Geldi, Agdan Okuma Basliyor ---")

	// 3. Veriyi parça parça okumak için döngü
	for {
		// r.Body'nin Read metodu, işletim sisteminden sıradaki veriyi buffer'a doldurur
		n, err := r.Body.Read(buffer)

		// Eğer en az 1 bayt bile okunduysa, bunu işleyelim
		if n > 0 {
			chunkCount++
			totalBytesRead += n

			// Konsola sadece okunan parça bilgisini basıyoruz.
			// Eğer verinin kendisini de görmek istersen: fmt.Println(string(buffer[:n]))
			fmt.Printf("[Parca %d] Okunan: %d bayt\n", chunkCount, n)
		}

		// 4. Hata veya bitiş kontrolü
		if err != nil {
			if err == io.EOF {
				fmt.Println("-> Isletim Sistemi Sinyali: Okunacak veri bitti (EOF).")
				break // Döngüden çık
			}
			// EOF dışında gerçek bir ağ kopması/hata durumu
			fmt.Println("-> Beklenmeyen okuma hatasi:", err)
			break
		}
	}

	// 5. İşimiz bittiğinde network soketini kapatarak bellek sızıntısını (leak) önlüyoruz
	r.Body.Close()

	fmt.Printf("--- Islem Tamamlandi. Toplam: %d bayt, %d parcada okundu ---\n\n", totalBytesRead, chunkCount)

	// İstemciye başarılı yanıt dönelim
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Veri basariyla akis (stream) halinde okundu."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buffer", handle)

	fmt.Println("Sunucu 8080 portunda basladi...")
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Sunucu baslatilamadi:", err)
	}
}
