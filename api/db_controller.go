package api

import (
	"archive/zip"
	"bytes"
	"strconv"
	"time"

	"bitfan/core"
	"github.com/gin-gonic/gin"
	"github.com/vjeantet/jodaTime"
)

type DatabaseController struct {
}

func (b *DatabaseController) Download(c *gin.Context) {

	// Create a new zip archive.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	zipFile, err := zipWriter.Create("bitfan.bolt.db")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	_, err = core.Storage().CopyTo(zipFile)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Make sure to check the error on Close.
	err = zipWriter.Close()
	if err != nil {
		c.String(500, err.Error())
		return
	}

	filename := jodaTime.Format("'bitfan_db_'YYYYMMdd-HHmmss'.zip'", time.Now())
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", strconv.Itoa(buf.Len()))

	c.Data(200, "application/zip", buf.Bytes())
}

func isalnum(c byte) bool {
	return (c >= '0' && c <= '9') || isletter(c)
}
func isletter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// Create a url-safe slug for fragments
func slugify(ins string) string {
	in := []byte(ins)
	if len(in) == 0 {
		return ins
	}
	out := make([]byte, 0, len(in))
	sym := false

	for _, ch := range in {
		if isalnum(ch) {
			sym = false
			out = append(out, ch)
		} else if sym {
			continue
		} else {
			out = append(out, '-')
			sym = true
		}
	}
	var a, b int
	var ch byte
	for a, ch = range out {
		if ch != '-' {
			break
		}
	}
	for b = len(out) - 1; b > 0; b-- {
		if out[b] != '-' {
			break
		}
	}
	return string(out[a : b+1])
}
