package meta

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// The format to write the time of the last check for updates.
const lastUpdateCheckLayout string = "Mon Jan 2 15:04:05 2006"

// Internal data required for `zpm` functioning.
type metaDto struct {
	LastUpdateCheck       string
	UpdatesAvailable      int32
	InstallationsRequired int32
}

type Meta struct {
	// Time passed since the last check for updates.
	LastUpdateCheck time.Time
	// The number of available updates since the last check.
	UpdatesAvailable int32
	// The number of required installations since the last check.
	InstallationsRequired int32
}

func (m Meta) Marshal() ([]byte, error) {
	dto := metaDto{
		LastUpdateCheck:       m.LastUpdateCheck.Format(lastUpdateCheckLayout),
		UpdatesAvailable:      m.UpdatesAvailable,
		InstallationsRequired: m.InstallationsRequired,
	}

	result, err := json.Marshal(dto)
	if err != nil {
		errors.Wrap(err, "while encoding meta: ")
		return nil, err
	}

	return result, nil
}

func Unmarshal(data []byte) (*Meta, error) {
	var dto metaDto

	err := json.Unmarshal(data, &dto)
	if err != nil {
		errors.Wrap(err, "while decoding meta: ")
		return nil, err
	}

	lastUpdateCheck, err := time.Parse(lastUpdateCheckLayout, dto.LastUpdateCheck)
	if err != nil {
		errors.Wrap(err, "while decoding meta datetime: ")
		return nil, err
	}

	result := &Meta{
		LastUpdateCheck:       lastUpdateCheck,
		UpdatesAvailable:      dto.UpdatesAvailable,
		InstallationsRequired: dto.InstallationsRequired,
	}

	return result, nil
}
