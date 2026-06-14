package permissioncheck

import (
	"fmt"
	"os"

	"github.com/noedaka/go-config-parser/internal/service"
)

func CheckFilePermissions(path string) *service.Issue {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	mode := info.Mode().Perm()
	if mode&0x07 != 0 {
		return &service.Issue{
			Severity:       service.High,
			Message:        fmt.Sprintf("файл %s имеет излишне широкие права доступа (%04o)", path, mode),
			Recommendation: "Установите права 0600 или 0640, убрав доступ для посторонних",
		}
	}

	return nil
}
