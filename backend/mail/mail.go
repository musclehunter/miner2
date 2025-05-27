package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

// MailConfig はメール送信の設定
type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

// DefaultConfig はデフォルトのメール設定を返します
func DefaultConfig() MailConfig {
	// 環境変数から設定を取得（なければデフォルト値を使用）
	host := getEnv("MAIL_HOST", "smtp.gmail.com")
	port := getEnvAsInt("MAIL_PORT", 587)
	username := getEnv("MAIL_USERNAME", "")
	password := getEnv("MAIL_PASSWORD", "")
	fromName := getEnv("MAIL_FROM_NAME", "Muscle Miner")
	fromAddr := getEnv("MAIL_FROM_ADDRESS", "noreply@example.com")

	return MailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		FromName: fromName,
		FromAddr: fromAddr,
	}
}

// getEnv は環境変数を取得（なければデフォルト値を使用）
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt は環境変数を整数として取得
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue
	}
	return result
}

// Sender はメール送信を行うインターフェース
type Sender interface {
	Send(to string, subject string, body string) error
	SendHTML(to string, subject string, htmlBody string) error
}

// SMTPSender はSMTPを使用してメールを送信する
type SMTPSender struct {
	Config MailConfig
}

// NewSMTPSender は新しいSMTPSenderを作成
func NewSMTPSender(config MailConfig) *SMTPSender {
	return &SMTPSender{Config: config}
}

// Send はテキストメールを送信
func (s *SMTPSender) Send(to string, subject string, body string) error {
	// メールメッセージを構築
	msg := buildMessage(s.Config.FromName, s.Config.FromAddr, to, subject, body, "text/plain")

	// SMTPサーバーに接続
	auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Host)
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	// メールを送信
	err := smtp.SendMail(addr, auth, s.Config.FromAddr, []string{to}, msg)
	if err != nil {
		log.Printf("メール送信エラー: %v", err)
		return err
	}

	log.Printf("メール送信成功: %s へ '%s'", to, subject)
	return nil
}

// SendHTML はHTMLメールを送信
func (s *SMTPSender) SendHTML(to string, subject string, htmlBody string) error {
	// メールメッセージを構築
	msg := buildMessage(s.Config.FromName, s.Config.FromAddr, to, subject, htmlBody, "text/html")

	// SMTPサーバーに接続
	auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Host)
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)

	// メールを送信
	err := smtp.SendMail(addr, auth, s.Config.FromAddr, []string{to}, msg)
	if err != nil {
		log.Printf("HTMLメール送信エラー: %v", err)
		return err
	}

	log.Printf("HTMLメール送信成功: %s へ '%s'", to, subject)
	return nil
}

// buildMessage はメールメッセージを構築
func buildMessage(fromName, fromAddr, to, subject, body, contentType string) []byte {
	var buf bytes.Buffer

	// ヘッダー
	fmt.Fprintf(&buf, "From: %s <%s>\r\n", fromName, fromAddr)
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: %s; charset=UTF-8\r\n", contentType)
	fmt.Fprintf(&buf, "\r\n")

	// 本文
	fmt.Fprintf(&buf, "%s\r\n", body)

	return buf.Bytes()
}

// MockSender はテスト用のメール送信モック
type MockSender struct {
	SentMessages []MockMessage
}

// MockMessage はモックで送信されたメッセージ
type MockMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// NewMockSender は新しいMockSenderを作成
func NewMockSender() *MockSender {
	return &MockSender{SentMessages: []MockMessage{}}
}

// Send はモックのテキストメール送信
func (s *MockSender) Send(to string, subject string, body string) error {
	s.SentMessages = append(s.SentMessages, MockMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
	log.Printf("モックメール送信: %s へ '%s'", to, subject)
	return nil
}

// SendHTML はモックのHTMLメール送信
func (s *MockSender) SendHTML(to string, subject string, body string) error {
	s.SentMessages = append(s.SentMessages, MockMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
	log.Printf("モックHTMLメール送信: %s へ '%s'", to, subject)
	return nil
}

// EmailVerificationData はメール確認テンプレートのデータ
type EmailVerificationData struct {
	UserName string
	VerificationURL string
}

// SendVerificationEmail はメール確認メールを送信
func SendVerificationEmail(sender Sender, to string, userName string, token string, baseURL string) error {
	subject := "メールアドレスの確認"
	
	// テンプレートデータ
	data := EmailVerificationData{
		UserName: userName,
		VerificationURL: fmt.Sprintf("%s/verify-email?token=%s", baseURL, token),
	}
	
	// HTMLメールを生成
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>メールアドレスの確認</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; color: #333; }
        .container { max-width: 600px; margin: 0 auto; background-color: #f9f9f9; padding: 20px; border-radius: 5px; }
        .header { text-align: center; margin-bottom: 20px; }
        .button { display: inline-block; background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
        .footer { margin-top: 30px; font-size: 12px; color: #777; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>メールアドレスの確認</h2>
        </div>
        <p>こんにちは、{{.UserName}}さん</p>
        <p>アカウント登録いただきありがとうございます。以下のボタンをクリックして、メールアドレスの確認を完了してください。</p>
        <p style="text-align: center;">
            <a href="{{.VerificationURL}}" class="button">メールアドレスを確認する</a>
        </p>
        <p>または、以下のURLをブラウザに貼り付けてアクセスしてください。</p>
        <p>{{.VerificationURL}}</p>
        <p>このリンクは24時間有効です。</p>
        <div class="footer">
            <p>このメールは自動送信されています。返信はできませんのでご了承ください。</p>
        </div>
    </div>
</body>
</html>
`
	
	// テンプレートをパース
	tmpl, err := template.New("verification").Parse(htmlTemplate)
	if err != nil {
		log.Printf("テンプレートパースエラー: %v", err)
		return err
	}
	
	// テンプレートを実行
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("テンプレート実行エラー: %v", err)
		return err
	}
	
	// HTMLメールを送信
	return sender.SendHTML(to, subject, buf.String())
}
