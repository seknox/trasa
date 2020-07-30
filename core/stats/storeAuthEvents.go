package stats

import (
	"database/sql"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/huandu/go-sqlbuilder"
	"github.com/seknox/trasa/core/misc"
	"github.com/seknox/trasa/utils"
)

//select sum(array_length(string_to_array(managed_accounts,','),1)) from servicesv1;

func (s StatStore) GetAggregatedLoginFails(entityType, entityID, orgID, timezone, timeFilter string) (reasons []failedReasonsByType, err error) {
	reasons = make([]failedReasonsByType, 0)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return reasons, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(`count(*)`, `failed_reason`)
	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))
	sb.Where(sb.Equal(`status`, false))
	sb.GroupBy(`failed_reason`)
	sb.OrderBy(`failed_reason`)
	sb.Limit(7)

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	rows, err := s.DB.Query(sqlStr, args...)

	//rows, err := s.DB.Query(`select count(*),failed_reason from auth_logs where status=false AND org_id=$1 GROUP BY failed_reason ORDER BY failed_reason`, orgID)
	if err != nil {
		return reasons, err
	}
	for rows.Next() {
		var reason failedReasonsByType
		err = rows.Scan(&reason.Value, &reason.Name)
		reason.Label = reason.Name
		if err != nil {
			return reasons, err
		}
		reasons = append(reasons, reason)
	}
	return reasons, err
}

func (s StatStore) GetAggregatedLoginHours(entityType, entityID, timezone, orgID, timeFilter, statusFilter string) (logins []loginsByHour, err error) {
	logins = make([]loginsByHour, 0)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return logins, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()

	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))
	sb.Select(sb.As(`EXTRACT('hour',timezone(timezone((login_time/1000000000)::int::timestamp,'UTC'),`+sb.Var(timezone)+`))`, `hour `), sb.As(`count(*)`, `c`))
	//sb.Select(sb.As(`EXTRACT('hour',(login_time/1000000000)::int::timestamp)`, `hour `), sb.As(`count(*)`, `c`))
	sb.GroupBy(`hour`)
	sb.OrderBy(`hour`)

	switch statusFilter {
	case "Success":
		sb.Where(sb.Equal(`status`, true))
	case "Failed":
		sb.Where(sb.Equal(`status`, false))
	}

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	rows, err := s.DB.Query(sqlStr, args...)

	if err != nil {
		return logins, err
	}
	for rows.Next() {

		//TODO break the logic into new function
		var login loginsByHour
		err = rows.Scan(&login.Hour, &login.Count)
		if err != nil {
			return logins, err
		}
		logins = append(logins, login)
	}
out:
	for i := 0; i < 24; i++ {

		for j := i; j >= 0; j-- {
			if len(logins) < i+1 {
				break
			}
			if strconv.Itoa(i) == logins[j].Hour {
				continue out
			}
		}

		temp1 := make([]loginsByHour, len(logins[:i]))
		copy(temp1, logins[:i])
		temp := append(temp1, loginsByHour{Hour: strconv.Itoa(i), Count: "0"})

		logins = append(temp, logins[i:]...)

	}

	return logins, err
}

func (s StatStore) GetAggregatedIPs(entityType, entityID, orgID, timezone, timeFilter, statusFilter string) (aggIps, error) {

	var ippool aggIps

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return ippool, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()

	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))
	sb.Select(`count(*)`, `user_ip`)
	sb.GroupBy(`user_ip`)
	sb.OrderBy(`user_ip`)

	switch statusFilter {
	case "Success":
		sb.Where(sb.Equal(`status`, true))
	case "Failed":
		sb.Where(sb.Equal(`status`, false))
	}

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	rows, err := s.DB.Query(sqlStr, args...)

	if err != nil {
		return ippool, err
	}

	//TODO move this logic outside

	var tmpArr []ipcount
	for rows.Next() {

		var tmp ipcount

		err = rows.Scan(&tmp.Count, &tmp.IP)
		if err != nil {
			return ippool, err
		}
		tmpArr = append(tmpArr, tmp)
	}

	ippool.Children = sortIps(tmpArr)

	ippool.Name = "TRASA"
	return ippool, nil
}

type ipcount struct {
	Count int
	IP    string
}

func sortIps(arr []ipcount) []firstOctet {
	var firstOctets = make([]firstOctet, 0)
	for _, t := range arr {

		var temp1, temp2, temp3, temp4 string
		userIP := t.IP
		count := t.Count

		splitted := strings.Split(userIP, ".")
		if len(splitted) != 4 {
			continue
		}
		temp1 = splitted[0]
		temp2 = splitted[1]
		temp3 = splitted[2]
		temp4 = splitted[3]

		fourth := fourthOctet{Key: temp4, Value: count, Name: temp1 + "." + temp2 + "." + temp3 + "." + temp4}
		third := thirdOctet{Key: temp3, Name: temp1 + "." + temp2 + "." + temp3 + ".0/24"}
		second := secondOctet{Key: temp2, Name: temp1 + "." + temp2 + ".0.0/16"}
		first := firstOctet{Key: temp1, Name: temp1 + ".0.0.0/8"}

		firstOctetEndIndex := len(firstOctets) - 1
		if (len(firstOctets) == 0) || (firstOctets[firstOctetEndIndex].Key != first.Key) {
			third.Children = append(third.Children, fourth)
			second.Children = append(second.Children, third)
			first.Children = append(first.Children, second)
			firstOctets = append(firstOctets, first)
			continue
		}

		secondOctetEndIndex := len(firstOctets[firstOctetEndIndex].Children) - 1
		if firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Key != second.Key {
			third.Children = append(third.Children, fourth)
			second.Children = append(second.Children, third)
			firstOctets[firstOctetEndIndex].Children = append(firstOctets[firstOctetEndIndex].Children, second)
			continue
		}

		thirdOctetEndIndex := len(firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children) - 1
		if firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children[thirdOctetEndIndex].Key != third.Key {
			third.Children = append(third.Children, fourth)
			firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children = append(firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children, third)
			continue
		}

		fourthOctetEndIndex := len(firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children[thirdOctetEndIndex].Children) - 1
		if firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children[thirdOctetEndIndex].Children[fourthOctetEndIndex].Key != fourth.Key {
			firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children[thirdOctetEndIndex].Children = append(firstOctets[firstOctetEndIndex].Children[secondOctetEndIndex].Children[thirdOctetEndIndex].Children, fourth)
			continue
		}

	}
	return firstOctets
}

func (s StatStore) GetLoginsByType(entityType, entityID, orgID, timezone, timeFilter, statusFilter string) (logins []nameValue, err error) {
	logins = make([]nameValue, 0)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return logins, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(`count(*)`, `COALESCE(service_type,'Other')`)
	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}

	if statusFilter == "success" {
		sb.Where(sb.Equal("status", true))
	} else if statusFilter == "failed" {
		sb.Where(sb.Equal("status", false))
	}
	sb.Where(sb.Equal(`org_id`, orgID))

	sb.GroupBy(`service_type`)

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	rows, err := s.DB.Query(sqlStr, args...)

	//rows, err := s.DB.Query(`select count(*),failed_reason from auth_logs where status=false AND org_id=$1 GROUP BY failed_reason ORDER BY failed_reason`, orgID)
	if err != nil {
		return logins, err
	}

	//TODO seperate function
	for rows.Next() {
		var reason nameValue
		err = rows.Scan(&reason.Value, &reason.Name)
		if err != nil {
			return logins, err
		}
		logins = append(logins, reason)
	}
out:
	for _, typ := range []string{"rdp", "ssh", "db", "dashboard", "http"} {
		for _, l := range logins {
			if l.Name == typ {
				continue out
			}
		}
		logins = append(logins, nameValue{Name: typ, Value: 0})
	}
	return logins, err
}

func (s StatStore) SortLoginByCity(entityType, entityID, orgID, timezone, timeFilter, statusFilter string) ([]geoDataType, error) {
	var dat geoDataType
	var docs []geoDataType = make([]geoDataType, 0)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return docs, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()

	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))
	sb.Select(`count(*)`, `user_ip`)
	sb.GroupBy(`user_ip`)
	sb.OrderBy(`user_ip`)

	switch statusFilter {
	case "Success":
		sb.Where(sb.Equal(`status`, true))
	case "Failed":
		sb.Where(sb.Equal(`status`, false))
	}

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	rows, err := s.DB.Query(sqlStr, args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ipStr string
		var count int64
		err := rows.Scan(&count, &ipStr)
		if err != nil {
			return docs, err
		}
		ip := net.ParseIP(ipStr)

		if utils.IsPrivateIP(ip) == false {
			loc, err := misc.Store.GetGeoLocation(ipStr)
			if err != nil {
				logrus.Debugf("get geo location: %v", err)
				continue
			}
			//country, city, coordinates := loc.
			dat.Name = loc.Country
			dat.City = loc.City
			dat.Value = count
			dat.Coordinates = loc.Location
			docs = append(docs, dat)

		}

	}
	return docs, err
	//docsStr, err := json.Marshal(docs)
	//return string(docsStr), err
}

func (s StatStore) GetAllAuthEventsByEntityType(entityType, entityID, timeFilter, timezone string) (totalEventsAuthEvents, error) {
	var events totalEventsAuthEvents

	// events.TotalLogins = &new(int64)
	// events.SuccessfulLogins = new(int64)
	// events.FailedLogins = new(int64)

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return events, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("count(status)", "status")
	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}

	*sb = addTimeFilter(timeFilter, now, *sb)

	sb.GroupBy("status")

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)

	rows, err := s.DB.Query(sqlStr, args...)
	if err != nil {
		return events, err
	}

	defer rows.Close()
	for rows.Next() {
		var count int64
		var status bool
		err = rows.Scan(&count, &status)
		if err != nil {
			return events, err
		}
		if status {
			events.SuccessfulLogins = count
		} else if !status {
			events.FailedLogins = count
		}
	}

	events.TotalLogins = events.SuccessfulLogins + events.FailedLogins

	return events, err
}

func (s StatStore) GetTodayHexaLoginEvents(entityType, entityID, orgID, statusFilter, timezone string) ([]todayHexa, error) {
	var hexas = make([]todayHexa, 0)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return hexas, err
	}
	now := time.Now().In(loc)
	todayStart := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute()))

	sb := sqlbuilder.NewSelectBuilder()

	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}
	sb.Where(sb.Equal(`org_id`, orgID))

	sb.Select(`org_id`, `email`, `service_name`, `service_id`, `login_time`, `status`)
	sb.Where(sb.GreaterThan(`login_time`, todayStart.UnixNano()))

	switch statusFilter {
	case "Success":
		sb.Where(sb.Equal(`status`, true))
	case "Failed":
		sb.Where(sb.Equal(`status`, false))
	}
	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)

	rows, err := s.DB.Query(sqlStr, args...)

	if err != nil {
		return hexas, err
	}

	for rows.Next() {
		var hexa todayHexa
		//var loginTime int64
		err := rows.Scan(&hexa.OrgID, &hexa.User, &hexa.AppName, &hexa.serviceID, &hexa.LoginTime, &hexa.Status)
		if err != nil {
			return hexas, err
		}
		hexas = append(hexas, hexa)
	}

	return hexas, nil

}

// GetTotalLoginsByDate generates time range based on current time and delta (start time for organization)
// and returns json array of total login events per day.
func (s StatStore) GetTotalLoginsByDate(entityType, entityID, orgID, timezone string) ([]totalEventsByDate, error) {
	var events []totalEventsByDate
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return events, err
	}
	now := time.Now().In(loc)

	limitTime := now.Add(-time.Hour * 24 * 46)

	var rows *sql.Rows

	if entityType == "service" {
		rows, err = s.DB.Query(`select count(*) as c,
       login_day,
       sum(case when status  then 1 else 0 end) as success,
       sum(case when status then 0 else 1 end) as fail
from (
    SELECT floor(login_time/86400000000000) as login_day ,status from auth_logs WHERE   login_time>$1 AND org_id=$2 AND service_id=$3
    )
GROUP BY login_day  ORDER BY login_day DESC LIMIT 45;`, limitTime.UnixNano(), orgID, entityID)

	} else if entityType == "user" {
		rows, err = s.DB.Query(`select count(*) as c,
       login_day,
       sum(case when status  then 1 else 0 end) as success,
       sum(case when status then 0 else 1 end) as fail
from (
    SELECT floor(login_time/86400000000000) as login_day ,status from auth_logs WHERE   login_time>$1 AND org_id=$2 AND service_id=$3
    )
GROUP BY login_day  ORDER BY login_day DESC LIMIT 45;`, limitTime.UnixNano(), orgID, entityID)

	} else {
		rows, err = s.DB.Query(`select count(*) as c,
       login_day,
       sum(case when status  then 1 else 0 end) as success,
       sum(case when status then 0 else 1 end) as fail
from (
    SELECT floor(login_time/86400000000000) as login_day ,status from auth_logs WHERE   login_time>$1 AND org_id=$2 
    )
GROUP BY login_day  ORDER BY login_day DESC LIMIT 45;`, limitTime.UnixNano(), orgID)

	}
	if err != nil {
		return events, err
	}
	var nextDay int64 = 0

	for rows.Next() {
		var event totalEventsByDate
		var day int64
		err = rows.Scan(&event.TotalLogins, &day, &event.SuccessfulLogins, &event.FailedLogins)
		if err != nil {
			return events, err
		}

		//TODO move this logic to different testable function
		event.Date = time.Unix(0, day*86400000000000).Format(time.RFC3339)
		event.TotalLogins = event.SuccessfulLogins + event.FailedLogins

		//Filling Missing days

		if nextDay != 0 {
			gap := nextDay - day - 1

			//		logger.Trace(gap)
			for i := gap; i > 0; i-- {
				events = append(events, totalEventsByDate{
					Date:             time.Unix(0, (day+i)*86400000000000).Format(time.RFC3339),
					TotalLogins:      0,
					SuccessfulLogins: 0,
					FailedLogins:     0,
				})
			}
		}

		events = append(events, event)
		nextDay = day
	}

	return events, err
}

func (s StatStore) GetRemoteAppCount(entityType, entityID, orgID, timezone, timeFilter, statusFilter string) (count int, err error) {

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return count, err
	}
	now := time.Now().In(loc)

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(`count(*)`)
	sb.From("auth_logs")
	if entityType == "service" {
		sb.Where(sb.Equal("service_id", entityID))
	} else if entityType == "user" {
		sb.Where(sb.Equal("user_id", entityID))
	}

	if statusFilter == "success" {
		sb.Where(sb.Equal("status", true))
	} else if statusFilter == "failed" {
		sb.Where(sb.Equal("status", false))
	}
	sb.Where(sb.Equal(`org_id`, orgID))

	sb.Where(sb.Equal(`endpoint`, "remoteApp"))
	sb.Where(sb.Equal(`service_type`, "rdp"))

	*sb = addTimeFilter(timeFilter, now, *sb)

	sqlStr, args := sb.Build()
	//change ? into $
	sqlStr = utils.SqlReplacer(sqlStr)
	err = s.DB.QueryRow(sqlStr, args...).Scan(&count)

	//rows, err := s.DB.Query(`select count(*),failed_reason from orgloginsv1 where status=false AND org_id=$1 GROUP BY failed_reason ORDER BY failed_reason`, orgID)
	if err != nil {
		return count, err
	}

	return count, err
}
