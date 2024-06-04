package common

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

var Templates map[string]*template.Template

func InitTemplate() {
	// init
	// load template file
	Templates = make(map[string]*template.Template)
	// register template
	if ChatTemplateDir == "" {
		Logger.Warnf("InitTemplate empty TemplateDir")
		return
	}
	err := filepath.Walk(ChatTemplateDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}
		var tempContent map[string]string
		fileBody, err := os.ReadFile(path)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("load failed: %s", info.Name()))
			return err
		}
		err = yaml.Unmarshal(fileBody, &tempContent)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("load failed: %s", info.Name()))
			return err
		}

		for k, v := range tempContent {
			if strings.TrimSpace(v) == "" {
				continue
			}
			subTemp, err := template.New(info.Name() + k).Parse(v)
			if err != nil {
				return err
			}
			Templates[k] = subTemp
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
