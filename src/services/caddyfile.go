package services

import (
	"app/config"
	"app/di"
	"app/models"
	"app/types"
	"app/utils/helper"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"
)

type CaddyfileService struct {
}

// Generate Caddyfile content from database
func (c *CaddyfileService) GenCaddyfile() (string, error) {
	var serverList []models.Server
	results := di.Container.DB.Find(&serverList)

	if results.RowsAffected == 0 {
		return "", fmt.Errorf("There is no server records.")
	}

	var content []string
	for _, serverItem := range serverList {
		routes := ""
		config := fmt.Sprintf("%s {\n\t%s\n}", serverItem.GetAddress(), routes)
		content = append(content, config)
	}
  return strings.Join(content, "\n"), nil
}

// caddy service reload config
func (c *CaddyfileService) ReloadFile(file string) (bool, error) {
	ret, err := c.Validate(file)
	if ret != true {
		return false, err
	}

	if config.Config.Caddy.ReloadCMD == "" {
		cmd := exec.Command(config.Config.Caddy.BinPath, "reload", "--config", file)
		output, err := cmd.CombinedOutput()
		if err != nil {
			results := strings.Split(strings.TrimSpace(string(output)), "\n")
			return false, fmt.Errorf("%s", results[len(results) - 1])
		}
		return true, nil
	} else {
		args := strings.Split(config.Config.Caddy.ReloadCMD, " ")
		args = append(args, file)
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return false, fmt.Errorf("%s", output)
		}
		return true, nil
	}
}

func (c *CaddyfileService) Reload(content string) (bool, error) {
	f, err := os.CreateTemp("", "Caddyfile")
	if err != nil {
		return false, err
	}
	f.WriteString(content)

	defer os.Remove(f.Name())

  ret, err := c.ReloadFile(f.Name())
	return ret, err
}

func (c *CaddyfileService) Validate(file string) (bool, error) {
	cmd := exec.Command(config.Config.Caddy.BinPath, "validate", "--config", file)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false, err
	}

	if err := cmd.Start(); err != nil {
		return false, err
	}

	output, _ := io.ReadAll(stderr)
	err = cmd.Wait()

	if err != nil {
		results := strings.Split(strings.TrimSpace(string(output)), "\n")
		return false, fmt.Errorf("%s", results[len(results) - 1])
	}
	return true, nil
}

func (c *CaddyfileService) HashPassword (password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
func (c *CaddyfileService) ListCertificate() []types.CertificateInfo {
  var data []types.CertificateInfo

	files, err := filepath.Glob(filepath.Join(config.Config.Caddy.DataPath + "/certificates/**/**/*.json"))
	if err != nil {
		return data
	}
	for _, f := range files {
		basename := filepath.Base(f)
		name := basename[0 : len(basename) - 5]

		var p fastjson.Parser
		rawData, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		item, err := p.ParseBytes(rawData)
		if err != nil {
			continue
		}
		var parser helper.JSONParser

		var sans []string
		_sans := item.GetArray("sans")
		for _, s := range _sans {
			parser.Value = s
			v, _ := parser.GetJSONString("")
			sans = append(sans, v)
		}
		parser.Value = item

		start, _ := parser.GetJSONString("issuer_data.renewal_info._retryAfter")
		end, _ := parser.GetJSONString("issuer_data.renewal_info.suggestedWindow.end")
		data = append(data, types.CertificateInfo{
			Name: name,
			Sans: sans,
			ValidityStart: start,
			ValidityEnd: end,

		})
	}
	return data
}