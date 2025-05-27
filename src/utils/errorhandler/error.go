package errorhandler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	NotFoundError       = errors.New("Not Found")
	InternalServerError = errors.New("Internal server error")
	BadRequestError     = errors.New("Bad request error")
)

func TranslateTogRPCError(err error) error {
	if errors.Is(err, NotFoundError) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, BadRequestError) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
