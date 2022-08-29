package route

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func echoCORSConfig(_config *config.Config) echo.MiddlewareFunc {
	CORSConfig := middleware.DefaultCORSConfig //use default cors config
	//cors allow orign
	if 0 < len(_config.CORSConfig.AllowOrigins) {
		origins := strings.Split(_config.CORSConfig.AllowOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}

		CORSConfig.AllowOrigins = origins
	}
	//cors allow method
	if 0 < len(_config.CORSConfig.AllowMethods) {
		methods := strings.Split(_config.CORSConfig.AllowMethods, ",")
		for i := range methods {
			methods[i] = strings.TrimSpace(methods[i]) //trim space
			methods[i] = strings.ToUpper(methods[i])   //to upper
		}

		CORSConfig.AllowMethods = methods
	}

	fmt.Fprintf(os.Stdout, "ECHO CORS Config:\n")

	tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	tabwrite.Write([]byte(strings.Join([]string{
		"", "allow-origins",
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"-", strings.Join(CORSConfig.AllowOrigins, ", "),
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"", "allow-methods",
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"-", strings.Join(CORSConfig.AllowMethods, ", "),
	}, "\t") + "\n"))

	tabwrite.Flush()

	fmt.Fprintln(os.Stdout, strings.Repeat("_", 40))

	// fmt.Fprintf(os.Stdout, "-   allow-origins: %v\n", strings.Join(CORSConfig.AllowOrigins, ", "))
	// fmt.Fprintf(os.Stdout, "-   allow-methods: %v\n", strings.Join(CORSConfig.AllowMethods, ", "))
	// fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("_", 40))

	return middleware.CORSWithConfig(CORSConfig)

}
