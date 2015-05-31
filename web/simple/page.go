package simple

import (
	"html/template"
)

var page = template.Must(template.New("simple.html").Parse(`<!DOCTYPE html>
<html lang="en">
<body onload='setFocus()'>
    <div id="story">{{ range . }}
        <p>{{ . }}</p>{{ end }}
    </div>
    <div id="input">
        <form action="run" id="f" method="POST">
            <input id="q" name="q"　type="text">
        </form>
    </div>
    <script>
			function setFocus(){
			    document.getElementById("q").focus();
			}
		</script>
</body>
</html>`))
