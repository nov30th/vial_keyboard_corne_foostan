package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/ncruces/zenity"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

//go:embed keyboard_conf_mapping.txt
var mappingFS embed.FS

// KeyboardMapping 存储有线和无线键盘之间的映射关系
type KeyboardMapping struct {
	WiredToWireless map[int]int
	WirelessToWired map[int]int
}

// 键盘配置结构
type KeyboardConfig struct {
	Version       int                      `json:"version"`
	UID           int64                    `json:"uid"`
	Layout        [][][]interface{}        `json:"layout"`
	EncoderLayout [][]interface{}          `json:"encoder_layout"`
	LayoutOptions interface{}              `json:"layout_options"`
	Macro         [][]interface{}          `json:"macro"`
	VialProtocol  int                      `json:"vial_protocol"`
	ViaProtocol   int                      `json:"via_protocol"`
	TapDance      [][]interface{}          `json:"tap_dance"`
	Combo         [][]interface{}          `json:"combo"`
	KeyOverride   []map[string]interface{} `json:"key_override"`
	Settings      map[string]interface{}   `json:"settings"`
}

// 固定的空位位置 (-1 值)
var (
	// 无线键盘的空位位置
	wirelessEmptyPositions = [][]int{
		{2, 6},                         // Row 3, Column 7
		{3, 0}, {3, 1}, {3, 2}, {3, 6}, // Row 4 empty positions
		{6, 6},                         // Row 7, Column 7
		{7, 0}, {7, 1}, {7, 2}, {7, 6}, // Row 8 empty positions
	}

	// 有线键盘的空位位置
	wiredEmptyPositions = [][]int{
		{2, 2}, // Row 3, Column 3
		{3, 2}, // Row 4, Column 3
	}
)

func main() {
	// 显示欢迎信息
	fmt.Println("键盘配置转换器 - 在有线和无线键盘配置之间转换")
	fmt.Println("=============================================")

	// 加载映射文件
	mapping, err := loadMapping()
	if err != nil {
		log.Fatalf("加载映射文件失败: %v", err)
	}

	// 选择源配置文件
	zenity.Info("请选择源配置文件，这是您希望转换其布局的文件。\n\n例如：如果要将无线键盘配置转换为有线格式，请选择无线键盘配置文件。",
		zenity.Title("选择源文件"))

	sourceFile, err := zenity.SelectFile(
		zenity.Title("选择源配置文件（包含要转换的布局）"),
		zenity.FileFilter{
			Name:     "键盘配置文件",
			Patterns: []string{"*.vil"},
		},
	)
	if err != nil {
		log.Fatalf("选择源文件失败: %v", err)
	}

	// 选择目标配置文件
	zenity.Info("请选择目标配置文件，这将提供输出文件的结构格式。\n\n例如：如果要将无线键盘配置转换为有线格式，请选择有线键盘配置文件。",
		zenity.Title("选择目标文件"))

	targetFile, err := zenity.SelectFile(
		zenity.Title("选择目标配置文件（提供目标结构）"),
		zenity.FileFilter{
			Name:     "键盘配置文件",
			Patterns: []string{"*.vil"},
		},
	)
	if err != nil {
		log.Fatalf("选择目标文件失败: %v", err)
	}

	// 加载源配置
	sourceConfig, err := loadConfig(sourceFile)
	if err != nil {
		log.Fatalf("加载源配置失败: %v", err)
	}

	// 加载目标配置
	targetConfig, err := loadConfig(targetFile)
	if err != nil {
		log.Fatalf("加载目标配置失败: %v", err)
	}

	// 确定键盘类型
	sourceType := determineKeyboardType(sourceConfig)
	targetType := determineKeyboardType(targetConfig)

	if sourceType == "unknown" || targetType == "unknown" {
		zenity.Error("无法确定键盘类型，请确保选择了正确的配置文件",
			zenity.Title("错误"))
		return
	}

	if sourceType == targetType {
		zenity.Error("源文件和目标文件是同一种键盘类型，无需转换",
			zenity.Title("错误"))
		return
	}

	var newLayout [][][]interface{}
	var message string

	// 根据键盘类型转换布局
	if sourceType == "wireless" && targetType == "wired" {
		// 将无线布局转换为有线布局
		newLayout = convertWirelessToWired(sourceConfig.Layout, mapping)
		message = "无线布局已转换为有线布局"
	} else { // sourceType == "wired" && targetType == "wireless"
		// 将有线布局转换为无线布局
		newLayout = convertWiredToWireless(sourceConfig.Layout, mapping)
		message = "有线布局已转换为无线布局"
	}

	// 确保布局层数与目标配置匹配
	adjustLayoutLayers(&newLayout, targetConfig.Layout)

	// 应用转换后的布局到目标配置
	targetConfig.Layout = newLayout

	// 保存新配置
	outputFile := createOutputFileName(sourceFile, targetFile)
	err = saveConfig(targetConfig, outputFile)
	if err != nil {
		zenity.Error(fmt.Sprintf("保存配置失败: %v", err),
			zenity.Title("错误"))
		return
	}

	zenity.Info(fmt.Sprintf("%s\n输出文件: %s", message, outputFile),
		zenity.Title("成功"))
}

// 加载键盘映射
func loadMapping() (KeyboardMapping, error) {
	mapping := KeyboardMapping{
		WiredToWireless: make(map[int]int),
		WirelessToWired: make(map[int]int),
	}

	// 从嵌入的文件系统加载映射文件
	content, err := mappingFS.ReadFile("keyboard_conf_mapping.txt")
	if err != nil {
		return mapping, fmt.Errorf("读取映射文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			var wiredIndex, wirelessIndex int
			fmt.Sscanf(parts[0], "%d", &wiredIndex)
			fmt.Sscanf(parts[1], "%d", &wirelessIndex)

			mapping.WiredToWireless[wiredIndex] = wirelessIndex
			mapping.WirelessToWired[wirelessIndex] = wiredIndex
		}
	}

	return mapping, nil
}

// 加载键盘配置文件
func loadConfig(filePath string) (KeyboardConfig, error) {
	var config KeyboardConfig

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// 保存键盘配置文件
func saveConfig(config KeyboardConfig, filePath string) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// 确定键盘类型
func determineKeyboardType(config KeyboardConfig) string {
	if len(config.Layout) >= 1 {
		firstArray := config.Layout[0]
		// 无线键盘有 8 行，每行 7 列
		if len(firstArray) == 8 && allRowsHaveLength(firstArray, 7) {
			return "wireless"
		}
		// 有线键盘有 4 行，每行 12 列
		if len(firstArray) == 4 && allRowsHaveLength(firstArray, 12) {
			return "wired"
		}
	}
	return "unknown"
}

// 检查所有行是否有指定的长度
func allRowsHaveLength(rows [][]interface{}, length int) bool {
	for _, row := range rows {
		if len(row) != length {
			return false
		}
	}
	return true
}

// 将无线键盘布局转换为有线布局
func convertWirelessToWired(wirelessLayout [][][]interface{}, mapping KeyboardMapping) [][][]interface{} {
	var wiredLayout [][][]interface{}

	for _, wirelessLayer := range wirelessLayout {
		wiredLayer := make([][]interface{}, 4)
		for i := range wiredLayer {
			wiredLayer[i] = make([]interface{}, 12)
			for j := range wiredLayer[i] {
				wiredLayer[i][j] = "KC_NO" // 默认值
			}
		}

		// 将无线布局映射到有线布局
		for rowIdx, row := range wirelessLayer {
			for colIdx, key := range row {
				// 计算无线数组中的线性索引
				wirelessLinearIdx := rowIdx*7 + colIdx

				// 检查此无线索引是否有映射
				if wiredLinearIdx, ok := mapping.WirelessToWired[wirelessLinearIdx]; ok {
					wiredRowIdx := wiredLinearIdx / 12
					wiredColIdx := wiredLinearIdx % 12

					// 仅在索引在范围内时分配
					if wiredRowIdx >= 0 && wiredRowIdx < 4 && wiredColIdx >= 0 && wiredColIdx < 12 {
						wiredLayer[wiredRowIdx][wiredColIdx] = key
					}
				}
			}
		}

		// 设置固定的 -1 值在空位置
		for _, pos := range wiredEmptyPositions {
			rowIdx, colIdx := pos[0], pos[1]
			if rowIdx < len(wiredLayer) && colIdx < len(wiredLayer[0]) {
				wiredLayer[rowIdx][colIdx] = float64(-1) // JSON 中的 -1 会解析为 float64
			}
		}

		wiredLayout = append(wiredLayout, wiredLayer)
	}

	return wiredLayout
}

// 将有线键盘布局转换为无线布局
func convertWiredToWireless(wiredLayout [][][]interface{}, mapping KeyboardMapping) [][][]interface{} {
	var wirelessLayout [][][]interface{}

	for _, wiredLayer := range wiredLayout {
		wirelessLayer := make([][]interface{}, 8)
		for i := range wirelessLayer {
			wirelessLayer[i] = make([]interface{}, 7)
			for j := range wirelessLayer[i] {
				wirelessLayer[i][j] = "KC_NO" // 默认值
			}
		}

		// 将有线布局映射到无线布局
		for rowIdx, row := range wiredLayer {
			for colIdx, key := range row {
				// 计算有线数组中的线性索引
				wiredLinearIdx := rowIdx*12 + colIdx

				// 检查此有线索引是否有映射
				if wirelessLinearIdx, ok := mapping.WiredToWireless[wiredLinearIdx]; ok {
					wirelessRowIdx := wirelessLinearIdx / 7
					wirelessColIdx := wirelessLinearIdx % 7

					// 仅在索引在范围内时分配
					if wirelessRowIdx >= 0 && wirelessRowIdx < 8 && wirelessColIdx >= 0 && wirelessColIdx < 7 {
						wirelessLayer[wirelessRowIdx][wirelessColIdx] = key
					}
				}
			}
		}

		// 设置固定的 -1 值在空位置
		for _, pos := range wirelessEmptyPositions {
			rowIdx, colIdx := pos[0], pos[1]
			if rowIdx < len(wirelessLayer) && colIdx < len(wirelessLayer[0]) {
				wirelessLayer[rowIdx][colIdx] = float64(-1) // JSON 中的 -1 会解析为 float64
			}
		}

		wirelessLayout = append(wirelessLayout, wirelessLayer)
	}

	return wirelessLayout
}

// 调整布局层数以匹配目标配置
func adjustLayoutLayers(newLayout *[][][]interface{}, targetLayout [][][]interface{}) {
	targetLayerCount := len(targetLayout)

	// 如果新布局层数过多，截断
	if len(*newLayout) > targetLayerCount {
		*newLayout = (*newLayout)[:targetLayerCount]
	} else if len(*newLayout) < targetLayerCount {
		// 如果新布局层数不足，从目标布局添加层
		for i := len(*newLayout); i < targetLayerCount; i++ {
			(*newLayout) = append((*newLayout), targetLayout[i])
		}
	}
}

// 创建输出文件名
func createOutputFileName(sourceFile, targetFile string) string {
	sourceBase := filepath.Base(sourceFile)
	targetBase := filepath.Base(targetFile)
	sourceExt := filepath.Ext(sourceFile)
	sourceNameWithoutExt := strings.TrimSuffix(sourceBase, sourceExt)

	return filepath.Join(filepath.Dir(sourceFile), sourceNameWithoutExt+"_converted_to_"+targetBase)
}
