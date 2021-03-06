package authentication

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates *template.Template

func initTemplates() {
	log.Println("initialising templates")
	templates = template.Must(template.New("providerSelection.gohtml").Funcs(template.FuncMap{
		"Title": strings.Title,
	}).Parse(`
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.12.1/css/all.min.css">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
		<script>
			function goto(url) {
			window.location.href = url
			}
		</script>
	</head>

	<body>
		<div style="min-height:100%;min-height:100vh;display:flex;align-items:center;">
			<div class="login-form" style="width:340px;margin:30px auto;">
				<form style="margin-bottom:15px;background:#f7f7f7;box-shadow:0px 2px 2px rgba(0,0,0,0. 3);padding:30px;">
					<div class="social-btn" style="">
						{{range .}}
							{{/* <a onclick="goto({{.}});return false" class="btn btn-primary btn-block" style="margin:10px 0;font-size:15px;text-align:left;line-height:36px;"><i class="fab fa-{{.}}" style="float:left;margin:4px 15px 0 5px;min-width:15px;font-size:28px;"></i>Authenticate using <b>{{. | Title}}</b></a> */}}
							<a href="{{.}}" class="btn btn-primary btn-block" style="margin:10px 0;font-size:15px;text-align:left;line-height:36px;"><i class="fab fa-{{.}}" style="float:left;margin:4px 15px 0 5px;min-width:15px;font-size:28px;"></i>Authenticate using <b>{{. | Title}}</b></a>
						{{end}}
					</div>
				</form>
			</div>
		</div>
	</body>
	`))
}

func providerSelectionTemplate(w http.ResponseWriter, r *http.Request, providers []string) {
	handleError(w, r, templates.ExecuteTemplate(w, "providerSelection.gohtml", Providers()))
}
