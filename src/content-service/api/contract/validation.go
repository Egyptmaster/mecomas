package contract

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
)

func (x *Id) Validate() error {
	details := make(map[string]string)
	if x.Id == "" {
		details["id"] = "must not be empty"
	}
	if len(details) == 0 {
		return nil
	}
	return convert[Id](details)
}

// convert takes the found validation issues and creates a gRPC status error including bad request details
func convert[T any](details map[string]string) error {
	if len(details) == 0 {
		return nil
	}
	st := status.Newf(codes.InvalidArgument, "validation error on %s", reflect.TypeFor[T]().Name())
	br := &errdetails.BadRequest{}
	for field, desc := range details {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		})
	}
	var err error
	st, err = st.WithDetails(br)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st.Err()
}
