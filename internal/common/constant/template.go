package constant

const (

	// NewUserRegistrationEmailTemplate 新用户注册邮件模板
	//
	//	@author centonhuang
	//	@update 2026-02-04 16:19:02
	//	@update 2026-02-04 16:19:02
	NewUserRegistrationEmailTemplate = `
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
		.content { background-color: #f9f9f9; padding: 20px; border: 1px solid #ddd; border-radius: 0 0 5px 5px; }
		.user-info { background-color: white; padding: 15px; margin: 15px 0; border-left: 4px solid #4CAF50; }
		.user-info p { margin: 8px 0; }
		.label { font-weight: bold; color: #555; }
		.footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; font-size: 12px; color: #888; text-align: center; }
		.action-needed { background-color: #fff3cd; padding: 10px; border-left: 4px solid #ffc107; margin: 15px 0; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h2>New User Registration</h2>
		</div>
		<div class="content">
			<p>A new user has registered on the platform and is pending approval.</p>

			<div class="user-info">
				<p><span class="label">User Name:</span> %s</p>
				<p><span class="label">Email:</span> %s</p>
				<p><span class="label">Avatar:</span> <a href="%s" target="_blank">View Avatar</a></p>
				<p><span class="label">Registration Time:</span> %s</p>
				<p><span class="label">User ID:</span> %d</p>
			</div>

			<div class="action-needed">
				<p><strong>Action Required:</strong> This user currently has <strong>pending</strong> permission and cannot access the system until approved by an administrator.</p>
			</div>

			<p>Please log in to the admin panel to review and approve this user's registration.</p>
		</div>
		<div class="footer">
			<p>This is an automated notification from Aris Member API.</p>
			<p>Server: %s</p>
		</div>
	</div>
</body>
</html>
`
	// DocumentationTemplate 文档模板
	//
	//	@author centonhuang
	//	@update 2026-02-04 16:30:00
	DocumentationTemplate = `<!doctype html>
<html>
  <head>
    <title>Aris Mem API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
)
