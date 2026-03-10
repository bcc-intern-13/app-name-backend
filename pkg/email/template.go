package email

import "fmt"

func verificationEmailTemplate(toEmail, link string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:'Helvetica Neue',Helvetica,Arial,sans-serif;">
    <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:40px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">

                    <!-- Header -->
                    <tr>
                        <td align="center" style="background-color:#111827;padding:32px 40px;">
                            <img src="YOUR_LOGO_URL_HERE" alt="Logo" width="120" style="display:block;" />
                        </td>
                    </tr>

                    <!-- Body -->
                    <tr>
                        <td style="padding:40px 40px 24px;">
                            <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#111827;">
                                Verify your email address
                            </h1>
                            <p style="margin:0 0 24px;font-size:15px;line-height:1.6;color:#6b7280;">
                                Thanks for signing up! Please verify your email address by clicking the button below. This link will expire in <strong>24 hours</strong>.
                            </p>

                            <!-- Button -->
                            <table cellpadding="0" cellspacing="0" style="margin:0 0 32px;">
                                <tr>
                                    <td align="center" style="background-color:#111827;border-radius:6px;">
                                        <a href="%s"
                                           style="display:inline-block;padding:14px 32px;font-size:15px;font-weight:600;color:#ffffff;text-decoration:none;letter-spacing:0.3px;">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin:0 0 8px;font-size:13px;color:#9ca3af;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin:0;font-size:12px;color:#6b7280;word-break:break-all;">
                                <a href="%s" style="color:#111827;">%s</a>
                            </p>
                        </td>
                    </tr>

                    <!-- Divider -->
                    <tr>
                        <td style="padding:0 40px;">
                            <hr style="border:none;border-top:1px solid #e5e7eb;margin:0;" />
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="padding:24px 40px 32px;">
                            <p style="margin:0;font-size:12px;color:#9ca3af;line-height:1.6;">
                                If you did not create an account, you can safely ignore this email.<br/>
                                This email was sent to <strong>%s</strong>.
                            </p>
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, link, link, link, toEmail)
}
