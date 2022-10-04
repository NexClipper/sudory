package control

import "github.com/NexClipper/sudory/pkg/server/macro"

func genUuidString(uuid string) string {
	if 0 < len(uuid) {
		return uuid
	}
	return macro.NewUuidString()
}
