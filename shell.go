package main

import (
	"github.com/Mrs4s/argv"
	httpDownloader "github.com/Mrs4s/go-http-downloader"
	"github.com/Mrs4s/six-cli/models"
	"github.com/Mrs4s/six-cli/six_cloud"
	"strconv"
	"time"

	"github.com/Mrs4s/ishell"
	"github.com/fatih/color"
	"os"
	"reflect"
	"strings"
)

type handler struct {
}

var (
	shell *ishell.Shell
	helps map[string]string
)

func init() {
	// 介绍 | 别名
	helps = map[string]string{
		"Login":     "登录6 Pan",
		"ListFiles": "列出当前目录文件列表 | ls",
		"JoinPath":  "进入目录 | cd",
		//"CopyFile":  "复制文件或目录 | cp",
		"Download": "下载文件或文件夹 | down",
	}
}

func runAsShell() {
	shell = ishell.New()
	shell.SetHistoryPath("history.shell")
	//shell.CustomCompleter(&readline.TabCompleter{})
	shell.Println("欢迎使用6 Pan命令行客户端!")
	shell.Println("开发By Mrs4s")
	shell.Println("")
	shell.SetPrompt("guest@six-pan:/$ ")
	// delete default commands
	shell.DeleteCmd("exit")
	shell.DeleteCmd("clear")
	shell.DeleteCmd("help")
	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Func: func(c *ishell.Context) {
			c.Print(c.HelpText())
		},
		Help: "获取帮助信息",
	})
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count < 2 {
			c.Println("再次键入 Ctrl+C 以确认退出")
			return
		}
		models.DefaultConf.SaveFile("config.json")
		os.Exit(0)
	})
	initCommands()
	if currentUser != nil {
		shell.Println("自动登录完成: " + currentUser.Username)
		refreshPrompt()
	}
	shell.Run()
}

func initCommands() {
	h := &handler{}
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		f, ok := v.Method(i).Interface().(func(*ishell.Context, *argv.Args))
		if ok {
			help := helps[m.Name]
			name := m.Name
			if strings.Contains(help, "|") {
				name = strings.Split(help, "|")[1]
				help = strings.Split(help, "|")[0]
			}
			var compFunc func([]string) []string
			comp := v.MethodByName(m.Name + "Completer")
			if comp.IsValid() {
				compFunc = comp.Interface().(func([]string) []string)
			}
			shell.AddCmd(&ishell.Cmd{
				Name: strings.TrimSpace(strings.ToLower(name)),
				Func: func(c *ishell.Context) {
					args := argv.Parse(append([]string{name}, c.Args...))
					f(c, &args)
				},
				Help:      strings.TrimSpace(help),
				Completer: compFunc,
			})
		}
	}
}

func (handler) Login(c *ishell.Context, args *argv.Args) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)
	c.Print("请输入用户名: ")
	username := c.ReadLine()
	c.Print("请输入密码: ")
	password := c.ReadPassword()
	user, err := six_cloud.LoginWithUsernameOrPhone(username, password)
	if err != nil {
		c.Println("登录失败: " + err.Error())
		return
	}
	currentUser = user
	currentPath = "/"
	models.DefaultConf.QingzhenToken = currentUser.Client.QingzhenToken
	c.Println("登录完成, 欢迎: " + currentUser.Username)
	refreshPrompt()
}

func (handler) ListFiles(c *ishell.Context, args *argv.Args) {
	var (
		blue       = color.New(color.FgBlue).SprintFunc()
		green      = color.New(color.FgGreen).SprintFunc()
		files, err = currentUser.GetFilesByPath(currentPath)
	)
	if len(args.Nokeys) >= 1 {
		tar := strings.Join(args.Nokeys, " ")
		if len(tar) == 0 || tar[0:1] != "/" {
			c.Println("非法操作")
			return
		}
		files, err = currentUser.GetFilesByPath(tar)
	}
	if err != nil {
		c.Println("获取文件失败: " + err.Error())
		return
	}
	if len(files) == 0 {
		c.Println("当前目录下无任何文件")
		return
	}
	printColor := func(files []*six_cloud.SixFile) (strs []string) {
		for _, file := range files {
			if file.IsDir {
				strs = append(strs, blue(file.Name))
				continue
			}
			strs = append(strs, green(file.Name))
		}
		return
	}
	filteredFiles := files
	if _, ok := args.Keys["R"]; ok {
		c.Println(".:")
		PrintColumns(printColor(filteredFiles), 2)
		for _, dir := range filterDirs(files) {
			subFiles, err := currentUser.GetFilesByPath(models.CombinePaths(currentPath, dir, "/"))
			if err == nil {
				c.Println("\n./" + dir + ":")
				PrintColumns(printColor(append(subFiles)), 2)
			}
		}
		return
	}
	if _, ok := args.Keys["d"]; ok {
		var dirs []*six_cloud.SixFile
		for _, file := range files {
			if file.IsDir {
				dirs = append(dirs, file)
			}

		}
		filteredFiles = dirs
	}

	if _, ok := args.Keys["a"]; ok {
		table := [][]string{
			{"序号", "创建时间", "文件大小", "文件名"},
		}
		for i, file := range filteredFiles {
			table = append(table, []string{
				strconv.FormatInt(int64(i), 10),
				time.Unix(file.CreateTime/1000, 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
				models.ConvertSizeString(file.Size),
				file.Name,
			})
		}
		PrintTables(table, 2)
		return
	}

	PrintColumns(printColor(filteredFiles), 2)
}

/*
func (handler) CopyFile(c *ishell.Context, args *argv.Args) {

}
*/

func (handler) Download(c *ishell.Context, args *argv.Args) {
	if len(args.Nokeys) != 1 || len(args.Nokeys[0]) == 0 {
		c.Println("使用方法: down <文件/目录>")
		return
	}
	targetPath := currentPath
	if args.Nokeys[0][0:1] == "/" {
		targetPath = models.GetParentPath(args.Nokeys[0])
	}
	files, err := currentUser.GetFilesByPath(targetPath)
	if err != nil {
		c.Println("错误:", err)
		return
	}
	var target *six_cloud.SixFile
	for _, file := range files {
		if file.Name == models.GetFileName(args.Nokeys[0]) {
			target = file
			continue
		}
	}
	if target == nil {
		c.Println("错误: 目标文件/目录不存在")
		return
	}
	for key, file := range target.GetLocalTree(models.DefaultConf.DownloadPath) {
		c.Println("添加下载", file.Path, "...")
		addr, err := file.GetDownloadAddress()
		if err != nil {
			c.Println("获取文件", file.Name, "的下载链接失败:", err)
			continue
		}
		info, err := httpDownloader.NewDownloaderInfo([]string{addr}, key, models.DefaultConf.DownloadBlockSize, int(models.DefaultConf.DownloadThread),
			map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:67.0) Gecko/20100101 Firefox/67.0"})
		models.DownloaderList = append(models.DownloaderList, httpDownloader.NewClient(info))
	}
}

func (handler) DownloadCompleter([]string) []string {
	return filterCurrentFiles()
}

func (handler) JoinPath(c *ishell.Context, args *argv.Args) {
	if len(c.Args) == 0 {
		c.Println("使用方法: cd <目录>")
	}
	arg := strings.Join(c.Args, " ")
	defer refreshPrompt()
	switch {
	case len(arg) == 0:
	case arg[0:1] == "/":
		_, err := currentUser.GetFilesByPath(arg)
		if err != nil {
			c.Println("切换失败: " + err.Error())
			return
		}
		currentPath = arg
	case arg == "..":
		if currentPath == "/" {
			return
		}
		currentPath = models.GetParentPath(currentPath)
	case strings.Contains(arg, "../"):
		if currentPath == "/" {
			return
		}
		for i := 0; i < strings.Count(arg, "../"); i++ {
			currentPath = models.GetParentPath(currentPath)
		}
	default:
		newPath := models.CombinePaths(currentPath, arg, "/")
		_, err := currentUser.GetFilesByPath(newPath)
		if err != nil {
			c.Println("切换失败: " + err.Error())
			return
		}
		currentPath = newPath
	}
}

func (handler) JoinPathCompleter([]string) []string {
	return filterCurrentDirs()
}

func refreshPrompt() {
	shell.SetPrompt(currentUser.Username + "@six-pan:" + currentPath + "$ ")
}

func filterCurrentFiles() []string {
	files, err := currentUser.GetFilesByPath(currentPath)
	if err != nil {
		return []string{}
	}
	return filterFiles(files)
}

func filterCurrentDirs() []string {
	files, err := currentUser.GetFilesByPath(currentPath)
	if err != nil {
		return []string{}
	}
	return filterDirs(files)
}

func filterFiles(files []*six_cloud.SixFile) (res []string) {
	for _, file := range files {
		if !file.IsDir {
			res = append(res, file.Name)
		}
	}
	return
}

func filterDirs(files []*six_cloud.SixFile) (res []string) {
	for _, file := range files {
		if file.IsDir {
			res = append(res, file.Name)
		}
	}
	return
}
