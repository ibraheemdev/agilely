{{ define "content" }}
<form action="{{mountpathed "register"}}" method="post">
	{{with .errors}}
		{{with (index . "")}}
			{{range .}}
				<span>{{.}}</span><br/>{{end}}
		{{end}}
	{{end -}}
	<div class="form-group">
		<label for="email">E-mail:</label>
		<input name="email" type="text" value="{{with .preserve}}{{with .email}}{{.}}{{end}}{{end}}" placeholder="E-mail"/><br/>
		{{with .errors}}
			{{range .email}}
				<span>{{.}}</span><br/>{{end}}
		{{end -}}
	</div>
	<label for="password">Password:</label>
	<input name="password" type="password" placeholder="Password"/><br/>
	{{with .errors}}
		{{range .password}}
			<span>{{.}}</span><br/>{{end}}
	{{end -}}
	<label for="confirm_password">Confirm Password:</label>
	<input name="confirm_password" type="password" placeholder="Confirm Password"/><br/>
	{{with .errors}}
		{{range .confirm_password}}
			<span>{{.}}</span><br/>{{end}}
	{{end -}}
	<input type="submit" value="Register"><br/>
	<a href="/">Cancel</a>

	{{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}"/>{{end}}
</form>
{{end}}
