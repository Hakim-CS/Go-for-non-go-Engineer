package main

import "fmt"

// Convert event type into readable message
func formatEvent(e Event) string {
	switch e.Type {

	case "PushEvent":
		return fmt.Sprintf("Pushed %d commits to %s", e.Payload.Size, e.Repo.Name)

	case "IssuesEvent":
		return fmt.Sprintf("Opened a new issue in %s", e.Repo.Name)

	case "WatchEvent":
		return fmt.Sprintf("Starred %s", e.Repo.Name)

	case "ForkEvent":
		return fmt.Sprintf("Forked %s", e.Repo.Name)

	default:
		return fmt.Sprintf("Did some activity in %s", e.Repo.Name)
	}
}
