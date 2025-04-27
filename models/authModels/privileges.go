package authModels

import (
	"fmt"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/icons"
)

type PermissionCache map[int64][]model.Privileges
type LevelNameCache map[string][]int64

type PrivilegeLevelsPrivilegesPrimaryKey struct {
	PrivilegeID      int64
	PrivilegeLevelID int64
}

type JoinedPrivilegesRaw struct {
	LevelID            int64
	LevelName          string
	PrivilegeID        int64
	PrivilegeName      string
	PrivilegeCreatedAt *time.Time
}

type JoinedPrivilegeLevel struct {
	LevelID    int64
	LevelName  string
	Privileges []model.Privileges
}

func UserPrivilegeLevelJoinAsRowData(upl []UserPrivilegeLevelJoin) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(upl))
	for i, j := range upl {
		rows[i] = datadisplay.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.FirstName, models.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(strconv.FormatInt(j.PrivilegeLevelID, 10), models.MD),
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.CreatedAt.Format("2006-01-02"), models.MD),
				},
			},
		}
	}
	return rows
}

func JoinedPrivilegeLevelAsRowData(jpl []JoinedPrivilegeLevel) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, j := range jpl {
		rows[i] = datadisplay.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.LevelName, models.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(strconv.FormatInt(j.LevelID, 10), models.MD),
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.Privileges[0].CreatedAt.Format("2006-01-02"), models.MD),
				},
			},
		}
	}
	return rows
}

func JoinedPrivilegesAsRowData(jpl []JoinedPrivilegesRaw) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, p := range jpl {
		ca := p.PrivilegeCreatedAt
		caStr := "No Created At"
		if ca != nil {
			caStr = ca.String()
		}
		xAttrs := templ.Attributes{
			"class":      "fill-red-400 size-6 p-1 rounded-xs mx-auto cursor-pointer hover:bg-[#FFFFFF44]",
			"hx-delete":  fmt.Sprintf("/privilege-levels-privileges/level/%d/privilege/%d", p.LevelID, p.PrivilegeID),
			"hx-trigger": "click",
			"hx-swap":    "none",
			"_":          "on htmx:beforeRequest remove closest <tr/>",
		}
		rows[i].ID = "row-" + strconv.Itoa(i)
		rows[i].Data = []datadisplay.CellData{
			{
				ID:    "lvl-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(p.LevelName, models.SM),
			},
			{
				ID:    "pr-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(p.PrivilegeName, models.SM),
			},
			{
				ID:    "ca-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(caStr, models.MD),
			},
			{
				ID:    "del-" + strconv.Itoa(i),
				Width: 1,
				Body:  icons.X(xAttrs),
			},
		}
	}
	return rows
}
