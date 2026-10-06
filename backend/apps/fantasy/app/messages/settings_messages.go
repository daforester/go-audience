package messages

func SubmissionsEnabledMessage(roomKey string, enabled bool) error {
	return sendToAll(roomKey, "SUBMISSIONS_ENABLED", map[string]interface{}{
		"Enabled": enabled,
	})
}

func PresentingTeamMessage(roomKey string, teamID string) error {
	return sendToAll(roomKey, "PRESENTING_TEAM", map[string]interface{}{
		"TeamID": teamID,
	})
}
