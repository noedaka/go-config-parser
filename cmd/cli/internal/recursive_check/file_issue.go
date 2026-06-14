package recursivecheck

import "github.com/noedaka/go-config-parser/internal/service"

type FileIssue struct {
	Issue    service.Issue
	FilePath string
}
