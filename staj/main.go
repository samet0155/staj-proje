package main

import (
	"fmt"
	"os"
	"staj/reporter"
	"staj/scanner"
)

func main() {
	target := "configs"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	fmt.Println("[*] Tarama başlatılıyor:", target)

	results, err := scanner.ScanDirectory(target)
	if err != nil {
		fmt.Println("Tarama hatası:", err)
		return
	}

	if len(results) == 0 {
		fmt.Println("[+] Hiç zaafiyet bulunamadı.")
	} else {
		fmt.Println("[!] Bulunan zaafiyetler:")
		for _, v := range results {
			line := "-"
			if v.Line > 0 {
				line = fmt.Sprintf("%d", v.Line)
			}
			fmt.Printf("Dosya: %s | Zaafiyet: %s | Satır: %s\n", v.File, v.Rule, line)
		}

		if err := reporter.SaveAsText(results); err != nil {
			fmt.Println("TXT rapor kaydedilemedi:", err)
		}
		if err := reporter.SaveAsJSON(results); err != nil {
			fmt.Println("JSON rapor kaydedilemedi:", err)
		}
		if err := reporter.SaveAsHTML(results); err != nil {
			fmt.Println("HTML rapor kaydedilemedi:", err)
		}
	}

	fmt.Println("[+] Tarama tamamlandı.")
}
