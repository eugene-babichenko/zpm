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
	LastCheckSuccess      bool
	UpdatesAvailable      []string
	InstallationsRequired []string
}

type Meta struct {
	// Time passed since the last check for updates.
	LastUpdateCheck time.Time
	// Indicates if all plugins were checked without errors for the last time.
	LastCheckSuccess bool
	// The number of available updates since the last check.
	UpdatesAvailable []string
	// The number of required installations since the last check.
	InstallationsRequired []string
}

func (m Meta) Marshal() ([]byte, error) {
	dto := metaDto{
		LastUpdateCheck:       m.LastUpdateCheck.Format(lastUpdateCheckLayout),
		LastCheckSuccess:      m.LastCheckSuccess,
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
		LastCheckSuccess:      dto.LastCheckSuccess,
		UpdatesAvailable:      dto.UpdatesAvailable,
		InstallationsRequired: dto.InstallationsRequired,
	}

	return result, nil
}
