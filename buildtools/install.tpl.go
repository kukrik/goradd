//** This file was code generated by got. ***

package buildtools

import (
	"bytes"
)

func drawInstaller(buf *bytes.Buffer) (err error) {

	buf.WriteString(`<html>
<head>
</head>
<body>
    <h1>Goradd Installer</h1>
    <p>
    Installing...
    </p>
    <p>
    `)

	buf.WriteString(textToHtml(results))

	buf.WriteString(`
    </p>
    `)
	if stop {
		buf.WriteString(`<p><a href="/">Return to Start Page<a></p>`)
	}

	buf.WriteString(`</body>
</html>
`)

	if !stop {

		buf.WriteString(`<script>
setTimeout(function() {
    location = "/installer";  // Refresh
}, 2000);
</script>
`)

	}

	return
}
