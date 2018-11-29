package Config

import "fmt"
import "os"
import "io"
import "strings"
import "bufio"

type Config struct {
	conMap map[string]string
	sArea  string
}

func (con *Config) InitConfig(conPath string) bool {
	con.conMap = make(map[string]string)

	pFile, err := os.Open(conPath)
	if err != nil {
		fmt.Println("Can't Open File", conPath, err)
		return false
	}

	defer pFile.Close()

	pRead := bufio.NewReader(pFile)

	for {
		pData, _, err := pRead.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Can't Read File", conPath, err)
			return false
		}

		sLine := strings.TrimSpace(string(pData[:]))

		/* Find first letter is #, skip */
		if strings.Index(sLine, "#") == 0 {
			continue
		}

		/* if has title [], means title, get title and continue */
		nBegin := strings.Index(sLine, "[")
		nEnd := strings.LastIndex(sLine, "]")
		if nBegin >= 0 && nEnd > 0 && nEnd > nBegin+1 {
			con.sArea = strings.TrimSpace(sLine[nBegin+1 : nEnd])
			continue
		}

		/* is sub string is nil, skip */
		if len(con.sArea) == 0 {
			continue
		}

		/* Don't found = , skip too */
		nIndex := strings.Index(sLine, "=")
		if nIndex < 0 {
			continue
		}

		sFrist := strings.TrimSpace(sLine[:nIndex])
		if len(sFrist) == 0 {
			continue
		}

		sSecond := strings.TrimSpace(sLine[nIndex+1:])

		nPos := strings.Index(sSecond, "\t#")
		if nPos >= 0 {
			sSecond = sSecond[0:nPos]
		}

		nPos = strings.Index(sSecond, " #")
		if nPos >= 0 {
			sSecond = sSecond[0:nPos]
		}

		nPos = strings.Index(sSecond, "\t\\")
		if nPos >= 0 {
			sSecond = sSecond[0:nPos]
		}

		nPos = strings.Index(sSecond, " \\")
		if nPos >= 0 {
			sSecond = sSecond[0:nPos]
		}

		if len(sSecond) == 0 {
			continue
		}

		sKey := con.sArea + "." + sFrist
		con.conMap[sKey] = strings.TrimSpace(sSecond)
	}

	return true
}

func (con Config) Read(sNode, sKey string) string {
	sKey = sNode + "." + sKey
	sValue, ok := con.conMap[sKey]
	if !ok {
		return ""
	}
	return sValue
}
