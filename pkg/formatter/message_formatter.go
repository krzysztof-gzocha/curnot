package formatter

const htmlFormat = "text/html"
const plainTextFormat = "text/plain"

type MessageFormatStrategy interface {
	Format(message string) (contentType string, formattedMessage string)
}

type HtmlMessageFormatStrategy struct {
}

type TextMessageFormatStrategy struct {
}

func (f HtmlMessageFormatStrategy) Format(message string) (contentType string, formattedMessage string) {
	return htmlFormat, `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <title>Currnot message</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
<body>
<p>` + message + `</p>
</body>
</html> `
}

func (f TextMessageFormatStrategy) Format(message string) (contentType string, formattedMessage string) {
	return plainTextFormat, message
}
