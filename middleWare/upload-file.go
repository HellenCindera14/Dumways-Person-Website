package middleware

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		file, err := c.FormFile("input-image")

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer src.Close()

		tempFile, err := ioutil.TempFile("uploads", "image-*.png")

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer tempFile.Close()

		io.Copy(tempFile, src)

		data := tempFile.Name()
		filename := data[8:]

		c.Set("dataFile", filename)

		return next(c)
	}
}
