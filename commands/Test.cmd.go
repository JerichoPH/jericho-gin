package commands

import (
	"fmt"
	"time"

	"jericho-gin/database"
	"jericho-gin/tools"
	"jericho-gin/wrongs"

	uuid "github.com/satori/go.uuid"
)

// TestCmd 测试用
type TestCmd struct{}

// NewTestCmd 构造函数
func NewTestCmd() *TestCmd {
	return &TestCmd{}
}

func (receiver TestCmd) uuid() []string {
	c := make(chan string)
	go func(c chan string) {
		uuidStr := uuid.NewV4().String()
		c <- uuidStr
	}(c)
	go tools.NewTimer(5).Ticker()
	return []string{<-c}
}

func (receiver TestCmd) ls() []string {
	_, res := (&tools.Cmd{}).Process("ls", "-la")
	return []string{res}
}

func (receiver TestCmd) redis() []string {
	if _, err := database.NewRedis(0).SetValue("test", "AAA", 15*time.Minute); err != nil {
		wrongs.ThrowForbidden(err.Error())
	}

	for i := 0; i < 100000; i++ {
		if val, err := database.NewRedis(0).GetValue("test"); err != nil {
			wrongs.ThrowForbidden(err.Error())
		} else {
			fmt.Println(i, val)
		}
	}

	return []string{""}
}

func (receiver TestCmd) t() []string {
	var (
		// 电务段器材
		equipmentUniqueCodes = make([]string, 0)

		// 检修车间器材
		entireInstanceIdentityCodes = make([]string, 0)
	)

	std := tools.NewStdoutHelper("写入excel文件")

	conn1 := database.NewGormLauncher().NewConn("bi-b051-rnv")
	conn2 := database.NewGormLauncher().NewConn("fix-b051-rnv")

	conn1.Table("equipments as e").
		Select("e.unique_code").
		Joins("join equipment_sub_models sm on e.equipment_sub_model_unique_code = sm.unique_code").
		Joins("join equipment_models em on sm.equipment_model_unique_code = em.unique_code").
		Joins("join equipment_categories c on em.equipment_category_unique_code = c.unique_code").
		Where("e.deleted_at is null").
		Pluck("e.unique_code", &equipmentUniqueCodes)

	conn2.Table("entire_instances as ei").
		Select("ei.identity_code").
		Joins("join entire_models sm on ei.model_unique_code = sm.unique_code").
		Joins("join entire_models em on sm.parent_unique_code = em.unique_code").
		Joins("join categories c on em.category_unique_code = c.unique_code").
		Where("ei.status = 'FIXING'").
		Where("ei.deleted_at is null").
		Where("sm.is_sub_model is true").
		Where("em.is_sub_model is false").
		Where("sm.deleted_at is null").
		Where("em.deleted_at is null").
		Where("c.deleted_at is null").
		Pluck("ei.identity_code", &entireInstanceIdentityCodes)

	excelWriter := tools.NewExcelWriter("./static/test.xlsx")
	excelSheet := excelWriter.ActiveSheetByIndex(0)
	// 设置表头
	var (
		excelRow  uint64 = 1
		excelRow2 uint64 = 1
	)

	excelSheet.AddRow(
		tools.
			NewExcelRow().
			SetRowNumber(excelRow).
			SetCells([]*tools.ExcelCell{
				tools.NewExcelCellAny("器材编码"),
			}),
	)

	for _, uniqueCode := range equipmentUniqueCodes {
		excelRow++
		excelSheet.AddRow(
			tools.
				NewExcelRow().
				SetRowNumber(excelRow).
				SetCells([]*tools.ExcelCell{
					tools.NewExcelCellAny(uniqueCode),
				}),
		)
		std.EchoLineDebug(fmt.Sprintf("「电子车间」写入器材编码 %s", uniqueCode))
	}

	excelSheet2 := excelWriter.CreateSheet("Sheet 2")
	excelSheet2.AddRow(
		tools.
			NewExcelRow().
			SetRowNumber(excelRow).
			SetCells([]*tools.ExcelCell{
				tools.NewExcelCellAny("器材编码"),
			}),
	)

	for _, identityCode := range entireInstanceIdentityCodes {
		excelRow2++
		excelSheet2.AddRow(
			tools.
				NewExcelRow().
				SetRowNumber(excelRow2).
				SetCells([]*tools.ExcelCell{
					tools.NewExcelCellAny(identityCode),
				}),
		)
		std.EchoLineDebug(fmt.Sprintf("「检修车间」写入器材编码 %s", identityCode))
	}

	excelWriter.Save()

	return []string{}
}

type (
	SyncTasks struct {
		Name                string  `gorm:"type:varchar(50);not null;default:'';comment:名称;"`
		Status              uint8   `gorm:"type:tinyint unsigned;not null;default:2;comment:状态;"`
		ParagraphUniqueCode string  `gorm:"type:char(4);not null;default:'';comment:段编码;"`
		Remark              string  `gorm:"type:varchar(255);not null;default:'';comment:备注;"`
		Project             string  `gorm:"type:varchar(50);default '';not null;comment:项目;"`
		RequestUrl          string  `gorm:"type:varchar(100);default '';not null;comment:请求url;"`
		RequestMethod       string  `gorm:"type:varchar(20);default '';not null;comment:请求method;"`
		RequestContent      *string `gorm:"type:longtext;comment:请求内容;"`
		ResponseContent     *string `gorm:"type:longtext;comment:响应内容;"`
		BatchCode           string  `gorm:"type:char(36);default:'';not null;comment:批次号;"`
	}
	RequestContent struct {
		UpdateEquipments []UpdateEquipment `json:"update_equipments"`
	}

	UpdateEquipment struct {
		Status                       string  `json:"status"`
		SerialNumber                 string  `json:"serial_number"`
		InstalledAt                  string  `json:"installed_at"`
		UpdatedAt                    string  `json:"updated_at"`
		uniqueCode                   string  `json:"unique_code"`
		WorkshopUniqueCode           string  `json:"workshop_unique_code"`
		StationUniqueCode            string  `json:"station_unique_code"`
		CenterUniqueCode             string  `json:"center_unique_code"`
		CrossingUniqueCode           string  `json:"crossing_unique_code"`
		LineUniqueCode               string  `json:"line_unique_code"`
		InstallLocationUniqueCode    string  `json:"install_location_unique_code"`
		OutdoorInstallLocationUuid   string  `json:"outdoor_install_location_uuid"`
		OutdoorInstallLocationExtend *string `json:"outdoor_install_location_extend"`
		SceneWorkAreaUniqueCode      string  `json:"scene_work_area_unique_code"`
	}

	Person struct {
		Name string
		Age  uint64
		Rank uint64
	}
)

// Handle 执行命令
func (receiver TestCmd) Handle(params []string) []string {
	switch params[0] {
	case "uuid":
		return receiver.uuid()
	case "ls":
		return receiver.ls()
	case "redis":
		return receiver.redis()
	case "t":
		return receiver.t()
	default:
		return []string{"没有找到命令"}
	}
}
