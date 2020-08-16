{{ define "content" }}
<form action="{{mountpathed "recover/end"}}" method="POST">
    {{with .errors}}
        {{with (index . "")}}
            {{range .}}
                <span>{{.}}</span><br/>{{end}}
        {{end}}
    {{end -}}
    <input type="password" name="password" placeholder="Password" value=""/><br/>
    {{with .errors}}
        {{range .password}}
            <span>{{.}}</span><br/>{{end}}
    {{end -}}
    <input type="password" name="confirm_password" placeholder="Confirm Password" value=""/><br/>
    {{with .errors}}
        {{range .confirm_password}}
            <span>{{.}}</span><br/>{{end}}
    {{end -}}
    <button type="submit">Recover</button><br/>
    <a href="/">Cancel</a>
    <input type="hidden" name="token" value="{{.recover_token}}"/>
    {{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}"/>{{end -}}
</form>
{{ end }}