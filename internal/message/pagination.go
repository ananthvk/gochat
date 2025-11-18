package message

import (
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

const defaultPageLimit = "30"

// This file implements cursor based pagination for efficient retrieval of messages
// Note: For now only backward pagination is supported (scroll up), forward pagination needs to be
// added after search/jump to message functionality is implemented

// The returned cursor will be the id of the last message the client has received, so in the subsequent request, messages < id will be returned

type Pagination struct {
	Before *ulid.ULID
	Limit  int
}

func readPagination(u url.Values) (Pagination, error) {
	before := u.Get("before")
	limit := u.Get("limit")
	if limit == "" {
		limit = defaultPageLimit
	}
	pagination := struct {
		Before string `validate:"omitempty,ulid"`
		Limit  string `validate:"number,min=1,max=100"`
	}{
		Before: before,
		Limit:  limit,
	}
	validator := validator.New(validator.WithRequiredStructEnabled())
	err := validator.Struct(pagination)
	if err != nil {
		return Pagination{}, err
	}
	i, err := strconv.Atoi(limit)
	if err != nil {
		return Pagination{}, err
	}
	var beforeId *ulid.ULID
	if pagination.Before == "" {
		beforeId = nil
	} else {
		id := ulid.MustParse(pagination.Before)
		beforeId = &id
	}
	return Pagination{
		Before: beforeId,
		Limit:  i,
	}, nil
}
