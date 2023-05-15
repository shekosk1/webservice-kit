package mid

import (
	"context"
	"net/http"

	"github.com/shekosk1/webservice-kit/business/web/auth"
	webV1 "github.com/shekosk1/webservice-kit/business/web/v1"
	"github.com/shekosk1/webservice-kit/foundation/web"
	"go.uber.org/zap"
)

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er webV1.ErrorResponse
				var status int

				switch {
				case webV1.IsRequestError(err):
					re := webV1.GetRequestError(err)
					er = webV1.ErrorResponse{
						Error: re.Error(),
					}
					status = re.Status

				case auth.IsAuthError(err):
					er = webV1.ErrorResponse{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					er = webV1.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
