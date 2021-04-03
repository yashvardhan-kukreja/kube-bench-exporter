package main

import (
	"fmt"
	"github.com/yashvardhan-kukreja/kube-bench-exporter/pkg/apis/aws"
	"github.com/yashvardhan-kukreja/kube-bench-exporter/pkg/global"
	"sync"
)

var targetDeserializers map[string]func(map[string]interface{}) (global.Target, error)
var wg sync.WaitGroup

func init() {
	targetDeserializers = map[string](func(map[string]interface{}) (global.Target, error)){
		"s3": aws.DeserializeInputJsonToS3Config,
	}
}

func main() {

	path := "/etc/config/target-config.json"
	inputConfigs, err := global.DecodeConfigFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	var targets []global.Target
	for _, input := range inputConfigs {
		destinationType := input["type"].(string)
		destinationConfig := input["config"].(map[string]interface{})
		destinationDeserializer := targetDeserializers[destinationType]

		currentTargetConfig, err := destinationDeserializer(destinationConfig)
		if err != nil {
			fmt.Println(err)
			return
		}
		targets = append(targets, currentTargetConfig)
	}

	wg.Add(len(targets))
	for _, target := range targets {
		go func(target global.Target, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := target.Export(); err != nil {
				fmt.Printf("\nerror occurred while exporting to the target (%+v) : %+v", target, err)
				return
			}
			fmt.Printf("\nSuccessfully exported the kube-bench report to the target %+v", target)
		}(target, &wg)
	}
	wg.Wait()
}
