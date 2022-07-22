package db

import (
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/logs"
	"github.com/Scalingo/go-scalingo/v4"
)

type LogsOpts struct {
	Follow bool
	Count  int
}

// Logs displays the addon logs.
// app may be an app UUID or name.
// addon may be a addon UUID or an addon type (e.g. MongoDB).
func Logs(app, addon string, opts LogsOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	addonUUID := addon
	// If addon does not contain a UUID, we consider it contains an addon type (e.g. MongoDB)
	if !strings.HasPrefix(addon, "ad-") {
		addonUUID, err = getAddonUUIDFromType(c, app, addon)
		if err != nil {
			return errgo.Notef(err, "fail to get the addon UUID based on its type")
		}
	}

	url, err := c.AddonLogsURL(app, addonUUID)
	if err != nil {
		return errgo.Notef(err, "fail to get log URL")
	}

	err = logs.Dump(url, opts.Count, "")
	if err != nil {
		return errgo.Notef(err, "fail to dump logs")
	}

	if opts.Follow {
		err := logs.Stream(url, "")
		if err != nil {
			return errgo.Notef(err, "fail to stream logs")
		}
	}
	return nil
}

func getAddonUUIDFromType(addonsClient scalingo.AddonsService, app, addonType string) (string, error) {
	addons, err := addonsClient.AddonsList(app)
	if err != nil {
		return "", errgo.Notef(err, "fail to list the addons to get the type UUID")
	}

	for _, addon := range addons {
		if strings.EqualFold(addonType, addon.AddonProvider.Name) {
			return addon.ID, nil
		}
	}

	return "", errgo.Newf("no '%s' addon exists", addonType)
}
