# otelslog
otelslog implementation to inject Datadog attributes to correlate tracing

Taken from Taken from https://darrenparkinson.uk/posts/2023-09-14-datadog-log-correlation-with-slog/

```
import (
    newslogin "github.com/dcatwoohoo/otelslog"
)

func main() {
    ...
    newslogin.DataDogFields = true

	logger := otelslog.NewLogger(os.Getenv("OTEL_SERVICE_NAME"))
	logger = newslogin.InitialiseLogging(logger.Handler())
}

```