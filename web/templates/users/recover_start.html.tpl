{{ define "content" }}
<form action="{{mountpathed "recover"}}" method="POST">
    <input type="text" name="email" placeholder="E-mail"/><br/>
    {{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}"/>{{end -}}
    <button type="submit">Recover</button><br/>
    <a href="/login">Cancel</a>
</form>
{{ end }}