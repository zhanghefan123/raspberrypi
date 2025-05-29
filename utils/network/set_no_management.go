package network

import (
	"fmt"
	"raspberrypi/utils/execute"
	"strings"
)

// SetNoManagement 将接口的模式设置为非节能模式
func SetNoManagement(InterfaceName string) error {
	err := execute.Command("nmcli", strings.Split(fmt.Sprintf("dev set %s managed no", InterfaceName), " "))
	if err != nil {
		return fmt.Errorf("failed to set %s to no efficient mode: %v", InterfaceName, err)
	}
	return nil
}
