package reporter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"staj/vulnerability"
)

func SaveAsText(vulns []vulnerability.Vulnerability) error {
	if err := os.MkdirAll("reports", os.ModePerm); err != nil {
		return err
	}

	fileName := fmt.Sprintf("reports/report_%s.txt", time.Now().Format("20060102_150405"))
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range vulns {
		line := "-"
		if v.Line > 0 {
			line = fmt.Sprintf("%d", v.Line)
		}
		_, err := f.WriteString(fmt.Sprintf("Dosya: %s | Zaafiyet: %s | Satır: %s\n", v.File, v.Rule, line))
		if err != nil {
			return err
		}
	}

	fmt.Println("[+] TXT raporu kaydedildi:", fileName)
	return nil
}

func SaveAsJSON(vulns []vulnerability.Vulnerability) error {
	if err := os.MkdirAll("reports", os.ModePerm); err != nil {
		return err
	}

	fileName := fmt.Sprintf("reports/report_%s.json", time.Now().Format("20060102_150405"))
	data, err := json.MarshalIndent(vulns, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("[+] JSON raporu kaydedildi:", fileName)
	return nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<title>Zaafiyet Raporu</title>
<style>
  body { font-family: Arial, sans-serif; background: #f4f4f9; color: #333; }
  table { border-collapse: collapse; width: 100%; max-width: 900px; margin: 20px auto; }
  th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
  th { background-color: #4CAF50; color: white; }
  tr:nth-child(even) { background-color: #f9f9f9; }
</style>
</head>
<body>
<h2 style="text-align:center;">Zaafiyet Raporu</h2>
<table>
<thead>
<tr>
  <th>Dosya</th>
  <th>Zaafiyet</th>
  <th>Satır</th>
</tr>
</thead>
<tbody>
{{range .Vulns}}
<tr>
  <td>{{.File}}</td>
  <td>{{.Rule}}</td>
  <td>{{if eq .Line 0}}-{{else}}{{.Line}}{{end}}</td>
</tr>
{{end}}
</tbody>
</table>
<p style="text-align:center; font-size: small; color: #666;">Rapor oluşturulma zamanı: {{.Timestamp}}</p>
</body>
</html>`

type htmlData struct {
	Timestamp string
	Vulns     []vulnerability.Vulnerability
}

func SaveAsHTML(results []vulnerability.Vulnerability) error {
	if err := os.MkdirAll("reports", os.ModePerm); err != nil {
		return err
	}

	fname := fmt.Sprintf("reports/report_%s.html", time.Now().Format("20060102_150405"))
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	data := htmlData{
		Timestamp: time.Now().Format("02 Jan 2006 15:04:05"),
		Vulns:     results,
	}

	return tmpl.Execute(f, data)
}
