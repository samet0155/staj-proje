# Config Dosyalarında Zaafiyet Tespit Aracı (Go)

Bu proje, sistemdeki çeşitli yapılandırma (config) dosyalarını tarayarak potansiyel güvenlik zaafiyetlerini tespit eden bir Go uygulamasıdır.  
Özellikle **hardcoded password**, **PermitRootLogin** gibi kritik yapılandırma hatalarını bulup raporlamayı amaçlar.

---

## Özellikler

- Desteklenen config formatları: `.conf`, `.cfg`, `.env`, `.ini`, `.json`, `.yaml`, `.yml`  
- Hardcoded şifre, gizli anahtar ve riskli yapılandırmaları tespit eder.  
- Taramayı belirtilen dizinde recursive (alt dizinlerle birlikte) yapar.  
- Tespit edilen zaafiyetleri TXT, JSON ve şık HTML formatında raporlar.  
- Modüler ve kolay genişletilebilir yapı.  

---

## Teknolojiler

- Programlama Dili: Go (Golang)  
- Kullanılan paketler:  
  - `bufio`, `os`, `filepath` (dosya ve dizin işlemleri)  
  - `encoding/json` (JSON parsing)  
  - `gopkg.in/yaml.v3` (YAML parsing)  
  - `html/template` (HTML raporlama)  

---

## Kurulum

1. Go ortamını kurun: [https://golang.org/dl/](https://golang.org/dl/)  
2. Projeyi klonlayın veya indirin.  
3. Gerekli bağımlılıkları yükleyin:  
   ```bash
   go get gopkg.in/yaml.v3
