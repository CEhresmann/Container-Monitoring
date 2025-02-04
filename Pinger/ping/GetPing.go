package ping

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type IPStatus struct {
	IP       string `json:"ip"`
	PingTime int    `json:"ping_time"`
	LastOK   string `json:"last_ok"`
}

func GetContainerIPs() ([]string, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка контейнеров: %w", err)
	}

	var ips []string
	containerIDs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, id := range containerIDs {
		cmdInspect := exec.Command("docker", "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", id)
		ip, err := cmdInspect.Output()
		if err == nil && len(ip) > 0 {
			ips = append(ips, strings.TrimSpace(string(ip)))
		}
	}
	return ips, nil
}

func Ping(ip string) int {
	start := time.Now()
	cmd := exec.Command("ping", "-c", "1", ip)
	if err := cmd.Run(); err != nil {
		return -1
	}
	return int(time.Since(start).Milliseconds())
}

func RunPinger(broker string, topic string) {
	for {
		ips, err := GetContainerIPs()
		if err != nil {
			log.Println("Ошибка получения IP контейнеров:", err)
		}
		for _, ip := range ips {
			pingTime := Ping(ip)
			ipStatus := IPStatus{
				IP:       ip,
				PingTime: pingTime,
				LastOK:   time.Now().Format("2006-01-02 15:04:05"),
			}
			ProduceToKafka(ipStatus, topic, broker)
		}
		time.Sleep(15 * time.Second)
	}
}
