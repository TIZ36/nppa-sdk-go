package nppa_test

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tiz36/nppa-sdk-go/domain"
)

// PackNppaPlayerBehaviorCollections 打包玩家行为上报数据
func PackNppaPlayerBehaviorCollections(
	no int,
	sessionId string,
	user domain.IdentityCsvData,
	behaviorType domain.NPPAPlayerBehaviorType,
	certType domain.NPPAPlayerCertType,
	deviceId string,
) domain.NPPAPlayerBehaviorReportCollection {
	return domain.NPPAPlayerBehaviorReportCollection{
		No: no,
		Si: sessionId,
		Bt: int(behaviorType),
		Ot: time.Now().Unix(),
		Ct: int(certType),
		Di: deviceId,
		Pi: user.Pi,
	}
}

// ReadAndParseJsonFile 读取并解析json文件
func ReadAndParseJsonFile(jsonFilePath string, target any) error {
	// load json file into struct domain.NPPAEndpointConfig
	// Open the JSON file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the JSON file into the struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

// ReadCsvAndParse 读取并解析csv文件
func ReadCsvAndParse(path string, offset, line int64) ([]domain.IdentityCsvData, error) {
	var result []domain.IdentityCsvData
	var readLine int64 = 0
	// read csv file
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error:", err)
			return // exit on error
		}
	}(file)

	csvReader := csv.NewReader(bufio.NewReader(file))

	// 跳过offset的行数
	for readLine < offset {
		_, e := csvReader.Read()
		if e != nil {
			fmt.Println("Error:", e)
			return nil, e
		}
		readLine++
	}

	for readLine < line+offset {
		row, e := csvReader.Read()
		if e != nil {
			fmt.Println("Error:", e)
			return nil, e
		}

		// parse csv data
		if len(row) == 6 {
			result = append(result, domain.IdentityCsvData{
				AppId:  row[0],
				AppUid: row[1],
				IdNum:  row[2],
				Name:   row[3],
				Pi:     row[4],
				Status: row[5],
			})
		}

		readLine++
	}

	return result, nil
}
