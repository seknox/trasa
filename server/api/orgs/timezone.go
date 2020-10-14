package orgs

import (
	"github.com/sirupsen/logrus"
	"time"
)

func GetTimeLocation(orgID string) *time.Location {

	orgDetail, err := Store.Get(orgID)
	if err != nil {
		logrus.Error(err)
		return time.UTC
	}

	loc, err := time.LoadLocation(orgDetail.Timezone)
	if err != nil {
		logrus.Errorf("load location: %v", err)
		return time.UTC
	}
	return loc
}
