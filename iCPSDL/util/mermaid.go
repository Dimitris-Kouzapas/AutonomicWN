package util

import "os/exec"

/*func CreateMermaidHTML(mermaid string) (string, error) {
    html := `
  	<!DOCTYPE html>
  	<html>
  	<head>
  	    <title>Mermaid Diagram</title>
  	    <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
  	    <script>mermaid.initialize({startOnLoad:true});</script>
  	</head>
  	<body>
  	    <div class="mermaid">`
    html += mermaid
    html += `
      </div>
    </body>
    </html>
    `
    fileName, err := CreateTmpFile( html, "html" )
    if err != nil { return "", err }
    return fileName, nil
}*/

func CreateMermaidSVG(mermaid string) (string, error) {
    fileName, err := CreateTmpFile( mermaid, "mmd" )
    if err != nil { return "", err }

    cmd1 := exec.Command("mmdc", "-i", fileName, "-o", fileName + ".svg")
    err = cmd1.Run()
    if err != nil { return "", err }
    return fileName + ".svg", nil
}
