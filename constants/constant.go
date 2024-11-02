package constants

const (
	PENDING  = "PENDING"
	ACCEPTED = "ACCEPTED"
	REJECTED = "REJECTED"
)

const (
	ADMIN   = "ADMIN"
	CAREER  = "CAREER"
	COMPANY = "COMPANY"
)
const (
	emailTemplate = `
	<!DOCTYPE html>
	<html lang="vi">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
			}
			.container {
				max-width: 600px;
				margin: auto;
				background: #ffffff;
				padding: 20px;
				border-radius: 5px;
				box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
			}
			.header {
				text-align: center;
				padding: 10px 0;
			}
			.header h1 {
				color: #4a4a4a;
			}
			.footer {
				margin-top: 20px;
				text-align: center;
				font-size: 0.8em;
				color: #666666;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Cảm ơn bạn đã ứng tuyển!</h1>
			</div>
			<p>Xin chào [Tên Ứng Viên],</p>
			<p>Cảm ơn bạn đã ứng tuyển vào vị trí <strong>[Tên Vị Trí]</strong> tại [Tên Công Ty]. Chúng tôi rất vui mừng khi nhận được hồ sơ của bạn.</p>
			<p>Đội ngũ tuyển dụng của chúng tôi sẽ xem xét hồ sơ của bạn và sẽ liên hệ trong thời gian sớm nhất. Nếu bạn có bất kỳ câu hỏi nào, đừng ngần ngại liên hệ với chúng tôi qua email này.</p>
			<p>Chúc bạn một ngày tuyệt vời!</p>
			<p>Trân trọng,</p>
			<p><em>Đội ngũ tuyển dụng tại [Tên Công Ty]</em></p>
			<div class="footer">
				<p>[Tên Công Ty] | [Địa chỉ Công Ty] | [Số điện thoại]</p>
			</div>
		</div>
	</body>
	</html>
	`
)
