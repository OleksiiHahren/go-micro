package main

import (
    "bytes"
    "github.com/vanng822/go-premailer/premailer"
    mail "github.com/xhit/go-simple-mail/v2"
    "html/template"
    "strconv"
    "time"
)

type Mail struct {
    Domain string
    Host string
    Port int
    Username string
    Password string
    Encryption string
    FromAdress string
    FromName string
}

type Message struct {
    From string
    FromName string
    To string
    Subject string
    Attachments []string
    Data any
    DataMap map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
    if msg.From == "" {
        msg.From = m.FromAdress
    }
    if msg.FromName == "" {
        msg.FromName = m.FromName
    }

    data := map[string]any {
        "message": msg.Data,
    }
    msg.DataMap = data
    formattedMessage, err := m.buildHTMLMessage(msg)
    if err!=nil {
        return err
    }
    plainMessage, err := m.buildPlainTextMessage(msg)
    if err!=nil {
        return err
    }
    server := mail.NewSMTPClient()
    server.Host = m.Host
    server.Port = m.Port
    server.Username = m.Username
    server.Password = m.Password
    server.Encryption = m.makeEncription(m.Encryption)
    server.KeepAlive = false
    server.ConnectTimeout = 10 *time.Second
    server.SendTimeout = 10 *time.Second

    smtpClient, err := server.Connect()
    if err!=nil {
        return  err
    }

    email := mail.NewMSG()
    email.SetFrom(msg.From).AddTo(msg.To)
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
    templateToRender := "./templates/mail.html.gohtml"

    t, err := template.New("email-html").ParseFiles(templateToRender)
    if err!=nil {
        return "", err
    }

    var tpl bytes.Buffer
    if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
        return "", err
    }
    formattedMessage := tpl.String()
    formattedMessage, err = m.inlineCSS(formattedMessage)
    if err !=nil {
        return "", err
    }
    return formattedMessage, err
}

func (m *Mail) inlineCSS(s string) (string, error) {
    options := premailer.Options{
    RemoveClasses: false,
    CssToAttributes: false,
    KeepBangImportant: true,
    }
    prem, err := premailer.NewPremailerFromString(s, &options)
    if err !=nil {
        return "", err
    }

    html, err := prem.Transform()
    if err !=nil {
        return "", err
    }
    return html, err
}

func (m *Mail ) buildPlainTextMessage(msg Message)  (string, error){
    templateToRender := "./templates/mail.plain.gohtml"

    t, err := template.New("email-plain").ParseFiles(templateToRender)
    if err!=nil {
        return "", err
    }

    var tpl bytes.Buffer
    if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
        return "", err
    }
    plainMessage := tpl.String()
    return plainMessage, err
}

func (m *Mail ) makeEncription(s string) mail.Encryption {
    switch s {
    case "tls":
        return mail.EncryptionSTARTTLS
    case "ssl":
        return mail.EncryptionSSL
    case "none", "":
        return mail.EncryptionNone
    default:
        return mail.EncryptionSTARTTLS
    }

}