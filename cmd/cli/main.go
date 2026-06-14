package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/noedaka/go-config-parser/cmd/cli/internal/config"
	permissioncheck "github.com/noedaka/go-config-parser/cmd/cli/internal/permission_check"
	recursivecheck "github.com/noedaka/go-config-parser/cmd/cli/internal/recursive_check"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
	"github.com/noedaka/go-config-parser/internal/service/rules"
)

func main() {
	cfg := config.Init()

	if cfg.IsSilent {
		log.Println("Активен флаг -s.")
		return
	}

	var targets []string
	if cfg.IsStdin {
		targets = []string{""}
	} else if cfg.IsRecursive {
		if flag.NArg() == 0 {
			log.Fatalf("Укажите директорию для рекурсивного анализа")
		}

		dir := flag.Arg(0)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			log.Fatalf("'%s' не является директорией", dir)
		}

		files, err := recursivecheck.WalkDir(dir)
		if err != nil {
			log.Fatalf("Ошибка обхода директории: %v", err)
		}

		if len(files) == 0 {
			log.Println("Не найдено конфигурационных файлов (.yaml/.json)")
			os.Exit(0)
		}

		targets = files
	} else {
		if flag.NArg() == 0 {
			log.Fatalf("Путь к файлу не указан")
		}

		targets = []string{flag.Arg(0)}
	}

	rulesList := []service.Rule{
		rules.DebugLogRule{},
		rules.PlaintextPasswordRule{},
		rules.ZeroHostRule{},
		rules.TLSDisabledRule{},
		rules.NewWeakAlgorithmRule(),
	}

	var allFileIssues []recursivecheck.FileIssue

	for _, target := range targets {
		var filePath string
		var input []byte
		var err error

		if target == "" { // stdin
			filePath = "<stdin>"
			input, err = io.ReadAll(os.Stdin)
		} else {
			filePath = target

			if permIssue := permissioncheck.CheckFilePermissions(target); permIssue != nil {
				allFileIssues = append(allFileIssues, recursivecheck.FileIssue{
					Issue:    *permIssue,
					FilePath: filePath,
				})
			}

			input, err = os.ReadFile(target)
		}

		if err != nil {
			allFileIssues = append(allFileIssues, recursivecheck.FileIssue{
				Issue: service.Issue{
					Severity:       service.High,
					Message:        "Не удалось прочитать файл: " + err.Error(),
					Recommendation: "Проверьте существование файла и права доступа",
				},
				FilePath: filePath,
			})
			continue
		}

		p := parser.Parser(parser.YamlJsonParser{})
		data, err := p.ParseConfig(input)
		if err != nil {
			allFileIssues = append(allFileIssues, recursivecheck.FileIssue{
				Issue: service.Issue{
					Severity:       service.High,
					Message:        "Невалидный YAML/JSON: " + err.Error(),
					Recommendation: "Исправьте синтаксис конфигурационного файла",
				},
				FilePath: filePath,
			})
			continue
		}

		for _, rule := range rulesList {
			for _, issue := range rule.Check(data) {
				allFileIssues = append(allFileIssues, recursivecheck.FileIssue{
					Issue:    issue,
					FilePath: filePath,
				})
			}
		}
	}

	for _, fi := range allFileIssues {
		log.Printf("[%s] %s\nФайл: %s\nРекомендация: %s\n\n",
			fi.Issue.Severity, fi.Issue.Message, fi.FilePath, fi.Issue.Recommendation)
	}

	if len(allFileIssues) > 0 && !cfg.IsSilent {
		log.Fatalf("Выходим с ошибкой, потому что найдены проблемы в конфигурации.")
	}

	log.Println("Проверка завершена, проблем не найдено.")
}
