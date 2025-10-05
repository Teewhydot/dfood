package service

import "fmt"

// EmailTemplates contains all email template functions
type EmailTemplates struct{}

// NewEmailTemplates creates a new instance of EmailTemplates
func NewEmailTemplates() *EmailTemplates {
	return &EmailTemplates{}
}

// WelcomeEmailTemplate generates welcome email content
func (t *EmailTemplates) WelcomeEmailTemplate(userName string) (subject, plainText, htmlContent string) {
	subject = "Welcome to dfood!"

	plainText = fmt.Sprintf(`Hi %s,

Welcome to dfood! We're excited to have you join our community of food lovers.

You can now:
- Discover amazing restaurants
- Share your dining experiences through reviews
- Help others find great places to eat

Start exploring and sharing your food journey with us!

Best regards,
The dfood Team`, userName)

	htmlContent = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to dfood</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">Welcome to dfood!</h1>
        
        <p>Hi %s,</p>
        
        <p>Welcome to dfood! We're excited to have you join our community of food lovers.</p>
        
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <h3 style="margin-top: 0;">You can now:</h3>
            <ul>
                <li>Discover amazing restaurants</li>
                <li>Share your dining experiences through reviews</li>
                <li>Help others find great places to eat</li>
            </ul>
        </div>
        
        <p>Start exploring and sharing your food journey with us!</p>
        
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`, userName)

	return subject, plainText, htmlContent
}

// PasswordResetEmailTemplate generates password reset email content
func (t *EmailTemplates) PasswordResetEmailTemplate(userName, resetToken string) (subject, plainText, htmlContent string) {
	subject = "Reset Your dfood Password"

	// You'll need to replace this with your actual frontend URL
	resetURL := fmt.Sprintf("https://dfood.com/reset-password?token=%s", resetToken)

	plainText = fmt.Sprintf(`Hi %s,

We received a request to reset your dfood password.

Click the link below to reset your password:
%s

This link will expire in 1 hour for security reasons.

If you didn't request this password reset, please ignore this email. Your password will remain unchanged.

Best regards,
The dfood Team`, userName, resetURL)

	htmlContent = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your dfood Password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">Reset Your Password</h1>
        
        <p>Hi %s,</p>
        
        <p>We received a request to reset your dfood password.</p>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #3498db; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
        </div>
        
        <p><strong>This link will expire in 1 hour</strong> for security reasons.</p>
        
        <p>If you didn't request this password reset, please ignore this email. Your password will remain unchanged.</p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        
        <p style="font-size: 12px; color: #666;">
            If the button doesn't work, copy and paste this link into your browser:<br>
            %s
        </p>
        
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`, userName, resetURL, resetURL)

	return subject, plainText, htmlContent
}

// ReviewNotificationTemplate generates review notification email for restaurant owners
func (t *EmailTemplates) ReviewNotificationTemplate(restaurantName, reviewerName string, rating float64, comment string) (subject, plainText, htmlContent string) {
	subject = fmt.Sprintf("New Review for %s", restaurantName)

	plainText = fmt.Sprintf(`Hello,

You have received a new review for %s!

Reviewer: %s
Rating: %.1f/5.0
Comment: %s

Log in to your dfood dashboard to respond to this review.

Best regards,
The dfood Team`, restaurantName, reviewerName, rating, comment)

	htmlContent = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>New Review for %s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2c3e50;">New Review!</h1>
        
        <p>Hello,</p>
        
        <p>You have received a new review for <strong>%s</strong>!</p>
        
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <p><strong>Reviewer:</strong> %s</p>
            <p><strong>Rating:</strong> %.1f/5.0 ⭐</p>
            <p><strong>Comment:</strong></p>
            <p style="font-style: italic;">"%s"</p>
        </div>
        
        <p>Log in to your dfood dashboard to respond to this review.</p>
        
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`, restaurantName, restaurantName, reviewerName, rating, comment)

	return subject, plainText, htmlContent
}

// OrderConfirmationTemplate generates order confirmation email
func (t *EmailTemplates) OrderConfirmationTemplate(userName, restaurantName, orderID string, totalAmount float64) (subject, plainText, htmlContent string) {
	subject = fmt.Sprintf("Order Confirmation - %s", orderID)

	plainText = fmt.Sprintf(`Hi %s,

Your order has been confirmed!

Order ID: %s
Restaurant: %s
Total Amount: $%.2f

Your food is being prepared and will be ready soon. You'll receive another email when it's ready for pickup/delivery.

Thank you for choosing dfood!

Best regards,
The dfood Team`, userName, orderID, restaurantName, totalAmount)

	htmlContent = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #27ae60;">Order Confirmed! 🎉</h1>
        
        <p>Hi %s,</p>
        
        <p>Your order has been confirmed!</p>
        
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <p><strong>Order ID:</strong> %s</p>
            <p><strong>Restaurant:</strong> %s</p>
            <p><strong>Total Amount:</strong> $%.2f</p>
        </div>
        
        <p>Your food is being prepared and will be ready soon. You'll receive another email when it's ready for pickup/delivery.</p>
        
        <p>Thank you for choosing dfood!</p>
        
        <p>Best regards,<br>The dfood Team</p>
    </div>
</body>
</html>`, userName, orderID, restaurantName, totalAmount)

	return subject, plainText, htmlContent
}
