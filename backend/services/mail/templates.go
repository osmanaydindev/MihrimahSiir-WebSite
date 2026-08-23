package mail

import (
	"bytes"
	"html/template"
	"strings"
)

// Şablonlar html/template ile derleniyor: tüm alanlar otomatik escape
// edilir. Kullanıcı girdisi veya Open Library verisi asla template.HTML
// ile sarmalanmaz.
const baseLayout = `
<div style="font-family:-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;background:#f5f5f5;padding:24px">
  <div style="max-width:560px;margin:0 auto;background:#ffffff;border-radius:12px;overflow:hidden;border:1px solid #e0e0e0">
    <div style="background:#1e1e1e;color:#ffffff;padding:20px 24px">
      <h1 style="margin:0;font-size:18px;font-weight:600">MihrimahSiir</h1>
    </div>
    <div style="padding:24px;color:#212121;font-size:14px;line-height:1.6">
      {{template "content" .}}
    </div>
    <div style="padding:16px 24px;background:#fafafa;border-top:1px solid #eeeeee;color:#757575;font-size:12px">
      Bu e-posta MihrimahSiir tarafından otomatik gönderilmiştir.
    </div>
  </div>
</div>`

const bookDetailsBlock = `
<table style="width:100%;border-collapse:collapse;margin:16px 0">
  <tr><td style="padding:6px 0;color:#757575;width:120px">Kitap</td><td style="padding:6px 0;font-weight:600">{{.Title}}</td></tr>
  {{if .Authors}}<tr><td style="padding:6px 0;color:#757575">Yazar</td><td style="padding:6px 0">{{.Authors}}</td></tr>{{end}}
  <tr><td style="padding:6px 0;color:#757575">ISBN</td><td style="padding:6px 0">{{.ISBN}}</td></tr>
  {{if .Publisher}}<tr><td style="padding:6px 0;color:#757575">Yayınevi</td><td style="padding:6px 0">{{.Publisher}}</td></tr>{{end}}
  {{if .PublishDate}}<tr><td style="padding:6px 0;color:#757575">Yayın Tarihi</td><td style="padding:6px 0">{{.PublishDate}}</td></tr>{{end}}
  {{if .Pages}}<tr><td style="padding:6px 0;color:#757575">Sayfa</td><td style="padding:6px 0">{{.Pages}}</td></tr>{{end}}
</table>`

var (
	adminNewRequestTmpl = mustParse("adminNewRequest", `
{{define "content"}}
  <p style="margin:0 0 12px"><strong>{{.Username}}</strong> yeni bir kitap talebi oluşturdu.</p>
  {{if not .MetadataFound}}
    <p style="margin:0 0 12px;padding:10px 12px;background:#fff3e0;border-radius:6px;color:#e65100">
      Open Library'de bu ISBN için kayıt bulunamadı. Kitap bilgilerini panelden elle girmeniz gerekiyor.
    </p>
  {{end}}
  `+bookDetailsBlock+`
  {{if .CoverURL}}<img src="{{.CoverURL}}" alt="Kapak" style="max-width:140px;border-radius:6px;margin-bottom:12px" />{{end}}
  {{if .UserNote}}<p style="margin:0 0 12px"><span style="color:#757575">Kullanıcı notu:</span> {{.UserNote}}</p>{{end}}
  {{if .PanelURL}}<p style="margin:16px 0 0"><a href="{{.PanelURL}}" style="display:inline-block;background:#1e1e1e;color:#ffffff;text-decoration:none;padding:10px 18px;border-radius:6px">Panelde İncele</a></p>{{end}}
{{end}}`)

	userApprovedTmpl = mustParse("userApproved", `
{{define "content"}}
  <p style="margin:0 0 12px">Merhaba {{.Username}},</p>
  <p style="margin:0 0 12px">İstediğin kitap kütüphaneye eklendi.</p>
  `+bookDetailsBlock+`
  {{if .BookURL}}<p style="margin:16px 0 0"><a href="{{.BookURL}}" style="display:inline-block;background:#1e1e1e;color:#ffffff;text-decoration:none;padding:10px 18px;border-radius:6px">Kitaba Git</a></p>{{end}}
{{end}}`)

	userRejectedTmpl = mustParse("userRejected", `
{{define "content"}}
  <p style="margin:0 0 12px">Merhaba {{.Username}},</p>
  <p style="margin:0 0 12px">Kitap talebin bu kez onaylanmadı.</p>
  `+bookDetailsBlock+`
  {{if .AdminNote}}<p style="margin:0 0 12px"><span style="color:#757575">Gerekçe:</span> {{.AdminNote}}</p>{{end}}
{{end}}`)
)

func mustParse(name, content string) *template.Template {
	return template.Must(template.New(name).Parse(content + baseLayout))
}

// templateData, üç şablonun da beslendiği düz veri yapısı.
// *fiber.Ctx buraya asla girmez; worker handler'dan sonra çalışıyor.
type templateData struct {
	Username      string
	Title         string
	Authors       string
	ISBN          string
	Publisher     string
	PublishDate   string
	Pages         int
	CoverURL      string
	UserNote      string
	AdminNote     string
	PanelURL      string
	BookURL       string
	MetadataFound bool
}

func render(tmpl *template.Template, data templateData) (string, error) {
	// Kapak yalnızca Open Library CDN'inden geliyorsa gösterilir.
	if !strings.HasPrefix(data.CoverURL, "https://covers.openlibrary.org/") {
		data.CoverURL = ""
	}
	if data.Title == "" {
		data.Title = "(başlık bulunamadı)"
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
