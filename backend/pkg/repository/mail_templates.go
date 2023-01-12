package repository

const (
	MailTemplateResetPassword = `
	<html xmlns="http://www.w3.org/1999/xhtml">
	<head>
		<title>HTML email template</title>
		<meta name="viewport" content="width = 375, initial-scale = -1">
	</head>

	<body style="background-color: #ffffff; font-size: 16px;">
		<center>
		<table align="center" border="0" cellpadding="0" cellspacing="0" style="width:600px;">
			<!-- BEGIN EMAIL -->
			<tr>
				<td align="center" bgcolor="#ffffff" style="padding:30px">
				<p style="text-align:left">Hello,<br><br> We received a request to reset the password for your account for this email address. To initiate the password reset process for your account, click the link below.
				</p>
				<p>
					<a target="_blank" style="text-decoration:none; background-color: black; border: black 1px solid; color: #fff; padding:10px 10px; display:block;" href="%s">
					<strong>Reset Password</strong></a>
				</p>
				<br>
				<p style="text-align:left">
					If you did not make this request, you can simply ignore this email.
				</p>
				<p style="text-align:left">
					Sincerely,<br>
					QR Menu Team
				</p>
				</td>
			</tr>
			</tbody>
		</table>
		</center>
	</body>
	</html>
	` // #nosec
)
