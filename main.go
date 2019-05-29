package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/discordapp/lilliput"
	"github.com/valyala/fasthttp"
)

func main() {
	server()
}

// server starts the fasthttp server on port 8080.
func server() {
	fn := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			buf, err := extract(ctx)
			if err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}

			resized, err := resize(buf)
			if err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}

			if err := write(resized); err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}
			ctx.Write([]byte(`File compressed and resized successfully`))

		case "/upload":
			ctx.Response.Header.Set("Content-Type", "text/html")
			ctx.Write([]byte(`<form action="/" method="POST" enctype="multipart/form-data">
  <input type="file" name="image" multiple="true" accept="image/*" required/>
  <input type="submit" value="Upload image"/>
</form>
`))

		default:
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.Write([]byte("404 not found"))
		}
	}

	log.Println("[info] Server started on port 8080")
	fasthttp.ListenAndServe(":8080", fn)
}

// extract extracts an image from a fasthttp request body as a byte slice.
//
// I'm pretty sure this will break if the image is over some arbitrary limit
// that I'm not aware of.
func extract(ctx *fasthttp.RequestCtx) ([]byte, error) {
	form, err := ctx.FormFile("image")
	if err != nil {
		return nil, err
	}

	file, err := form.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, form.Size)
	_, err = file.Read(buf)
	return buf, err
}

// resize resizes an image.
func resize(buf []byte) ([]byte, error) {
	decoder, err := lilliput.NewDecoder(buf)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	header, err := decoder.Header()
	if err != nil {
		return nil, err
	}

	width := header.Width()
	height := header.Height()

	// Prepare to resize the image using 8192x8192 maximum resize buffer size.
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	opts := &lilliput.ImageOptions{
		FileType:             ".jpeg",
		Width:                width / 2,
		Height:               height / 2,
		ResizeMethod:         lilliput.ImageOpsResize,
		NormalizeOrientation: true,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: 85},
	}

	// Create a buffer to store the output image, 5MB in this case.
	out := make([]byte, 5*1024*1024)
	out, err = ops.Transform(decoder, opts, out)
	return out, err
}

// write writes the given bytes to a file.
func write(buf []byte) error {
	path := "images/compressed.jpeg"
	if err := ioutil.WriteFile(path, buf, 0400); err != nil {
		return err
	}
	log.Println("[info] File written to", path)
	return nil
}
