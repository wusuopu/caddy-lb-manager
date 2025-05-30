package services

import (
	"app/config"
	"app/di"
	"app/models"
	"app/types"
	"app/utils/helper"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"
)

type CaddyfileService struct {
	lastReloadedAt int64
}


func (c *CaddyfileService) generateRouteCaddyfile(route models.Route) (string) {
	handle := "handle"
	path := route.Path
	if route.StripPath && path != "" {
		handle = "handle_path"
		// handle_path is an exact match by default, not a prefix match. You must append a * for a fast prefix match
		if !strings.HasSuffix(path, "*") {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			path += "*"
		}
	} else {
		// 以 `path` 开头的路由
		if strings.HasSuffix(path, "/") {
			path += "*"
		}
	}
	auth := "\n"
	if route.Authentication.Username != "" && route.Authentication.HashedPw != "" {
		auth = fmt.Sprintf("\tbasic_auth {\n\t\t\t%s %s\n\t\t}\n", route.Authentication.Username, route.Authentication.HashedPw)
	}

	// handle_path <path> {
	//     basic_auth {
	//         <username> <hashed_password>
	//     }
	//     reverse_proxy <upstream> {
	//         header_up
	//         header_down
	//     }
	// }
	var headerConfig []string
	if route.HeaderUp != nil {
		var HeaderUp []map[string]interface{}
		json.Unmarshal(route.HeaderUp, &HeaderUp)
		for _, item := range HeaderUp {
			headerConfig = append(headerConfig, fmt.Sprintf("header_up %s \"%s\"", item["key"], item["value"]))
		}
	}
	if route.HeaderDown != nil {
		var HeaderDown []map[string]interface{}
		json.Unmarshal(route.HeaderDown, &HeaderDown)
		for _, item := range HeaderDown {
			headerConfig = append(headerConfig, fmt.Sprintf("header_down %s \"%s\"", item["key"], item["value"]))
		}
	}
	headerContent := "\t\t\t" + strings.Join(headerConfig, "\n\t\t\t") + "\n"

	content := fmt.Sprintf("\t%s %s {\n\t%s\t\treverse_proxy %s {\n%s\t\t}\n\t}", handle, path, auth, route.UpStream.GetAddress(), headerContent)
	return content
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
		if !serverItem.Enable {
			continue
		}

		var routeList []models.Route
		di.Container.DB.Preload("UpStream").Preload("Authentication").Where("server_id = ?", serverItem.ID).Order("sort ASC").Order("id ASC").Find(&routeList)
		var routeContent []string
		for _, routeItem := range routeList {
			if !routeItem.Enable {
				continue
			}
			if routeItem.UpStream.ID == 0 {
				continue
			}

			routeContent = append(routeContent, c.generateRouteCaddyfile(routeItem))
		}

		routes := strings.Join(routeContent, "\n")
		config := fmt.Sprintf("%s {\n%s\n}", serverItem.GetAddress(), routes)
		content = append(content, config)
	}

	globalOptions := ""
	if config.Config.Caddy.TLSEmail != "" {
		globalOptions = fmt.Sprintf("{\n\temail %s\n}\n", config.Config.Caddy.TLSEmail)
	}
  return globalOptions + strings.Join(content, "\n"), nil
}

func (c *CaddyfileService) TouchReloadTime() (bool, error) {
	var setting models.SystemSetting
	results := di.Container.DB.Where("name = ?", "caddyfile_reload_time").First(&setting)

	setting.Value, _ = json.Marshal(time.Now().UTC().UnixMicro())

	if results.RowsAffected == 0 {
		setting.Name = "caddyfile_reload_time"
		results = di.Container.DB.Create(&setting)
	} else {
		results = di.Container.DB.Save(&setting)
	}

	if results.Error != nil {
		return false, results.Error
	}
	c.lastReloadedAt = time.Now().UTC().UnixMicro()
  return true, nil
}

func (c *CaddyfileService) PullConfigAndReload() (bool, error) {
	if c.lastReloadedAt == 0 {
		// 服务启动时直接使用上次保存的配置
		fi, err := os.Stat(config.Config.Caddy.ConfigPath)
		if err == nil && fi.Size() > 0 {
			c.lastReloadedAt = time.Now().UTC().UnixMicro()
			return true, nil
		}
	}

	row := di.Container.DB.Model(&models.SystemSetting{}).Where("name = ?", "caddyfile_reload_time").Select("value").Row()
	var value int64
	err := row.Scan(&value)
	if err != nil {
		return false, fmt.Errorf("caddyfile_reload_time not found\n")
	}

	if c.lastReloadedAt >= value {
		return true, nil
	}

	content, err := di.Service.CaddyfileService.GenCaddyfile()
	if err != nil {
		return false, err
	}

	ret, err := di.Service.CaddyfileService.Reload(content)
	if ret != true {
		return false, err
	}
	c.lastReloadedAt = time.Now().UTC().UnixMicro()
	fmt.Printf("pull caddyfile reload success\n")

  return true, nil
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