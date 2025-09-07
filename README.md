# otelslog
otelslog implementation to inject Datadog attributes to correlate tracing

Taken from Taken from https://darrenparkinson.uk/posts/2023-09-14-datadog-log-correlation-with-slog/

I put this into a repository because I was using the otel slog logging bridge but DataDog needed its own attributes and some modifiers on the ids.  I found the above article but it seems the otel libraries broke that compatibility with what the author posted.

```
import (
    newslogin "github.com/dchapman992000/otelslog"
)

func main() {
    ...
    newslogin.DataDogFields = true

	logger := otelslog.NewLogger(os.Getenv("OTEL_SERVICE_NAME"))
	logger = newslogin.InitialiseLogging(logger.Handler())
}

```