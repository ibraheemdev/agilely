<!DOCTYPE html>
{{ safe "<!-- Application Layout -->" }}
<html>
    <head>
        <title>Poller</title>
    </head>
    <body>
    {{with .flash_success}}<div class="alert alert-success">{{.}}</div>{{end}}
	{{with .flash_error}}<div class="alert alert-danger">{{.}}</div>{{end}}
    {{ template "content" . }}
    </body>
</html>