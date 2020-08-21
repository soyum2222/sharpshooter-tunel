package common

import (
	"encoding/json"
	"os"
)

func GenConf(i interface{}) {

	for _, v := range os.Args {

		if v == "genconf" {

			data, err := json.MarshalIndent(i, "", "	")
			if err != nil {
				panic(err)
			}
			file, err := os.Create("config.json")
			if err != nil {
				panic(err)
			}

			_, err = file.Write(data)
			if err != nil {
				panic(err)
			}

			os.Exit(0)

		}
	}
}
