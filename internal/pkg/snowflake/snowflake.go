package snowflake

import (
	"time"

	"github.com/sony/sonyflake"
)

var (
	client    *sonyflake.Sonyflake
	machineId uint16
)

func getMachineID() (uint16, error) {
	return machineId, nil
}

func Init(machineID uint16) {

	machineId = machineID
	startTime, err := time.Parse("2006-01-02 15:04:05", "2021-01-01 00:00:00")
	if err != nil {
		panic(err)
	}

	client = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: startTime,
		MachineID: getMachineID,
	})
}

func GenSnowflakeID() uint64 {
	if client == nil {
		return 0
	}

	id, err := client.NextID()
	if err != nil {
		return 0
	}

	return id
}
