package systemd

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/samber/lo"
)

type Repository struct{}

func NewRepository() Repository { return Repository{} }

func (r Repository) ListUserUnits(ctx context.Context, optFns ...filterOption) ([]dbus.UnitStatus, error) {
	var options filterOptions

	for _, fn := range optFns {
		fn(&options)
	}

	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create new sys connection: %w", err)
	}

	defer conn.Close()

	var units []dbus.UnitStatus

	if options.unitStatus != "" {
		units, err = conn.ListUnitsFilteredContext(ctx, []string{options.unitStatus})
		if err != nil {
			return nil, fmt.Errorf("unable to list sys units: %w", err)
		}
	} else {
		units, err = conn.ListUnitsContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to list sys units: %w", err)
		}
	}

	if options.unitType != "" {
		units = lo.Filter(units, func(unit dbus.UnitStatus, _ int) bool {
			return strings.HasSuffix(unit.Name, fmt.Sprintf(".%s", options.unitType))
		})
	}

	return units, nil
}

func (r Repository) ListSystemUnits(ctx context.Context, optFns ...filterOption) ([]dbus.UnitStatus, error) {
	var options filterOptions

	for _, fn := range optFns {
		fn(&options)
	}

	conn, err := dbus.NewUserConnectionContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create new user connection: %w", err)
	}

	defer conn.Close()

	var units []dbus.UnitStatus

	if options.unitStatus != "" {
		units, err = conn.ListUnitsFilteredContext(ctx, []string{options.unitStatus})
		if err != nil {
			return nil, fmt.Errorf("unable to list user units: %w", err)
		}
	} else {
		units, err = conn.ListUnitsContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to list user units: %w", err)
		}
	}

	if options.unitType != "" {
		units = lo.Filter(units, func(unit dbus.UnitStatus, _ int) bool {
			return strings.HasSuffix(unit.Name, fmt.Sprintf(".%s", options.unitType))
		})
	}

	return units, nil
}

type (
	filterOption func(o *filterOptions)

	filterOptions struct {
		unitStatus string
		unitType   string
	}
)

func WithUnitStatus(unitStatus string) filterOption {
	return func(o *filterOptions) {
		o.unitStatus = unitStatus
	}
}

func WithUnitType(unitType string) filterOption {
	return func(o *filterOptions) {
		o.unitType = unitType
	}
}
