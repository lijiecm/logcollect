package store

import (
	"github.com/weihualiu/logcollect/util"
	"log"
	"os"
	"path"
)

func Parse(data []byte) {
	pack, err := NewPackCommon(data)
	if err != nil {
		log.Fatal(err)
	}

	if pack.Type == byte(0x01) {
		//api
		//tag1:项目名称,tag2:项目环境,tag3:接口名称
		//创建存储目录
		basepath := "/var/log/logcollect/api/"
		err = os.MkdirAll(basepath, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if int(pack.TagNum) < 3 {
			log.Fatal("data tag number error!")
		}
		filePath := path.Join(basepath, string(pack.TagList[0].Name), string(pack.TagList[1].Name), string(pack.TagList[2].Name))
		//写入文件
		err = util.AppendToFile(filePath, pack.Body)
		if err != nil {
			log.Fatal(err)
		}

	}

}
