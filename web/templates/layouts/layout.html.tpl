{{ define "base" }}
    <!DOCTYPE html>
    {{ safe "<!-- Application Layout -->" }}
    <html>
        <head>
            <title>Poller</title>
        </head>
        <body>
            {{ template "content" . }}
        </body>
    </html>
{{ end }}